package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/dnakolan/trail-data-service/internal/models"
)

type TrailStorage interface {
	Save(ctx context.Context, trail *models.Trail) error
	FindAll(ctx context.Context, filter *models.TrailFilter) ([]*models.Trail, error)
	FindById(ctx context.Context, uid string) (*models.Trail, error)
	Clear(ctx context.Context) error
}

type trailStorage struct {
	sync.RWMutex
	data map[string]*models.Trail
}

func NewTrailStorage() *trailStorage {
	return &trailStorage{
		data: make(map[string]*models.Trail),
	}
}

func (s *trailStorage) Save(ctx context.Context, trail *models.Trail) error {
	s.Lock()
	defer s.Unlock()
	s.data[trail.UID.String()] = trail
	return nil
}

func (s *trailStorage) FindAll(ctx context.Context, filter *models.TrailFilter) ([]*models.Trail, error) {
	s.RLock()
	defer s.RUnlock()
	trails := make([]*models.Trail, 0, len(s.data))
	for _, trail := range s.data {
		if filter == nil || trail.MatchesFilter(filter) {
			trails = append(trails, trail)
		}
	}
	return trails, nil
}

func (s *trailStorage) FindById(ctx context.Context, uid string) (*models.Trail, error) {
	s.RLock()
	defer s.RUnlock()
	trail, ok := s.data[uid]
	if !ok {
		return nil, errors.New("trail not found")
	}
	return trail, nil
}

func (s *trailStorage) Clear(ctx context.Context) error {
	s.Lock()
	defer s.Unlock()
	s.data = make(map[string]*models.Trail)
	return nil
}
