package models

import "gorm.io/gorm"

type User struct {
	gorm.DB
	Username string `gorm:"unique"`
	Password string `json:"-"`
}