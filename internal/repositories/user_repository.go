package repositories

import (
	"time"

	"github.com/sigit14ap/simple-dating-app-backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserById(userID uint) (*models.User, error)
	FindNextUser(userIDs []uint) (*models.User, error)
	UpdateUnlimitedSwipe(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where(&models.User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserById(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.Where(&models.User{ID: userID}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindNextUser(userIDs []uint) (*models.User, error) {
	var user models.User

	err := r.db.Not("id", userIDs).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUnlimitedSwipe(userID uint) error {
	expiredTime := time.Now().AddDate(0, 0, 30)
	err := r.db.Model(&models.User{}).Where(&models.User{ID: userID}).Update("unlimited_swipe_expired_at", expiredTime).Error

	if err != nil {
		return err
	}

	return nil
}
