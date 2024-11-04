package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type PositionService interface {
	GetAllPositions() ([]models.Position, error)
	GetPositionById(id uint) (*models.Position, error)
	CreatePosition(position *models.Position) error
	UpdatePosition(position *models.Position) error
	DeletePosition(id uint) error
}

type positionService struct {
	positionRepo repositories.PositionRepository
}

func NewPositionService(positionRepo repositories.PositionRepository) PositionService {
	return &positionService{
		positionRepo: positionRepo,
	}
}

func (s *positionService) GetAllPositions() ([]models.Position, error) {
	return s.positionRepo.GetAllPositions()
}

func (s *positionService) GetPositionById(id uint) (*models.Position, error) {
	return s.positionRepo.GetPositionById(id)
}

func (s *positionService) CreatePosition(position *models.Position) error {
	if position.Name == "" {
		return errors.New("position name is required")
	}
	return s.positionRepo.CreatePosition(position)
}

func (s *positionService) UpdatePosition(position *models.Position) error {
	return s.positionRepo.UpdatePosition(position)
}

func (s *positionService) DeletePosition(id uint) error {
	position, err := s.positionRepo.GetPositionById(id)
	if err != nil {
		return err
	}
	if position == nil {
		return errors.New("position not found")
	}
	return s.positionRepo.DeletePosition(position)
}
