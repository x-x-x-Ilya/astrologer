package services

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type NasaClientI interface {
	Picture(date time.Time) (models.Picture, error)
}

type NasaClient struct {
	client ClientServiceI
	apiKey string
	url    string
}

func NewNasaClient(apiKey string, client ClientServiceI) (NasaClientI, error) {
	if client == nil {
		return nil, nilErr("client")
	}

	return &NasaClient{
		client,
		apiKey,
		"https://api.nasa.gov",
	}, nil
}

type pictureResponse struct {
	URL          string `json:"url"`
	ResponseType string `json:"media_type"`
}

func (n *NasaClient) Picture(date time.Time) (models.Picture, error) {
	year, month, day := date.Date()
	queryParams := map[string][]string{
		"date":    {fmt.Sprintf("%d-%d-%d", year, month, day)},
		"api_key": {n.apiKey},
	}

	response, err := n.client.Get(n.url+"/planetary/apod/", queryParams)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Panicf("can't close db rows %+v", err)
		}
	}()

	if err != nil {
		return models.Picture{}, errors.Wrapf(err, "can't get response from %s for date: %s", n.url, fmt.Sprintf("%d-%d-%d", year, month, day))
	}

	var responseStruct pictureResponse

	err = json.NewDecoder(response.Body).Decode(&responseStruct)
	if err != nil {
		return models.Picture{}, errors.Wrapf(err, "can't decode body to struct: %T %+v", responseStruct, responseStruct)
	}

	if responseStruct.ResponseType != "image" {
		return models.Picture{}, errors.Errorf("nasa response type is not image, response type is: %s", responseStruct.ResponseType)
	}

	imgResponse, err := n.client.Get(responseStruct.URL, nil)
	defer func() {
		err := imgResponse.Body.Close()
		if err != nil {
			log.Panicf("can't close db rows %+v", err)
		}
	}()

	if err != nil {
		return models.Picture{}, errors.Wrapf(err, "can't get: %s", responseStruct.URL)
	}

	buffer := make([]byte, imgResponse.ContentLength)

	_, err = io.ReadFull(imgResponse.Body, buffer)
	if err != nil {
		return models.Picture{}, errors.Wrapf(err, "can't read response body from %s", responseStruct.URL)
	}

	return models.NewPicture(date, buffer), nil
}
