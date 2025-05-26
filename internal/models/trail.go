package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TrailDifficulty string

const (
	TrailDifficultyEasy   TrailDifficulty = "easy"
	TrailDifficultyMedium TrailDifficulty = "medium"
	TrailDifficultyHard   TrailDifficulty = "hard"
)

type Trail struct {
	UID        uuid.UUID        `json:"trail_id"`
	Name       *string          `json:"trail_name"`
	LatStart   *float64         `json:"lat_start"`
	LonStart   *float64         `json:"lon_start"`
	Difficulty *TrailDifficulty `json:"difficulty"`
	LengthKm   *float64         `json:"length_km"`
	CreatedAt  *time.Time       `json:"created_at"`
}

type CreateTrailRequest struct {
	Name       *string          `json:"trail_name"`
	LatStart   *float64         `json:"lat_start"`
	LonStart   *float64         `json:"lon_start"`
	Difficulty *TrailDifficulty `json:"difficulty"`
	LengthKm   *float64         `json:"length_km"`
}

func (t *Trail) Validate() error {
	if t.Name == nil || *t.Name == "" {
		return errors.New("trail name is required")
	}
	if t.LatStart == nil {
		return errors.New("trail start latitude is required")
	}
	if *t.LatStart < -90 || *t.LatStart > 90 {
		return errors.New("trail start latitude must be between -90 and 90")
	}
	if t.LonStart == nil {
		return errors.New("trail start longitude is required")
	}
	if *t.LonStart < -180 || *t.LonStart > 180 {
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
		Name:       t.Name,
		LatStart:   t.LatStart,
		LonStart:   t.LonStart,
		Difficulty: t.Difficulty,
		LengthKm:   t.LengthKm,
	}
	return trail.Validate()
}
