package app

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

func migrationsUp(address string, port int64, userName string, password string, dbName string) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id:   "2",
				Up:   []string{"CREATE TABLE pictures (id int, apod_date timestamp)"},
				Down: []string{"DROP TABLE pictures"},
			},
		},
	}

	connectionString := ConnectionString(address, port, userName, password, dbName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
