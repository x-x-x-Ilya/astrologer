package services

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"github.com/x-x-x-Ilya/astrologer/internal/helpers"
	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type PicturesServiceI interface {
	Pictures(limit, offset int64) (models.Pictures, error)
	PictureOfTheDay(date time.Time) (models.Picture, error)
}

type PicturesService struct {
	storageService     StorageServiceI
	picturesRepository database.PicturesRepositoryI
	nasaClient         NasaClientI
	transactionService TransactionServiceI
}

func NewPicturesService(
	storageService StorageServiceI, picturesRepository database.PicturesRepositoryI,
	nasaClient NasaClientI, transactionService TransactionServiceI,
) (PicturesServiceI, error) {
	err := helpers.IsNotNil(nasaClient, picturesRepository, storageService, transactionService)
	if err != nil {
		return nil, errors.Wrapf(err, "err NewPicturesService")
	}

	return PicturesService{
		storageService,
		picturesRepository,
		nasaClient,
		transactionService,
	}, nil
}

func (p PicturesService) Pictures(limit, offset int64) (models.Pictures, error) {
	dbPictures, err := p.picturesRepository.Pictures(limit, offset)
	if err != nil {
		return models.Pictures{}, errors.Wrapf(err, "can't get pictures from db with params: limit = %d, offset = %d", limit, offset)
	}

	response := make(models.Pictures, 0, len(dbPictures))

	for _, dbPicture := range dbPictures {
		response = append(response, models.NewPicture(dbPicture.AsTime(), nil))
	}

	return response, nil
}

func (p PicturesService) PictureOfTheDay(date time.Time) (models.Picture, error) {
	dbPicture, err := p.picturesRepository.Picture(date)
	if err != nil {
		return models.Picture{}, errors.Wrapf(err, "can't get picture from db for date: %s", date.String())
	}

	if dbPicture != nil {
		file, err := p.storageService.Read(dbPicture.Date())
		if err != nil {
			return models.Picture{}, errors.Wrapf(err, "can't read picture file for date: %s", dbPicture.Date())
		}

		return models.NewPicture(dbPicture.AsTime(), file), nil
	}

	var newPicture models.Picture

	err = p.transactionService.Run(func(tx *sql.Tx) error {
		newPicture, err = p.nasaClient.Picture(date)
		if err != nil {
			return errors.Wrapf(err, "can't request picture for date: %s", date.String())
		}

		err = p.picturesRepository.Add(tx, newPicture)
		if err != nil {
			return errors.Wrapf(err, "can't add picture: %+v", newPicture)
		}

		err = p.storageService.Save(newPicture.Date(), newPicture.File())
		if err != nil {
			return errors.Wrapf(err, "can't safe file with date: %s", newPicture.Date())
		}

		return nil
	})
	if err != nil {
		return models.Picture{}, errors.WithStack(err)
	}

	return newPicture, nil
}
