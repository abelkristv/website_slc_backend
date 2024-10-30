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
			Username:    "BL23-2",
			Password:    hashedPassword,
			Role:        models.AssistantRole.String(),
			AssistantId: 1,
		},
		{
			Username:    "DT23-2",
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
	seedAssistant(db)
	seedUsers(db)
}

func ClearDatabase(db *gorm.DB) {
	tables := []string{"teaching_histories", "periods", "courses", "users", "assistants"}

	for _, table := range tables {
		err := db.Exec("DELETE FROM " + table).Error
		if err != nil {
			log.Fatalf("Failed to clear table %s: %v", table, err)
		}
		log.Printf("Cleared table %s successfully", table)

	}
	err := db.Exec("ALTER SEQUENCE teaching_histories_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}

	err = db.Exec("ALTER SEQUENCE assistants_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}
	err = db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}
	err = db.Exec("ALTER SEQUENCE periods_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}
	err = db.Exec("ALTER SEQUENCE courses_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}

	// err = db.Exec("ALTER SEQUENCE assistants_id_seq RESTART WITH 1").Error
	// if err != nil {
	// 	log.Fatalf("Failed to reset sequence for assistants: %v", err)
	// }

	log.Println("Database cleared and sequences reset successfully")
}
