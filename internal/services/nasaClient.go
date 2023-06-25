package services

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/internal/helpers"
	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

const (
	nasaAddress = "https://api.nasa.gov"
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
	err := helpers.IsNotNil(client)
	if err != nil {
		return nil, errors.Wrapf(err, "err NewNasaClient")
	}

	return &NasaClient{
		client,
		apiKey,
		nasaAddress,
	}, nil
}

type pictureResponse struct {
	URL          string `json:"url"`
	ResponseType string `json:"media_type"`
}

func (n *NasaClient) Picture(date time.Time) (models.Picture, error) {
	year, month, day := date.Date()
	queryDate := fmt.Sprintf("%d-%d-%d", year, month, day)

	queryParams := map[string][]string{
		"date":    {queryDate},
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
		return models.Picture{}, errors.Wrapf(err, "can't get response from %s for date: %s", n.url, queryDate)
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
