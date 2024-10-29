package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type CarrerPathRepository interface {
	GetAllCarrerPaths() ([]models.CarrerPath, error)
	GetCarrerPathById(id uint) (*models.CarrerPath, error)
	CreateCarrerPath(event *models.CarrerPath) error
	UpdateCarrerPath(event *models.CarrerPath) error
	DeleteCarrerPath(event *models.CarrerPath) error
}

type carrerPathRepository struct {
	db *gorm.DB
}

func NewCarrerPathRepository(db *gorm.DB) CarrerPathRepository {
	return &carrerPathRepository{
		db: db,
	}
}

func (r *carrerPathRepository) GetAllCarrerPaths() ([]models.CarrerPath, error) {
	var carrerPath []models.CarrerPath
	err := r.db.Find(&carrerPath).Error
	if err != nil {
		return nil, err
	}
	return carrerPath, nil
}

func (r *carrerPathRepository) GetCarrerPathById(id uint) (*models.CarrerPath, error) {
	var carrerPath models.CarrerPath
	err := r.db.First(&carrerPath, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &carrerPath, nil
}

func (r *carrerPathRepository) CreateCarrerPath(carrerPath *models.CarrerPath) error {
	return r.db.Create(carrerPath).Error
}

func (r *carrerPathRepository) UpdateCarrerPath(carrerPath *models.CarrerPath) error {
	return r.db.Save(carrerPath).Error
}

func (r *carrerPathRepository) DeleteCarrerPath(carrerPath *models.CarrerPath) error {
	return r.db.Delete(carrerPath).Error
}
