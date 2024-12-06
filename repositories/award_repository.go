// repositories/award_repository.go
package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type AwardRepository interface {
	CreateAward(award *models.Award) error
	GetAwardByID(id uint) (*models.Award, error)
	UpdateAward(award *models.Award) error
	DeleteAward(id uint) error
	GetAllAwards() ([]models.Award, error)
	GetAllAwardsGroupedByPeriod() ([]models.Period, error)
}

type awardRepository struct {
	db *gorm.DB
}

func NewAwardRepository(db *gorm.DB) AwardRepository {
	return &awardRepository{db}
}

func (r *awardRepository) CreateAward(award *models.Award) error {
	return r.db.Create(award).Error
}

func (r *awardRepository) GetAwardByID(id uint) (*models.Award, error) {
	var award models.Award
	if err := r.db.First(&award, id).Error; err != nil {
		return nil, err
	}
	return &award, nil
}

func (r *awardRepository) UpdateAward(award *models.Award) error {
	return r.db.Save(award).Error
}

func (r *awardRepository) DeleteAward(id uint) error {
	return r.db.Delete(&models.Award{}, id).Error
}

func (r *awardRepository) GetAllAwards() ([]models.Award, error) {
	var awards []models.Award
	if err := r.db.Find(&awards).Error; err != nil {
		return nil, err
	}
	return awards, nil
}

func (r *awardRepository) GetAllAwardsGroupedByPeriod() ([]models.Period, error) {
	var periods []models.Period
	err := r.db.Debug().
		Preload("AssistantAwards.Award").
		Preload("AssistantAwards.Assistant").
		Find(&periods).Error
	return periods, err
}
