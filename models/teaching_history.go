package models

import "gorm.io/gorm"

type TeachingHistory struct {
	gorm.Model
	AssistantId int
	CourseId    int
	PeriodId    int
	Assistant   Assistant `gorm:"foreignKey:AssistantId"`
	Course      Course    `gorm:"foreignKey:CourseId"`
	Period      Award     `gorm:"foreignKey:PeriodId"`
}
