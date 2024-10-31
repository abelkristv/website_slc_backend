package api_service

import (
	"log"
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/wiredsync/api"
	api_repositories "github.com/abelkristv/slc_website/wiredsync/api/repositories"
	"gorm.io/gorm"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AssistantService struct {
	userRepository api_repositories.AssistantRepository
}

func NewUserService(repo api_repositories.AssistantRepository) *AssistantService {
	return &AssistantService{userRepository: repo}
}
func (s *AssistantService) FetchAssistant(db *gorm.DB, authToken TokenResponse) {
	assistant_data, err := s.userRepository.FetchDataFromAPI()
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	for _, assistant := range assistant_data.Active {
		log.Print(assistant.Username)

		if !isValidUsername(assistant.Username) {
			log.Printf("Username %s does not match the required format. Skipping...\n", assistant.Username)
			continue
		}

		if assistant.Name == "Admin Lab Mass Comm" {
			continue
		}

		if !processUser(db, s, assistant, "active") {
			log.Printf("User %s already exists, skipping...\n", assistant.Username)
			continue
		}
	}
	for _, assistant := range assistant_data.Inactive {
		log.Print(assistant.Username)

		if !isValidUsername(assistant.Username) {
			log.Printf("Username %s does not match the required format. Skipping...\n", assistant.Username)
			continue
		}

		if assistant.Name == "Admin Lab Mass Comm" {
			continue
		}

		if !processUser(db, s, assistant, "inactive") {
			log.Printf("User %s already exists, skipping...\n", assistant.Username)
			continue
		}
	}
	periods := insertPeriod(db)
	for _, assistant := range assistant_data.Active {
		for _, period := range periods {
			if !isValidUsername(assistant.Username) {
				log.Printf("Username %s does not match the required format. Skipping...\n", assistant.Username)
				continue
			}

			schedules, err := api.FetchTeachingHistory(assistant.BinusianID, period.SemesterID, authToken.AccessToken, assistant.Name, period.Description)
			if err != nil {
				log.Printf("Failed to fetch teaching history for assistant %s in semester %s: %v", assistant.Username, period.SemesterID, err)
				continue
			}

			for _, schedule := range schedules {
				log.Printf("Assistant %s teaches %s - %s\n", assistant.Username, schedule.CourseCode, schedule.CourseTitle)

				var course models.Course
				if err := db.Where("course_code = ?", schedule.CourseCode).First(&course).Error; err != nil {
					log.Printf("Course not found for code %s: %v", schedule.CourseCode, err)
					continue
				}

				var periodModel models.Period
				if err := db.Where("period_title = ?", period.Description).First(&periodModel).Error; err != nil {
					log.Printf("Period not found for title %s: %v", period.Description, err)
					continue
				}

				var foundAssistant models.Assistant
				initial := assistant.Username[:2]
				generation := assistant.Username[2:]
				if err := db.Where("initial = ? AND generation = ?", initial, generation).First(&foundAssistant).Error; err != nil {
					log.Printf("Failed to find assistant %s in database: %v", assistant.Username, err)
					continue
				}

				var existingHistory models.TeachingHistory
				if err := db.Where("assistant_id = ? AND course_id = ? AND period_id = ?", foundAssistant.ID, course.ID, periodModel.ID).First(&existingHistory).Error; err == nil {
					log.Printf("Teaching history for assistant %s for course %s in period %s already exists, skipping...\n", assistant.Username, schedule.CourseCode, period.Description)
					continue
				}

				teachingHistory := models.TeachingHistory{
					AssistantId: int(foundAssistant.ID),
					CourseId:    int(course.ID),
					PeriodId:    int(periodModel.ID),
				}

				if err := db.Create(&teachingHistory).Error; err != nil {
					log.Printf("Failed to create teaching history: %v", err)
				} else {
					log.Printf("Inserted teaching history for assistant %s for course %s in period %s", assistant.Username, schedule.CourseCode, period.Description)
				}
			}
		}
	}
}

func (s *AssistantService) FetchAssistantRoles(username string) []string {
	roles, err := s.userRepository.FetchAssistantRoles(username)
	if err != nil {
		log.Printf("Failed to fetch roles for %s: %v\n", username, err)
		return []string{"Assistant"}
	}
	return roles
}

func insertPeriod(db *gorm.DB) []api.Period {
	periods, err := api.FetchPeriods()
	if err != nil {
		log.Fatalf("Failed to fetch periods: %v", err)
	}

	layouts := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
	}

	var createdPeriods []models.Period

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
			// continue
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

		createdPeriods = append(createdPeriods, periodModel)
	}

	return periods
}
