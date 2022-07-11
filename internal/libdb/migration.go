package db

import (
	"database/sql"
	_ "embed"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var migrationFunctions = []func(trans *sql.Tx) error{
	migrate0to1,
}

func (db *DB) Migrate() error {
	log.Info().Msg("running migrations")

	// list tables
	tx, err := db.pool.Begin()
	if err != nil {
		return errors.WithMessage(err, "could not begin transaction")
	}
	defer smartRollback(tx)

	rows, err := db.pool.Query(`SELECT "table_name" FROM "information_schema"."tables" WHERE "table_schema" = 'public';`)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	existingTables := make(map[string]struct{})
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return errors.WithStack(err)
		}
		existingTables[tableName] = struct{}{}
	}

	var databaseVersion int

	if _, found := existingTables["version"]; found {
		err := db.pool.QueryRow(`SELECT "version" FROM "version";`).Scan(&databaseVersion)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.WithStack(err)
		}
	}

	if x := len(migrationFunctions); databaseVersion == x {
		log.Info().Msg("migrations up-to-date without any changes")
		return nil
	} else if databaseVersion > x {
		return errors.New("corrupt database: database version too high")
	}

	for _, f := range migrationFunctions[databaseVersion:] {
		if err := f(tx); err != nil {
			return errors.WithStack(err)
		}
	}

	log.Info().Msg("committing migrations")
	return errors.WithStack(
		tx.Commit(),
	)
}

//go:embed migrations/0to1.sql
var migrate0to1SQL string

func migrate0to1(trans *sql.Tx) error {
	log.Info().Msg("migrating new database to v1")

	_, err := trans.Exec(migrate0to1SQL)
	if err != nil {
		return errors.Wrap(err, "failed to migrate database version from v0 to v1")
	}

	return nil
}
