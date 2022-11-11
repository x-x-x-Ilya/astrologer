package services

import (
	"os"
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
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (StorageService) Read(fileName string) ([]byte, error) {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return dat, nil
}
