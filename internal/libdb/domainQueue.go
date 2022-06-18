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

type QueueItem struct {
	ID uuid.UUID
	CreatedAt time.Time
	Domain string
}

func (db *DB) DomainQueueFetch(id uuid.UUID) (*QueueItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stmt, err := db.pool.PrepareContext(ctx, `SELECT "id", "created_at", "domain" FROM "domain_queue" WHERE "id" = $1;`)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	qi := new(QueueItem)

	if err := stmt.QueryRow(id).Scan(&qi.ID, &qi.CreatedAt, &qi.Domain); err != nil {
		return nil, errors.Wrap(err, "could not scan")
	}

	return qi, nil
}

type CurrentJob struct {
	ID uuid.UUID
	QueueItem uuid.UUID
	WorkerID string
	LastChecKInTime time.Time
}

var (
	ErrNoQueuedDomains = errors.New("db: no queued domains")
	ErrWorkerIDInUse = errors.New("db: worker ID in use")
)

func (db *DB) RequestJob(workerID string) (*CurrentJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer smartRollback(tx)

	newID := uuid.New()

	stmt, err := tx.Prepare(`INSERT INTO "current_jobs"("id", "queue_item", "worker_id")
	VALUES ($1,
			(SELECT "id"
			 FROM "domain_queue"
			 WHERE "id" NOT IN (SELECT "queue_item" FROM "current_jobs")
			 ORDER BY "created_at"
			 LIMIT 1), $2)
	RETURNING "current_jobs"."queue_item";`)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	var queueItem uuid.UUID
	row := stmt.QueryRow(newID, workerID)
	if err := row.Scan(&queueItem); err != nil {
		if e, ok := err.(*pq.Error); ok {
			// If the subquery returns no results, it'll fail with this error
			// because it returns null
			switch e.Code{
			case errorCodeNotNullViolation:
				return nil, ErrNoQueuedDomains
			case errorCodeUniqueViolation:
				if e.Constraint == "current_jobs_worker_id_key" {
					return nil, ErrWorkerIDInUse
				}
			}
		}
		return nil, errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &CurrentJob{
		ID: newID,
		WorkerID: workerID,
		QueueItem: queueItem,
	}, nil
}