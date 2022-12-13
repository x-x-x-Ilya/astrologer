package config

import (
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

	const (
		UserName = `DB_USER`
		Password = `DB_PASSWORD`
		Name     = `DB_NAME`
		Address  = `DB_ADDRESS`
		Port     = `DB_PORT`
	)

	db.user = os.Getenv(UserName)
	if db.user == "" {
		return getEnvErr(UserName)
	}

	db.password = os.Getenv(Password)
	if db.password == "" {
		return getEnvErr(Password)
	}

	db.name = os.Getenv(Name)
	if db.name == "" {
		return getEnvErr(Name)
	}

	db.address = os.Getenv(Address)
	if db.address == "" {
		return getEnvErr(Address)
	}

	db.port, err = strconv.ParseInt(os.Getenv(Port), 10, 64)
	if err != nil {
		return errors.Wrap(err, Port+" is invalid")
	}

	return nil
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
