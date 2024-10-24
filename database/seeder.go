package database

import (
	"log"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/services"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) {
	hashedPassword, err := services.HashPassword("hehe")
	if err != nil {
		log.Fatal(err)
	}
	users := []models.User{
		{
			Username:    "abel",
			Password:    hashedPassword,
			Role:        models.AssistantRole.String(),
			AssistantId: 1,
		},
		{
			Username:    "jason",
			Password:    hashedPassword,
			Role:        models.Admin.String(),
			AssistantId: 2,
		},
	}

	for _, user := range users {
		err := db.Create(&user).Error
		if err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
	}

	log.Println("Users table seeded successfully")
}

func seedAssistant(db *gorm.DB) {
	assistant := []models.Assistant{
		{
			Email:          "abel@gmail.com",
			Bio:            "hehe",
			ProfilePicture: "not found",
			Initial:        "BL",
			Generation:     "23-2",
		},
		{
			Email:          "jason@gmail.com",
			Bio:            "hehe",
			ProfilePicture: "not found",
			Initial:        "DT",
			Generation:     "23-2",
		},
	}

	for _, assistant := range assistant {
		err := db.Create(&assistant).Error
		if err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
	}

	log.Println("Assistant table seeded successfully")
}

func SeedDatabase(db *gorm.DB) {
	seedUsers(db)
	seedAssistant(db)
}
