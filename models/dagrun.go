package models

import (
	"time"

	"github.com/estenssoros/goflow/state"
)

type DagRun struct {
	ID        int
	DagID     string
	StartDate time.Time
	EndDate   time.Time
	State     state.State
}
