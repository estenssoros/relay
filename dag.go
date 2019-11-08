package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/estenssoros/dasorm/nulls"
	"github.com/estenssoros/relay/db"
	"github.com/estenssoros/relay/models"
	"github.com/estenssoros/relay/state"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DAG Directed Acylcical Graph
type DAG struct {
	ID                   string
	Description          string
	ScheduleInterval     string
	StartDate            time.Time
	EndDate              time.Time
	DefaultArgs          []interface{}
	Params               map[string]interface{}
	Concurrency          int
	MaxActiveRuns        int
	DagRunTimeout        time.Duration
	DefaultView          string
	Orientation          string
	Catchup              bool
	OnFailureCallBack    func() error
	OnSuccessCallBack    func() error
	AccessControl        map[string]string
	IsPausedUponCreation bool

	tasks map[string]TaskInterface
}

// FormattedID formatted dag id
func (d DAG) FormattedID() string {
	return fmt.Sprintf("DAG[%s]", d.ID)
}

func (d DAG) String() string {
	ju, _ := json.Marshal(d)
	return string(ju)
}

// Roots are the first tasks that need tp be run
func (d *DAG) Roots() []TaskInterface {
	tasks := []TaskInterface{}
	for _, t := range d.tasks {
		if t.IsRoot() {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func printSeparator(sep string, num int) {
	fmt.Println(strings.Repeat(sep, num))
}

// Run runs a dag
func (d *DAG) Run(ctx context.Context, dagRun *models.DagRun) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	start := time.Now()

	runner := NewTaskRunner(d.tasks)
	if err := runner.Check(); err != nil {
		return errors.Wrap(err, "runner check")
	}
	var w sync.WaitGroup
	w.Add(1)

	go runner.Run(ctx, &w)

	dagRun.UpdateState(state.Running)

	for _, task := range d.tasks {
		if task.IsRoot() {
			task.SetState(state.Queued)
		}
		taskModel := &models.TaskInstance{
			TaskID:    task.GetID(),
			DagRunID:  dagRun.ID,
			StartDate: time.Now().UTC(),
			EndDate:   nulls.Time{},
			State:     task.GetState(),
			Operator:  task.OperatorType(),
		}
		if err := taskModel.Create(); err != nil {
			return errors.Wrap(err, "create task model")
		}
		task.SetModel(taskModel)
		runner.evalQueue <- task
	}

	w.Wait()

	dagRun.Finish(runner.FinalState())

	logrus.Infof("dag took %v", time.Since(start))
	return nil
}

// DagConfig basic config for a new dag
type DagConfig struct {
	ID               string
	Description      string
	ScheduleInterval string
}

// NewDag creats a new dag
func NewDag(input *DagConfig) (*DAG, error) {
	return &DAG{
		ID:               input.ID,
		Description:      input.Description,
		ScheduleInterval: input.ScheduleInterval,
		tasks:            map[string]TaskInterface{},
	}, nil
}

// AddTask sets the task state to pending and adds a task to a dag
func (d *DAG) AddTask(t TaskInterface) error {
	_, ok := d.tasks[t.GetID()]
	if ok {
		return errors.Errorf("task %s already exists in dag", t.GetID())
	}
	t.SetState(state.Pending)
	d.tasks[t.GetID()] = t
	t.SetDag(d)
	return nil
}

// NewBash creates a new bash operators on a dag
func (d *DAG) NewBash(o *BashOperator) (*BashOperator, error) {
	return o, d.AddTask(o)
}

// NewGo creates a new go operator on a dag
func (d *DAG) NewGo(o *GoOperator) (*GoOperator, error) {
	return o, d.AddTask(o)
}

func (d *DAG) getTask(taskID string) (TaskInterface, error) {
	t, ok := d.tasks[taskID]
	if !ok {
		return nil, errors.Errorf("missing task %s", taskID)
	}
	return t, nil
}

func getUpstream(task TaskInterface, level int) error {
	fmt.Println(strings.Repeat("\t", level) + task.String())
	level++
	for _, t := range task.downstreamList() {
		getUpstream(t, level)
	}
	return nil
}

// TreeView shows an ascii tree representation of the DAG
func (d *DAG) TreeView() error {
	printSeparator("-", 50)
	fmt.Println(d.FormattedID(), "TREE VIEW")
	printSeparator("-", 50)

	for _, t := range d.Roots() {
		if err := getUpstream(t, 0); err != nil {
			return errors.Wrap(err, "get upstream")
		}
	}
	printSeparator("-", 50)
	return nil
}

// NextRun returns the next run time based on the chron expression
func (d *DAG) NextRun() (time.Time, error) {
	cron, err := cronexpr.Parse(d.ScheduleInterval)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "cron parse")
	}
	return cron.Next(time.Now().UTC()), nil
}

func (d *DAG) getOrCreateDagModel() error {
	conn := db.Connection
	dagModel := &models.DAG{}
	conn.Where(models.DAG{ID: d.ID}).Attrs(models.DAG{
		Description:      d.Description,
		ScheduleInterval: d.ScheduleInterval,
	}).FirstOrCreate(&dagModel)
	return conn.Error
}

func (d *DAG) updateNextScheduled() error {
	conn := db.Connection
	dagModel := &models.DAG{}
	conn.Where(models.DAG{ID: d.ID}).Assign(models.DAG{
		LastSchedulerRun: time.Now().UTC(),
	}).FirstOrCreate(&dagModel)
	return conn.Error
}

// DagRun creeates a dag run model from a dag
func (d *DAG) DagRun() *models.DagRun {
	return &models.DagRun{
		DagID:         d.ID,
		ExecutionDate: time.Now().UTC(),
		State:         state.Pending,
		StartDate:     time.Now().UTC(),
	}
}
