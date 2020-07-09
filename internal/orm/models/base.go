package models

import (
	"time"
)

type BaseModel struct {
	// ID should use uuid_generate_v4() for the pk's
	ID        int        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"index;not null;default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt *time.Time `gorm:"index" json:"-"`
}

type BaseModelSoftDelete struct {
	BaseModel
	DeletedAt *time.Time `sql:"index"`
}
