package repositories

import (
	"log"

	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssistantRepository interface {
	GetAllAssistants() ([]models.Assistant, error)
	GetAssistantById(id uint) (*models.Assistant, error)
	CreateAssistant(user *models.Assistant) error
	UpdateAssistant(user *models.Assistant) error
	DeleteAssistant(user *models.Assistant) error
	GetAllGenerations() ([]string, error)
	GetAssistantsByGeneration(generation string) ([]models.Assistant, error)
	SearchAssistantsByName(name string) ([]models.Assistant, error)
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
	err := r.db.Preload("TeachingHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("period_id")
	}).
		Preload("AssistantAward").
		Preload("AssistantAward.Award").
		Preload("AssistantExperience").
		Preload("AssistantExperience.Position").
		Preload("AssistantExperience.Position.Company").
		Preload("TeachingHistory.Period").
		Preload("TeachingHistory.Course").
		Preload("AssistantSocialMedia").
		Preload("AssistantSocialMedia").First(&assistant, id).Error

	log.Print(assistant.AssistantExperience)
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

func (r *assistantRepository) GetAllGenerations() ([]string, error) {
	var generations []string
	err := r.db.Model(&models.Assistant{}).
		Distinct("generation").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "generation"}, Desc: false}).
		Pluck("generation", &generations).Error
	if err != nil {
		return nil, err
	}
	return generations, nil
}

func (r *assistantRepository) GetAssistantsByGeneration(generation string) ([]models.Assistant, error) {
	var assistants []models.Assistant
	err := r.db.Where("generation = ?", generation).Find(&assistants).Error
	if err != nil {
		return nil, err
	}
	return assistants, nil
}

func (r *assistantRepository) SearchAssistantsByName(name string) ([]models.Assistant, error) {
	var assistants []models.Assistant
	err := r.db.Where("full_name LIKE ?", "%"+name+"%").Find(&assistants).Error
	if err != nil {
		return nil, err
	}
	return assistants, nil
}
