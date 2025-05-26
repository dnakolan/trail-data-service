package handlers

import (
	"net/http"

	"github.com/dnakolan/trail-data-service/internal/services"
	"github.com/gin-gonic/gin"
)

type loginHandler struct {
	service services.LoginService
}

func NewLoginHandler(service services.LoginService) *loginHandler {
	return &loginHandler{service: service}
}

func (l *loginHandler) LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	token, err := l.service.Login(c.Request.Context(), credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to login"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"token": token})
}
