package models

import "gorm.io/gorm"

type EventImage struct {
	gorm.Model
	Image   string `gorm:"type:text"`
	EventId int    `gorm:"uniqueIndex"`
	Event   Event  `gorm:"foreignKey:EventId"`
}
