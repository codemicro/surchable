package db

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type CurrentJob struct {
	ID              uuid.UUID
	QueueItem       uuid.UUID
	WorkerID        string
	LastChecKInTime time.Time
}

var (
	ErrNoQueuedDomains = errors.New("db: no queued domains")
	ErrCrawlerIDInUse  = errors.New("db: crawler ID in use")
)

func (db *DB) RequestJobForCrawler(workerID string) (*CurrentJob, error) {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer smartRollback(tx)

	newID := uuid.New()

	var queueItem uuid.UUID
	row := tx.QueryRow(`INSERT INTO "current_jobs"("id", "queue_item", "crawler_id")
	VALUES ($1,
			(SELECT "id"
			 FROM "domain_queue"
			 WHERE "id" NOT IN (SELECT "queue_item" FROM "current_jobs")
			 ORDER BY "priority" DESC, "created_at"
			 LIMIT 1), $2)
	RETURNING "current_jobs"."queue_item";`, newID, workerID)
	if err := row.Scan(&queueItem); err != nil {
		if e, ok := err.(*pq.Error); ok {
			// If the subquery returns no results, it'll fail with this error
			// because it returns null
			switch e.Code {
			case errorCodeNotNullViolation:
				return nil, ErrNoQueuedDomains
			case errorCodeUniqueViolation:
				if e.Constraint == "current_jobs_crawler_id_key" {
					return nil, ErrCrawlerIDInUse
				}
			}
		}
		return nil, errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &CurrentJob{
		ID:        newID,
		WorkerID:  workerID,
		QueueItem: queueItem,
	}, nil
}

func (db *DB) UpdateTimeForJobByWorkerID(workerID string, checkInTime time.Time) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	if _, err := tx.Exec(
		`UPDATE "current_jobs" SET "last_check_in" = $1 WHERE "crawler_id" = $2;`,
		checkInTime,
		workerID,
	); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

func (db *DB) RemoveTimedOutJobs() error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	rows, err := tx.Query(`SELECT "id", "queue_item" FROM "current_jobs" WHERE "last_check_in" < now()-'10 minute'::interval;`)
	if err != nil {
		return errors.WithStack(err)
	}

	var ids [][2]uuid.UUID

	for rows.Next() {
		var jobID, queueItemID uuid.UUID
		if err := rows.Scan(&jobID, &queueItemID); err != nil {
			return errors.WithStack(err)
		}

		ids = append(ids, [2]uuid.UUID{jobID, queueItemID})
	}

	for _, item := range ids {
		if _, err := tx.Exec(`UPDATE "domain_queue" SET "priority" = 5 WHERE "id" = $1;`, item[1]); err != nil {
			return errors.WithStack(err)
		}

		if _, err := tx.Exec(`DELETE FROM "current_jobs" WHERE "id" = $1;`, item[0]); err != nil {
			return errors.WithStack(err)
		}
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

var ErrNoActiveJob = errors.New("db: no active jobs with a matching crawler ID")

func (db *DB) CompleteJobByCrawlerID(crawlerID string) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	mods, err := tx.Exec(
		`WITH job_info AS (
			DELETE FROM "current_jobs" WHERE "crawler_id" = $1 RETURNING "queue_item"
			) DELETE FROM "domain_queue" WHERE "id" = (SELECT "queue_item" FROM "job_info");`,
		crawlerID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if n, err := mods.RowsAffected(); err != nil {
		return errors.WithStack(err)
	} else if n == 0 {
		return ErrNoActiveJob
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

func (db *DB) CancelJobByCrawlerID(crawlerID string) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	mods, err := tx.Exec(
		`WITH job_info AS (
			DELETE FROM "current_jobs" WHERE "crawler_id" = $1 RETURNING "queue_item"
			) UPDATE "domain_queue" SET "priority" = 3 WHERE "id" = (SELECT "queue_item" FROM "job_info");`,
		crawlerID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if n, err := mods.RowsAffected(); err != nil {
		return errors.WithStack(err)
	} else if n == 0 {
		return ErrNoActiveJob
	}

	return errors.WithStack(
		tx.Commit(),
	)
}
