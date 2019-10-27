package goflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/estenssoros/goflow/db"
	"github.com/estenssoros/goflow/models"
	"github.com/estenssoros/goflow/state"
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

	taskDict  map[string]TaskInterface
	taskCount int
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
	for _, t := range d.taskDict {
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
func (d *DAG) Run(ctx context.Context, dagRun *models.DagRun) error {
	start := time.Now()

	runner := NewTaskRunner(d.taskDict)

	go runner.Run(ctx)

	dagRun.UpdateState(state.Running)

	for _, task := range d.taskDict {
		if task.IsRoot() {
			task.SetState(state.Queued)
		}
		runner.evalQueue <- task
	}
	<-runner.Done

	dagRun.UpdateState(runner.FinalState())

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
		taskDict:         map[string]TaskInterface{},
	}, nil
}

// AddTask adds a task to a dag
func (d *DAG) AddTask(t TaskInterface) error {
	_, ok := d.taskDict[t.GetID()]
	if ok {
		return errors.Errorf("task %s already exists in dag", t.GetID())
	}
	t.SetState(state.Pending)
	d.taskDict[t.GetID()] = t
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
	t, ok := d.taskDict[taskID]
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
		DefaultView:      d.DefaultView,
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
