package database

import (
	"database/sql"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type PicturesRepositoryI interface {
	Add(models.Picture) (int64, error)
	Picture(date time.Time) (*models.Picture, error)
	Pictures(params any) (models.Pictures, error)
}

type PicturesRepository struct {
	dbr *dbr.Connection
}

func NewPicturesRepository(db *sql.DB) (PicturesRepositoryI, error) {
	if db == nil {
		return nil, errors.New("DB can't be nil")
	}

	dbc := GetDbc(db)
	return &PicturesRepository{
		dbc,
	}, nil
}

func (n PicturesRepository) Add(picture models.Picture) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (n PicturesRepository) Pictures(params any) (models.Pictures, error) {
	//TODO implement me
	panic("implement me")
}

func (n PicturesRepository) Picture(date time.Time) (*models.Picture, error) {
	//TODO implement me
	panic("implement me")
}

func GetDbc(db *sql.DB) *dbr.Connection {
	d := dialect.MySQL

	er := &dbr.NullEventReceiver{}

	dbc := dbr.Connection{
		DB:            db,
		Dialect:       d,
		EventReceiver: er,
	}

	return &dbc
}
