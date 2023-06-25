package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type ClientServiceI interface {
	Get(url string, queryParameters map[string][]string) (*http.Response, error)
	Post(url string, body any) (*http.Response, error)
}

type ClientService struct {
	http.Client
}

func NewClientService(timeout time.Duration) ClientServiceI {
	return &ClientService{
		http.Client{
			Timeout: timeout,
		},
	}
}

func (s ClientService) Get(url string, queryParameters map[string][]string) (*http.Response, error) {
	return s.doRequest(http.MethodGet, url, nil, queryParameters)
}

func (s ClientService) Post(url string, body any) (*http.Response, error) {
	return s.doRequest(http.MethodGet, url, body, nil)
}

func (s ClientService) doRequest(method, url string, body any, queryParameters map[string][]string) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query := req.URL.Query()

	for key, values := range queryParameters {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	req.URL.RawQuery = query.Encode()

	resp, err := s.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
