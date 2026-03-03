package middleware

import (
	"net/http"
	"strings"
	JsonWebToken "whatsapp-clone-api/JWT"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Authorization header required"})
			ctx.Abort()
			return 
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return JsonWebToken.JWT_SECRET_KEY, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid token"})
			ctx.Abort()
			return 
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("id", claims["id"])
		}

		ctx.Next()
	}
}
