package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Gallery struct {
    gorm.Model
    AssistantId				int
    GalleryTitle			string
    GalleryDescription		string
	GalleryStatus 			string
	GalleryNotes			string
    Assistant				*Assistant     `gorm:"foreignKey:AssistantId"`
    GalleryImages			pq.StringArray `gorm:"type:text[]"`
}

