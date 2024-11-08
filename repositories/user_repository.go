package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	LoginByUserInitial(initial string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetUserCount() (int, error)
	GetPaginatedUsers(offset, limit int) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Assistant").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	// log.Print(id)

	var assistant models.Assistant
	err = r.db.Preload("TeachingHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("period_id")
	}).
		Preload("AssistantPosition").
		Preload("AssistantAward").
		Preload("AssistantAward.Award").
		Preload("AssistantPosition.Position").
		Preload("TeachingHistory.Period").
		Preload("TeachingHistory.Course").
		Preload("AssistantSocialMedia").
		Preload("AssistantSocialMedia").First(&assistant, user.AssistantId).Error
	// err := r.db.First(&socialMedia, id).Error
	// log.Print(s)

	user.Assistant = assistant

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Assistant").Preload("Assistant.AssistantAward").Preload("Assistant.AssistantSocialMedia").Preload("Assistant.TeachingHistory").Preload("Assistant.TeachingHistory.Course").Preload("Assistant.TeachingHistory.Period").Preload("Assistant.AssistantPosition").Preload("Assistant.AssistantPosition.Position").Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) LoginByUserInitial(initial string) (*models.User, error) {
	var user models.User
	err := r.db.Joins("JOIN assistants ON users.assistant_id = assistants.id").
		Where("assistants.initial || assistants.generation = ?", initial).
		First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *userRepository) GetPaginatedUsers(offset, limit int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Assistant").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserCount() (int, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
