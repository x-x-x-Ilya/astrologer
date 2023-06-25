package app

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

func migrationsUp(connectionString string) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id:   "2",
				Up:   []string{"CREATE TABLE pictures (id int, apod_date timestamp)"},
				Down: []string{"DROP TABLE pictures"},
			},
		},
	}

	db, err := sql.Open(dbDriver, connectionString)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = migrate.Exec(db, dbDriver, migrations, migrate.Up)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
