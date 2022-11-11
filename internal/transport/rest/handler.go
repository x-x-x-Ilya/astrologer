package rest

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/x-x-x-Ilya/astrologer/internal/services"
	"net/http"
	"net/url"
	"time"

	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest/transformers"
)

type PicturesController struct {
	picturesService services.PicturesServiceI
}

func NewPicturesController(picturesService services.PicturesServiceI) (*PicturesController, error) {
	if picturesService == nil {
		return nil, errors.New("nil")
	}

	return &PicturesController{
		picturesService,
	}, nil
}

func (c *PicturesController) pictures(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	var reqParams transformers.PicturesParams

	err = schema.NewDecoder().Decode(&reqParams, values)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	pictures, err := c.picturesService.Pictures(reqParams.Limit, reqParams.Offset)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	res := transformers.ToRests(pictures)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (c *PicturesController) pictureOfTheDay(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	var reqParams transformers.PictureParams

	reqParams.Date, err = time.Parse("2006-01-02", values.Get("date"))
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	pictures, err := c.picturesService.PictureOfTheDay(reqParams.Date)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	res := transformers.ToRest(pictures)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != "" {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Errorf("%+v", errors.WithStack(err))
		}

		_, err = w.Write(response)
		if err != nil {
			log.Errorf("%+v", errors.WithStack(err))
		}
	}
}
