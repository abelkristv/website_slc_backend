package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type EventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.eventRepo.GetAllEvents()
}

func (s *EventService) GetEventById(id uint) (*models.Event, error) {
	return s.eventRepo.GetEventById(id)
}

func (s *EventService) CreateEvent(eventTitle string, eventDescription string, writerId int, eventtype string, periodId int) (*models.Event, error) {
	if eventTitle == "" || eventDescription == "" || writerId < 0 || eventtype == "" || periodId < 0 {
		return nil, errors.New("all fields are required")
	}

	newEvent := &models.Event{
		EventTitle:       eventTitle,
		EventDescription: eventDescription,
		WriterId:         writerId,
		Type:             eventtype,
		PeriodId:         periodId,
	}

	err := s.eventRepo.CreateEvent(newEvent)
	if err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	existingEvent, err := s.eventRepo.GetEventById(event.ID)
	if err != nil {
		return err
	}
	if existingEvent == nil {
		return errors.New("event not found")
	}

	existingEvent.EventTitle = event.EventTitle
	existingEvent.EventDescription = event.EventDescription
	existingEvent.Type = event.Type
	existingEvent.WriterId = event.WriterId
	existingEvent.PeriodId = event.PeriodId

	return s.eventRepo.UpdateEvent(existingEvent)
}

func (s *EventService) DeleteEvent(id uint) error {
	event, err := s.eventRepo.GetEventById(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}
	return s.eventRepo.DeleteEvent(event)
}
