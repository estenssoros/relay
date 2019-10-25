package models

import (
	"time"

	"github.com/estenssoros/goflow/state"
)

type TaskRun struct {
	ID             int
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
