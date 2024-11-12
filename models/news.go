package models

import "gorm.io/gorm"

type News struct {
	gorm.Model
	AssistantId     int
	NewsTitle       string
	NewsDescription string
	Assistant       *Assistant  `gorm:"foreignKey:AssistantId"`
	NewsImages      []NewsImage `gorm:"foreignKey:NewsId"`
}
