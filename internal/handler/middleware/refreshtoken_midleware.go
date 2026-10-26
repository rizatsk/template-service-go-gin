package middleware

import (
	"log"
	"service_go_gin/internal/config"
	"service_go_gin/pkg/response"
	"service_go_gin/pkg/tokenize"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshTokenMiddleware(env *config.Environment) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("refreshToken")
		if err != nil {
			log.Printf("Refresh token is not found from cookie: %v", err)
			response.Failed(
				ctx,
				response.NewAppError(401, "Refresh token is not valid", nil),
			)
			ctx.Abort()
			return
		}

		claims := &tokenize.JWTRefreshTokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(env.JWTRefreshTokenSecretKey), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Refresh token is not valid failed: %v", err)
			response.Failed(
				ctx,
				response.NewAppError(401, "Refresh token is not valid", nil),
			)
			ctx.Abort()
			return
		}

		ctx.Set("auth_id", claims.AuthID)
		ctx.Set("user_id", claims.UserID)
		ctx.Next()
	}
}
