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
	"sync"

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

	semaphore := make(chan struct{}, 500)
	var wg sync.WaitGroup

	for _, courseOutline := range courseOutlines {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(courseOutline api.CourseOutline) {
			defer wg.Done()
			defer func() { <-semaphore }()

			parts := strings.SplitN(courseOutline.Name, "-", 2)
			if len(parts) < 2 {
				log.Printf("Unexpected course format for %s. Skipping...", courseOutline.Name)
				return
			}

			courseCode := parts[0]
			courseTitle := parts[1]

			courseDescription, err := api.FetchCourseDescription(courseOutline.CourseOutlineId, authToken.AccessToken)
			if err != nil {
				log.Printf("Failed to fetch course description for course %s: %v", courseCode, err)
				return
			}

			course := models.Course{
				CourseCode:        courseCode,
				CourseTitle:       courseTitle,
				CourseDescription: courseDescription.CourseDescription,
			}

			var existingCourse models.Course
			if err := db.Where("course_code = ?", course.CourseCode).First(&existingCourse).Error; err != nil && err != gorm.ErrRecordNotFound {
				log.Fatalf("Error querying course: %v", err)
			}

			if existingCourse.ID != 0 {
				log.Printf("Course %s already exists, skipping...", course.CourseTitle)
				return
			}

			if err := db.Create(&course).Error; err != nil {
				log.Fatalf("Failed to create course: %v", err)
			}
			log.Printf("Created course: %s with code: %s", course.CourseTitle, course.CourseCode)
		}(courseOutline)
	}

	wg.Wait()
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
