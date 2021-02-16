// nolint
package holding

import (
	"context"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s service) Update() httperrors.HTTPError {
	ctx := context.Background()
	tx, err := s.pg.NewTXContext(ctx)
	if err != nil {
		return httperrors.NewInternalServerError(errors.Wrap(err, "start transaction"))
	}

	defer tx.Rollback() // nolint

	limit := 20
	offset := 0

	c, err := tx.GetCountOfETF()
	if err != nil {
		return httperrors.NewInternalServerError(err)
	}

	for l := 0; limit < c; l += 20 {
		res, _, err := tx.GetPaginationETFs(l, offset)
		if err != nil {
			return httperrors.NewInternalServerError(err)
		}

		for i := range res {
			for j := range res[i].SectorWeights {
				errD := tx.DeleteSectorWeightByID(res[j].SectorWeights[j].ID)
				if errD != nil {
					return httperrors.NewInternalServerError(errD)
				}
			}

			for j := range res[i].Holdings {
				errD := tx.DeleteHoldingByID(res[j].Holdings[j].ID)
				if errD != nil {
					return httperrors.NewInternalServerError(errD)
				}
			}

			httpErr := s.getAndSave(tx, res[i].FundURI, res[i].ID)
			if httpErr != nil {
				logrus.Debug("couldn't get and save holdings", "uri", res[i].FundURI)
				continue
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return httperrors.NewInternalServerError(err)
	}

	return nil
}

func (s service) updateSchedule() {
	for {
		select {
		case <-s.ticker.C:
			if err := s.Update(); err != nil {
				logrus.Error("couldn't update holdings", "err", err)
				continue
			}
		}
	}
}
