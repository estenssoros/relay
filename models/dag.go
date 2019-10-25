package models

import (
	"time"
)

type Dag struct {
	ID                 string
	IsPausedAtCreation bool
	IsPaused           bool
	IsSubDag           bool
	IsActive           bool
	LastSchedulerRun   time.Time
	Owners             string
	Description        string
	DefaultView        string
	ScheduleInterval   string
}
