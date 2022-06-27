package db

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strconv"
)

type IndexClassification uint8

// There's a hard limit of 5 index classes maximum due to the database schema.
// Change the number of bits used for index classification if you need more than 5.
const (
	IndexClassPageBody IndexClassification = 1 << iota
	IndexClassPageTitle
	IndexClassPageDescription
)

const bitStringLength = 5

func (i IndexClassification) ToBitString() string {
	x := strconv.FormatInt(int64(i), 2)
	for len(x) < bitStringLength {
		x = "0" + x
	}
	return "B" + x + ""
}

func (i IndexClassification) Value() (driver.Value, error) {
	return i.ToBitString(), nil
}

func (i *IndexClassification) Scan(inp any) error {
	var x string
	switch y := inp.(type) {
	case string:
		x = y
	case []byte:
		x = string(y)
	case int64:
		*i = IndexClassification(y)
		return nil
	case float64:
		*i = IndexClassification(y)
		return nil
	default:
		return errors.WithStack(fmt.Errorf("cannot scan %T into IndexClassification: unknown type", inp))
	}
	newVal, err := strconv.ParseInt(x, 2, 8)
	if err != nil {
		return err
	}

	*i = IndexClassification(newVal)
	return nil
}

type TokenMap map[string]IndexClassification

func (tm TokenMap) Add(tokens []string, classification IndexClassification) {
	for _, token := range tokens {
		if existingClassification, found := tm[token]; found {
			tm[token] = existingClassification | classification
		} else {
			tm[token] = classification
		}
	}
}

type TokenSet struct {
	PageID uuid.UUID
	Tokens TokenMap
}

func (db *DB) SearchIndexQueryByPageID(pageID uuid.UUID) (*TokenSet, error) {
	ts := new(TokenSet)
	ts.PageID = pageID
	ts.Tokens = make(map[string]IndexClassification)

	ctx, cancel := db.newContext()
	defer cancel()

	rows, err := db.pool.QueryContext(
		ctx,
		`SELECT "token", "classification" from "search_index" WHERE "page_id" = $1`,
		ts.PageID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			token          string
			classification IndexClassification
		)
		if err := rows.Scan(&token, &classification); err != nil {
			return nil, errors.WithStack(err)
		}
		ts.Tokens[token] = classification
	}

	return ts, nil
}

func (db *DB) SearchIndexUpsert(ts *TokenSet) error {
	ctx, cancel := db.newContext()
	defer cancel()

	// Query this page ID, then diff that token set and the token set we have here.
	// Then apply the modifications required.

	existingTokenSet, err := db.SearchIndexQueryByPageID(ts.PageID)
	if err != nil {
		return errors.WithStack(err)
	}

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	for tokenString, newClassification := range ts.Tokens {

		existingClassification, alreadyExists := existingTokenSet.Tokens[tokenString]

		if !alreadyExists {
			_, err := tx.Exec(
				`INSERT INTO "search_index"("token", "page_id", "classification") VALUES ($1, $2, $3);`,
				tokenString,
				ts.PageID,
				newClassification,
			)
			if err != nil {
				return errors.WithStack(err)
			}
		} else if existingClassification != newClassification {
			_, err := tx.Exec(
				`UPDATE "search_index" SET "classification" = $1 WHERE "page_id" = $2 AND "token" = $3;`,
				newClassification,
				ts.PageID,
				tokenString,
			)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	for tokenString := range existingTokenSet.Tokens {
		if _, found := ts.Tokens[tokenString]; !found {
			_, err := tx.Exec(
				`DELETE FROM "search_index" WHERE "page_id" = $1 AND "token" = $2;`,
				ts.PageID,
				tokenString,
			)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return errors.WithStack(
		tx.Commit(),
	)
}
