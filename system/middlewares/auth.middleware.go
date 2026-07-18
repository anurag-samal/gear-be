package middlewares

import (
	"github/anurag/altar-be/system/packages"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const ContextUserKey contextKey = "user"

func AuthMiddleware(jwtManager *pkg.JWTManager) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := jwtManager.ParseAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid access token",
			})
			return
		}

		c.Set(ContextUserKey, claims)

		c.Next()
	}
}

func GetClaims(c *gin.Context) *pkg.Claims {

	value, exists := c.Get(ContextUserKey)
	if !exists {
		return nil
	}

	claims, ok := value.(*pkg.Claims)
	if !ok {
		return nil
	}

	return claims
}
