package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, rawAPIKey, err := h.service.RegisterBusiness(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, ErrWeakPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	var dispName string
	if user.DisplayName != nil {
		dispName = *user.DisplayName
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "registration successful",
		"api_key":      rawAPIKey,
		"user_id":      user.ID,
		"username":     user.Username,
		"display_name": dispName,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginUserDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.LoginBusiness(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return

	}

	var dispName string
	if user.DisplayName != nil {
		dispName = *user.DisplayName
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"user_id":      user.ID,
		"username":     user.Username,
		"display_name": dispName,
		"status":       user.Status,
	})

}
