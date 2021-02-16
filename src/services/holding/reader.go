package holding

import (
	"context"

	"github.com/orbis-challenge/src/models"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
)

func (s service) GetAll(ctx context.Context, limit, offset int) ([]models.EtfResponse, int, httperrors.HTTPError) {
	tx := s.pg.QueryContext(ctx)
	etfs, count, err := tx.GetPaginationETFs(limit, offset)
	if err != nil {
		return nil, -1, httperrors.NewInternalServerError(err)
	}
	res := make([]models.EtfResponse, len(etfs))
	for i := range etfs {
		res[i] = etfs[i].Convert2Response()
	}

	return res, count, nil
}

func (s service) GetByTicker(ctx context.Context, ticker string) (models.EtfResponse, httperrors.HTTPError) {
	tx := s.pg.QueryContext(ctx)
	etfs, err := tx.GetETFByTicker(ticker)
	if err != nil {
		return models.EtfResponse{}, httperrors.NewInternalServerError(err)
	}

	return etfs.Convert2Response(), nil
}
