package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type SLCPositionService interface {
	GetAllPositions() ([]models.SLCPosition, error)
	GetPositionByID(id uint) (*models.SLCPosition, error)
	CreatePosition(position *models.SLCPosition) error
	UpdatePosition(id uint, position *models.SLCPosition) error
	DeletePosition(id uint) error
}

type slcPositionService struct {
	repo repositories.SLCPositionRepository
}

func NewSLCPositionService(repo repositories.SLCPositionRepository) SLCPositionService {
	return &slcPositionService{repo: repo}
}

func (s *slcPositionService) GetAllPositions() ([]models.SLCPosition, error) {
	return s.repo.GetAll()
}

func (s *slcPositionService) GetPositionByID(id uint) (*models.SLCPosition, error) {
	position, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("position not found")
	}
	return position, nil
}

func (s *slcPositionService) CreatePosition(position *models.SLCPosition) error {
	return s.repo.Create(position)
}

func (s *slcPositionService) UpdatePosition(id uint, position *models.SLCPosition) error {
	existingPosition, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("position not found")
	}

	existingPosition.PositionName = position.PositionName
	return s.repo.Update(existingPosition)
}

func (s *slcPositionService) DeletePosition(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("position not found")
	}
	return s.repo.Delete(id)
}
