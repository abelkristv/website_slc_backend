package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	"github.com/abelkristv/slc_website/wiredsync/database"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var authToken TokenResponse
var client *http.Client

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("USERNAME_WIREDSYNC")
	password := os.Getenv("PASSWORD_WIREDSYNC")
	if username == "" || password == "" {
		log.Fatalf("USERNAME_WIREDSYNC or PASSWORD_WIREDSYNC is not set in .env")
	}
	var err error
	client, err = createAuthenticatedClient(username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	db := setupDatabase()
	insertCourseOutlines(db)
	insertPeriod(db)
	data := fetchUserData()

	for _, user := range data.Active {
		log.Print(user.Username)

		if !isValidUsername(user.Username) {
			log.Printf("Username %s does not match the required format. Skipping...\n", user.Username)
			continue
		}

		if user.Name == "Admin Lab Mass Comm" {
			continue
		}

		if !processUser(db, user) {
			log.Printf("User %s already exists, skipping...\n", user.Username)
			continue
		}
	}
	for _, user := range data.Inactive {
		log.Print(user.Username)

		if !isValidUsername(user.Username) {
			log.Printf("Username %s does not match the required format. Skipping...\n", user.Username)
			continue
		}

		if user.Name == "Admin Lab Mass Comm" {
			continue
		}

		if !processUser(db, user) {
			log.Printf("User %s already exists, skipping...\n", user.Username)
			continue
		}
	}
}

func insertPeriod(db *gorm.DB) {
	periods, err := api.FetchPeriods()
	if err != nil {
		log.Fatalf("Failed to fetch periods: %v", err)
	}

	layouts := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
	}

	for _, period := range periods {
		var start, end time.Time
		var parseErr error

		for _, layout := range layouts {
			start, parseErr = time.Parse(layout, period.Start)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			log.Printf("Failed to parse start date %s: %v", period.Start, parseErr)
			continue
		}

		for _, layout := range layouts {
			end, parseErr = time.Parse(layout, period.End)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			log.Printf("Failed to parse end date %s: %v", period.End, parseErr)
			continue
		}

		periodModel := models.Period{
			PeriodTitle: period.Description,
			StartDate:   start,
			EndDate:     end,
		}

		var existingPeriod models.Period
		if err := db.Where("period_title = ?", periodModel.PeriodTitle).First(&existingPeriod).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error querying period: %v", err)
		}

		if existingPeriod.ID != 0 {
			log.Printf("Period %s already exists, skipping...", periodModel.PeriodTitle)
			continue
		}

		if err := db.Create(&periodModel).Error; err != nil {
			log.Fatalf("Failed to create period: %v", err)
		}
		log.Printf("Created period: %s", periodModel.PeriodTitle)
	}
}
func insertCourseOutlines(db *gorm.DB) {
	courseOutlines, err := api.FetchCourseOutlines(authToken.AccessToken)
	if err != nil {
		log.Fatalf("Failed to fetch course outlines: %v", err)
	}

	for _, courseOutline := range courseOutlines {
		parts := strings.SplitN(courseOutline.Name, "-", 2)
		if len(parts) < 2 {
			log.Printf("Unexpected course format for %s. Skipping...", courseOutline.Name)
			continue
		}

		courseCode := parts[0]
		courseTitle := parts[1]

		course := models.Course{
			CourseCode:  courseCode,
			CourseTitle: courseTitle,
		}

		var existingCourse models.Course
		if err := db.Where("course_code = ?", course.CourseCode).First(&existingCourse).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error querying course: %v", err)
		}

		if existingCourse.ID != 0 {
			log.Printf("Course %s already exists, skipping...", course.CourseTitle)
			continue
		}

		if err := db.Create(&course).Error; err != nil {
			log.Fatalf("Failed to create course: %v", err)
		}
		log.Printf("Created course: %s with code: %s", course.CourseTitle, course.CourseCode)
	}
}

func createAuthenticatedClient(username, password string) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	loginURL := "https://bluejack.binus.ac.id/lapi/api/Account/LogOn"
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	authToken = tokenResponse
	log.Println("Login successful, token stored.")
	return client, nil
}

func setupDatabase() *gorm.DB {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	return db
}

func fetchUserData() api.UserDataResponse {
	data, err := api.FetchDataFromAPI()
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	return data
}

func isValidUsername(username string) bool {
	var usernamePattern = regexp.MustCompile(`^[A-Z]{2}\d{2}-\d{1}$`)
	return usernamePattern.MatchString(username)
}

func processUser(db *gorm.DB, user api.User) bool {
	var foundUser models.User
	result := db.Where("username = ?", user.Username).First(&foundUser)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Fatalf("Error querying user: %v", result.Error)
	}

	if result.RowsAffected > 0 {
		return false // User already exists
	}

	email := fetchEmail(user.BinusianID)
	roles := fetchRoles(user.Username)

	assistant := createAssistant(db, user, email)
	createUser(db, user, roles, int(assistant.ID))

	return true
}

func fetchEmail(binusianID string) string {
	email, err := api.FetchEmail(binusianID)
	if err != nil {
		log.Printf("Failed to fetch email for BinusianID %s: %v\n", binusianID, err)
		return ""
	}
	return email
}

func fetchRoles(username string) []string {
	roles, err := api.FetchAssistantRoles(username)
	if err != nil {
		log.Printf("Failed to fetch roles for %s: %v\n", username, err)
		return []string{"Assistant"}
	}
	return roles
}

func createAssistant(db *gorm.DB, user api.User, email string) models.Assistant {
	initial := user.Username[:2]
	generation := user.Username[2:]
	profilePictureURL := fmt.Sprintf("https://bluejack.binus.ac.id/lapi/api/Account/GetThumbnail?id=%s", user.PictureID)

	assistant := models.Assistant{
		Initial:        initial,
		Generation:     generation,
		Email:          email,
		ProfilePicture: profilePictureURL,
		FullName:       user.Name,
	}

	if err := db.Create(&assistant).Error; err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}

	return assistant
}

func createUser(db *gorm.DB, user api.User, roles []string, assistantID int) {
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
