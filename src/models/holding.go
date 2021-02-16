package models

import (
	"time"

	"github.com/google/uuid"
)

// Etf model.
type Holding struct {
	tableName struct{} `sql:"holdings" pg:",discard_unknown_columns"` // nolint

	ID     uuid.UUID `json:"id,omitempty" pg:"id,pk"`
	Name   string    `json:"name"`
	Weight float64   `json:"weight"`
	EtfsID uuid.UUID `json:"etfs_id"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
