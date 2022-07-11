package db

import (
	"database/sql"

	"github.com/pkg/errors"
)

var ErrDomainAlreadyInBlocklist = errors.New("db: domain already in blocklist")

type BlocklistEntry struct {
	Domain string
	Reason string
}

func (db *DB) AddDomainToBlocklist(domain, reason string) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.Exec(`INSERT INTO "domain_blocklist"("domain", "reason") VALUES($1, $2)`, domain, reason)
	if err != nil {
		if isPostgresErrWithCode(err, errorCodeUniqueViolation) {
			return ErrDomainAlreadyInBlocklist
		}
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

func (db *DB) QueryDomainBlocklistByDomain(domain string) (*BlocklistEntry, error) {
	o := new(BlocklistEntry)

	if err := db.pool.QueryRow(`SELECT * FROM "domain_blocklist" WHERE "domain" = $1`, domain).Scan(&o.Domain, &o.Reason); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return o, nil
}

func (db *DB) RemoveDomainFromBlocklist(domain string) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.Exec(`DELETE FROM "domain_blocklist" WHERE "domain" = $1`, domain)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}