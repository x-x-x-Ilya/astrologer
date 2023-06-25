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

type db struct {
	user     string
	password string
	name     string
	address  string
	port     int64
}

func newDBConfig() (DBI, error) {
	var db db

	err := db.getEnv()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *db) getEnv() error {
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

func (db db) ConnectionString() string {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Address(), db.Port(), db.User(), db.Password(), db.Name())

	return connectionString
}

func (db db) User() string {
	return db.user
}

func (db db) Password() string {
	return db.password
}

func (db db) Name() string {
	return db.name
}

func (db db) Address() string {
	return db.address
}

func (db db) Port() int64 {
	return db.port
}
