// services/award_service.go
package services

import (
	"errors"
	"sort"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantResponse struct {
	ID             uint
	FullName       string
	Email          string
	ProfilePicture string
	Initial        string
	Generation     string
	Status         string
}

type AwardResponse struct {
	AwardTitle string
	Assistants []AssistantResponse
}

type PeriodResponse struct {
	PeriodTitle string
	Awards      []AwardResponse
}

type AwardService interface {
	CreateAward(award *models.Award) error
	GetAwardByID(id uint) (*models.Award, error)
	UpdateAward(award *models.Award) error
	DeleteAward(id uint) error
	GetAllAwards() ([]models.Award, error)
	GetAwardsGroupedByPeriod() ([]PeriodResponse, error)
}

type awardService struct {
	repo repositories.AwardRepository
}

func NewAwardService(repo repositories.AwardRepository) AwardService {
	return &awardService{repo}
}

func (s *awardService) CreateAward(award *models.Award) error {
	if award.AwardTitle == "" {
		return errors.New("AwardTitle is required")
	}
	if award.AwardDescription == "" {
		return errors.New("AwardDescription is required")
	}

	return s.repo.CreateAward(award)
}
func (s *awardService) GetAwardByID(id uint) (*models.Award, error) {
	return s.repo.GetAwardByID(id)
}

func (s *awardService) UpdateAward(award *models.Award) error {
	return s.repo.UpdateAward(award)
}

func (s *awardService) DeleteAward(id uint) error {
	return s.repo.DeleteAward(id)
}

func (s *awardService) GetAllAwards() ([]models.Award, error) {
	return s.repo.GetAllAwards()
}

func (s *awardService) GetAwardsGroupedByPeriod() ([]PeriodResponse, error) {
	periods, err := s.repo.GetAllAwardsGroupedByPeriod()
	if err != nil {
		return nil, err
	}

	awardOrder := map[string]int{
		"Best Performing Assistant":          1,
		"Best TPA":                           2,
		"Best Teaching Index Assistant":      3,
		"Best Qualification":                 4,
		"Best Assistant Diploma":             5,
		"Guider of Best RIG Performing Team": 6,
		"Best RIG":                           7,
		"Best Assistant Candidate":           8,
	}

	var result []PeriodResponse
	for _, period := range periods {
		periodResponse := PeriodResponse{
			PeriodTitle: period.PeriodTitle,
		}

		awardMap := make(map[string]*AwardResponse)
		for _, assistantAward := range period.AssistantAwards {
			awardTitle := assistantAward.Award.AwardTitle
			if _, exists := awardMap[awardTitle]; !exists {
				awardMap[awardTitle] = &AwardResponse{
					AwardTitle: awardTitle,
				}
			}

			assistant := AssistantResponse{
				ID:             assistantAward.Assistant.ID,
				FullName:       assistantAward.Assistant.FullName,
				Email:          assistantAward.Assistant.Email,
				ProfilePicture: assistantAward.Assistant.ProfilePicture,
				Initial:        assistantAward.Assistant.Initial,
				Generation:     assistantAward.Assistant.Generation,
				Status:         assistantAward.Assistant.Status,
			}

			awardMap[awardTitle].Assistants = append(awardMap[awardTitle].Assistants, assistant)
		}

		var awards []AwardResponse
		for _, award := range awardMap {
			awards = append(awards, *award)
		}

		sort.Slice(awards, func(i, j int) bool {
			orderI := awardOrder[awards[i].AwardTitle]
			orderJ := awardOrder[awards[j].AwardTitle]

			if orderI == 0 && orderJ == 0 {
				return awards[i].AwardTitle < awards[j].AwardTitle
			}
			if orderI == 0 {
				return false
			}
			if orderJ == 0 {
				return true
			}
			return orderI < orderJ
		})

		periodResponse.Awards = awards
		result = append(result, periodResponse)
	}

	return result, nil
}
