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
	ContentSHA1   *string
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
		`SELECT "id", "url", "normalised_url", "content_sha1", "loaded_at", "not_load_before" FROM "page_loads" WHERE "normalised_url" = $1;`,
		normalisedURL,
	).Scan(&pageLoad.ID, &pageLoad.URL, &pageLoad.NormalisedURL, &pageLoad.ContentSHA1, &pageLoad.LoadedAt, &pageLoad.NotLoadBefore)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoMatchingPageLoad
		}
		return nil, errors.WithStack(err)
	}

	return pageLoad, nil
}
