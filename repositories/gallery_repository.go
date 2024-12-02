package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type GalleryRepository interface {
	CreateGallery(gallery *models.Gallery) error
	GetGalleryByID(id uint) (*models.Gallery, error)
	GetAllGalleries() ([]models.Gallery, error)
	UpdateGallery(gallery *models.Gallery) error
	DeleteGallery(id uint) error
	GetByStatus(status string, galleries *[]models.Gallery) error
	GetByAssistantID(assistantID uint, galleries *[]models.Gallery) error
}

type galleryRepository struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) GalleryRepository {
	return &galleryRepository{db: db}
}

func (r *galleryRepository) CreateGallery(gallery *models.Gallery) error {
	return r.db.Create(gallery).Error
}

func (r *galleryRepository) GetGalleryByID(id uint) (*models.Gallery, error) {
	var gallery models.Gallery
	if err := r.db.Preload("Assistant").First(&gallery, id).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *galleryRepository) GetAllGalleries() ([]models.Gallery, error) {
	var galleries []models.Gallery
	if err := r.db.Preload("Assistant").Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *galleryRepository) UpdateGallery(gallery *models.Gallery) error {
	return r.db.Save(gallery).Error
}

func (r *galleryRepository) DeleteGallery(id uint) error {
	return r.db.Delete(&models.Gallery{}, id).Error
}

func (r *galleryRepository) GetByStatus(status string, galleries *[]models.Gallery) error {
    return r.db.Preload("Assistant").Where("gallery_status = ?", status).Find(galleries).Error
}

func (r *galleryRepository) GetByAssistantID(assistantID uint, galleries *[]models.Gallery) error {
    return r.db.Preload("Assistant").Where("assistant_id = ?", assistantID).Find(galleries).Error
}
