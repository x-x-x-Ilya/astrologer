package services

import (
	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"github.com/x-x-x-Ilya/astrologer/internal/helpers"
)

type TransactionServiceI interface {
	Run(fn database.TransactionFn) error
	RunManualCommit(fn database.TransactionFn) error
}

type TransactionService struct {
	transaction TransactionServiceI
}

func NewTransactionService(transaction TransactionServiceI) (TransactionServiceI, error) {
	err := helpers.IsNotNil(transaction)
	if err != nil {
		return nil, errors.Wrapf(err, "err NewTransactionService")
	}

	return &TransactionService{
		transaction,
	}, nil
}

func (t *TransactionService) Run(fn database.TransactionFn) error {
	err := t.transaction.Run(fn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *TransactionService) RunManualCommit(fn database.TransactionFn) error {
	err := t.transaction.RunManualCommit(fn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
