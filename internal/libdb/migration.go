package db

import (
	"database/sql"

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
	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Warn().Stack().Err(errors.WithStack(err)).Msg("failed to rollback transaction")
		}
	}()

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

	if _, found := existingTables[tableNameVersion]; found {
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

func migrate0to1(trans *sql.Tx) error {
	log.Info().Msg("migrating new database to v1")

	_, err := trans.Exec(`CREATE TABLE "domain_queue"
	(
		id         uuid      not null primary key,
		created_at timestamp not null default now(),
		domain     varchar(253),
		subdomain  varchar(253)
	);`)
	if err != nil {
		return errors.Wrap(err, "failed to create `domain_queue` table")
	}

	_, err = trans.Exec(`CREATE TABLE "page_loads"
	(
		id              uuid      not null primary key,
		url             text      not null,
		content         text,
		loaded_at       timestamp not null default now(),
		not_load_before timestamp
	);`)
	if err != nil {
		return errors.Wrap(err, "failed to create `page_loads` table")
	}

	_, err = trans.Exec(`CREATE TABLE "current_jobs"
	(
		id            uuid      not null primary key,
		queue_item    uuid      not null references domain_queue (id),
		worker_id     uuid      not null,
		last_check_in timestamp not null
	);`)
	if err != nil {
		return errors.Wrap(err, "failed to create `current_jobs` table")
	}

	_, err = trans.Exec(`CREATE TABLE "version"
	(
		version int
	);
	INSERT INTO "version"(version)
	VALUES (1);`)
	if err != nil {
		return errors.Wrap(err, "failed to create `version` table and insert version number")
	}

	return nil
}
