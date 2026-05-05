package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secrets [][]byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		var token *jwt.Token
		var err error

		for _, secret := range secrets {
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secret, nil

			})

			if err == nil && token.Valid {
				break
			}
		}

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		userID := int(userIDFloat)

		c.Set("userID", userID)

		c.Next()
	}
}
