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
	}

	for _, pattern := range usernamePatterns {
		if matched, _ := regexp.MatchString(pattern, username); matched {
			return true
		}
	}
	return false
}

func processUser(db *gorm.DB, s *AssistantService, user api_models.Assistant, status string) bool {
	var foundUser models.User
	result := db.Where("username = ?", user.Username).First(&foundUser)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Fatalf("Error querying user: %v", result.Error)
	}

	if result.RowsAffected > 0 {
		return false // User already exists
	}

	email := fetchEmail(user.BinusianID)
	roles := s.FetchAssistantRoles(user.Username)

	assistant := createAssistant(db, user, email, status)
	createUser(db, user, roles, int(assistant.ID))

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
		Status:         status,
	}

	if err := db.Create(&assistant).Error; err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}

	return assistant
}

func createUser(db *gorm.DB, user api_models.Assistant, roles []string, assistantID int) {
	hashedPassword := generatePassword(user.Username)

	role := "Assistant"
	if len(roles) > 0 {
		role = roles[0]
	}

	newUser := models.User{
		Username:    user.Username,
		Password:    hashedPassword,
		Role:        role,
		AssistantId: assistantID,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("Created new user: %s with email and role: %s\n", newUser.Username, newUser.Role)
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
