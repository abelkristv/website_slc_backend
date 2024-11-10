package models

import (
	"time"

	"gorm.io/gorm"
)

type AssistantExperience struct {
	gorm.Model
	AssistantId         int
	CompanyName         string `gorm:"size:100"`
	CompanyLogo 		string `gorm:"type:text"`
	PositionName        string `gorm:"size:100"`
	PositionDescription string `gorm:"size:text"`
	StartDate           time.Time
	EndDate             time.Time

	Assistant Assistant `gorm:"foreignKey:AssistantId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
