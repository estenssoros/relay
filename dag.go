package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/estenssoros/goflow/state"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Dag struct {
	ID                   string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Description          string
	ScheduleInterval     string
	StartDate            time.Time
	EndDate              time.Time
	FullFilePath         string
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

func (d Dag) FormattedID() string {
	return fmt.Sprintf("DAG[%s]", d.ID)
}

func (d Dag) String() string {
	ju, _ := json.Marshal(d)
	return string(ju)
}

// Roots are the first tasks that need tp be run
func (d *Dag) Roots() []TaskInterface {
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
func (d *Dag) Run(ctx context.Context) error {
	start := time.Now()

	runner := NewTaskRunner(d.taskDict)

	go runner.Run(ctx)

	for _, task := range d.taskDict {
		if task.IsRoot() {
			task.SetState(state.Queued)
		}
		runner.evalQueue <- task
	}
	<-runner.Done

	logrus.Infof("dag took %v", time.Since(start))
	return nil
}

func NewDag(input *DagConfig) (*Dag, error) {
	return &Dag{
		ID:               input.ID,
		Description:      input.Description,
		ScheduleInterval: input.ScheduleInterval,
		taskDict:         map[string]TaskInterface{},
	}, nil
}

// AddTask adds a task to a dag
func (d Dag) AddTask(t TaskInterface) error {
	_, ok := d.taskDict[t.GetID()]
	if ok {
		return errors.Errorf("task %s already exists in dag", t.GetID())
	}
	t.SetState(state.Pending)
	d.taskDict[t.GetID()] = t
	return nil
}

// NewBash creates a new bash operators on a dag
func (d *Dag) NewBash(o *BashOperator) (*BashOperator, error) {
	d.AddTask(o)
	return o, nil
}

func (d *Dag) getTask(taskID string) (TaskInterface, error) {
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
func (d *Dag) TreeView() error {
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
func (d *Dag) NextRun() (time.Time, error) {
	cron, err := cronexpr.Parse(d.ScheduleInterval)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "cron parse")
	}
	return cron.Next(time.Now().UTC()), nil
}
