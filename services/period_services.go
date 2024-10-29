package services

import (
	"errors"
	"time"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type PeriodService struct {
	periodRepo repositories.PeriodRepository
}

func NewPeriodService(periodRepo repositories.PeriodRepository) *PeriodService {
	return &PeriodService{
		periodRepo: periodRepo,
	}
}

func (s *PeriodService) GetAllPeriods() ([]models.Period, error) {
	return s.periodRepo.GetAllPeriods()
}

func (s *PeriodService) GetPeriodById(id uint) (*models.Period, error) {
	return s.periodRepo.GetPeriodById(id)
}

func (s *PeriodService) CreatePeriod(periodTitle string, startDate time.Time, endDate time.Time) (*models.Period, error) {
	if periodTitle == "" {
		return nil, errors.New("all fields are required")
	}

	newPeriod := &models.Period{
		PeriodTitle: periodTitle,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err := s.periodRepo.CreatePeriod(newPeriod)
	if err != nil {
		return nil, err
	}

	return newPeriod, nil
}

func (s *PeriodService) UpdatePeriod(period *models.Period) error {
	existingPeriod, err := s.periodRepo.GetPeriodById(period.ID)
	if err != nil {
		return err
	}
	if existingPeriod == nil {
		return errors.New("event not found")
	}

	existingPeriod.PeriodTitle = period.PeriodTitle
	existingPeriod.StartDate = period.StartDate
	existingPeriod.EndDate = period.EndDate

	return s.periodRepo.UpdatePeriod(existingPeriod)
}

func (s *PeriodService) DeletePeriod(id uint) error {
	period, err := s.periodRepo.GetPeriodById(id)
	if err != nil {
		return err
	}
	if period == nil {
		return errors.New("period not found")
	}
	return s.periodRepo.DeletePeriod(period)
}
