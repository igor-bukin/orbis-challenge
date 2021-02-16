package postgres

import (
	"github.com/google/uuid"
	"github.com/orbis-challenge/src/models"
)

func (q DBQuery) SaveETF(etf *models.Etf) (models.Etf, error) {
	etf.ID = uuid.UUID{}
	_, err := q.Model(etf).
		Returning("*").
		Insert()

	return *etf, err
}

func (q DBQuery) GetPaginationETFs(limit, offset int) ([]*models.Etf, int, error) {
	res := []*models.Etf{}
	c, err := q.Model(&res).
		Column("*").
		Relation("Holdings").
		Relation("SectorWeights").
		Limit(limit).
		Offset(offset).
		SelectAndCount()
	if err != nil {
		return nil, -1, err
	}

	return res, c, nil
}

func (q DBQuery) GetETFByTicker(ticker string) (res models.Etf, err error) {
	err = q.Model(&res).
		Column("*").
		Relation("Holdings").
		Relation("SectorWeights").
		Where(`ticker = ?`, ticker).
		First()
	return res, err
}

func (q DBQuery) GetCountOfETF() (c int, err error) {
	e := models.Etf{}
	return q.Model(&e).
		Count()
}
