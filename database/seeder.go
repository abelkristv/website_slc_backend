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

// func seedSocialMedia(db *gorm.DB) {
// 	socialMedias := []models.SocialMedia{
// 		{
// 			SocialMediaName:  "Instagram",
// 			SocialMediaImage: "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e7/Instagram_logo_2016.svg/2048px-Instagram_logo_2016.svg.png",
// 		},
// 		{
// 			SocialMediaName:  "LinkedIn",
// 			SocialMediaImage: "https://upload.wikimedia.org/wikipedia/commons/c/ca/LinkedIn_logo_initials.png",
// 		},
// 		{
// 			SocialMediaName:  "Github",
// 			SocialMediaImage: "https://cdn-icons-png.flaticon.com/512/25/25231.png",
// 		},
// 		{
// 			SocialMediaName:  "Line",
// 			SocialMediaImage: "https://upload.wikimedia.org/wikipedia/commons/thumb/4/41/LINE_logo.svg/2048px-LINE_logo.svg.png",
// 		},
// 		{
// 			SocialMediaName:  "Whatsapp",
// 			SocialMediaImage: "https://static.vecteezy.com/system/resources/thumbnails/018/930/746/small/whatsapp-logo-whatsapp-icon-whatsapp-transparent-free-png.png",
// 		},
// 	}

// 	for _, sm := range socialMedias {
// 		err := db.Create(&sm).Error
// 		if err != nil {
// 			log.Fatalf("Failed to seed social media: %v", err)
// 		}
// 	}

// 	log.Println("Social media table seeded successfully")
// }

func SeedDatabase(db *gorm.DB) {
	// seedSocialMedia(db)
}

func ClearDatabase(db *gorm.DB) {
	tables := []string{"teaching_histories", "periods", "courses", "users", "assistants", "social_medias"}

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

	err = db.Exec("ALTER SEQUENCE social_medias_id_seq RESTART WITH 1").Error
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
