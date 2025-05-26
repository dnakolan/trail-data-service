package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTrail(t *testing.T) {
	name := "Test Trail"
	lat := 45.5231
	lon := -122.6765
	difficulty := TrailDifficultyMedium
	lengthKm := 10.5

	trail := NewTrail(name, lat, lon, difficulty, lengthKm)

	assert.NotNil(t, trail)
	assert.NotNil(t, trail.UID)
	assert.Equal(t, name, *trail.Name)
	assert.Equal(t, lat, *trail.Lat)
	assert.Equal(t, lon, *trail.Lon)
	assert.Equal(t, difficulty, *trail.Difficulty)
	assert.Equal(t, lengthKm, *trail.LengthKm)
	assert.Equal(t, time.Time{}, *trail.CreatedAt)
}

func TestNewTrailFromRequest(t *testing.T) {
	name := "Test Trail"
	lat := 45.5231
	lon := -122.6765
	difficulty := TrailDifficultyMedium
	lengthKm := 10.5

	req := &CreateTrailRequest{
		Name:       &name,
		Lat:        &lat,
		Lon:        &lon,
		Difficulty: &difficulty,
		LengthKm:   &lengthKm,
	}

	trail := NewTrailFromRequest(req)

	assert.NotNil(t, trail)
	assert.NotNil(t, trail.UID)
	assert.Equal(t, name, *trail.Name)
	assert.Equal(t, lat, *trail.Lat)
	assert.Equal(t, lon, *trail.Lon)
	assert.Equal(t, difficulty, *trail.Difficulty)
	assert.Equal(t, lengthKm, *trail.LengthKm)
	assert.Equal(t, time.Time{}, *trail.CreatedAt)
}

func TestTrail_Validate(t *testing.T) {
	validName := "Test Trail"
	validLat := 45.5231
	validLon := -122.6765
	validDifficulty := TrailDifficultyMedium
	validLength := 10.5

	tests := []struct {
		name          string
		trail         *Trail
		expectedError string
	}{
		{
			name: "valid trail",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "",
		},
		{
			name: "missing name",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Lat:        &validLat,
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail name is required",
		},
		{
			name: "empty name",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       stringPtr(""),
					Lat:        &validLat,
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail name is required",
		},
		{
			name: "missing latitude",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail start latitude is required",
		},
		{
			name: "invalid latitude (too high)",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        float64Ptr(91.0),
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail start latitude must be between -90 and 90",
		},
		{
			name: "invalid latitude (too low)",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        float64Ptr(-91.0),
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail start latitude must be between -90 and 90",
		},
		{
			name: "missing longitude",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail start longitude is required",
		},
		{
			name: "invalid longitude (too high)",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Lon:        float64Ptr(181.0),
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail start longitude must be between -180 and 180",
		},
		{
			name: "invalid longitude (too low)",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Lon:        float64Ptr(-181.0),
					Difficulty: &validDifficulty,
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail start longitude must be between -180 and 180",
		},
		{
			name: "missing difficulty",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:     &validName,
					Lat:      &validLat,
					Lon:      &validLon,
					LengthKm: &validLength,
				},
			},
			expectedError: "trail difficulty is required",
		},
		{
			name: "invalid difficulty",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Lon:        &validLon,
					Difficulty: trailDifficultyPtr("invalid"),
					LengthKm:   &validLength,
				},
			},
			expectedError: "trail difficulty must be easy, medium, or hard",
		},
		{
			name: "missing length",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Lon:        &validLon,
					Difficulty: &validDifficulty,
				},
			},
			expectedError: "trail length is required",
		},
		{
			name: "negative length",
			trail: &Trail{
				CreateTrailRequest: CreateTrailRequest{
					Name:       &validName,
					Lat:        &validLat,
					Lon:        &validLon,
					Difficulty: &validDifficulty,
					LengthKm:   float64Ptr(-1.0),
				},
			},
			expectedError: "trail length must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.trail.Validate()
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestTrailFilter_Validate(t *testing.T) {
	validLat := 45.5231
	validLon := -122.6765
	validRadius := 10.0

	tests := []struct {
		name          string
		filter        *TrailFilter
		expectedError string
	}{
		{
			name:          "empty filter",
			filter:        &TrailFilter{},
			expectedError: "",
		},
		{
			name: "valid proximity filter",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: &validLat,
					Lon: &validLon,
				},
				RadiusKm: &validRadius,
			},
			expectedError: "",
		},
		{
			name: "missing longitude",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: &validLat,
				},
				RadiusKm: &validRadius,
			},
			expectedError: "invalid missing lon filter",
		},
		{
			name: "missing latitude",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lon: &validLon,
				},
				RadiusKm: &validRadius,
			},
			expectedError: "invalid missing lat filter",
		},
		{
			name: "missing radius",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: &validLat,
					Lon: &validLon,
				},
			},
			expectedError: "invalid missing radius filter",
		},
		{
			name: "invalid latitude (too high)",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: float64Ptr(91.0),
					Lon: &validLon,
				},
				RadiusKm: &validRadius,
			},
			expectedError: "invalid lat filter outside of bounds -90 to 90",
		},
		{
			name: "invalid longitude (too high)",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: &validLat,
					Lon: float64Ptr(181.0),
				},
				RadiusKm: &validRadius,
			},
			expectedError: "invalid lon filter outside of bounds -180 to 180",
		},
		{
			name: "negative radius",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: &validLat,
					Lon: &validLon,
				},
				RadiusKm: float64Ptr(-1.0),
			},
			expectedError: "invalid radius filter must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestTrail_MatchesFilter(t *testing.T) {
	trail := NewTrail("Test Trail", 45.5231, -122.6765, TrailDifficultyMedium, 10.5)

	mediumDifficulty := TrailDifficultyMedium
	hardDifficulty := TrailDifficultyHard

	tests := []struct {
		name     string
		filter   *TrailFilter
		expected bool
	}{
		{
			name:     "nil filter",
			filter:   nil,
			expected: true,
		},
		{
			name: "matching name",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Name: stringPtr("Test Trail"),
				},
			},
			expected: true,
		},
		{
			name: "non-matching name",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Name: stringPtr("Different Trail"),
				},
			},
			expected: false,
		},
		{
			name: "matching difficulty",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Difficulty: &mediumDifficulty,
				},
			},
			expected: true,
		},
		{
			name: "non-matching difficulty",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Difficulty: &hardDifficulty,
				},
			},
			expected: false,
		},
		{
			name: "within radius",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: float64Ptr(45.5231),
					Lon: float64Ptr(-122.6765),
				},
				RadiusKm: float64Ptr(1.0),
			},
			expected: true,
		},
		{
			name: "outside radius",
			filter: &TrailFilter{
				CreateTrailRequest: CreateTrailRequest{
					Lat: float64Ptr(46.5231),
					Lon: float64Ptr(-123.6765),
				},
				RadiusKm: float64Ptr(1.0),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trail.MatchesFilter(tt.filter)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidTrailDifficulty(t *testing.T) {
	tests := []struct {
		difficulty string
		expected   bool
	}{
		{"easy", true},
		{"medium", true},
		{"hard", true},
		{"invalid", false},
		{"", false},
		{"EASY", false},
		{"Medium", false},
	}

	for _, tt := range tests {
		t.Run(tt.difficulty, func(t *testing.T) {
			result := IsValidTrailDifficulty(tt.difficulty)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func trailDifficultyPtr(d TrailDifficulty) *TrailDifficulty {
	return &d
}
