package domain

import (
	"context"

	"gorm.io/gorm"
)

type AuthToken struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type AuthTokenRepository interface {
	GetAuthTokenByUserId(ctx context.Context, userId string) ([]*AuthToken, error)
	InsertToken(ctx context.Context, tx *gorm.DB, id string, userId string) error
	DeleteTokenByIds(ctx context.Context, ids []string) error
	DeleteTokenById(ctx context.Context, tx *gorm.DB, id string) error
	GetAuthTokenById(ctx context.Context, authId string) (*AuthToken, error)
}

func (AuthToken) TableName() string {
	return "service_go_gin.auth_tokens"
}
