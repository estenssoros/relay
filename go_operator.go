package main

import (
	"fmt"

	"github.com/estenssoros/goflow/state"
)

type GoOperator struct {
	TaskID            string
	Dag               *Dag `json:"-"` // avoid recursion
	Retries           int
	Message           string
	GoFunc            func() error
	upstreamTaskIDs   []string
	downstreamTaskIDs []string
	State             state.State
}

// GetID returns the tag id for an operator
func (o *GoOperator) GetID() string {
	return o.TaskID
}

// GetDag returns the dag for an operator
func (o *GoOperator) GetDag() *Dag {
	return o.Dag
}

// HasDag checks to see if the operators dag is nil
func (o *GoOperator) HasDag() bool {
	return o.Dag != nil
}

// SetDag sets the dag on an operator
func (o *GoOperator) SetDag(dag *Dag) {
	o.Dag = dag
}

// addDownstreamTask adds a task id to the downstream list
func (o *GoOperator) addDownstreamTask(taskID string) {
	o.downstreamTaskIDs = append(o.downstreamTaskIDs, taskID)
}

// addUpstreamTask adds a task to the upstream list
func (o *GoOperator) addUpstreamTask(taskID string) {
	o.upstreamTaskIDs = append(o.upstreamTaskIDs, taskID)
}

// SetUpstream creates relationships between tasks
func (o *GoOperator) SetUpstream(task TaskInterface) {
	setRelatives(o, task, true)
}

// SetDownStream creates relationships between tasks
func (o *GoOperator) SetDownStream(task TaskInterface) {
	setRelatives(o, task, false)
}

func (o GoOperator) String() string {
	return o.TaskID
}

func (o *GoOperator) FormattedID() string {
	return fmt.Sprintf("[TASK] %s", o.TaskID)
}

// hasUpstream returns true if the operators has upstream tasks
func (o *GoOperator) hasUpstream() bool {
	return len(o.upstreamTaskIDs) > 0
}

// downstreamList returns the list of downstream tasks
func (o *GoOperator) downstreamList() []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range o.downstreamTaskIDs {
		task, err := o.Dag.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

func (o *GoOperator) upstreamList() []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range o.upstreamTaskIDs {
		task, err := o.Dag.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

func (o *GoOperator) IsRoot() bool {
	return !o.hasUpstream()
}

// Run run the bash operator
func (o *GoOperator) Run() error {
	return o.GoFunc()
}

func (o *GoOperator) SetState(s state.State) {
	o.State = s
}

func (o *GoOperator) GetState() state.State {
	return o.State
}
