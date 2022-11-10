package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type ClientServiceI interface {
	DoRequest(method string, url string, body interface{}, queryParameters map[string][]string) (*http.Response, error)
}

type ClientService struct {
	client http.Client
}

func NewClientService(timeout time.Duration) ClientServiceI {
	return &ClientService{
		http.Client{
			Timeout: timeout,
		},
	}
}

func (ClientService) DoRequest(method string, url string, body interface{}, queryParameters map[string][]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, values := range queryParameters {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Error(err)
	}
}
