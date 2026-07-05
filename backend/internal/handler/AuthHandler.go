package handler

import (
	"errors"
	"net/http"

	"reservation-system/backend/internal/middleware"
	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email and password are required",
		})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims, ok := c.Get(middleware.AuthClaimsContextKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication context not found",
		})
		return
	}

	authClaims, ok := claims.(*service.AuthClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authentication context",
		})
		return
	}

	user, err := h.authService.Me(c.Request.Context(), authClaims.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authToken, ok := c.Get("authToken")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication token not found",
		})
		return
	}

	token, ok := authToken.(string)
	if !ok || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authentication token",
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrTokenRevoked) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "logged out",
	})
}
