package services

import (
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type ContactUsService interface {
	GetAllContacts() ([]models.ContactUs, error)
	GetContactById(id uint) (*models.ContactUs, error)
	CreateContact(contact *models.ContactUs) error
	UpdateContact(contact *models.ContactUs) error
	DeleteContact(id uint) error
}

type contactUsService struct {
	repo repositories.ContactUsRepository
}

func NewContactUsService(repo repositories.ContactUsRepository) ContactUsService {
	return &contactUsService{
		repo: repo,
	}
}

func (s *contactUsService) GetAllContacts() ([]models.ContactUs, error) {
	return s.repo.GetAllContacts()
}

func (s *contactUsService) GetContactById(id uint) (*models.ContactUs, error) {
	return s.repo.GetContactById(id)
}

func (s *contactUsService) CreateContact(contact *models.ContactUs) error {
	contact.IsRead = false
	return s.repo.CreateContact(contact)
}

func (s *contactUsService) UpdateContact(contact *models.ContactUs) error {
	return s.repo.UpdateContact(contact)
}

func (s *contactUsService) DeleteContact(id uint) error {
	contact, err := s.repo.GetContactById(id)
	if err != nil {
		return err
	}
	if contact == nil {
		return nil
	}
	return s.repo.DeleteContact(contact)
}
