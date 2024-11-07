// repositories/assistant_award_repository.go
package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type AssistantAwardRepository interface {
	CreateAssistantAward(assistantAward *models.AssistantAward) error
	GetAssistantAwardByID(id uint) (*models.AssistantAward, error)
	GetAssistantAwardsByAssistantID(assistantID int) ([]models.AssistantAward, error)
	UpdateAssistantAward(assistantAward *models.AssistantAward) error
	DeleteAssistantAward(id uint) error
}

type assistantAwardRepository struct {
	db *gorm.DB
}

func NewAssistantAwardRepository(db *gorm.DB) AssistantAwardRepository {
	return &assistantAwardRepository{db}
}

func (r *assistantAwardRepository) CreateAssistantAward(assistantAward *models.AssistantAward) error {
	return r.db.Create(assistantAward).Error
}

func (r *assistantAwardRepository) GetAssistantAwardByID(id uint) (*models.AssistantAward, error) {
	var assistantAward models.AssistantAward
	if err := r.db.Preload("Assistant").Preload("Award").Preload("Period").First(&assistantAward, id).Error; err != nil {
		return nil, err
	}
	return &assistantAward, nil
}

func (r *assistantAwardRepository) GetAssistantAwardsByAssistantID(assistantID int) ([]models.AssistantAward, error) {
	var assistantAwards []models.AssistantAward
	if err := r.db.Where("assistant_id = ?", assistantID).Preload("Award").Preload("Period").Find(&assistantAwards).Error; err != nil {
		return nil, err
	}
	return assistantAwards, nil
}

func (r *assistantAwardRepository) UpdateAssistantAward(assistantAward *models.AssistantAward) error {
	return r.db.Save(assistantAward).Error
}

func (r *assistantAwardRepository) DeleteAssistantAward(id uint) error {
	return r.db.Delete(&models.AssistantAward{}, id).Error
}
