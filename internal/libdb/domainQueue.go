package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const DefaultDomainQueuePriority = 0

var (
	ErrDomainAlreadyQueued = errors.New("db: domain already queued")
	ErrDomainInBlocklist   = errors.New("db: domain in blocklist")
)

func (db *DB) AddDomainToQueue(domain string, path string, priority int) (*uuid.UUID, error) {
	bl, err := db.QueryDomainBlocklistByDomain(domain)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if bl != nil {
		return nil, ErrDomainInBlocklist
	}

	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer smartRollback(tx)

	id := uuid.New()
	if _, err := tx.Exec(`INSERT INTO "domain_queue"("id", "domain", "priority", "start") VALUES($1, $2, $3, $4)`, id, domain, priority, path); err != nil {
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
	ID        uuid.UUID
	CreatedAt time.Time
	Domain    string
	Priority  int
	StartPath string
}

func (db *DB) QueryDomainQueue(id uuid.UUID) (*QueueItem, error) {
	ctx, cancel := db.newContext()
	defer cancel()

	qi := new(QueueItem)

	if err := db.pool.QueryRowContext(
		ctx,
		`SELECT "id", "created_at", "domain", "priority", "start" FROM "domain_queue" WHERE "id" = $1;`,
		id,
	).Scan(&qi.ID, &qi.CreatedAt, &qi.Domain, &qi.Priority, &qi.StartPath); err != nil {
		return nil, errors.Wrap(err, "could not scan")
	}

	return qi, nil
}
