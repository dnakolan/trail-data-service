package services

import (
	"context"
	"errors"

	"github.com/dnakolan/trail-data-service/internal/models"
	"github.com/dnakolan/trail-data-service/internal/storage"
)

const (
	DUPLICATE_TRAIL_RADIUS_KM = 25.0
)

type TrailsService interface {
	CreateTrail(ctx context.Context, trail *models.Trail) error
	GetTrail(ctx context.Context, uid string) (*models.Trail, error)
	UpdateTrail(ctx context.Context, trail *models.Trail) error
	DeleteTrail(ctx context.Context, uid string) error
	GetAllTrails(ctx context.Context, filter *models.TrailFilter) ([]*models.Trail, error)
	GetNearbyTrails(ctx context.Context, lat float64, lon float64, radiusKm float64) ([]*models.Trail, error)
}

type trailsService struct {
	storage storage.TrailStorage
}

func NewTrailsService(storage storage.TrailStorage) *trailsService {
	return &trailsService{storage: storage}
}

func (s *trailsService) CreateTrail(ctx context.Context, trail *models.Trail) error {
	radiusKm := DUPLICATE_TRAIL_RADIUS_KM
	filter := &models.TrailFilter{
		CreateTrailRequest: models.CreateTrailRequest{
			Name: trail.Name,
			Lat:  trail.Lat,
			Lon:  trail.Lon,
		},
		RadiusKm: &radiusKm,
	}

	if duplicates, _ := s.GetAllTrails(ctx, filter); len(duplicates) > 0 {
		return errors.New("trail already exists")
	}

	return s.storage.Save(ctx, trail)
}

func (s *trailsService) GetTrail(ctx context.Context, uid string) (*models.Trail, error) {
	return s.storage.FindById(ctx, uid)
}

func (s *trailsService) UpdateTrail(ctx context.Context, trail *models.Trail) error {
	return s.storage.Save(ctx, trail)
}

func (s *trailsService) DeleteTrail(ctx context.Context, uid string) error {
	return s.storage.Delete(ctx, uid)
}

func (s *trailsService) GetAllTrails(ctx context.Context, filter *models.TrailFilter) ([]*models.Trail, error) {
	return s.storage.FindAll(ctx, filter)
}

func (s *trailsService) GetNearbyTrails(ctx context.Context, lat float64, lon float64, radiusKm float64) ([]*models.Trail, error) {

	return nil, errors.New("not implemented")
}
