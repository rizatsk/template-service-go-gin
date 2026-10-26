package service

import (
	"context"
	"log"
	"service_go_gin/internal/config"
	"service_go_gin/internal/domain"
	"service_go_gin/pkg/asynclock"
	"service_go_gin/pkg/response"
	"service_go_gin/pkg/tokenize"
	"uuid"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo domain.UserRepository
	authRepo domain.AuthTokenRepository
	env      *config.Environment
	syncLock *asynclock.SyncLock
	db       *gorm.DB
}

type DependsAuthService struct {
	UserRepo domain.UserRepository
	AuthRepo domain.AuthTokenRepository
	Env      *config.Environment
	SyncLock *asynclock.SyncLock
	DB       *gorm.DB
}

func NewAuthService(depends *DependsAuthService) domain.AuthService {
	return &AuthService{
		userRepo: depends.UserRepo,
		authRepo: depends.AuthRepo,
		env:      depends.Env,
		syncLock: depends.SyncLock,
		db:       depends.DB,
	}
}

func (authService *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.Login, error) {
	user, err := authService.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("User is not registered, error: %v", err)
		return nil, response.NewAppError(401, "Email or password is wrong", nil)
	}

	err_validate_pass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err_validate_pass != nil {
		log.Printf("Password is wrong, err : %v", err_validate_pass)
		return nil, response.NewAppError(401, "Email or password is wrong", nil)
	}

	tokenString, err := tokenize.GenerateAccessToken(authService.env, user)
	if err != nil {
		log.Printf("Failed generate token jwt error: %v", err)
		return nil, response.NewAppError(500, "", nil)
	}

	authTokenId := uuid.NewV7().String()
	refreshTokenString, err := tokenize.GenerateRefreshToken(authService.env, authTokenId, user.ID)
	if err != nil {
		log.Printf("Failed generate refersh token jwt error: %v", err)
		return nil, response.NewAppError(500, "", nil)
	}

	// Locking
	authService.syncLock.Lock(req.Email)
	defer authService.syncLock.Unlock(req.Email)

	authTokens, err := authService.authRepo.GetAuthTokenByUserId(ctx, user.ID)
	if err != nil {
		log.Printf("Failed get tokens user from table : %v", err)
		return nil, response.NewAppError(500, "", nil)
	}

	if len(authTokens) >= authService.env.UserMaxDevice {
		var idsTokenOld []string
		for i := authService.env.UserMaxDevice - 1; i < len(authTokens); i++ {
			idsTokenOld = append(idsTokenOld, authTokens[i].ID)
		}

		if err := authService.authRepo.DeleteTokenByIds(ctx, idsTokenOld); err != nil {
			log.Printf("Failed delete token by id : %v", err)
			return nil, response.NewAppError(500, "", nil)
		}
	}

	if err := authService.authRepo.InsertToken(ctx, nil, authTokenId, user.ID); err != nil {
		log.Printf("Failed create token from table : %v", err)
		return nil, response.NewAppError(500, "", nil)
	}

	return &domain.Login{TokenString: tokenString, RefreshTokenString: refreshTokenString}, nil
}

func (authService *AuthService) RefreshToken(ctx *gin.Context) (*domain.Login, error) {
	// Locking
	dataToken := &domain.DataRefreshToken{
		AuthId: ctx.MustGet("auth_id").(string),
		UserId: ctx.MustGet("user_id").(string),
	}

	// Locking
	authService.syncLock.Lock(dataToken.UserId)
	defer authService.syncLock.Unlock(dataToken.UserId)

	// Get token by auth id
	_, err := authService.authRepo.GetAuthTokenById(ctx.Request.Context(), dataToken.AuthId)
	if err != nil {
		log.Printf("Token is not have in table : %v", err)
		return nil, response.NewAppError(401, "Token is not valid", nil)
	}

	// Get data user by user id
	user, err := authService.userRepo.FindByUserId(ctx.Request.Context(), dataToken.UserId)
	if err != nil {
		log.Printf("Token is not valid email is not fonund in user : %v", err)
		return nil, response.NewAppError(401, "Token is not valid", nil)
	}

	// Generate access token
	tokenString, err := tokenize.GenerateAccessToken(authService.env, user)
	if err != nil {
		log.Printf("Failed generate token jwt error: %v", err)
		return nil, response.NewAppError(500, "", nil)
	}

	// Generate refresh token
	authTokenId := uuid.NewV7().String()
	refreshTokenString, err := tokenize.GenerateRefreshToken(authService.env, authTokenId, user.ID)
	if err != nil {
		log.Printf("Failed generate refersh token jwt error: %v", err)
		return nil, response.NewAppError(500, "", nil)
	}

	if err := authService.db.Transaction(func(tx *gorm.DB) error {
		if err := authService.authRepo.DeleteTokenById(ctx, tx, dataToken.AuthId); err != nil {
			log.Printf("Failed delete token by id : %v", err)
			return err
		}

		if err := authService.authRepo.InsertToken(ctx, tx, authTokenId, user.ID); err != nil {
			log.Printf("Failed create token from table : %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, response.NewAppError(500, "", nil)
	}

	return &domain.Login{TokenString: tokenString, RefreshTokenString: refreshTokenString}, nil
}
