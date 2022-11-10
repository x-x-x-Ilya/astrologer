package services

import (
	"encoding/json"
	"net/http"
	"time"

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

func NewNasaClient(client ClientServiceI) (NasaClientI, error) {
	if client == nil {
		return nil, nilErr("client")
	}

	nasaClient := &NasaClient{
		client,
		"",
		"api.nasa.gov",
	}

	err := nasaClient.refreshKey()
	if err != nil {
		return nil, err
	}

	return nasaClient, nil
}

func (n *NasaClient) Picture(date time.Time) (models.Picture, error) {
	queryParams := map[string][]string{"date": {date.String()}}

	response, err := n.client.DoRequest(http.MethodGet, n.url+"/", nil, queryParams)
	if err != nil {
		return models.Picture{}, err
	}

	defer closeBody(response.Body)

	var responseStruct any
	err = json.NewDecoder(response.Body).Decode(responseStruct)
	if err != nil {
		return models.Picture{}, err
	}

	return models.Picture{}, nil
}

func (n *NasaClient) refreshKey() error {
	return nil
}
