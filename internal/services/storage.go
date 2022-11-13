package services

import (
	"os"

	"github.com/pkg/errors"
)

type StorageServiceI interface {
	Save(fileName string, data []byte) error
	Read(fileName string) ([]byte, error)
}

type StorageService struct{}

func NewStorageService() StorageServiceI {
	return StorageService{}
}

func (StorageService) Save(fileName string, data []byte) error {
	f, err := os.Create("./" + fileName)
	if err != nil {
		return errors.Wrapf(err, "creating file with name %s failed", fileName)
	}

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrapf(err, "writing data into file with name %s failed", fileName)
	}

	return nil
}

func (StorageService) Read(fileName string) ([]byte, error) {
	dat, err := os.ReadFile("./" + fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file with name %s failed", fileName)
	}

	return dat, nil
}
