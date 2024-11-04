package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type PositionRepository interface {
	GetAllPositions() ([]models.Position, error)
	GetPositionById(id uint) (*models.Position, error)
	CreatePosition(position *models.Position) error
	UpdatePosition(position *models.Position) error
	DeletePosition(position *models.Position) error
}

type positionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepository{
		db: db,
	}
}

func (r *positionRepository) GetAllPositions() ([]models.Position, error) {
	var positions []models.Position
	err := r.db.Find(&positions).Error
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *positionRepository) GetPositionById(id uint) (*models.Position, error) {
	var position models.Position
	err := r.db.First(&position, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *positionRepository) CreatePosition(position *models.Position) error {
	return r.db.Create(position).Error
}

func (r *positionRepository) UpdatePosition(position *models.Position) error {
	return r.db.Save(position).Error
}

func (r *positionRepository) DeletePosition(position *models.Position) error {
	return r.db.Delete(position).Error
}
