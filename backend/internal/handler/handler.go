package handler

import (
	"errors"
	"net/http"

	"reservation-system/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	repo *repository.Repository
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": h.repo.GetRootMessage(),
	})
}

func (h *Handler) Health(c *gin.Context) {
	if err := h.repo.CheckHealth(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "database unavailable",
			"error":  err.Error(),
		})
		return
	}

	if err := h.repo.CheckRedisHealth(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "redis unavailable",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": h.repo.GetHealthStatus(),
	})
}

func (h *Handler) Login(c *gin.Context) {
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

	token, err := h.repo.CreateSession(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email": req.Email,
		},
	})
}

func (h *Handler) Me(c *gin.Context) {
	token := repository.ExtractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing authorization token",
		})
		return
	}

	email, err := h.repo.GetSessionEmail(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "session not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"email": email,
		},
	})
}

func (h *Handler) Logout(c *gin.Context) {
	token := repository.ExtractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing authorization token",
		})
		return
	}

	if err := h.repo.DeleteSession(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "logged out",
	})
}
