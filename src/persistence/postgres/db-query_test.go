package postgres

import (
	"context"
	"testing"

	"github.com/go-pg/pg"
	"github.com/orbis-challenge/src/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRollback(t *testing.T) {
	tc := []struct {
		name    string
		query   DBQuery
		isValid bool
	}{
		{
			name: "Empty model",
			query: DBQuery{
				DBModel:   &pg.Tx{},
				completed: true,
			},
			isValid: true,
		},
		{
			name:  "Empty model",
			query: DBQuery{},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.query.Rollback()
			if test.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tc := []struct {
		name    string
		query   DBQuery
		isValid bool
	}{
		{
			name: "Not completed transaction",
			query: DBQuery{
				DBModel:   &pg.Tx{},
				completed: true,
			},
			isValid: true,
		},
		{
			name: "Completed transaction",
			query: DBQuery{
				completed: true,
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.query.Commit()
			assert.NoError(t, err)
		})
	}
}

func TestNewTxContext(t *testing.T) {
	err := Load(&config.Config.PostgresTest, logrus.New())
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()
	actual, err := postgres.NewTXContext(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func TestQueryContext(t *testing.T) {
	err := Load(&config.Config.PostgresTest, logrus.New())
	if !assert.NoError(t, err) {
		return
	}

	actual := postgres.QueryContext(context.Background())
	assert.NotNil(t, actual)
}
