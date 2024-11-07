package services

import (
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type AssistantSocialMediaService interface {
	CreateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error
	GetAssistantSocialMediaByID(id uint) (*models.AssistantSocialMedia, error)
	UpdateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error
	DeleteAssistantSocialMedia(id uint) error
	FindByAssistantId(assistantId int) (*models.AssistantSocialMedia, error)
}

type assistantSocialMediaService struct {
	repository repositories.AssistantSocialMediaRepository
}

func NewAssistantSocialMediaService(repository repositories.AssistantSocialMediaRepository) AssistantSocialMediaService {
	return &assistantSocialMediaService{repository}
}

func (s *assistantSocialMediaService) CreateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error {
	return s.repository.CreateAssistantSocialMedia(socialMedia)
}

func (s *assistantSocialMediaService) GetAssistantSocialMediaByID(id uint) (*models.AssistantSocialMedia, error) {
	return s.repository.GetAssistantSocialMediaByID(id)
}

func (s *assistantSocialMediaService) UpdateAssistantSocialMedia(socialMedia *models.AssistantSocialMedia) error {
	return s.repository.UpdateAssistantSocialMedia(socialMedia)
}

func (s *assistantSocialMediaService) DeleteAssistantSocialMedia(id uint) error {
	return s.repository.DeleteAssistantSocialMedia(id)
}

func (s *assistantSocialMediaService) FindByAssistantId(assistantId int) (*models.AssistantSocialMedia, error) {
	return s.repository.FindByAssistantId(assistantId)
}
