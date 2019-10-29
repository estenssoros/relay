package models

import (
	"time"
)

// DAG (directed acyclic graph) a collection of tasks with directional
// dependencies
type DAG struct {
	ID               string `gorm:"PRIMARY_KEY"`
	IsPaused         bool
	LastSchedulerRun time.Time
	Description      string
	ScheduleInterval string
}
