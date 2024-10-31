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
	"strings"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	api_repositories "github.com/abelkristv/slc_website/wiredsync/api/repositories"
	api_service "github.com/abelkristv/slc_website/wiredsync/api/services"
	"github.com/abelkristv/slc_website/wiredsync/database"
	"github.com/joho/godotenv"
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

const BaseURL = "https://bluejack.binus.ac.id/lapi/api"

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

	userRepository := api_repositories.NewAssistantRepository() // Concrete repository implementation
	userService := api_service.NewUserService(userRepository)
	userService.FetchAssistant(db, api_service.TokenResponse(authToken))

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
