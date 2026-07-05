package middleware

import (
	"errors"
	"net/http"

	"reservation-system/backend/internal/repository"
	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

const AuthClaimsContextKey = "authClaims"

func NewAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization token",
			})
			return
		}

		claims, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			status := http.StatusUnauthorized
			message := "invalid token"

			if !errors.Is(err, service.ErrInvalidToken) && !errors.Is(err, service.ErrTokenRevoked) {
				status = http.StatusInternalServerError
				message = "failed to validate token"
			}

			c.AbortWithStatusJSON(status, gin.H{
				"error": message,
			})
			return
		}

		c.Set(AuthClaimsContextKey, claims)
		c.Set("authToken", token)
		c.Next()
	}
}

func extractBearerToken(headerValue string) string {
	return repository.ExtractBearerToken(headerValue)
}
