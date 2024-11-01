package database

import (
	"log"
	"os"

	"github.com/abelkristv/slc_website/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL database:", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.Assistant{}, &models.Award{}, &models.SocialMedia{}, &models.AssistantSocialMedia{}, &models.Event{}, &models.Period{}, &models.User{})
	if err != nil {
		log.Fatal("failed to migrate PostgreSQL database:", err)
		return nil, err
	}

	return db, nil
}
