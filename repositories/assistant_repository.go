package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type AssistantRepository interface {
	GetAllAssistants() ([]models.Assistant, error)
	GetAssistantById(id uint) (*models.Assistant, error)
	CreateAssistant(user *models.Assistant) error
	UpdateAssistant(user *models.Assistant) error
	DeleteAssistant(user *models.Assistant) error
}

type assistantRepository struct {
	db *gorm.DB
}

func NewAssistantRepository(db *gorm.DB) AssistantRepository {
	return &assistantRepository{
		db: db,
	}
}

func (r *assistantRepository) GetAllAssistants() ([]models.Assistant, error) {
	var assistant []models.Assistant
	err := r.db.Find(&assistant).Error
	if err != nil {
		return nil, err
	}
	return assistant, nil
}

func (r *assistantRepository) GetAssistantById(id uint) (*models.Assistant, error) {
	var assistant models.Assistant
	err := r.db.First(&assistant, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &assistant, nil
}

func (r *assistantRepository) CreateAssistant(user *models.Assistant) error {
	return r.db.Create(user).Error
}

func (r *assistantRepository) UpdateAssistant(user *models.Assistant) error {
	return r.db.Save(user).Error
}

func (r *assistantRepository) DeleteAssistant(user *models.Assistant) error {
	return r.db.Delete(user).Error
}
