package models

import (
	"gorm.io/gorm"
)

type SLCPosition struct {
	gorm.Model
	PositionName string `gorm:"size:100"`
}
