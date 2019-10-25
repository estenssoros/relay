package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Base base model type for id, created at, and updated at
type Base struct {
	ID        uuid.UUID `gorm:"type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
