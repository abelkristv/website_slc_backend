package models

import (
	"time"

	"gorm.io/gorm"
)

type AssistantExperience struct {
	gorm.Model
	AssistantId         int
	CompanyName         string `gorm:"size:100"`
	PositionName        string `gorm:"size:100"`
	PositionDescription string `gorm:"size:100"`
	StartDate           time.Time
	EndDate             time.Time

	Assistant Assistant `gorm:"foreignKey:AssistantId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
