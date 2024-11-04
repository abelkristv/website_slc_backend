package services

import (
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantPositionService interface {
	GetAllAssistantPositions() ([]models.AssistantPosition, error)
	GetAssistantPositionById(id uint) (*models.AssistantPosition, error)
	CreateAssistantPosition(assistantPosition *models.AssistantPosition) error
	UpdateAssistantPosition(assistantPosition *models.AssistantPosition) error
	DeleteAssistantPosition(assistantPosition *models.AssistantPosition) error
	CreatePositionByAssistant(assistantId int, positionName, positionDesc string, startDate, endDate time.Time) (*models.AssistantPosition, error)
}

type assistantPositionService struct {
	repo         repositories.AssistantPositionRepository
	positionRepo repositories.PositionRepository
}

func NewAssistantPositionService(repo repositories.AssistantPositionRepository, positionRepo repositories.PositionRepository) AssistantPositionService {
	return &assistantPositionService{
		repo:         repo,
		positionRepo: positionRepo,
	}
}

func (s *assistantPositionService) GetAllAssistantPositions() ([]models.AssistantPosition, error) {
	return s.repo.GetAllAssistantPositions()
}

func (s *assistantPositionService) GetAssistantPositionById(id uint) (*models.AssistantPosition, error) {
	return s.repo.GetAssistantPositionById(id)
}

func (s *assistantPositionService) CreateAssistantPosition(assistantPosition *models.AssistantPosition) error {
	return s.repo.CreateAssistantPosition(assistantPosition)
}

func (s *assistantPositionService) UpdateAssistantPosition(assistantPosition *models.AssistantPosition) error {
	return s.repo.UpdateAssistantPosition(assistantPosition)
}

func (s *assistantPositionService) DeleteAssistantPosition(assistantPosition *models.AssistantPosition) error {
	return s.repo.DeleteAssistantPosition(assistantPosition)
}

func (s *assistantPositionService) CreatePositionByAssistant(assistantId int, positionName, positionDesc string, startDate, endDate time.Time) (*models.AssistantPosition, error) {

	position := models.Position{
		Name: positionName,
	}

	if err := s.positionRepo.CreatePosition(&position); err != nil {
		return nil, err
	}

	assistantPosition := models.AssistantPosition{
		AssistantId: assistantId,
		PositionId:  int(position.ID),
		Description: positionDesc,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.CreateAssistantPosition(&assistantPosition); err != nil {
		return nil, err
	}

	return &assistantPosition, nil
}
