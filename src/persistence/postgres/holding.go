package postgres

import (
	"github.com/google/uuid"
	"github.com/orbis-challenge/src/models"
)

func (q DBQuery) SaveHolding(holding *models.Holding) (models.Holding, error) {
	_, err := q.Model(holding).
		Returning("*").
		Insert()

	return *holding, err
}

func (q DBQuery) DeleteHoldingByID(id uuid.UUID) (err error) {
	h := models.Holding{ID: id}
	_, err = q.Model(&h).
		Delete()
	return err
}
