package services

import (
	"os"

	"github.com/pkg/errors"
)

type StorageServiceI interface {
	Save(fileName string, data []byte) error
	Read(fileName string) ([]byte, error)
}

type StorageService struct {
	storagePath string
}

func NewStorageService(storagePath string) StorageServiceI {
	return StorageService{
		storagePath,
	}
}

func (s StorageService) Save(fileName string, data []byte) error {
	f, err := os.Create(s.storagePath + fileName)
	if err != nil {
		return errors.Wrapf(err, "creating file with name %s failed", fileName)
	}

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrapf(err, "writing data into file with name %s failed", fileName)
	}

	return nil
}

func (s StorageService) Read(fileName string) ([]byte, error) {
	dat, err := os.ReadFile(s.storagePath + fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file with name %s failed", fileName)
	}

	return dat, nil
}
