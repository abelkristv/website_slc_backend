package models

import "gorm.io/gorm"

type Award struct {
	gorm.Model
	AwardTitle       string    `gorm:"size:100"`
	AwardDescription string    `gorm:"type:text"`
	AwardImage       string    `gorm:"type:text"`
	PeriodId         string    `gorm:"uniqueIndex"`
	Period           Period    `gorm:"foreignKey:PeriodId"`
	AssistantId      int       `gorm:"uniqueIndex"`
	Assistant        Assistant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
