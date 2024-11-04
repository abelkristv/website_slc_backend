package repositories

import (
	"github.com/abelkristv/slc_website/models"
	"gorm.io/gorm"
)

type ContactUsRepository interface {
	GetAllContacts() ([]models.ContactUs, error)
	GetContactById(id uint) (*models.ContactUs, error)
	CreateContact(contact *models.ContactUs) error
	UpdateContact(contact *models.ContactUs) error
	DeleteContact(contact *models.ContactUs) error
}

type contactUsRepository struct {
	db *gorm.DB
}

func NewContactUsRepository(db *gorm.DB) ContactUsRepository {
	return &contactUsRepository{
		db: db,
	}
}

func (r *contactUsRepository) GetAllContacts() ([]models.ContactUs, error) {
	var contacts []models.ContactUs
	err := r.db.Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *contactUsRepository) GetContactById(id uint) (*models.ContactUs, error) {
	var contact models.ContactUs
	err := r.db.First(&contact, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *contactUsRepository) CreateContact(contact *models.ContactUs) error {
	return r.db.Create(contact).Error
}

func (r *contactUsRepository) UpdateContact(contact *models.ContactUs) error {
	return r.db.Save(contact).Error
}

func (r *contactUsRepository) DeleteContact(contact *models.ContactUs) error {
	return r.db.Delete(contact).Error
}
