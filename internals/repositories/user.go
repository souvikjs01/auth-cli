package repository

import (
	"errors"

	"github.com/souvikjs01/auth-cli/internals/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User

	err := r.db.
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
