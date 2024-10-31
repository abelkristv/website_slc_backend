// services/teaching_history_service.go
package services

import (
	"github.com/abelkristv/slc_website/models"
	"github.com/abelkristv/slc_website/repositories"
)

type TeachingHistoryService interface {
	GetTeachingHistoryByAssistantAndPeriod(assistantUsername, periodName string) ([]models.TeachingHistory, error)
	GetTeachingHistoryGroupedByPeriod(assistantUsername string) (map[string][]models.TeachingHistory, error)
}

type teachingHistoryService struct {
	repo repositories.TeachingHistoryRepository
}

func NewTeachingHistoryService(repo repositories.TeachingHistoryRepository) TeachingHistoryService {
	return &teachingHistoryService{repo: repo}
}

func (s *teachingHistoryService) GetTeachingHistoryByAssistantAndPeriod(assistantUsername, periodName string) ([]models.TeachingHistory, error) {
	return s.repo.GetByAssistantAndPeriod(assistantUsername, periodName)
}

func (s *teachingHistoryService) GetTeachingHistoryGroupedByPeriod(assistantUsername string) (map[string][]models.TeachingHistory, error) {
	return s.repo.GetByUsernameGroupedByPeriod(assistantUsername)
}
