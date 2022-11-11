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
	pictures, err := p.picturesRepository.Pictures(limit, offset)
	if err != nil {
		return models.Pictures{}, err
	}

	return pictures, nil
}

func (p PicturesService) PictureOfTheDay(date time.Time) (models.Picture, error) {
	/*picture, err := p.picturesRepository.Picture(date)
	if err != nil {
		return models.Picture{}, err
	}
	*/
	var picture *models.Picture
	if picture == nil {
		newPicture, err := p.nasaClient.Picture(date)
		if err != nil {
			return models.Picture{}, err
		}

		_, err = p.picturesRepository.Add(newPicture)
		if err != nil {
			return models.Picture{}, err
		}

		err = p.storageService.Save(newPicture.Date().String(), newPicture.File())
		if err != nil {
			return models.Picture{}, err
		}

		return newPicture, nil
	}

	file, err := p.storageService.Read(picture.Date().String())
	if err != nil {
		return models.Picture{}, err
	}

	return models.NewPicture(date, file), nil
}
