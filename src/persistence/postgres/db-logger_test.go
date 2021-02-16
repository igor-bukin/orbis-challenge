package postgres

import (
	"context"
	"testing"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type (
	log     struct{}
	logType byte
)

const (
	_ logType = iota
	none
	before
	after
)

var actual = none

func (log) Debug(args ...interface{})              { actual = before }
func (log) Debugf(fmt string, args ...interface{}) { actual = after }

func TestBeforeQuery(t *testing.T) {
	log := dbLogger{
		logger: log{},
	}

	input := &pg.QueryEvent{
		Ctx: context.Background(),
	}

	log.BeforeQuery(input)
	assert.EqualValues(t, before, actual)
}

func TestAfterQuery(t *testing.T) {
	log := dbLogger{
		logger: log{},
	}

	input := &pg.QueryEvent{
		Error: errors.New("test error"),
		Ctx:   context.Background(),
	}

	log.BeforeQuery(input)
	log.AfterQuery(input)
	assert.EqualValues(t, after, actual)
}
