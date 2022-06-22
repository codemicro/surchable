package db

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PageInformation struct {
	ID                      uuid.UUID
	LoadID                  uuid.UUID
	PageTitle               *string
	PageMetaDescriptionText *string
	PageContentText         *string
	PageRawHTML             string
	RawHTMLSHA1             [sha1.Size]byte
	OutboundLinks           []string
}

func (db *DB) UpsertPageInformation(pi *PageInformation) (uuid.UUID, error) {
	if pi.ID == uuid.Nil {
		pi.ID = uuid.New()
	}

	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}
	defer smartRollback(tx)

	row := tx.QueryRow(
		`INSERT INTO "page_information"("id", "load_id", "page_title", "page_meta_description_text", "page_content_text",
                               "page_raw_html", "raw_html_sha1", "outbound_links")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT ("load_id") DO UPDATE SET "page_title"                 = $3,
                          "page_meta_description_text" = $4,
                          "page_content_text"          = $5,
                          "page_raw_html"              = $6,
                          "raw_html_sha1"              = $7,
                          "outbound_links"             = $8
RETURNING "page_information"."id";`,

		pi.ID,
		pi.LoadID,
		pi.PageTitle,
		pi.PageMetaDescriptionText,
		pi.PageContentText,
		pi.PageRawHTML,
		hex.EncodeToString(pi.RawHTMLSHA1[0:20]),
		stringSliceToPGArray(pi.OutboundLinks),
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
