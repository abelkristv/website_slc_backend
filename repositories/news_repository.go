package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type NewsRepository interface {
	CreateNews(news *models.News) error
	GetNewsByID(id uint) (*models.News, error)
	GetAllNews() ([]models.News, error)
	UpdateNews(news *models.News) error
	DeleteNews(id uint) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) CreateNews(news *models.News) error {
	return r.db.Create(news).Error
}

func (r *newsRepository) GetNewsByID(id uint) (*models.News, error) {
	var news models.News
	if err := r.db.First(&news, id).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) GetAllNews() ([]models.News, error) {
	var news []models.News
	if err := r.db.Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}


func (r *newsRepository) UpdateNews(news *models.News) error {
	return r.db.Save(news).Error
}

func (r *newsRepository) DeleteNews(id uint) error {
	return r.db.Delete(&models.News{}, id).Error
}
