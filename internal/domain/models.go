package domain

import (
	"time"

	"github.com/google/uuid"
)

type Junction struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type Signal struct {
	ID              uuid.UUID `json:"id"`
	JunctionID      uuid.UUID `json:"junction_id"`
	Direction       string    `json:"direction"`
	State           string    `json:"state"`
	DurationSeconds int       `json:"duration_seconds"`
	UpdatedAt       time.Time `json:"updated_at"`
}
