package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type PeriodRepository interface {
	GetAllPeriods() ([]models.Period, error)
	GetPeriodById(id uint) (*models.Period, error)
	CreatePeriod(user *models.Period) error
	UpdatePeriod(user *models.Period) error
	DeletePeriod(user *models.Period) error
}

type periodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) PeriodRepository {
	return &periodRepository{
		db: db,
	}
}

func (r *periodRepository) GetAllPeriods() ([]models.Period, error) {
	var period []models.Period
	err := r.db.Find(&period).Error
	if err != nil {
		return nil, err
	}
	return period, nil
}

func (r *periodRepository) GetPeriodById(id uint) (*models.Period, error) {
	var period models.Period
	err := r.db.First(&period, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) CreatePeriod(period *models.Period) error {
	return r.db.Create(period).Error
}

func (r *periodRepository) UpdatePeriod(period *models.Period) error {
	return r.db.Save(period).Error
}

func (r *periodRepository) DeletePeriod(period *models.Period) error {
	return r.db.Delete(period).Error
}
