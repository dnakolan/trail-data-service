package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
		UID:                uuid.New(),
		CreateTrailRequest: req,
		CreatedAt:          &now,
	}

	if err := h.service.CreateTrail(c.Request.Context(), trail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, trail)
}

func (h *TrailsHandler) GetTrailsHandler(c *gin.Context) {
	uid := c.Param("uid")
	trail, err := h.service.GetTrail(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, errors.New("trail not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, trail)
}

func (h *TrailsHandler) ListTrailsHandler(c *gin.Context) {
	filter, err := parseFilter(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trails, err := h.service.GetAllTrails(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, trails)
}

func parseFilter(query url.Values) (*models.TrailFilter, error) {
	var name *string
	var lat *float64
	var lon *float64
	var radiusKm *float64
	var difficulty *models.TrailDifficulty
	var lengthKm *float64

	nameStr := query.Get("name")
	if nameStr != "" {
		name = &nameStr
	}

	latStr := query.Get("lat")
	if latStr != "" {
		val, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid lat: %w", err)
		}
		lat = &val
	}

	lonStr := query.Get("lon")
	if lonStr != "" {
		val, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min_lng: %w", err)
		}
		lon = &val
	}

	radiusKmStr := query.Get("radius-km")
	if radiusKmStr != "" {
		val, err := strconv.ParseFloat(radiusKmStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid radius-km: %w", err)
		}
		radiusKm = &val
	}

	difficultyStr := query.Get("difficulty")
	if difficultyStr != "" {
		if !models.IsValidTrailDifficulty(difficultyStr) {
			return nil, fmt.Errorf("invalid difficulty: %s", difficultyStr)
		}
		dif := models.TrailDifficulty(difficultyStr)
		difficulty = &dif
	}

	lengthKmStr := query.Get("length-km")
	if lengthKmStr != "" {
		val, err := strconv.ParseFloat(lengthKmStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid length-km: %w", err)
		}
		lengthKm = &val
	}

	filter := &models.TrailFilter{
		CreateTrailRequest: models.CreateTrailRequest{
			Name:       name,
			Lat:        lat,
			Lon:        lon,
			Difficulty: difficulty,
			LengthKm:   lengthKm,
		},
		RadiusKm: radiusKm,
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	return filter, nil
}
