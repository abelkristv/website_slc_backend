// services/award_service.go
package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Initial        string `json:"initial"`
	Generation     string `json:"generation"`
	Status         string `json:"status"`
}

type AwardResponse struct {
	AwardTitle string              `json:"awardTitle"`
	Assistants []AssistantResponse `json:"assistants"`
}

type PeriodResponse struct {
	PeriodTitle string          `json:"periodTitle"`
	Awards      []AwardResponse `json:"awards"`
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

	var result []PeriodResponse
	for _, period := range periods {
		periodResponse := PeriodResponse{
			PeriodTitle: period.PeriodTitle,
		}

		// Group awards by award title
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

		for _, award := range awardMap {
			periodResponse.Awards = append(periodResponse.Awards, *award)
		}

		result = append(result, periodResponse)
	}

	return result, nil
}
