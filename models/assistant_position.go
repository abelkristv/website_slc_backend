package models

import (
	"time"

	"gorm.io/gorm"
)

type AssistantPosition struct {
	gorm.Model
	AssistantId int `gorm:"uniqueIndex"`
	PositionId  int `gorm:"uniqueIndex"`
	StartDate   time.Time
	EndDate     time.Time

	Assistant Assistant `gorm:"foreignKey:AssistantId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Position  Position  `gorm:"foreignKey:PositionId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
