package services

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type PicturesServiceI interface {
	Pictures(limit int64, offset int64) (models.Pictures, error)
	PictureOfTheDay(date time.Time) (models.Picture, error)
}

type PicturesService struct {
	storageService     StorageServiceI
	picturesRepository database.PicturesRepositoryI
	nasaClient         NasaClientI
	transactionService TransactionServiceI
}

func nilErr(entityName string) error {
	return errors.Errorf("error, %s is nil", entityName)
}

func NewPicturesService(
	storageService StorageServiceI, picturesRepository database.PicturesRepositoryI, nasaClient NasaClientI, transactionService TransactionServiceI,
) (PicturesServiceI, error) {
	if nasaClient == nil {
		return nil, nilErr("nasaClient")
	}

	if picturesRepository == nil {
		return nil, nilErr("picturesRepository")
	}

	return PicturesService{
		storageService,
		picturesRepository,
		nasaClient,
		transactionService,
	}, nil
}

func (p PicturesService) Pictures(limit int64, offset int64) (models.Pictures, error) {
	dbPictures, err := p.picturesRepository.Pictures(limit, offset)
	if err != nil {
		return models.Pictures{}, errors.Wrapf(err, "can'transaction get pictures from db with params: limit = %d, offset = %d", limit, offset)
	}

	response := make(models.Pictures, 0, len(dbPictures))

	for _, dbPicture := range dbPictures {
		response = append(response, models.NewPicture(dbPicture.Date(), nil))
	}

	return response, nil
}

func (p PicturesService) PictureOfTheDay(date time.Time) (models.Picture, error) {
	dbPicture, err := p.picturesRepository.Picture(date)
	if err != nil {
		return models.Picture{}, errors.Wrapf(err, "can'transaction get picture from db for date: %s", date.String())
	}

	if dbPicture != nil {
		file, err := p.storageService.Read(dbPicture.Date().Format("2006-01-02 15:04:05"))
		if err != nil {
			return models.Picture{}, errors.Wrapf(err, "can'transaction read picture file for date: %s", dbPicture.Date().Format("2006-01-02 15:04:05"))
		}

		return models.NewPicture(dbPicture.Date(), file), nil
	}

	var newPicture models.Picture

	err = p.transactionService.Run(func(tx *sql.Tx) error {
		newPicture, err = p.nasaClient.Picture(date)
		if err != nil {
			return errors.Wrapf(err, "can'transaction request picture for date: %s", date.String())
		}

		err = p.picturesRepository.Add(tx, newPicture)
		if err != nil {
			return errors.Wrapf(err, "can'transaction add picture: %+v", newPicture)
		}

		err = p.storageService.Save(newPicture.Date().Format("2006-01-02 15:04:05"), newPicture.File())
		if err != nil {
			return errors.Wrapf(err, "can'transaction safe file with date: %s", newPicture.Date().Format("2006-01-02 15:04:05"))
		}

		return nil
	})

	return newPicture, nil
}
