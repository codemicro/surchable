package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var ErrDomainAlreadyQueued = errors.New("db: domain already queued")

func (db *DB) DomainQueueInsert(domain string) (*uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer smartRollback(tx)

	stmt, err := tx.Prepare(`INSERT INTO "domain_queue"("id", "domain") VALUES($1, $2)`)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	id := uuid.New()
	if _, err := stmt.Exec(id, domain); err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == errorCodeUniqueViolation {
				return nil, ErrDomainAlreadyQueued
			}
		}
		return nil, errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &id, nil
}
