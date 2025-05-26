package handlers

import (
	"net/http"
	"time"

	"github.com/dnakolan/trail-data-service/internal/models"
	"github.com/dnakolan/trail-data-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrailsHandler struct {
	service services.TrailsService
}

func NewTrailsHandler(service services.TrailsService) *TrailsHandler {
	return &TrailsHandler{service: service}
}

func (h *TrailsHandler) CreateTrailHandler(c *gin.Context) {
	var req models.CreateTrailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	trail := &models.Trail{
		UID:        uuid.New(),
		Name:       req.Name,
		LatStart:   req.LatStart,
		LonStart:   req.LonStart,
		Difficulty: req.Difficulty,
		LengthKm:   req.LengthKm,
		CreatedAt:  &now,
	}

	if err := h.service.CreateTrail(c.Request.Context(), trail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, trail)
}
