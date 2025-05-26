package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/umahmood/haversine"
)

type TrailDifficulty string

const (
	TrailDifficultyEasy   TrailDifficulty = "easy"
	TrailDifficultyMedium TrailDifficulty = "medium"
	TrailDifficultyHard   TrailDifficulty = "hard"
)

type CreateTrailRequest struct {
	Name       *string          `json:"name"`
	Lat        *float64         `json:"lat"`
	Lon        *float64         `json:"lon"`
	Difficulty *TrailDifficulty `json:"difficulty"`
	LengthKm   *float64         `json:"length_km"`
}

type Trail struct {
	CreateTrailRequest
	UID       uuid.UUID  `json:"trail_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type TrailFilter struct {
	CreateTrailRequest
	RadiusKm *float64 `json:"radius_km"`
}

func (t *Trail) Validate() error {
	if t.Name == nil || *t.Name == "" {
		return errors.New("trail name is required")
	}
	if t.Lat == nil {
		return errors.New("trail start latitude is required")
	}
	if *t.Lat < -90 || *t.Lat > 90 {
		return errors.New("trail start latitude must be between -90 and 90")
	}
	if t.Lon == nil {
		return errors.New("trail start longitude is required")
	}
	if *t.Lon < -180 || *t.Lon > 180 {
		return errors.New("trail start longitude must be between -180 and 180")
	}
	if t.Difficulty == nil {
		return errors.New("trail difficulty is required")
	}
	if *t.Difficulty != TrailDifficultyEasy && *t.Difficulty != TrailDifficultyMedium && *t.Difficulty != TrailDifficultyHard {
		return errors.New("trail difficulty must be easy, medium, or hard")
	}
	if t.LengthKm == nil {
		return errors.New("trail length is required")
	}
	if *t.LengthKm < 0 {
		return errors.New("trail length must be positive")
	}
	return nil
}

func (t *CreateTrailRequest) Validate() error {
	trail := &Trail{
		CreateTrailRequest: *t,
	}
	return trail.Validate()
}

func (t *TrailFilter) Validate() error {
	if t.Name != nil && *t.Name == "" {
		return errors.New("invalid empty name filter")
	}

	if t.Lat != nil && t.Lon != nil && t.RadiusKm != nil {
		if *t.Lat < -90 || *t.Lat > 90 {
			return errors.New("invalid lat filter outside of bounds -90 to 90")
		} else if *t.Lon < -180 || *t.Lon > 180 {
			return errors.New("invalid lon filter outside of bounds -180 to 180")
		} else if *t.RadiusKm < 0 {
			return errors.New("invalid radius filter must be positive")
		}
	} else if t.Lat != nil && t.Lon != nil {
		return errors.New("invalid missing radius filter")
	} else if t.Lat != nil {
		return errors.New("invalid missing lon filter")
	} else if t.Lon != nil {
		return errors.New("invalid missing lat filter")
	} else if t.RadiusKm != nil {
		return errors.New("invalid missing lat or lon filter")
	}

	if t.Difficulty != nil {
		if *t.Difficulty != TrailDifficultyEasy && *t.Difficulty != TrailDifficultyMedium && *t.Difficulty != TrailDifficultyHard {
			return errors.New("trail difficulty must be easy, medium, or hard")
		}
	}

	if t.LengthKm != nil {
		if *t.LengthKm < 0 {
			return errors.New("invalid length filter must be positive")
		}
	}

	return nil
}

func (t *Trail) MatchesFilter(filter *TrailFilter) bool {
	if filter.Name != nil && *filter.Name != "" && *filter.Name != *t.Name {
		return false
	}
	if filter.Lat != nil && filter.Lon != nil && filter.RadiusKm != nil {
		_, distance := haversine.Distance(
			haversine.Coord{Lat: *filter.Lat, Lon: *filter.Lon},
			haversine.Coord{Lat: *t.Lat, Lon: *t.Lon},
		)
		if distance > *filter.RadiusKm {
			return false
		}
	}
	if filter.Difficulty != nil && *filter.Difficulty != *t.Difficulty {
		return false
	}
	if filter.LengthKm != nil && *filter.LengthKm != *t.LengthKm {
		return false
	}
	return true
}

func IsValidTrailDifficulty(s string) bool {
	switch TrailDifficulty(s) {
	case TrailDifficultyEasy, TrailDifficultyMedium, TrailDifficultyHard:
		return true
	default:
		return false
	}
}
