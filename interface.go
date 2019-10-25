package main

import (
	"errors"

	"github.com/estenssoros/goflow/state"
)

type TaskInterface interface {
	String() string
	GetID() string
	FormattedID() string
	hasUpstream() bool
	downstreamList() []TaskInterface
	upstreamList() []TaskInterface
	GetDag() *Dag
	SetDag(*Dag)
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
