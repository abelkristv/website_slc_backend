package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantService struct {
	assistantRepo repositories.AssistantRepository
}

func NewAssistantService(assistantRepo repositories.AssistantRepository) *AssistantService {
	return &AssistantService{
		assistantRepo: assistantRepo,
	}
}

func (s *AssistantService) GetAllAssistants() ([]models.Assistant, error) {
	return s.assistantRepo.GetAllAssistants()
}

type TeachingHistoryEntry struct {
	PeriodTitle string
	Courses     []map[string]interface{}
}

type AssistantAwardEntry struct {
	AwardTitle       string
	AwardDescription string
}

func (s *AssistantService) GetAssistantById(id uint) (map[string]interface{}, error) {
	assistant, err := s.assistantRepo.GetAssistantById(id)
	if err != nil {
		return nil, err
	}
	if assistant == nil {
		return nil, nil
	}

	groupedHistory := make(map[string]interface{})
	groupedHistory["ID"] = assistant.ID
	groupedHistory["Email"] = assistant.Email
	groupedHistory["Bio"] = assistant.Bio
	groupedHistory["FullName"] = assistant.FullName
	groupedHistory["ProfilePicture"] = assistant.ProfilePicture
	groupedHistory["Initial"] = assistant.Initial
	groupedHistory["Generation"] = assistant.Generation
	groupedHistory["Status"] = assistant.Status

	type SocialMediaResponse struct {
		AssistantId         int
		GithubLink          string
		InstagramLink       string
		LinkedInLink        string
		WhatsappLink        string
		PersonalWebsiteLink string
	}

	socialMediaResponse := &SocialMediaResponse{
		GithubLink:          assistant.AssistantSocialMedia.GithubLink,
		InstagramLink:       assistant.AssistantSocialMedia.InstagramLink,
		LinkedInLink:        assistant.AssistantSocialMedia.LinkedInLink,
		WhatsappLink:        assistant.AssistantSocialMedia.WhatsappLink,
		PersonalWebsiteLink: assistant.AssistantSocialMedia.PersonalWebsiteLink,
	}

	groupedHistory["SocialMedia"] = socialMediaResponse

	var teachingHistoryEntries []TeachingHistoryEntry
	var assistantAwardEntries []AssistantAwardEntry

	for _, history := range assistant.TeachingHistory {
		periodTitle := history.Period.PeriodTitle
		courseData := map[string]interface{}{
			"CourseTitle":       history.Course.CourseTitle,
			"CourseCode":        history.Course.CourseCode,
			"CourseDescription": history.Course.CourseDescription,
		}

		found := false
		for i := range teachingHistoryEntries {
			if teachingHistoryEntries[i].PeriodTitle == periodTitle {
				teachingHistoryEntries[i].Courses = append(teachingHistoryEntries[i].Courses, courseData)
				found = true
				break
			}
		}
		if !found {
			teachingHistoryEntries = append(teachingHistoryEntries, TeachingHistoryEntry{
				PeriodTitle: periodTitle,
				Courses:     []map[string]interface{}{courseData},
			})
		}
	}

	for _, award := range assistant.AssistantAward {
		awardTitle := award.Award.AwardTitle
		AwardDescription := award.Award.AwardDescription

		found := false

		if !found {
			assistantAwardEntries = append(assistantAwardEntries, AssistantAwardEntry{
				AwardTitle:       awardTitle,
				AwardDescription: AwardDescription,
			})
		}
	}

	sortedTeachingHistory := make([]map[string]interface{}, len(teachingHistoryEntries))
	for i, entry := range teachingHistoryEntries {
		sortedTeachingHistory[i] = map[string]interface{}{
			"PeriodTitle": entry.PeriodTitle,
			"Courses":     entry.Courses,
		}
	}

	groupedHistory["TeachingHistories"] = sortedTeachingHistory
	groupedHistory["Awards"] = assistantAwardEntries

	var experienceEntries []map[string]interface{}
	for _, experience := range assistant.AssistantExperience {
		var startDate string
		if experience.StartDate.Format("2006-01-02 15:04:05-07") != "0001-01-01 00:00:00+00" {
			startDate = experience.StartDate.Format("2006-01-02 15:04:05-07")
		} else {
			startDate = ""
		}
		var endDate string
		if experience.EndDate.Format("2006-01-02 15:04:05-07") != "0001-01-01 00:00:00+00" {
			endDate = experience.EndDate.Format("2006-01-02 15:04:05-07")
		} else {
			endDate = ""
		}

		experienceData := map[string]interface{}{
			"CompanyName":         experience.CompanyName,
			"CompanyLogo":         experience.CompanyLogo,
			"PositionName":        experience.PositionName,
			"PositionDescription": experience.PositionDescription,
			"StartDate":           startDate,
			"EndDate":             endDate,
		}
		experienceEntries = append(experienceEntries, experienceData)
	}
	groupedHistory["AssistantExperience"] = experienceEntries

	return groupedHistory, nil
}

func (s *AssistantService) CreateAssistant(email, bio, profile_picture, initial, generation string) (*models.Assistant, error) {
	if email == "" || bio == "" || initial == "" || generation == "" {
		return nil, errors.New("all fields are required")
	}

	newAssistant := &models.Assistant{
		Email:          email,
		Bio:            bio,
		ProfilePicture: profile_picture,
		Initial:        initial,
		Generation:     generation,
	}

	err := s.assistantRepo.CreateAssistant(newAssistant)
	if err != nil {
		return nil, err
	}

	return newAssistant, nil
}

func (s *AssistantService) UpdateAssistant(assistant *models.Assistant) error {
	existingAssistant, err := s.assistantRepo.GetAssistantById(assistant.ID)
	if err != nil {
		return err
	}
	if existingAssistant == nil {
		return errors.New("user not found")
	}

	existingAssistant.Email = assistant.Email
	existingAssistant.Bio = assistant.Bio
	existingAssistant.ProfilePicture = assistant.ProfilePicture
	existingAssistant.Initial = assistant.Initial
	existingAssistant.Generation = assistant.Generation

	return s.assistantRepo.UpdateAssistant(existingAssistant)
}

func (s *AssistantService) DeleteAssistant(id uint) error {
	assistant, err := s.assistantRepo.GetAssistantById(id)
	if err != nil {
		return err
	}
	if assistant == nil {
		return errors.New("user not found")
	}
	return s.assistantRepo.DeleteAssistant(assistant)
}

func (s *AssistantService) GetAllGenerations() ([]string, error) {
	return s.assistantRepo.GetAllGenerations()
}

func (s *AssistantService) GetAssistantsByGeneration(generation string) ([]models.Assistant, error) {
	return s.assistantRepo.GetAssistantsByGeneration(generation)
}

func (s *AssistantService) SearchAssistantsByName(name string) ([]models.Assistant, error) {
	return s.assistantRepo.SearchAssistantsByName(name)
}
