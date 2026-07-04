package handler

import (
	"net/http"
	"strings"

	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/dto"
	"reservation-system/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

type createUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type createUserResponse struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.PasswordHash == "" || req.Role == "" {
		h.errorResponse(c, http.StatusBadRequest, "name, email, password_hash, and role are required")
		return
	}

	userRole := domain.UserRole(strings.ToLower(req.Role))
	if userRole != domain.UserRoleUser && userRole != domain.UserRoleAdmin {
		h.errorResponse(c, http.StatusBadRequest, "role must be user or admin")
		return
	}

	requestDTO := dto.CreateUserRequest{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		Role:         userRole,
	}

	user, err := h.userService.CreateUser(c.Request.Context(), requestDTO)
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	c.JSON(http.StatusCreated, createUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(timeLayout),
		UpdatedAt: user.UpdatedAt.Format(timeLayout),
	})
}

func (h *UserHandler) errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

const timeLayout = "2006-01-02T15:04:05Z07:00"
