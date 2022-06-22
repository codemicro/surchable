package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/codemicro/surchable/internal/util"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PageLoad struct {
	ID            uuid.UUID
	URL           string
	NormalisedURL string
	LoadedAt      time.Time
	NotLoadBefore *time.Time
}

var ErrNoMatchingPageLoad = errors.New("db: no page loads matching that URL")

func (db *DB) QueryPageLoadsByURL(url string) (*PageLoad, error) {
	normalisedURL, err := util.NormaliseURL(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pageLoad := new(PageLoad)

	// url is unique, hence QueryRowContext and not QueryContext
	err = db.pool.QueryRowContext(
		ctx,
		`SELECT "id", "url", "normalised_url", "loaded_at", "not_load_before" FROM "page_loads" WHERE "normalised_url" = $1;`,
		normalisedURL,
	).Scan(&pageLoad.ID, &pageLoad.URL, &pageLoad.NormalisedURL, &pageLoad.LoadedAt, &pageLoad.NotLoadBefore)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoMatchingPageLoad
		}
		return nil, errors.WithStack(err)
	}

	return pageLoad, nil
}

func (db *DB) UpsertPageLoad(pl *PageLoad) (uuid.UUID, error) {
	if pl.ID == uuid.Nil {
		pl.ID = uuid.New()
	}

	x, err := util.NormaliseURL(pl.URL)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	pl.NormalisedURL = x

	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}
	defer smartRollback(tx)

	row := tx.QueryRow(
		`INSERT INTO "page_loads"("id", "url", "normalised_url", "loaded_at", "not_load_before")
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ("normalised_url") DO UPDATE SET "url"             = $2,
                          "loaded_at"       = $4,
                          "not_load_before" = $5
RETURNING "page_loads"."id";`,
		pl.ID,
		pl.URL,
		pl.NormalisedURL,
		pl.LoadedAt,
		pl.NotLoadBefore,
	)

	var newID uuid.UUID
	if err := row.Scan(&newID); err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	return newID, nil
}
