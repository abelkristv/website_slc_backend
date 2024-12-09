package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	"github.com/abelkristv/slc_website/wiredsync/api/course"
	api_repositories "github.com/abelkristv/slc_website/wiredsync/api/repositories"
	api_service "github.com/abelkristv/slc_website/wiredsync/api/services"
	"github.com/abelkristv/slc_website/wiredsync/api/token"
	"github.com/abelkristv/slc_website/wiredsync/config"
	"github.com/abelkristv/slc_website/wiredsync/database"
	"gorm.io/gorm"
)

var client *http.Client

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <USERNAME_WIREDSYNC> <PASSWORD_WIREDSYNC>", os.Args[0])
	}

	username := os.Args[1]
	password := os.Args[2]
	if username == "" || password == "" {
		log.Fatalf("USERNAME_WIREDSYNC or PASSWORD_WIREDSYNC is not set")
	}
	var err error
	client, err = token.CreateAuthenticatedClient(username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	db := setupDatabase()
	insertCourseOutlines(db)

	userRepository := api_repositories.NewAssistantRepository()
	userService := api_service.NewUserService(userRepository)
	userService.FetchAssistant(db, api_service.TokenResponse(config.AuthToken))

}

func insertCourseOutlines(db *gorm.DB) {
	courseOutlines, err := api.FetchCourseOutlines(config.AuthToken.AccessToken)
	if err != nil {
		log.Fatalf("Failed to fetch course outlines: %v", err)
	}

	semaphore := make(chan struct{}, 500)
	var wg sync.WaitGroup

	for _, courseOutline := range courseOutlines {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(courseOutline course.GetCourseOutlineResponse) {
			defer wg.Done()
			defer func() { <-semaphore }()

			parts := strings.SplitN(courseOutline.Name, "-", 2)
			if len(parts) < 2 {
				log.Printf("Unexpected course format for %s. Skipping...", courseOutline.Name)
				return
			}

			courseCode := parts[0]
			courseTitle := parts[1]

			courseDescription, err := api.FetchCourseDescription(courseOutline.CourseOutlineId, config.AuthToken.AccessToken)
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

func setupDatabase() *gorm.DB {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	return db
}
