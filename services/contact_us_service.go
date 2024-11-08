package services

import (
	"fmt"
	"net/smtp"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type ContactUsService interface {
	GetAllContacts() ([]models.ContactUs, error)
	GetContactById(id uint) (*models.ContactUs, error)
	CreateContact(contact *models.ContactUs) error
	UpdateContact(contact *models.ContactUs) error
	DeleteContact(id uint) error
	MarkContactAsRead(id uint, isRead bool) error
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
	if err := s.sendNotificationEmail(contact); err != nil {
		return fmt.Errorf("failed to send notification email: %w", err)
	}
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

func (s *contactUsService) MarkContactAsRead(id uint, isRead bool) error {
	return s.repo.UpdateContactIsRead(id, isRead)
}

func (s *contactUsService) sendNotificationEmail(contact *models.ContactUs) error {
	auth := smtp.PlainAuth(
		"",
		"dteamslc@gmail.com",
		"czer ojen exze vkxi",
		"smtp.gmail.com",
	)

	to := "abel.kristanto@binus.edu"
	subject := "New Contact Message Received"
	body := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s", contact.Name, contact.Email, contact.Message)
	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"dteamslc@gmail.com",
		[]string{to},
		msg,
	)
}
