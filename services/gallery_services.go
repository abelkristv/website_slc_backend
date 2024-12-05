package services

import (
	"fmt"

	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type GalleryService interface {
	CreateGallery(gallery *models.Gallery) error
	GetGalleryByID(id uint) (*models.Gallery, error)
	GetAllGalleries() ([]models.Gallery, error)
	UpdateGallery(userID uint, updatedGallery *models.Gallery) error
	DeleteGallery(userID, galleryID uint) error
	GetGalleriesByStatus(status string) ([]models.Gallery, error)
	GetGalleriesByAssistantID(assistantID uint) ([]models.Gallery, error)
	AcceptGallery(gallery *models.Gallery) error
	RejectGallery(gallery *models.Gallery) error
}

type galleryService struct {
	repo repositories.GalleryRepository
}

func NewGalleryService(repo repositories.GalleryRepository) GalleryService {
	return &galleryService{repo: repo}
}

func (s *galleryService) CreateGallery(gallery *models.Gallery) error {
	return s.repo.CreateGallery(gallery)
}

func (s *galleryService) GetGalleryByID(id uint) (*models.Gallery, error) {
	return s.repo.GetGalleryByID(id)
}

func (s *galleryService) GetAllGalleries() ([]models.Gallery, error) {
	return s.repo.GetAllGalleries()
}

func (s *galleryService) UpdateGallery(userID uint, updatedGallery *models.Gallery) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	if user.AssistantId == 0 {
		return fmt.Errorf("no associated Assistant ID for user")
	}

	existingGallery, err := s.repo.GetGalleryByID(updatedGallery.ID)
	if err != nil {
		return fmt.Errorf("gallery not found")
	}

	if existingGallery.AssistantId != user.AssistantId {
		return fmt.Errorf("unauthorized to update this gallery")
	}

	updatedGallery.AssistantId = existingGallery.AssistantId
	if user.Assistant.SLCPosition.PositionName == "Operations Management Officer" {
		updatedGallery.GalleryStatus = "accepted"
	} else {
		updatedGallery.GalleryStatus = "pending"
	}

	err = s.repo.UpdateGallery(updatedGallery)
	if err != nil {
		return fmt.Errorf("failed to update gallery: %w", err)
	}

	return nil
}

func (s *galleryService) DeleteGallery(userID, galleryID uint) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	if user.AssistantId == 0 {
		return fmt.Errorf("no associated Assistant ID for user")
	}

	gallery, err := s.repo.GetGalleryByID(galleryID)
	if err != nil {
		return fmt.Errorf("gallery not found")
	}

	if gallery.AssistantId != user.AssistantId {
		return fmt.Errorf("unauthorized to delete this gallery")
	}

	err = s.repo.DeleteGallery(galleryID)
	if err != nil {
		return fmt.Errorf("failed to delete gallery: %w", err)
	}

	return nil
}

func (s *galleryService) GetGalleriesByStatus(status string) ([]models.Gallery, error) {
	var galleries []models.Gallery
	err := s.repo.GetByStatus(status, &galleries)
	return galleries, err
}

func (s *galleryService) GetGalleriesByAssistantID(assistantID uint) ([]models.Gallery, error) {
	var galleries []models.Gallery
	err := s.repo.GetByAssistantID(assistantID, &galleries)
	return galleries, err
}

func (s *galleryService) AcceptGallery(gallery *models.Gallery) error {
	gallery.GalleryStatus = "accepted"
	return s.repo.UpdateGallery(gallery)
}

func (s *galleryService) RejectGallery(gallery *models.Gallery) error {
	gallery.GalleryStatus = "rejected"
	return s.repo.UpdateGallery(gallery)
}
