package models

import (
	"gorm.io/gorm"
)

type AssistantAward struct {
	gorm.Model
	AssistantId int
	AwardId     int
	PeriodId    int
	AwardImage  string `gorm:"type:text"`

	Assistant Assistant `gorm:"foreignKey:AssistantId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Award     Award     `gorm:"foreignKey:AwardId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Period    Period    `gorm:"foreignKey:PeriodId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
