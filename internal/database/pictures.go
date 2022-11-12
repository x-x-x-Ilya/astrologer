package database

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/jmoiron/sqlx"
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

	dbc := GetDbc(db)
	return &PicturesRepository{
		dbc,
	}, nil
}

type DbrSession interface {
	Select(column ...string) *dbr.SelectStmt
	SelectBySql(query string, value ...interface{}) *dbr.SelectStmt
	Update(table string) *dbr.UpdateStmt
	UpdateBySql(query string, value ...interface{}) *dbr.UpdateStmt
	InsertInto(table string) *dbr.InsertStmt
	InsertBySql(query string, value ...interface{}) *dbr.InsertStmt
	DeleteFrom(table string) *dbr.DeleteStmt
	DeleteBySql(query string, value ...interface{}) *dbr.DeleteStmt
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	sqlx.Execer
	sqlx.ExecerContext
}

func GetDbrSession(db *sql.DB, tx *sql.Tx) DbrSession {
	if reflect.ValueOf(tx).IsNil() {
		return GetDbc(db).NewSession(nil)
	} else {
		return GetDbrTransaction(GetDbc(db), tx)
	}
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

func (rep PicturesRepository) Add(tx *sql.Tx, picture models.Picture) error {
	_, err := GetDbrSession(rep.dbr.DB, tx).
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

	var date time.Time
	/*rows, err := rep.dbr.NewSession(nil).
	Select("apod_date").
	From("pictures").
	Limit(uint64(limit)).
	Offset(uint64(offset)).
	Load(&pictures)
	*/
	rows, err := rep.dbr.Query("SELECT apod_date FROM pictures LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for rows.Next() {
		err := rows.Scan(&date)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		pictures = append(pictures, Picture{date: date})
	}

	err = rows.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return pictures.toDomains(), nil
	}

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
	d := dialect.PostgreSQL

	er := &dbr.NullEventReceiver{}

	dbc := dbr.Connection{
		DB:            db,
		Dialect:       d,
		EventReceiver: er,
	}

	return &dbc
}
