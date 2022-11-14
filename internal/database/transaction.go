package database

import (
	"database/sql"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	db *sql.DB
}

type TransactionFn = func(tx *sql.Tx) error

func NewTransaction(db *sql.DB) (*Transaction, error) {
	if db == nil {
		return nil, errors.New("DB can't be nil")
	}

	return &Transaction{db: db}, nil
}

func (t *Transaction) Run(fn TransactionFn) error {
	return RunTx(t.db, fn)
}

func (t *Transaction) RunManualCommit(fn TransactionFn) error {
	return RunTxManualCommit(t.db, fn)
}

func RunTx(db *sql.DB, transactionFn TransactionFn) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				log.Errorf("%+v", errors.Wrap(err, "can't rollback transaction"))
			}

			panic(p)
		}
	}()

	err = transactionFn(tx)
	if err != nil {
		return HandleTxError(tx, err)
	}

	err = tx.Commit()
	if err != nil {
		return HandleTxError(tx, err)
	}

	return nil
}

func RunTxManualCommit(db *sql.DB, transactionFn TransactionFn) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				log.Errorf("%+v", errors.WithStack(err))
			}

			panic(p)
		}
	}()

	err = transactionFn(tx)
	if err != nil {
		return HandleTxError(tx, err)
	}

	return err
}

func HandleTxError(tx *sql.Tx, err error) error {
	if txErr := tx.Rollback(); txErr != nil {
		return errors.Wrap(err, txErr.Error())
	}

	return err
}

func GetDbrTransaction(dbc *dbr.Connection, tx *sql.Tx) *dbr.Tx {
	sess := dbc.NewSession(nil)

	dbrTx := dbr.Tx{
		EventReceiver: sess.EventReceiver,
		Dialect:       sess.Dialect,
		Tx:            tx,
		Timeout:       sess.Timeout,
	}

	return &dbrTx
}
