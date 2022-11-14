package database

import (
	"database/sql"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
)

func GetDbc(db *sql.DB) *dbr.Connection {
	dbc := dbr.Connection{
		DB:            db,
		Dialect:       dialect.PostgreSQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}

	return &dbc
}
