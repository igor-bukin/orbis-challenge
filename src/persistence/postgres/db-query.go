package postgres

import (
	"context"
	"errors"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// DBModel interface.
type DBModel interface {
	Model(model ...interface{}) *orm.Query
	Exec(query interface{}, params ...interface{}) (pg.Result, error)
	Query(model, query interface{}, params ...interface{}) (pg.Result, error)
}

// DBQuery model.
type DBQuery struct {
	DBModel
	completed bool
}

// Rollback rollbacks query if it was transaction or returns error
// nolint:gocritic
func (q DBQuery) Rollback() error {
	switch t := q.DBModel.(type) {
	case *pg.Tx:
		if !q.completed {
			return t.Rollback()
		}
		return nil
	}

	return errors.New("rollback failed: not in Tx")
}

// Commit makes commit on transaction or do nothing
// nolint:gocritic
func (q *DBQuery) Commit() error {
	switch t := q.DBModel.(type) {
	case *pg.Tx:
		if !q.completed {
			q.completed = true

			return t.Commit()
		}
	}

	return nil
}

// NewTXContext returns DBQuery instance with new transaction
// DBQuery.Commit() must be called to run transaction
func (p Client) NewTXContext(ctx context.Context) (DBQuery, error) {
	tx, err := p.conn.WithContext(ctx).Begin()
	return DBQuery{DBModel: tx}, err
}

// QueryContext returns DBQuery instance of current db pool
func (p Client) QueryContext(ctx context.Context) DBQuery {
	return DBQuery{DBModel: p.conn.WithContext(ctx)}
}
