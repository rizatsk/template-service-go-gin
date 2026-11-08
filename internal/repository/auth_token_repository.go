package repository

import (
	"context"
	"service_blog_go/internal/domain"
	"service_blog_go/pkg/time"

	"gorm.io/gorm"
)

type authTokenRepository struct {
	db *gorm.DB
}

func NewauthTokenRepository(db *gorm.DB) domain.AuthTokenRepository {
	return &authTokenRepository{db: db}
}

func (authTokenRepo *authTokenRepository) GetAuthTokenByUserId(ctx context.Context, userId string) ([]*domain.AuthToken, error) {
	var authToken []*domain.AuthToken

	err := authTokenRepo.db.WithContext(ctx).
		Where("user_id=?", userId).
		Order("created_at DESC").
		Find(&authToken).Error
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (authTokenRepo *authTokenRepository) InsertToken(ctx context.Context, tx *gorm.DB, id string, userId string) error {
	dbConn := authTokenRepo.db.WithContext(ctx)
	if tx != nil {
		dbConn = tx
	}

	newAuthToken := domain.AuthToken{
		ID:        id,
		UserId:    userId,
		CreatedAt: time.TimeNoewFormat(),
	}
	if err := dbConn.Create(newAuthToken).Error; err != nil {
		return err
	}

	return nil
}

func (authTokenRepo *authTokenRepository) DeleteTokenByIds(ctx context.Context, ids []string) error {
	if err := authTokenRepo.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.AuthToken{}).Error; err != nil {
		return err
	}

	return nil
}

func (authTokenRepo *authTokenRepository) DeleteTokenById(ctx context.Context, tx *gorm.DB, id string) error {
	dbConn := authTokenRepo.db.WithContext(ctx)
	if tx != nil {
		dbConn = tx
	}

	if err := dbConn.Where("id = ?", id).Delete(&domain.AuthToken{}).Error; err != nil {
		return err
	}

	return nil
}

func (authTokenRepo *authTokenRepository) GetAuthTokenById(ctx context.Context, authId string) (*domain.AuthToken, error) {
	var token domain.AuthToken

	if err := authTokenRepo.db.WithContext(ctx).Where("id = ?", authId).Take(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}
