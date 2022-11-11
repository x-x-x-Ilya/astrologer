package app

import (
	"database/sql"
	log "github.com/sirupsen/logrus"

	migrate "github.com/rubenv/sql-migrate"
)

func migrationsUp(address string, port int64, userName string, password string, dbName string) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id:   "2",
				Up:   []string{"CREATE TABLE pictures (id int, date date)"},
				Down: []string{"DROP TABLE pictures"},
			},
		},
	}

	connectionString := ConnectionString(address, port, userName, password, dbName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Infof("Applied %d migrations!\n", n)

	return nil
}
