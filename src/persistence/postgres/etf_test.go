package postgres

import (
	"context"
	"testing"

	"github.com/orbis-challenge/src/models"
	"github.com/orbis-challenge/src/persistence/postgres/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDBQuery_SaveETF(t *testing.T) {
	fixtures.PrepareFixtures()
	etf := models.Etf{
		Name:    "Test",
		Ticker:  "TickerTest",
		FundURI: "www.test.com",
	}

	tx := GetDB().QueryContext(context.Background())

	etfN, err := tx.SaveETF(&etf)
	assert.NoError(t, err)
	assert.Equal(t, etfN.Name, etf.Name)
	assert.Equal(t, etfN.Ticker, etf.Ticker)
	assert.Equal(t, etfN.FundURI, etf.FundURI)
}
