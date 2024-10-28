package services

import (
	"errors"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantService struct {
	assistantRepo repositories.AssistantRepository
}

func NewAssistantService(assistantRepo repositories.AssistantRepository) *AssistantService {
	return &AssistantService{
		assistantRepo: assistantRepo,
	}
}

func (s *AssistantService) GetAllAssistants() ([]models.Assistant, error) {
	return s.assistantRepo.GetAllAssistants()
}

func (s *AssistantService) GetAssistantById(id uint) (*models.Assistant, error) {
	return s.assistantRepo.GetAssistantById(id)
}

func (s *AssistantService) CreateAssistant(email, bio, profile_picture, initial, generation string) (*models.Assistant, error) {
	if email == "" || bio == "" || initial == "" || generation == "" {
		return nil, errors.New("all fields are required")
	}

	newAssistant := &models.Assistant{
		Email:          email,
		Bio:            bio,
		ProfilePicture: profile_picture,
		Initial:        initial,
		Generation:     generation,
	}

	err := s.assistantRepo.CreateAssistant(newAssistant)
	if err != nil {
		return nil, err
	}

	return newAssistant, nil
}

func (s *AssistantService) UpdateAssistant(assistant *models.Assistant) error {
	existingAssistant, err := s.assistantRepo.GetAssistantById(assistant.ID)
	if err != nil {
		return err
	}
	if existingAssistant == nil {
		return errors.New("user not found")
	}

	existingAssistant.Email = assistant.Email
	existingAssistant.Bio = assistant.Bio
	existingAssistant.ProfilePicture = assistant.ProfilePicture
	existingAssistant.Initial = assistant.Initial
	existingAssistant.Generation = assistant.Generation

	return s.assistantRepo.UpdateAssistant(existingAssistant)
}

func (s *AssistantService) DeleteAssistant(id uint) error {
	assistant, err := s.assistantRepo.GetAssistantById(id)
	if err != nil {
		return err
	}
	if assistant == nil {
		return errors.New("user not found")
	}
	return s.assistantRepo.DeleteAssistant(assistant)
}
