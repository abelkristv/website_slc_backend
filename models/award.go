package models

import "gorm.io/gorm"

type Award struct {
	gorm.Model
	AwardTitle       string `gorm:"size:100"`
	AwardDescription string `gorm:"type:text"`
}
