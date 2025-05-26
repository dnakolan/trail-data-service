package storage

import (
	"context"
	"testing"
	"time"

	"github.com/dnakolan/trail-data-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrailStorage_Save(t *testing.T) {
	storage := NewTrailStorage()
	ctx := context.Background()
	uid := uuid.New()
	now := time.Now()

	trail := &models.Trail{
		CreateTrailRequest: models.CreateTrailRequest{
			Name:       stringPtr("Lamar River Trail"),
			Lat:        float64Ptr(44.8472),
			Lon:        float64Ptr(-109.6278),
			Difficulty: trailDifficultyPtr(models.TrailDifficulty("hard")),
			LengthKm:   float64Ptr(53),
		},
		UID:       uid,
		CreatedAt: &now,
	}

	err := storage.Save(ctx, trail)
	require.NoError(t, err)

	// Verify the trail was saved
	saved, err := storage.FindById(ctx, uid.String())
	require.NoError(t, err)
	assert.Equal(t, trail, saved)
}

func TestTrailStorage_FindById(t *testing.T) {
	storage := NewTrailStorage()
	ctx := context.Background()
	uid := uuid.New()
	notFoundUID := uuid.New()
	now := time.Now()

	trail := &models.Trail{
		CreateTrailRequest: models.CreateTrailRequest{
			Name:       stringPtr("Lamar River Trail"),
			Lat:        float64Ptr(44.8472),
			Lon:        float64Ptr(-109.6278),
			Difficulty: trailDifficultyPtr(models.TrailDifficulty("hard")),
			LengthKm:   float64Ptr(53),
		},
		UID:       uid,
		CreatedAt: &now,
	}

	// Save a trail first
	err := storage.Save(ctx, trail)
	require.NoError(t, err)

	tests := []struct {
		name          string
		uid           string
		expectError   bool
		expectedError string
	}{
		{
			name:        "successful retrieval",
			uid:         uid.String(),
			expectError: false,
		},
		{
			name:          "trail not found",
			uid:           notFoundUID.String(),
			expectError:   true,
			expectedError: "trail not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := storage.FindById(ctx, tt.uid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, found)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, found)
				assert.Equal(t, trail, found)
			}
		})
	}
}

func TestTrailStorage_FindAll(t *testing.T) {
	storage := NewTrailStorage()
	ctx := context.Background()
	now := time.Now()

	// Create test trails
	trails := []*models.Trail{
		{
			CreateTrailRequest: models.CreateTrailRequest{
				Name:       stringPtr("Lamar River Trail"),
				Lat:        float64Ptr(44.8472),
				Lon:        float64Ptr(-109.6278),
				Difficulty: trailDifficultyPtr(models.TrailDifficulty("hard")),
				LengthKm:   float64Ptr(53),
			},
			UID:       uuid.New(),
			CreatedAt: &now,
		},
		{
			CreateTrailRequest: models.CreateTrailRequest{
				Name:       stringPtr("Trail of Ten Falls"),
				Lat:        float64Ptr(43.8242),
				Lon:        float64Ptr(-121.5654),
				Difficulty: trailDifficultyPtr(models.TrailDifficulty("medium")),
				LengthKm:   float64Ptr(10),
			},
			UID:       uuid.New(),
			CreatedAt: &now,
		},
		{
			CreateTrailRequest: models.CreateTrailRequest{
				Name:       stringPtr("Angel's Rest"),
				Lat:        float64Ptr(45.6789),
				Lon:        float64Ptr(-122.3456),
				Difficulty: trailDifficultyPtr(models.TrailDifficulty("medium")),
				LengthKm:   float64Ptr(10),
			},
			UID:       uuid.New(),
			CreatedAt: &now,
		},
	}

	// Save all trails
	for _, w := range trails {
		err := storage.Save(ctx, w)
		require.NoError(t, err)
	}

	tests := []struct {
		name           string
		filter         *models.TrailFilter
		expectedCount  int
		expectedNames  []string
		expectedFilter func(*models.Trail) bool
	}{
		{
			name:          "no filter",
			filter:        nil,
			expectedCount: 3,
			expectedNames: []string{"Lamar River Trail", "Trail of Ten Falls", "Angel's Rest"},
		},
		{
			name: "filter by difficulty",
			filter: &models.TrailFilter{
				CreateTrailRequest: models.CreateTrailRequest{
					Difficulty: trailDifficultyPtr(models.TrailDifficulty("medium")),
				},
			},
			expectedCount: 2,
			expectedNames: []string{"Trail of Ten Falls", "Angel's Rest"},
			expectedFilter: func(w *models.Trail) bool {
				return *w.Difficulty == models.TrailDifficulty("medium")
			},
		},
		{
			name: "filter by proximity",
			filter: &models.TrailFilter{
				CreateTrailRequest: models.CreateTrailRequest{
					Lat: float64Ptr(45.6789),
					Lon: float64Ptr(-122.3456),
				},
				RadiusKm: float64Ptr(10),
			},
			expectedCount: 1,
			expectedNames: []string{"Angel's Rest"},
			expectedFilter: func(w *models.Trail) bool {
				return *w.Lat == 45.6789 && *w.Lon == -122.3456
			},
		},
		{
			name: "filter with no matches",
			filter: &models.TrailFilter{
				CreateTrailRequest: models.CreateTrailRequest{
					Name: stringPtr("nonexistent"),
				},
			},
			expectedCount: 0,
			expectedNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := storage.FindAll(ctx, tt.filter)
			require.NoError(t, err)
			assert.Len(t, found, tt.expectedCount)

			// Verify the names of returned trails
			names := make([]string, len(found))
			for i, w := range found {
				names[i] = *w.Name
			}
			assert.ElementsMatch(t, tt.expectedNames, names)

			// If there's a specific filter function, verify each trail matches it
			if tt.expectedFilter != nil {
				for _, w := range found {
					assert.True(t, tt.expectedFilter(w))
				}
			}
		})
	}
}

// Helper functions to create pointers
func float64Ptr(v float64) *float64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func trailDifficultyPtr(v models.TrailDifficulty) *models.TrailDifficulty {
	return &v
}
