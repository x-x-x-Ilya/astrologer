package services

import (
	"github.com/pkg/errors"
	"github.com/x-x-x-Ilya/astrologer/internal/database"
	"time"

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
}

func nilErr(entityName string) error {
	return errors.Errorf("error, %s is nil", entityName)
}

func NewPicturesService(storageService StorageServiceI, picturesRepository database.PicturesRepositoryI, nasaClient NasaClientI) (PicturesServiceI, error) {
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
	}, nil
}

func (p PicturesService) Pictures(limit int64, offset int64) (models.Pictures, error) {
	dbPictures, err := p.picturesRepository.Pictures(limit, offset)
	if err != nil {
		return models.Pictures{}, err
	}

	var response = make(models.Pictures, 0, len(dbPictures))
	if dbPictures != nil {
		for _, dbPicture := range dbPictures {
			response = append(response, models.NewPicture(dbPicture.Date(), nil))
		}
	}

	return response, nil
}

func (p PicturesService) PictureOfTheDay(date time.Time) (models.Picture, error) {
	dbPicture, err := p.picturesRepository.Picture(date)
	if err != nil {
		return models.Picture{}, err
	}

	if dbPicture != nil {
		file, err := p.storageService.Read(dbPicture.Date().Format("2006-01-02 15:04:05"))
		if err != nil {
			return models.Picture{}, err
		}

		return models.NewPicture(dbPicture.Date(), file), nil
	}

	newPicture, err := p.nasaClient.Picture(date)
	if err != nil {
		return models.Picture{}, err
	}

	err = p.picturesRepository.Add(nil, newPicture)
	if err != nil {
		return models.Picture{}, err
	}

	err = p.storageService.Save(newPicture.Date().Format("2006-01-02 15:04:05"), newPicture.File())
	if err != nil {
		return models.Picture{}, err
	}

	return newPicture, nil
}
