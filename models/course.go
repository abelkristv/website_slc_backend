package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	CourseTitle       string `gorm:"size:300"`
	CourseCode        string `gorm:"size:100"`
	CourseDescription string `gorm:"size:text"`
}
