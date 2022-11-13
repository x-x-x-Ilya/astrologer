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
	Add(tx *sql.Tx, picture models.Picture) error
	Picture(date time.Time) (*models.Picture, error)
	Pictures(limit int64, offset int64) (models.Pictures, error)
}

type PicturesRepository struct {
	dbr *dbr.Connection
}

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

func NewPicturesRepository(db *sql.DB) (PicturesRepositoryI, error) {
	if db == nil {
		return nil, errors.New("DB can't be nil")
	}

	return &PicturesRepository{
		GetDbc(db),
	}, nil
}

/*
func GetDbrSession(db *sql.DB, tx *sql.Tx)  *dbr.Tx {
	if reflect.ValueOf(tx).IsNil() {
		return GetDbrTransaction(GetDbc(db), tx)  // GetDbc(db).NewSession(nil)
	} else {
		return GetDbrTransaction(GetDbc(db), tx)
	}
}
*/

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

	_, err := rep.dbr.NewSession(nil).
		Select("apod_date").
		From("pictures").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		Load(&pictures)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pictures.toDomains(), nil
}

func (rep PicturesRepository) Picture(date time.Time) (*models.Picture, error) {
	var picture Picture

	err := rep.dbr.QueryRow("SELECT apod_date FROM pictures WHERE apod_date >= $1 and apod_date < $2", date, date.AddDate(0, 0, 1)).Scan(&picture.date)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		case err != nil:
			return nil, errors.WithStack(err)
		}
	}

	domainPicture := picture.toDomain()

	return &domainPicture, nil
}

func GetDbc(db *sql.DB) *dbr.Connection {
	dbc := dbr.Connection{
		DB:            db,
		Dialect:       dialect.PostgreSQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}

	return &dbc
}
