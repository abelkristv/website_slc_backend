package services

import (
	"errors"
	"log"
	"sort"
	"time"

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
	Period           string
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
	groupedHistory["SLCPosition"] = assistant.SLCPosition

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
		awardPeriod := award.Period

		found := false

		if !found {
			assistantAwardEntries = append(assistantAwardEntries, AssistantAwardEntry{
				AwardTitle:       awardTitle,
				AwardDescription: AwardDescription,
				Period:           awardPeriod.PeriodTitle,
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

	type AssistantExperienceEntry struct {
		PositionName        string
		PositionDescription string
		StartDate           *time.Time
		EndDate             *time.Time
		Location            string
	}

	type CompanyExperience struct {
		CompanyName string
		CompanyLogo string
		Experiences []AssistantExperienceEntry
	}

	var assistantExperienceByCompany []CompanyExperience

	companyExperienceMap := make(map[string]*CompanyExperience)

	for _, exp := range assistant.AssistantExperience {
		log.Print(exp.Position.Company.CompanyName)
		log.Print(exp.Position.PositionName)
		companyName := exp.Position.Company.CompanyName
		companyLogo := exp.Position.Company.CompanyLogo
		experienceData := AssistantExperienceEntry{
			PositionName:        exp.Position.PositionName,
			PositionDescription: exp.Position.PositionDescription,
			StartDate:           &exp.Position.StartDate,
			EndDate:             &exp.Position.EndDate,
			Location:            exp.Position.Location,
		}

		if companyExp, exists := companyExperienceMap[companyName]; exists {
			companyExp.Experiences = append(companyExp.Experiences, experienceData)
		} else {
			newCompanyExperience := &CompanyExperience{
				CompanyName: companyName,
				CompanyLogo: companyLogo,
				Experiences: []AssistantExperienceEntry{experienceData},
			}
			companyExperienceMap[companyName] = newCompanyExperience
		}

	}

	for _, companyExp := range companyExperienceMap {
		sort.Slice(companyExp.Experiences, func(i, j int) bool {
			// If both EndDates are nil, compare by StartDate (ongoing jobs are more recent)
			if companyExp.Experiences[i].EndDate.Equal(time.Time{}) && companyExp.Experiences[j].EndDate.Equal(time.Time{}) {
				return companyExp.Experiences[i].StartDate.After(*companyExp.Experiences[j].StartDate)
			}

			// If one EndDate is nil, the ongoing job (with nil EndDate) should come first
			if companyExp.Experiences[i].EndDate.Equal(time.Time{}) {
				return true // i is ongoing, so it comes first
			}
			if companyExp.Experiences[j].EndDate.Equal(time.Time{}) {
				return false // j is ongoing, so it comes first
			}

			// Both EndDates are non-nil, compare them (most recent first)
			if companyExp.Experiences[i].EndDate.After(*companyExp.Experiences[j].EndDate) {
				return true
			}
			if companyExp.Experiences[i].EndDate.Before(*companyExp.Experiences[j].EndDate) {
				return false
			}

			// If EndDates are equal, compare by StartDate (most recent first)
			return companyExp.Experiences[i].StartDate.After(*companyExp.Experiences[j].StartDate)
		})
	}

	for _, companyExp := range companyExperienceMap {
		assistantExperienceByCompany = append(assistantExperienceByCompany, *companyExp)
	}

	sort.Slice(assistantExperienceByCompany, func(i, j int) bool {
		iExperiences := assistantExperienceByCompany[i].Experiences
		jExperiences := assistantExperienceByCompany[j].Experiences

		// Separate ongoing companies and non-ongoing companies
		iHasOngoing := false
		jHasOngoing := false

		var iLatestExperience, jLatestExperience *AssistantExperienceEntry

		// Check if company i has ongoing positions (EndDate is empty)
		for _, experience := range iExperiences {
			if experience.EndDate.Equal(time.Time{}) {
				iHasOngoing = true
			}
			if iLatestExperience == nil ||
				(experience.EndDate != nil && !experience.EndDate.Equal(time.Time{}) && experience.EndDate.After(*iLatestExperience.EndDate)) {
				iLatestExperience = &experience
			}
		}

		// Check if company j has ongoing positions (EndDate is empty)
		for _, experience := range jExperiences {
			if experience.EndDate.Equal(time.Time{}) {
				jHasOngoing = true
			}
			if jLatestExperience == nil ||
				(experience.EndDate != nil && !experience.EndDate.Equal(time.Time{}) && experience.EndDate.After(*jLatestExperience.EndDate)) {
				jLatestExperience = &experience
			}
		}

		// If both companies are ongoing (both have time.Time{} as EndDate)
		if iHasOngoing && jHasOngoing {
			// Compare based on the most recent StartDate for ongoing jobs
			if iLatestExperience.StartDate.After(*jLatestExperience.StartDate) {
				return true
			}
			if iLatestExperience.StartDate.Before(*jLatestExperience.StartDate) {
				return false
			}
		}

		// Otherwise, compare based on EndDate (most recent first) for non-ongoing companies
		if iHasOngoing {
			return true // i is ongoing, so it comes first
		}
		if jHasOngoing {
			return false // j is ongoing, so it comes first
		}

		// Compare based on EndDate for non-ongoing positions
		if iLatestExperience.EndDate.After(*jLatestExperience.EndDate) {
			return true
		}
		if iLatestExperience.EndDate.Before(*jLatestExperience.EndDate) {
			return false
		}

		// If both EndDates are equal, compare by the most recent StartDate
		return iLatestExperience.StartDate.After(*jLatestExperience.StartDate)
	})

	groupedHistory["AssistantExperiences"] = assistantExperienceByCompany
	log.Print(assistantExperienceByCompany)

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
	if assistant.SLCPositionID != 0 {
		positionExists, err := s.assistantRepo.CheckPositionExists(assistant.SLCPositionID)
		if err != nil {
			return err
		}
		if !positionExists {
			return errors.New("invalid position ID")
		}
	}
	if existingAssistant == nil {
		return errors.New("user not found")
	}

	if assistant.Email != "" {
		existingAssistant.Email = assistant.Email
	}
	if assistant.Bio != "" {
		existingAssistant.Bio = assistant.Bio
	}
	if assistant.ProfilePicture != "" {
		existingAssistant.ProfilePicture = assistant.ProfilePicture
	}
	if assistant.Initial != "" {
		existingAssistant.Initial = assistant.Initial
	}
	if assistant.Generation != "" {
		existingAssistant.Generation = assistant.Generation
	}
	if assistant.FullName != "" {
		existingAssistant.FullName = assistant.FullName
	}
	if assistant.SLCPositionID != 0 {
		existingAssistant.SLCPositionID = assistant.SLCPositionID
	}

	log.Print(existingAssistant.SLCPositionID)

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
