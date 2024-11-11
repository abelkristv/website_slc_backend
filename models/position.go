package models

import (
	"time"

	"gorm.io/gorm"
)

type Position struct {
	gorm.Model
	PositionName      		string         `gorm:"size:100"`
	PositionDescription 	string       `gorm:"type:text"`
	StartDate         		time.Time      `gorm:"type:date"`
	EndDate           		time.Time      `gorm:"type:date"`
	Location				string    `gorm:"size:100"`
	CompanyId         		int
	Company           		Company       `gorm:"foreignKey:CompanyId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
