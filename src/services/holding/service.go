package holding

import (
	"context"
	"sync"
	"time"

	"github.com/orbis-challenge/src/config"
	"github.com/orbis-challenge/src/integrations/ssga"
	"github.com/orbis-challenge/src/models"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/orbis-challenge/src/persistence/postgres"
)

const (
	durationUpdateHour = 12
)

type Service interface {
	Load(ctx context.Context) httperrors.HTTPError
	GetAll(ctx context.Context, limit, offset int) ([]models.EtfResponse, int, httperrors.HTTPError)
	GetByTicker(ctx context.Context, ticker string) (models.EtfResponse, httperrors.HTTPError)
}

type service struct {
	pg     *postgres.Client
	parser ssga.Parser
	ticker *time.Ticker
}

var (
	srv  Service
	once = &sync.Once{}
)

// New returns service instance.
func New(pg *postgres.Client, cfg *config.SSGA) Service {
	once.Do(func() {
		parser := ssga.New(cfg)
		instance := service{
			pg:     pg,
			parser: parser,
			ticker: time.NewTicker(time.Hour * durationUpdateHour),
		}
		go instance.updateSchedule()

		srv = instance
	})
	return srv
}
