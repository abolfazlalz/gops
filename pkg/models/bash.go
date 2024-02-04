package models

import (
	"gorm.io/gorm"
	"time"
)

type Bash struct {
	ID        uint           `gorm:"primarykey"`
	Title     string         `json:"title"`
	File      string         `json:"file"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
