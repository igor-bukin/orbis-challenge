package etf

import (
	"net/http"

	"github.com/orbis-challenge/src/handlers/common"
	"github.com/orbis-challenge/src/models"
	"github.com/orbis-challenge/src/services"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// Load - load top of etfs
func Load(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	httpErr := services.Get().Holding().Load(ctx)
	if httpErr != nil {
		logrus.Error("couldn't load list of etfs",
			zap.Error(httpErr))
		common.SendHTTPError(w, httpErr)
		return
	}

	common.SendResponse(w, http.StatusNoContent, &struct{}{})
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	limit, offset, httpErr := common.GetLimitAndOffset(r.URL.Query())
	if httpErr != nil {
		common.SendHTTPError(w, httpErr)
		return
	}

	res, c, httpErr := services.Get().Holding().GetAll(ctx, limit, offset)
	if httpErr != nil {
		logrus.Error("couldn't get list of etfs",
			zap.Error(httpErr))
		common.SendHTTPError(w, httpErr)
		return
	}

	common.SendResponse(w, http.StatusOK, struct {
		Count int                  `json:"count"`
		Etfs  []models.EtfResponse `json:"etfs"`
	}{
		Count: c,
		Etfs:  res,
	})
}

func GetByTicker(w http.ResponseWriter, r *http.Request) {
	var (
		ticker = r.URL.Query().Get("ticker")
		ctx    = r.Context()
	)

	res, httpErr := services.Get().Holding().GetByTicker(ctx, ticker)
	if httpErr != nil {
		logrus.Error("couldn't get etf by ticker",
			zap.Error(httpErr))
		common.SendHTTPError(w, httpErr)
		return
	}

	common.SendResponse(w, http.StatusOK, res)
}
