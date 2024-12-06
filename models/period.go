package models

import (
	"time"

	"gorm.io/gorm"
)

type Period struct {
	gorm.Model
	PeriodTitle     string           `gorm:"size:100"`
	StartDate       time.Time        `gorm:"type:date"`
	EndDate         time.Time        `gorm:"type:date"`
	AssistantAwards []AssistantAward `gorm:"foreignKey:PeriodId;references:ID"`
}
