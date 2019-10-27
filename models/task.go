package models

import (
	"time"

	"github.com/estenssoros/goflow/state"
)

// TaskInstance stores the state of a task instance. This table is the
// authority and single source of truth around what tasks have run and the
// state they are in.
type TaskInstance struct {
	TaskID         string
	DagRunID       int
	ExecutionDate  time.Time
	StartDate      time.Time
	EndDate        time.Time
	Duration       float64
	State          state.State
	TryNumber      int
	MaxTries       int
	HostName       string
	UnixName       string
	PriorityWeight int
	Operator       string
}
