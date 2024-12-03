package services

import (
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type GalleryService interface {
	CreateGallery(gallery *models.Gallery) error
	GetGalleryByID(id uint) (*models.Gallery, error)
	GetAllGalleries() ([]models.Gallery, error)
	UpdateGallery(gallery *models.Gallery) error
	DeleteGallery(id uint) error
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

func (s *galleryService) UpdateGallery(gallery *models.Gallery) error {
	return s.repo.UpdateGallery(gallery)
}

func (s *galleryService) DeleteGallery(id uint) error {
	return s.repo.DeleteGallery(id)
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