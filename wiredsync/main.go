package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	"github.com/abelkristv/slc_website/wiredsync/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	data, err := api.FetchDataFromAPI()
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	var usernamePattern = regexp.MustCompile(`^[A-Z]{2}\d{2}-\d{1}$`)

	for _, user := range data.Active {
		log.Print(user.Username)

		if !usernamePattern.MatchString(user.Username) {
			log.Printf("Username %s does not match the required format. Skipping...\n", user.Username)
			continue
		}

		var foundUser models.User
		result := db.Where("username = ?", user.Username).First(&foundUser)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			log.Fatalf("Error querying user: %v", result.Error)
		}

		if result.RowsAffected == 0 { // User not found
			// Fetch email using BinusianID
			email, err := api.FetchEmail(user.BinusianID)
			if err != nil {
				log.Printf("Failed to fetch email for %s: %v\n", user.Username, err)
				continue
			}

			// Fetch assistant roles using username
			roles, err := api.FetchAssistantRoles(user.Username)
			if err != nil {
				log.Printf("Failed to fetch roles for %s: %v\n", user.Username, err)
				continue
			}

			// Fetch profile picture using PictureId
			profilePicture, err := api.FetchProfilePicture(user.PictureID)
			if err != nil {
				log.Printf("Failed to fetch profile picture for %s: %v\n", user.Username, err)
				continue
			}

			initial := user.Username[:2]    // Assuming the first two characters are the initials
			generation := user.Username[2:] // The rest is the generation
			assistant := models.Assistant{
				Initial:        initial,
				Generation:     generation,
				Email:          email,
				ProfilePicture: profilePicture,
			}

			if err := db.Create(&assistant).Error; err != nil {
				log.Fatalf("Failed to create assistant: %v", err)
			}

			rawPassword := user.Username + user.Username + user.Username

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Failed to hash password: %v", err)
			}

			newUser := models.User{
				Username:    user.Username,
				Password:    string(hashedPassword), // Store the hashed password
				Role:        roles[0],               // Use the first role from the fetched roles
				AssistantId: int(assistant.ID),
			}

			if err := db.Create(&newUser).Error; err != nil {
				log.Fatalf("Failed to create user: %v", err)
			}

			fmt.Printf("Created new user: %s with email: %s, role: %s, and profile picture\n", newUser.Username, email, newUser.Role)
		} else {
			fmt.Printf("User already exists: %s\n", foundUser.Username)
		}
	}
}
