package api_service

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	api_models "github.com/abelkristv/slc_website/wiredsync/api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func isValidUsername(username string) bool {
	var usernamePatterns = []string{
		`^[A-Z]{2}\d{2}-\d{1}$`,
		`^LC\d{3}$`,
		`^LS\d{3}$`,
		`^LB\d{3}$`,
	}

	for _, pattern := range usernamePatterns {
		if matched, _ := regexp.MatchString(pattern, username); matched {
			return true
		}
	}
	return false
}
func processUser(db *gorm.DB, s *AssistantService, user api_models.Assistant, status string) bool {

	var existingAssistant models.Assistant
	if err := db.Where("full_name = ?", user.Name).First(&existingAssistant).Error; err == nil {
		log.Printf("Assistant with FullName %s already exists. Skipping creation.", user.Name)
		return false
	}

	email := fetchEmail(user.BinusianID)

	assistant := createAssistant(db, user, email, status)
	createUser(db, user, int(assistant.ID))

	return true
}

func createAssistant(db *gorm.DB, user api_models.Assistant, email string, status string) models.Assistant {
	var initial, generation string
	profilePictureURL := fmt.Sprintf("https://bluejack.binus.ac.id/lapi/api/Account/GetThumbnail?id=%s", user.PictureID)

	if strings.HasPrefix(user.Username, "LC") && len(user.Username) == 5 {
		initial = user.Username
		generation = "PART-TIME"
	} else {
		initial = user.Username[:2]
		generation = user.Username[2:]
	}

	assistant := models.Assistant{
		Initial:        initial,
		Generation:     generation,
		Email:          email,
		ProfilePicture: profilePictureURL,
		FullName:       user.Name,

		Status: status,
	}

	if err := db.Create(&assistant).Error; err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}

	return assistant
}

func createUser(db *gorm.DB, user api_models.Assistant, assistantID int) {
	hashedPassword := generatePassword(user.Username)

	newUser := models.User{
		Username:    user.Username,
		Password:    hashedPassword,
		AssistantId: assistantID,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("Created new user: %s\n", newUser.Username)
}

func generatePassword(username string) string {
	rawPassword := username + username + username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

func fetchEmail(binusianID string) string {
	email, err := api.FetchEmail(binusianID)
	if err != nil {
		log.Printf("Failed to fetch email for BinusianID %s: %v\n", binusianID, err)
		return ""
	}
	return email
}
