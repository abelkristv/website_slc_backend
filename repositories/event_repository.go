package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetAllEvents() ([]models.Event, error)
	GetEventById(id uint) (*models.Event, error)
	CreateEvent(event *models.Event) error
	UpdateEvent(event *models.Event) error
	DeleteEvent(event *models.Event) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) GetAllEvents() ([]models.Event, error) {
	var event []models.Event
	err := r.db.Find(&event).Error
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *eventRepository) GetEventById(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) CreateEvent(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) UpdateEvent(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) DeleteEvent(event *models.Event) error {
	return r.db.Delete(event).Error
}
