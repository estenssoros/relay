package goflow

import (
	"errors"

	"github.com/estenssoros/goflow/state"
)

// TaskInterface an interface for all operators on a DAG
type TaskInterface interface {
	String() string
	GetID() string
	FormattedID() string
	hasUpstream() bool
	downstreamList() []TaskInterface
	upstreamList() []TaskInterface
	GetDag() *DAG
	SetDag(*DAG)
	HasDag() bool
	addDownstreamTask(string)
	addUpstreamTask(string)
	SetState(state.State)
	GetState() state.State
	IsRoot() bool
	Run() error
}

func setRelatives(task, other TaskInterface, upstream bool) error {
	if task.GetDag() != other.GetDag() {
		return errors.New("tried to set relationship between tasks in more than one DAG")
	}
	if !other.HasDag() {
		other.SetDag(task.GetDag())
	}
	if upstream {
		other.addDownstreamTask(task.GetID())
		task.addUpstreamTask(other.GetID())
		return nil
	}
	task.addDownstreamTask(other.GetID())
	other.addUpstreamTask(task.GetID())
	return nil
}

func upstreamList(task TaskInterface) []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range task.upstreamList() {
		task, err := task.GetDag().getTask(taskID.String())
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

func downstreamList(task TaskInterface) []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range task.downstreamList() {
		task, err := task.GetDag().getTask(taskID.String())
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

// ConnectionInterface interface for connection
type ConnectionInterface interface {
	Close() error
	Exec(string) error
}
