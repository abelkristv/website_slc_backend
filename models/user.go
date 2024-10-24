package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string    `gorm:"size:100"`
	Password    string    `gorm:"size:100"`
	Role        string    `gorm:"size:10"`
	AssistantId int       `gorm:"uniqueIndex"`
	Assistant   Assistant `gorm:"foreignKey:AssistantId"`
}

type UserRole int

const (
	Admin UserRole = iota + 1
	AssistantRole
)

func (role UserRole) String() string {
	switch role {
	case Admin:
		return "Admin"
	case AssistantRole:
		return "Assistant"
	default:
		return "Unknown"
	}
}
