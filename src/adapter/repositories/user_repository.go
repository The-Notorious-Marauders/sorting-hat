package repositories

import (
	"context"
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/adapter/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateLastLoginAt(ctx context.Context, userID uint, lastLoginAt time.Time) error
	Create(ctx context.Context, user *models.User) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) UpdateLastLoginAt(ctx context.Context, userID uint, lastLoginAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", lastLoginAt).Error
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
