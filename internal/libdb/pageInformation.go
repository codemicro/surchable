package db

import (
	"crypto/sha1"
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

func (db *DB) InsertPageInformation(pi *PageInformation) error {
	if pi.ID == uuid.Nil {
		pi.ID = uuid.New()
	}

	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.Exec(
		`INSERT INTO "page_information"("id", "load_id", "page_title", "page_meta_description_text", "page_content_text", "page_raw_html", "raw_html_sha1", "outbound_links")
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,

		pi.ID,
		pi.LoadID,
		pi.PageTitle,
		pi.PageMetaDescriptionText,
		pi.PageContentText,
		pi.PageRawHTML,
		pi.RawHTMLSHA1,
		pi.OutboundLinks,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}
