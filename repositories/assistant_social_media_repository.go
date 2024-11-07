package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type AssistantSocialMediaRepository interface {
	CreateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error
	GetAssistantSocialMediaByID(id uint) (*models.AssistantSocialMedia, error)
	UpdateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error
	DeleteAssistantSocialMedia(id uint) error
	FindByAssistantId(assistantId int) (*models.AssistantSocialMedia, error)
}

type assistantSocialMediaRepository struct {
	db *gorm.DB
}

func NewAssistantSocialMediaRepository(db *gorm.DB) AssistantSocialMediaRepository {
	return &assistantSocialMediaRepository{db}
}

func (r *assistantSocialMediaRepository) CreateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error {
	return r.db.Create(socialMedia).Error
}

func (r *assistantSocialMediaRepository) GetAssistantSocialMediaByID(id uint) (*models.AssistantSocialMedia, error) {
	var socialMedia models.AssistantSocialMedia
	if err := r.db.First(&socialMedia, id).Error; err != nil {
		return nil, err
	}
	return &socialMedia, nil
}

func (r *assistantSocialMediaRepository) UpdateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error {
	return r.db.Save(socialMedia).Error
}

func (r *assistantSocialMediaRepository) DeleteAssistantSocialMedia(id uint) error {
	return r.db.Delete(&models.AssistantSocialMedia{}, id).Error
}

func (r *assistantSocialMediaRepository) FindByAssistantId(assistantId int) (*models.AssistantSocialMedia, error) {
	var socialMedia models.AssistantSocialMedia
	if err := r.db.Where("assistant_id = ?", assistantId).First(&socialMedia).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &socialMedia, nil
}
