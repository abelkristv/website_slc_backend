package models

import (
	"gorm.io/gorm"
)

type AssistantExperience struct {
	gorm.Model
	AssistantId 	int
	PositionId  	int
	Position    	Position     `gorm:"foreignKey:PositionId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}