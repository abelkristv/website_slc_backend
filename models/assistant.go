package models

import "gorm.io/gorm"

type Assistant struct {
	gorm.Model
	Email                string               `gorm:"size:100"`
	Bio                  string               `gorm:"size:100"`
	FullName             string               `gorm:"size:100"`
	ProfilePicture       string               `gorm:"type:text"`
	Initial              string               `gorm:"size:6"`
	Generation           string               `gorm:"size:60"`
	Status               string               `gorm:"size:10"`
	TeachingHistory      []TeachingHistory    `gorm:"foreignKey:AssistantId"`
	AssistantSocialMedia AssistantSocialMedia `gorm:"foreignKey:AssistantId"` // Ensure the foreign key is set correctly
	AssistantPosition    []AssistantPosition  `gorm:"foreignKey:AssistantId"`
	AssistantAward       []AssistantAward     `gorm:"foreignKey:AssistantId"`
}
