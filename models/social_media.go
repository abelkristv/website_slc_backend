package models

import "gorm.io/gorm"

type SocialMedia struct {
	gorm.Model
	SocialMediaName  string `gorm:"size:100"`
	SocialMediaImage string `gorm:"size:text"`
}
