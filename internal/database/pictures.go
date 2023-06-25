package database

import (
	"database/sql"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	"github.com/x-x-x-Ilya/astrologer/internal/helpers"
	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type Picture struct {
	date time.Time `db:"apod_date"`
}

func (p Picture) toDomain() models.Picture {
	return models.NewPicture(p.date, nil)
}

type Pictures []Picture

func (p Pictures) toDomains() models.Pictures {
	domains := make(models.Pictures, 0, len(p))

	for _, picture := range p {
		domains = append(domains, models.NewPicture(picture.date, nil))
	}

	return domains
}

type PicturesRepositoryI interface {
	Add(tx *sql.Tx, picture models.Picture) error
	Picture(date time.Time) (*models.Picture, error)
	Pictures(limit int64, offset int64) (models.Pictures, error)
}

type PicturesRepository struct {
	dbr *dbr.Connection
}

func NewPicturesRepository(db *sql.DB) (PicturesRepositoryI, error) {
	err := helpers.IsNotNil(db)
	if err != nil {
		return nil, errors.Wrapf(err, "err NewPicturesRepository")
	}

	return &PicturesRepository{
		GetDbc(db),
	}, nil
}

func (rep PicturesRepository) Add(tx *sql.Tx, picture models.Picture) error {
	_, err := GetDbrTransaction(rep.dbr, tx).
		InsertInto("pictures").
		Columns("apod_date").
		Values(picture.Date()).
		Exec()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (rep PicturesRepository) Pictures(limit int64, offset int64) (models.Pictures, error) {
	var pictures Pictures

	rowsNum, err := rep.dbr.NewSession(nil).
		Select("apod_date").
		From("pictures").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		Load(&pictures)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if rowsNum == 0 {
		return nil, nil
	}

	return pictures.toDomains(), nil
}

func (rep PicturesRepository) Picture(date time.Time) (*models.Picture, error) {
	var picture Picture

	rowsNum, err := rep.dbr.NewSession(nil).
		Select("apod_date").
		From("pictures").
		Where(
			dbr.And(
				dbr.Gte("apod_date", date),
				dbr.Lt("apod_date", date.AddDate(0, 0, 1)),
			),
		).
		Load(&picture.date)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if rowsNum == 0 {
		return nil, nil
	}

	domainPicture := picture.toDomain()

	return &domainPicture, nil
}
