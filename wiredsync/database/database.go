package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/abelkristv/slc_website/models"
)

func SetupDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=abel password=hehe dbname=slc_website port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Assistant{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
