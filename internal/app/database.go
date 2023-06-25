package app

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/config"
)

const (
	dbDriver        = "postgres"
	retriesInterval = time.Second * 30
	retriesAmount   = 3
)

type DatabaseConnector interface {
	OpenDBConnect(dbConf config.DBI) *sql.DB
}

type postgresConnector struct{}

func NewPostgresConnector() DatabaseConnector {
	return postgresConnector{}
}

func (postgresConnector) OpenDBConnect(dbConf config.DBI) *sql.DB {
	connStr := dbConf.ConnectionString()

	db, err := sql.Open(dbDriver, connStr)
	for i := 0; i < retriesAmount && err != nil; i++ {
		log.Errorf("can't open db connect: %+v (attempt[%d])", err, i+1)

		time.Sleep(retriesInterval)

		db, err = sql.Open(dbDriver, connStr)
	}

	if err != nil {
		log.Panicf("can't open db connect: %+v", err)
	}

	err = db.Ping()
	for i := 0; i < retriesAmount && err != nil; i++ {
		log.Errorf("can't ping db connect: %+v (attempt[%d])", err, i+1)

		time.Sleep(retriesInterval)

		err = db.Ping()
	}

	if err != nil {
		log.Panicf("can't ping db connect: %+v", err)
	}

	return db
}
