package models

import "gorm.io/gorm"

type NewsImage struct {
	gorm.Model
	NewsId    int    `gorm:"index"`
	NewsImage string `gorm:"size:text"`
	News      *News  `gorm:"foreignKey:NewsId"`
}
