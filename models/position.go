package models

import "gorm.io/gorm"

type Position struct {
	gorm.Model
	Name string `gorm:"size:100"`
}
