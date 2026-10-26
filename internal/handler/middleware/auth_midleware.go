package middleware

import (
	"log"
	"service_go_gin/internal/config"
	"service_go_gin/pkg/response"
	"service_go_gin/pkg/tokenize"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(env *config.Environment) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("token")
		if err != nil {
			log.Printf("Token is not found from cookie: %v", err)
			response.Failed(
				ctx,
				response.NewAppError(401, "Bad Token is not valid", nil),
			)
			ctx.Abort()
			return
		}

		claims := &tokenize.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(env.JWTSecretKey), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Token is not valid err valdiate: %v", err)
			response.Failed(
				ctx,
				response.NewAppError(401, "Token is not valid", nil),
			)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("name", claims.Name)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
		ctx.Set("created_at", claims.CreatedAt)
		ctx.Set("updated_at", claims.UpdatedAt)
		ctx.Next()
	}
}
