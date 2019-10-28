package models

import (
	"time"

	"github.com/estenssoros/dasorm/nulls"
	"github.com/estenssoros/goflow/db"
	"github.com/estenssoros/goflow/state"
)

// TaskInstance stores the state of a task instance. This table is the
// authority and single source of truth around what tasks have run and the
// state they are in.
type TaskInstance struct {
	TaskInstanceID int `gorm:"PRIMARY_KEY"`
	TaskID         string
	DagRunID       int
	StartDate      time.Time
	EndDate        nulls.Time
	Duration       float64
	State          state.State
	TryNumber      int
	MaxTries       int
	HostName       string
	UnixName       string
	PriorityWeight int
	Operator       string
}

func (t *TaskInstance) Create() error {
	conn := db.Connection
	return conn.Create(t).Error
}
func (t *TaskInstance) Update() error {
	conn := db.Connection
	return conn.Save(t).Error
}

func (t *TaskInstance) Start() error {
	t.State = state.Running
	t.StartDate = time.Now().UTC()
	conn := db.Connection
	return conn.Save(t).Error
}

func (t *TaskInstance) Stop() error {
	now := time.Now().UTC()
	t.Duration = now.Sub(t.StartDate).Seconds()
	t.EndDate = nulls.NewTime(now)
	conn := db.Connection
	return conn.Save(t).Error
}
