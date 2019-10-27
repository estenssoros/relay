package models

import (
	"time"

	"github.com/estenssoros/goflow/db"
	"github.com/estenssoros/goflow/state"
)

// DagRun describes an instance of a Dag. It can be created by the scheduler or by an external trigger
type DagRun struct {
	ID            int
	DagID         string
	ExecutionDate time.Time
	State         state.State
	StartDate     time.Time
	EndDate       time.Time
}

func (d *DagRun) Create() error {
	conn := db.Connection
	return conn.Create(d).Error
}

func (d *DagRun) UpdateState(s state.State) error {
	conn := db.Connection
	conn.Model(d).Update("state", s)
	return conn.Error
}
