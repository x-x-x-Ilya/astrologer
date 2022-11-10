package services

import (
	"os"
	"strconv"
)

type StorageServiceI interface {
	Save(id int64, file []byte) error
	Read(id int64) ([]byte, error)
}

func Save(id int64) error {
	_, err := os.Create(strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}

	return nil
}

func Read(id int64) ([]byte, error) {
	dat, err := os.ReadFile(strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	return dat, nil
}
