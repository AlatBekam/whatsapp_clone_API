package middleware

import (
	"net/http"
	"strings"
	JsonWebToken "whatsapp-clone-api/JWT"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}


func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Authorization header required"})
			ctx.Abort()
			return 
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid Authorization format"})
			ctx.Abort()
			return 
		}

		tokenString := splitToken[1]
		if (tokenString == "") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Token required"})
			ctx.Abort()
			return 
		}
		
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
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

		// fmt.Println("CLAIMS ID:", claims.ID) // 🔥 DEBUG
		ctx.Set("userID", claims.ID)
		ctx.Next()
	}
}


		// authHeader := ctx.GetHeader("Authorization")
		
		// claims := &JWTClaims{}

		// if authHeader == "" {
		// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Authorization header required"})
		// 	ctx.Abort()
		// 	return 
		// }

		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// 	return JsonWebToken.JWT_SECRET_KEY, nil
		// })

		// if err != nil || !token.Valid {
		// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid token"})
		// 	ctx.Abort()
		// 	return 
		// }

		// if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 	ctx.Set("id", claims["id"])
		// }

		// fmt.Println("CLAIMS ID:", claims.ID) // 🔥 DEBUG
		// ctx.Set("userID", claims.ID)
		// ctx.Next()
