package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type SLCPositionRepository interface {
	GetAll() ([]models.SLCPosition, error)
	GetByID(id uint) (*models.SLCPosition, error)
	Create(position *models.SLCPosition) error
	Update(position *models.SLCPosition) error
	Delete(id uint) error
}

type slcPositionRepository struct {
	db *gorm.DB
}

func NewSLCPositionRepository(db *gorm.DB) SLCPositionRepository {
	return &slcPositionRepository{db: db}
}

func (r *slcPositionRepository) GetAll() ([]models.SLCPosition, error) {
	var positions []models.SLCPosition
	err := r.db.Find(&positions).Error
	return positions, err
}

func (r *slcPositionRepository) GetByID(id uint) (*models.SLCPosition, error) {
	var position models.SLCPosition
	err := r.db.First(&position, id).Error
	return &position, err
}

func (r *slcPositionRepository) Create(position *models.SLCPosition) error {
	return r.db.Create(position).Error
}

func (r *slcPositionRepository) Update(position *models.SLCPosition) error {
	return r.db.Save(position).Error
}

func (r *slcPositionRepository) Delete(id uint) error {
	return r.db.Delete(&models.SLCPosition{}, id).Error
}
