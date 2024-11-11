package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	CompanyName     string    `gorm:"size:100"`
	CompanyLogo     string    `gorm:"type:text"`
}
