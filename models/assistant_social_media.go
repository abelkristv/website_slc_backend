package models

import "gorm.io/gorm"

type AssistantSocialMedia struct {
	gorm.Model
	SocialMediaId   int
	AssistantId     int
	SocialMediaLink string
	Assistant       Assistant   `gorm:"foreignKey:AssistantId"`
	SocialMedia     SocialMedia `gorm:"foreignKey:SocialMediaId"`
}
