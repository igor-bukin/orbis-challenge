package holding

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbis-challenge/src/models"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/orbis-challenge/src/persistence/postgres"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const zeroWeight = 0.0

func (s service) Load(ctx context.Context) httperrors.HTTPError {
	etfs, err := s.parser.GetEtfs()
	if err != nil {
		return httperrors.NewInternalServerError(err)
	}
	tx, err := s.pg.NewTXContext(ctx)
	if err != nil {
		return httperrors.NewInternalServerError(errors.Wrap(err, "start transaction"))
	}

	for i := range etfs.Data.Funds.Etfs.Datas {
		etfRaw := etfs.Data.Funds.Etfs.Datas[i]
		etf := models.Etf{
			Name:    etfRaw.FundName,
			Ticker:  etfRaw.FundTicker,
			FundURI: etfRaw.FundURI,
		}
		savedETF, err := tx.SaveETF(&etf)
		if err != nil {
			return httperrors.NewInternalServerError(err)
		}

		httpErr := s.getAndSave(tx, etfRaw.FundURI, savedETF.ID)
		if httpErr != nil {
			logrus.Error("couldn't get and save holdings", "uri", etfRaw.FundURI)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return httperrors.NewInternalServerError(err)
	}

	return nil
}

func (s service) getAndSave(tx postgres.DBQuery, uri string, etfID uuid.UUID) httperrors.HTTPError {
	url := fmt.Sprintf("https://www.ssga.com%s", uri)
	h, w, err := s.parser.GetTopHoldingsAndWeight(url)
	if err != nil {
		logrus.Debug("couldn't get holdings", "url", url)
		return httperrors.NewInternalServerError(err)
	}
	for j := range h {
		h[j].EtfsID = etfID
		if h[j].Weight <= zeroWeight {
			_, errHolding := tx.SaveHolding(&h[j])
			if errHolding != nil {
				return httperrors.NewInternalServerError(errHolding)
			}
		}
	}

	for j := range w {
		w[j].EtfsID = etfID
		_, errSW := tx.SaveSectorWeight(&w[j])
		if errSW != nil {
			return httperrors.NewInternalServerError(errSW)
		}
	}

	return nil
}
