package models

import (
	"time"

	"gorm.io/gorm"
)

type CarrerPath struct {
	gorm.Model
	CarrerTitle       string    `gorm:"size:100"`
	CarrerDescription string    `gorm:"type:text"`
	CarrerImage       string    `gorm:"type:text"`
	StartDate         time.Time `gorm:"type:date"`
	EndDate           time.Time `gorm:"type:date"`
	AssistantId       int       `gorm:"uniqueIndex"`
	Assistant         Assistant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
