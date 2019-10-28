package goflow

import (
	"fmt"

	"github.com/estenssoros/goflow/models"
	"github.com/estenssoros/goflow/state"
)

// GoOperator operator for go functions
type GoOperator struct {
	TaskID            string
	DAG               *DAG `json:"-"` // avoid recursion
	Retries           int
	Message           string
	GoFunc            func() error
	upstreamTaskIDs   []string
	downstreamTaskIDs []string
	State             state.State
	model             *models.TaskInstance
}

// GetID returns the tag id for an operator
func (o *GoOperator) GetID() string {
	return o.TaskID
}

// GetDag returns the dag for an operator
func (o *GoOperator) GetDag() *DAG {
	return o.DAG
}

// HasDag checks to see if the operators dag is nil
func (o *GoOperator) HasDag() bool {
	return o.DAG != nil
}

// SetDag sets the dag on an operator
func (o *GoOperator) SetDag(dag *DAG) {
	o.DAG = dag
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

// FormattedID exports the formatted id for an operator
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
		task, err := o.DAG.getTask(taskID)
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
		task, err := o.DAG.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

// IsRoot checks to see if an operator has upstream tasks
func (o *GoOperator) IsRoot() bool {
	return !o.hasUpstream()
}

// Run run the bash operator
func (o *GoOperator) Run() error {
	return o.GoFunc()
}

// SetState sets the state on an operator
func (o *GoOperator) SetState(s state.State) {
	o.State = s
}

// GetState gets the state from an operator
func (o *GoOperator) GetState() state.State {
	return o.State
}

func (o *GoOperator) OperatorType() string { return `go` }

func (o *GoOperator) SetModel(m *models.TaskInstance) {
	o.model = m
}

func (o *GoOperator) GetModel() *models.TaskInstance { return o.model }
