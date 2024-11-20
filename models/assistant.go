package models

import (
	"time"

	"gorm.io/gorm"
)

type Assistant struct {
	gorm.Model
	Email                string `gorm:"size:100"`
	Bio                  string `gorm:"size:100"`
	FullName             string `gorm:"size:100"`
	ProfilePicture       string `gorm:"type:text"`
	Initial              string `gorm:"size:6"`
	Generation           string `gorm:"size:60"`
	Status               string `gorm:"size:10"`
	SLCPositionID        uint
	SLCPosition          SLCPosition           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TeachingHistory      []TeachingHistory     `gorm:"foreignKey:AssistantId"`
	AssistantSocialMedia AssistantSocialMedia  `gorm:"foreignKey:AssistantId"`
	AssistantExperience  []AssistantExperience `gorm:"foreignKey:AssistantId"`
	AssistantAward       []AssistantAward      `gorm:"foreignKey:AssistantId"`
	DOB                  time.Time
}
