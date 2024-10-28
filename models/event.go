package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	EventTitle       string `gorm:"size:100"`
	EventDescription string `gorm:"size:100"`
	WriterId         int    `gorm:"uniqueIndex"`
	Type             string `gorm:"size:10"`
	PeriodId         int    `gorm:"uniqueIndex"`
	User             User   `gorm:"foreignKey:WriterId"`
	Period           Period `gorm:"foreignKey:PeriodId"`
}
