package tokenize

import (
	"service_go_gin/internal/config"
	"service_go_gin/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	jwt.RegisteredClaims
}
type JWTRefreshTokenClaims struct {
	AuthID string `json:"id"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(env *config.Environment, param *domain.User) (string, error) {
	claims := JWTClaims{
		UserID:    param.ID,
		Name:      param.Name,
		Email:     param.Email,
		Role:      param.Role,
		CreatedAt: param.CreatedAt.Format(env.FormatDateTime),
		UpdatedAt: param.UpdatedAt.Format(env.FormatDateTime),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(env.JWTExpiredHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.JWTSecretKey))

	return tokenString, err
}

func GenerateRefreshToken(env *config.Environment, authTokenId string, userId string) (string, error) {
	claims := JWTRefreshTokenClaims{
		AuthID: authTokenId,
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(env.JWTRfreshTokenExpiredDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(env.JWTRefreshTokenSecretKey))

	return tokenString, err
}
