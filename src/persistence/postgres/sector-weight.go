package postgres

import (
	"github.com/google/uuid"
	"github.com/orbis-challenge/src/models"
)

func (q DBQuery) SaveSectorWeight(sectorWeight *models.SectorWeight) (*models.SectorWeight, error) {
	_, err := q.Model(sectorWeight).
		Returning("*").
		Insert()

	return sectorWeight, err
}

func (q DBQuery) DeleteSectorWeightByID(id uuid.UUID) (err error) {
	sw := models.SectorWeight{ID: id}
	_, err = q.Model(&sw).
		Delete()
	return err
}
