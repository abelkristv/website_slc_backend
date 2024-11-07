// services/award_service.go
package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AwardService interface {
	CreateAward(award *models.Award) error
	GetAwardByID(id uint) (*models.Award, error)
	UpdateAward(award *models.Award) error
	DeleteAward(id uint) error
	GetAllAwards() ([]models.Award, error)
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
