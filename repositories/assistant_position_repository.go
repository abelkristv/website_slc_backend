package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type AssistantPositionRepository interface {
	GetAllAssistantPositions() ([]models.AssistantPosition, error)
	GetAssistantPositionById(id uint) (*models.AssistantPosition, error)
	CreateAssistantPosition(assistantPosition *models.AssistantPosition) error
	UpdateAssistantPosition(assistantPosition *models.AssistantPosition) error
	DeleteAssistantPosition(assistantPosition *models.AssistantPosition) error
}

type assistantPositionRepository struct {
	db *gorm.DB
}

func NewAssistantPositionRepository(db *gorm.DB) AssistantPositionRepository {
	return &assistantPositionRepository{
		db: db,
	}
}

func (r *assistantPositionRepository) GetAllAssistantPositions() ([]models.AssistantPosition, error) {
	var assistantPositions []models.AssistantPosition
	err := r.db.Preload("Assistant").Preload("Position").Find(&assistantPositions).Error
	if err != nil {
		return nil, err
	}
	return assistantPositions, nil
}

func (r *assistantPositionRepository) GetAssistantPositionById(id uint) (*models.AssistantPosition, error) {
	var assistantPosition models.AssistantPosition
	err := r.db.Preload("Assistant").Preload("Position").First(&assistantPosition, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &assistantPosition, nil
}

func (r *assistantPositionRepository) CreateAssistantPosition(assistantPosition *models.AssistantPosition) error {
	return r.db.Create(assistantPosition).Error
}

func (r *assistantPositionRepository) UpdateAssistantPosition(assistantPosition *models.AssistantPosition) error {
	return r.db.Save(assistantPosition).Error
}

func (r *assistantPositionRepository) DeleteAssistantPosition(assistantPosition *models.AssistantPosition) error {
	return r.db.Delete(assistantPosition).Error
}
