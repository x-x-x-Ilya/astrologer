package app

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
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

func ConnectionString(address string, port int64, userName string, password string, dbName string) string {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		address, port, userName, password, dbName)

	return connectionString
}

func (postgresConnector) OpenDBConnect(dbConf config.DBI) *sql.DB {
	connectionString := ConnectionString(dbConf.Address(), dbConf.Port(), dbConf.User(), dbConf.Password(), dbConf.Name())

	db, err := sql.Open("postgres", connectionString)
	for i := 0; i < 3 && err != nil; i++ {
		log.Errorf("can't open db connect: %+v (attempt[%d])", err, i+1)

		time.Sleep(time.Second * 30)

		db, err = sql.Open("postgres", connectionString)
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
