package app

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/internal/config"
)

type DatabaseConnector interface {
	OpenDBConnect(dbConf config.DBI) *sql.DB
}

type postgresConnector struct{}

func NewPostgresConnector() DatabaseConnector {
	return postgresConnector{}
}

func (postgresConnector) OpenDBConnect(dbConf config.DBI) *sql.DB {
	db, err := sql.Open("postgres", dbConf.ConnectionString())
	for i := 0; i < 3 && err != nil; i++ {
		log.Errorf("can't open db connect: %+v (attempt[%d])", err, i+1)

		time.Sleep(time.Second * 30)

		db, err = sql.Open("postgres", dbConf.ConnectionString())
	}

	if err != nil {
		log.Panicf("can't open db connect: %+v", err)
	}

	err = db.Ping()
	for i := 0; i < 3 && err != nil; i++ {
		log.Errorf("can't ping db connect: %+v (attempt[%d])", err, i+1)

		time.Sleep(time.Second * 30)

		err = db.Ping()
	}

	if err != nil {
		log.Panicf("can't ping db connect: %+v", err)
	}

	return db
}
