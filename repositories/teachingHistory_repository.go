// repositories/teaching_history_repository.go
package repositories

import (
	"log"
	"strings"

	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type TeachingHistoryRepository interface {
	GetByAssistantAndPeriod(assistantUsername, periodName string) ([]models.TeachingHistory, error)
	GetByUsernameGroupedByPeriod(assistantUsername string) (map[string][]models.TeachingHistory, error)
}

type teachingHistoryRepository struct {
	db *gorm.DB
}

func NewTeachingHistoryRepository(db *gorm.DB) TeachingHistoryRepository {
	return &teachingHistoryRepository{db: db}
}

func (r *teachingHistoryRepository) GetByAssistantAndPeriod(assistantUsername, periodName string) ([]models.TeachingHistory, error) {
	var histories []models.TeachingHistory

	normalizedUsername := strings.ToLower(strings.TrimSpace(assistantUsername))
	normalizedPeriod := strings.ToLower(strings.TrimSpace(periodName))

	result := r.db.
		Table("teaching_histories").
		Select("teaching_histories.id, teaching_histories.assistant_id, teaching_histories.course_id").
		Joins("JOIN assistants ON assistants.id = teaching_histories.assistant_id").
		Joins("JOIN courses ON courses.id = teaching_histories.course_id").
		Joins("JOIN periods ON periods.id = teaching_histories.period_id").
		Where("LOWER(CONCAT(assistants.initial, assistants.generation)) = ? AND LOWER(periods.period_title) = ?", normalizedUsername, normalizedPeriod).
		Preload("Assistant").
		Preload("Course").
		Find(&histories)

	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
	}

	return histories, result.Error
}

func (r *teachingHistoryRepository) GetByUsernameGroupedByPeriod(assistantUsername string) (map[string][]models.TeachingHistory, error) {
	var histories []models.TeachingHistory

	normalizedUsername := strings.ToLower(strings.TrimSpace(assistantUsername))

	result := r.db.
		Table("teaching_histories").
		Select("teaching_histories.id, teaching_histories.assistant_id, teaching_histories.course_id, teaching_histories.period_id").
		Joins("JOIN assistants ON assistants.id = teaching_histories.assistant_id").
		Joins("JOIN courses ON courses.id = teaching_histories.course_id").
		Joins("JOIN periods ON periods.id = teaching_histories.period_id").
		Where("LOWER(CONCAT(assistants.initial, assistants.generation)) = ?", normalizedUsername).
		Preload("Assistant").
		Preload("Course").
		Preload("Period").
		Find(&histories)

	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
		return nil, result.Error
	}

	groupedHistories := make(map[string][]models.TeachingHistory)

	for _, history := range histories {
		periodTitle := history.Period.PeriodTitle // Assuming you have a Period field in the TeachingHistory model
		groupedHistories[periodTitle] = append(groupedHistories[periodTitle], history)
	}

	return groupedHistories, nil
}
