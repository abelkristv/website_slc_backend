package models

import "gorm.io/gorm"

type Assistant struct {
	gorm.Model
	Email          string       `gorm:"size:100"`
	Bio            string       `gorm:"size:100"`
	ProfilePicture string       `gorm:"type:text"`
	Initial        string       `gorm:"size:5"`
	Generation     string       `gorm:"size:5"`
	CarrerPath     []CarrerPath `gorm:"foreignKey:AssistantId"`
	Award          []Award      `gorm:"foreignKey:AssistantId"`
}
