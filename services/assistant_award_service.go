// services/assistant_award_service.go
package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantAwardService interface {
	CreateAssistantAward(assistantAward *models.AssistantAward) error
	GetAssistantAwardByID(id uint) (*models.AssistantAward, error)
	GetAssistantAwardsByAssistantID(assistantID int) ([]models.AssistantAward, error)
	UpdateAssistantAward(assistantAward *models.AssistantAward) error
	DeleteAssistantAward(id uint) error
}

type assistantAwardService struct {
	repo repositories.AssistantAwardRepository
}

func NewAssistantAwardService(repo repositories.AssistantAwardRepository) AssistantAwardService {
	return &assistantAwardService{repo}
}

func (s *assistantAwardService) CreateAssistantAward(assistantAward *models.AssistantAward) error {
	if assistantAward.AssistantId == 0 {
		return errors.New("AssistantId is required")
	}
	if assistantAward.AwardId == 0 {
		return errors.New("AwardId is required")
	}
	if assistantAward.PeriodId == 0 {
		return errors.New("PeriodId is required")
	}

	return s.repo.CreateAssistantAward(assistantAward)
}
func (s *assistantAwardService) GetAssistantAwardByID(id uint) (*models.AssistantAward, error) {
	return s.repo.GetAssistantAwardByID(id)
}

func (s *assistantAwardService) GetAssistantAwardsByAssistantID(assistantID int) ([]models.AssistantAward, error) {
	return s.repo.GetAssistantAwardsByAssistantID(assistantID)
}

func (s *assistantAwardService) UpdateAssistantAward(assistantAward *models.AssistantAward) error {
	return s.repo.UpdateAssistantAward(assistantAward)
}

func (s *assistantAwardService) DeleteAssistantAward(id uint) error {
	return s.repo.DeleteAssistantAward(id)
}
