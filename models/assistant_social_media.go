package models

import "gorm.io/gorm"

type AssistantSocialMedia struct {
	gorm.Model
	AssistantId         int
	GithubLink          string
	InstagramLink       string
	LinkedInLink        string
	WhatsappLink        string
	PersonalWebsiteLink string
	// Assistant           Assistant `gorm:"foreignKey:AssistantId"`
}
