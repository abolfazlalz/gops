package models

import (
	"gorm.io/gorm"
	"time"
)

type CommandLog struct {
	gorm.Model
	Command   string
	Args      string
	Completed bool `json:"completed" gorm:"default=false"`
	Duration  time.Duration
}
