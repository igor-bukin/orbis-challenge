package postgres

import (
	"context"
	"time"

	"github.com/go-pg/pg"
)

type startedKey string

const (
	execStartedKey startedKey = "exec_started"
)

// DebugLogger interface
type DebugLogger interface {
	Debug(args ...interface{})
	Debugf(fmt string, args ...interface{})
}

type dbLogger struct {
	logger DebugLogger
}

// BeforeQuery logs before query execution
func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
	if q.Ctx == nil {
		q.Ctx = context.Background()
	}
	q.Ctx = context.WithValue(q.Ctx, execStartedKey, time.Now())
	d.logger.Debug(q.FormattedQuery())
}

// AfterQuery logs after query execution
func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	if q.Ctx == nil {
		q.Ctx = context.Background()
	}
	execTime := time.Since(q.Ctx.Value(execStartedKey).(time.Time)).Round(time.Millisecond)
	if q.Error != nil {
		d.logger.Debugf("Query execution error: %s, exec time: %s", q.Error.Error(), execTime)
	}

	if q.Result != nil {
		d.logger.Debugf("Rows affected: %d, Rows returned: %d, exec time: %s",
			q.Result.RowsAffected(), q.Result.RowsReturned(), execTime)
	}
}
