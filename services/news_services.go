package services

import (
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type NewsService interface {
	CreateNews(news *models.News) error
	GetNewsByID(id uint) (*models.News, error)
	GetAllNews() ([]models.News, error)
	UpdateNews(news *models.News) error
	DeleteNews(id uint) error
}

type newsService struct {
	repo repositories.NewsRepository
}

func NewNewsService(repo repositories.NewsRepository) NewsService {
	return &newsService{repo: repo}
}

func (s *newsService) CreateNews(news *models.News) error {
	return s.repo.CreateNews(news)
}

func (s *newsService) GetNewsByID(id uint) (*models.News, error) {
	return s.repo.GetNewsByID(id)
}

func (s *newsService) GetAllNews() ([]models.News, error) {
	return s.repo.GetAllNews()
}

func (s *newsService) UpdateNews(news *models.News) error {
	return s.repo.UpdateNews(news)
}

func (s *newsService) DeleteNews(id uint) error {
	return s.repo.DeleteNews(id)
}
