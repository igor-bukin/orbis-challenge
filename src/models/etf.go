package models

import (
	"time"

	"github.com/google/uuid"
)

// Etf model.
type Etf struct {
	tableName struct{} `sql:"etfs" pg:",discard_unknown_columns"` // nolint

	ID      uuid.UUID `json:"id,omitempty" pg:"id,pk"`
	Name    string    `json:"name"`
	Ticker  string    `json:"ticker"`
	FundURI string    `json:"fund_uri"`

	Holdings      []*Holding      `json:"holdings" pg:"fk:etfs_id"`
	SectorWeights []*SectorWeight `json:"sector_weights" pg:"fk:etfs_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type EtfResponse struct {
	Name          string          `json:"name"`
	Ticker        string          `json:"ticker"`
	Holdings      []*Holding      `json:"holdings"`
	SectorWeights []*SectorWeight `json:"sectorWeights"`
}

func (e *Etf) Convert2Response() EtfResponse {
	return EtfResponse{
		Name:          e.Name,
		Ticker:        e.Ticker,
		Holdings:      e.Holdings,
		SectorWeights: e.SectorWeights,
	}
}
