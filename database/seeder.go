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
			AssistantId: 1,
		},
		{
			Username:    "DT23-2",
			Password:    hashedPassword,
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

//		log.Println("Social media table seeded successfully")
//	}
func seedAwards(db *gorm.DB) {
	awards := []models.Award{
		{
			AwardTitle:       "Best TPA",
			AwardDescription: "Awarded to the best TPA for their outstanding contribution to teaching.",
		},
		{
			AwardTitle:       "Best Rig",
			AwardDescription: "Recognizing the best teaching assistant for their excellent skills in managing teaching rigs.",
		},
		{
			AwardTitle:       "Best Qualif",
			AwardDescription: "Awarded for exceptional qualification achievements in the academic year.",
		},
		{
			AwardTitle:       "Best Performing Assistant",
			AwardDescription: "Given to the assistant with the best overall performance in supporting students and faculty.",
		},
		{
			AwardTitle:       "Best Teaching Index",
			AwardDescription: "Awarded for the highest teaching index based on student feedback and academic support.",
		},
		{
			AwardTitle:       "Best Qualification",
			AwardDescription: "Given to the assistant with the best academic qualification and contributions to the department.",
		},
		{
			AwardTitle:       "Best Performing Assistant (Diploma)",
			AwardDescription: "Awarded to the best-performing assistant in the diploma program, based on teaching excellence and commitment.",
		},
	}

	// Insert awards into the database
	for _, award := range awards {
		err := db.Create(&award).Error
		if err != nil {
			log.Fatalf("Failed to seed awards: %v", err)
		}
	}

	log.Println("Awards table seeded successfully")
}

func seedAssistantAwards(db *gorm.DB) {
	assistantAwards := []models.AssistantAward{
		{
			AssistantId: 68,                                    // Set to assistantId 68
			AwardId:     1,                                     // AwardId set to 1 (you can modify this based on your awards table)
			PeriodId:    1,                                     // PeriodId set to 1 (assuming there's a period with ID 1)
			AwardImage:  "https://example.com/award_image.png", // Example award image URL
		},
		// Add more entries if needed
	}

	for _, assistantAward := range assistantAwards {
		err := db.Create(&assistantAward).Error
		if err != nil {
			log.Fatalf("Failed to seed assistant awards: %v", err)
		}
	}

	log.Println("AssistantAwards table seeded successfully")
}

func seedSLCPositions(db *gorm.DB) {
	positions := []models.SLCPosition{
		{PositionName: "Junior Laboratory Assistant"},
		{PositionName: "Subject Coordinator Staff"},
		{PositionName: "Subject Coordinator Officer"},
		{PositionName: "Network Administration Staff"},
		{PositionName: "Network Administration Officer"},
		{PositionName: "Research & Development Staff"},
		{PositionName: "Research & Development Officer"},
		{PositionName: "Database Administration Staff"},
		{PositionName: "Database Administration Officer"},
		{PositionName: "Operation Manager Officer"},
	}

	for _, position := range positions {
		err := db.Create(&position).Error
		if err != nil {
			log.Fatalf("Failed to seed SLC positions: %v", err)
		}
	}

	log.Println("SLCPosition table seeded successfully")
}

func SeedDatabase(db *gorm.DB) {
	seedSLCPositions(db)
}

func ClearDatabase(db *gorm.DB) {
	tables := []string{"teaching_histories", "slc_positions", "positions", "periods", "courses", "users", "assistant_social_media", "assistants", "social_media", "contact_us"}

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

	err = db.Exec("ALTER SEQUENCE slc_positions_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}

	err = db.Exec("ALTER SEQUENCE social_media_id_seq RESTART WITH 1").Error
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
	// err = db.Exec("ALTER SEQUENCE assistant_positions_id_seq RESTART WITH 1").Error
	// if err != nil {
	// 	log.Fatalf("Failed to reset sequence for users: %v", err)
	// }
	err = db.Exec("ALTER SEQUENCE assistant_social_media_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}
	err = db.Exec("ALTER SEQUENCE contact_us_id_seq RESTART WITH 1").Error
	if err != nil {
		log.Fatalf("Failed to reset sequence for users: %v", err)
	}

	// err = db.Exec("ALTER SEQUENCE assistants_id_seq RESTART WITH 1").Error
	// if err != nil {
	// 	log.Fatalf("Failed to reset sequence for assistants: %v", err)
	// }

	log.Println("Database cleared and sequences reset successfully")
}
