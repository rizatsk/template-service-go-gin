package repository

import (
	"context"
	"service_go_gin/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (userRepo *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := userRepo.db.WithContext(ctx).Table(user.TableName()).Where("email=?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepository) FindByUserId(ctx context.Context, userId string) (*domain.User, error) {
	var user domain.User

	err := userRepo.db.WithContext(ctx).Table(user.TableName()).Where("id=?", userId).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
