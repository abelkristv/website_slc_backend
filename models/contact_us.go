package models

import "gorm.io/gorm"

type ContactUs struct {
	gorm.Model
	Name    string `gorm:"size:100"`
	Email   string `gorm:"size:100"`
	Phone   string `gorm:"size:100"`
	Message string `gorm:"size:text"`
}
