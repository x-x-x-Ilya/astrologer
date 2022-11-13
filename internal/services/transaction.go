package services

import (
	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/internal/database"
)

type TransactionServiceI interface {
	Run(fn database.TransactionFn) error
	RunManualCommit(fn database.TransactionFn) error
}

type TransactionService struct {
	transaction TransactionServiceI
}

func NewTransactionService(transaction TransactionServiceI) (TransactionServiceI, error) {
	if transaction == nil {
		return nil, errors.New("transaction can't be nil")
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
