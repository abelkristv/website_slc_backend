package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/abelkristv/slc_website/models"
	"github.com/joho/godotenv"
)

func SetupDatabase() (*gorm.DB, error) {
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

	err = db.AutoMigrate(&models.User{}, &models.Assistant{}, &models.Course{}, &models.Period{}, &models.TeachingHistory{}, &models.Position{}, &models.AssistantPosition{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
