package api_service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	api_models "github.com/abelkristv/slc_website/wiredsync/api/models"
	"github.com/abelkristv/slc_website/wiredsync/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BinusianResponse struct {
	BinusianId         string `json:"BinusianId"`
	BirthDate          string `json:"BirthDate"`
	Campus             string `json:"Campus"`
	Name               string `json:"Name"`
	Number             string `json:"Number"`
	PictureId          string `json:"PictureId"`
	ProgramDescription string `json:"ProgramDescription"`
	SeatNum            int    `json:"seatNum"`
	SeatNumber         string `json:"seatNumber"`
}

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
	if err := db.Where("full_name = ? AND generation != ?", user.Name, "PART-TIME").First(&existingAssistant).Error; err == nil {
		log.Printf("Assistant with FullName %s already exists with a Generation other than PART-TIME. Skipping creation.", user.Name)
		return false
	}

	if err := db.Where("full_name = ? AND generation = ?", user.Name, "PART-TIME").First(&existingAssistant).Error; err == nil {
		existingAssistant.Generation = user.Username[2:]
		existingAssistant.Initial = user.Username[:2]
		if updateErr := db.Save(&existingAssistant).Error; updateErr != nil {
			log.Printf("Failed to update Assistant with FullName %s and Generation PART-TIME: %v", user.Name, updateErr)
			return false
		}
		log.Printf("Updated Assistant with FullName %s and Generation PART-TIME.", user.Name)
		return true
	}

	binusianData, err := fetchBinusianData(user.BinusianID)
	if err != nil {
		log.Printf("Error fetching Binusian data for %s: %v", user.BinusianID, err)
		// return false
	}

	email := fetchEmail(user.BinusianID)

	assistant := createAssistant(db, user, email, status, binusianData.BirthDate)
	createUser(db, user, assistant, int(assistant.ID))

	return true
}

func fetchBinusianData(binusianID string) (BinusianResponse, error) {
	url := fmt.Sprintf("%s/Student/GetBinusianByIds?binusianIds=%s", config.BaseURL, binusianID)
	log.Printf("Fetching dob data of %s from api : %s", binusianID, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BinusianResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return BinusianResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BinusianResponse{}, fmt.Errorf("failed to fetch Binusian data: %s", resp.Status)
	}

	var binusianResponse map[string]BinusianResponse
	if err := json.NewDecoder(resp.Body).Decode(&binusianResponse); err != nil {
		return BinusianResponse{}, err
	}

	log.Print(binusianResponse)

	if binusianData, exists := binusianResponse[binusianID]; exists {
		return binusianData, nil
	}

	return BinusianResponse{}, fmt.Errorf("Binusian data not found for %s", binusianID)
}

func createAssistant(db *gorm.DB, user api_models.Assistant, email string, status string, birthDate string) models.Assistant {
	var initial, generation string
	profilePictureURL := fmt.Sprintf("https://bluejack.binus.ac.id/lapi/api/Account/GetThumbnail?id=%s", user.PictureID)

	if (strings.HasPrefix(user.Username, "LC") || strings.HasPrefix(user.Username, "LB") || strings.HasPrefix(user.Username, "LS")) && len(user.Username) == 5 {
		initial = user.Username
		generation = "PART-TIME"
	} else {
		initial = user.Username[:2]
		generation = user.Username[2:]
	}

	layouts := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	var parsedDob time.Time
	var parseErr error

	for _, layout := range layouts {
		parsedDob, parseErr = time.Parse(layout, birthDate)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		log.Printf("Error parsing birthDate '%s': %v", birthDate, parseErr)
		parsedDob = time.Time{}
	}

	assistant := models.Assistant{
		Initial:        initial,
		Generation:     generation,
		Email:          email,
		ProfilePicture: profilePictureURL,
		FullName:       user.Name,
		Status:         status,
		DOB:            parsedDob,
	}

	if err := db.Create(&assistant).Error; err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}

	return assistant
}

func createUser(db *gorm.DB, user api_models.Assistant, assistant models.Assistant, assistantID int) {
	hashedPassword := generatePassword(assistant)

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

func generatePassword(assistant models.Assistant) string {
	rawPassword := assistant.DOB.Format("02-01-2006")

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
