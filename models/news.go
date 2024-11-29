package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type News struct {
    gorm.Model
    AssistantId     int
    NewsTitle       string
    NewsDescription string
    Assistant       *Assistant     `gorm:"foreignKey:AssistantId"`
    NewsImages      pq.StringArray `gorm:"type:text[]"`
}

