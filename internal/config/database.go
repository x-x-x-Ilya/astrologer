package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type DBI interface {
	User() string
	Password() string
	Name() string
	Address() string
	Port() int64
	ConnectionString() string
}

type DB struct {
	user     string
	password string
	name     string
	address  string
	port     int64
}

func newDBConfig() (DBI, error) {
	var db DB

	err := db.getEnv()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) getEnv() error {
	var err error

	db.user = os.Getenv(`DB_USER`)
	if db.user == "" {
		return getEnvErr(`DB_USER`)
	}

	db.password = os.Getenv(`DB_PASSWORD`)
	if db.password == "" {
		return getEnvErr(`DB_PASSWORD`)
	}

	db.name = os.Getenv(`DB_NAME`)
	if db.name == "" {
		return getEnvErr(`DB_NAME`)
	}

	db.address = os.Getenv(`DB_HOST`)
	if db.address == "" {
		return getEnvErr(`DB_HOST`)
	}

	db.port, err = strconv.ParseInt(os.Getenv(`DB_PORT`), 10, 64)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (db DB) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Address(), db.Port(), db.User(), db.Password(), db.Name())
}

func (db DB) User() string {
	return db.user
}

func (db DB) Password() string {
	return db.password
}

func (db DB) Name() string {
	return db.name
}

func (db DB) Address() string {
	return db.address
}

func (db DB) Port() int64 {
	return db.port
}
