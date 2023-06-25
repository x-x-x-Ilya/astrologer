package rest

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/internal/helpers"
	"github.com/x-x-x-Ilya/astrologer/internal/services"
	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest/transformers"
)

type PicturesController struct {
	picturesService services.PicturesServiceI
}

func NewPicturesController(picturesService services.PicturesServiceI) (*PicturesController, error) {
	err := helpers.IsNotNil(picturesService)
	if err != nil {
		return nil, errors.Wrapf(err, "err picturesService")
	}

	return &PicturesController{
		picturesService,
	}, nil
}

func (c *PicturesController) pictures(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var reqParams transformers.PicturesParams

	err = schema.NewDecoder().Decode(&reqParams, values)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	pictures, err := c.picturesService.Pictures(reqParams.Limit, reqParams.Offset)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := transformers.ToRests(pictures)

	respondWithJSON(w, http.StatusOK, res)
}

func (c *PicturesController) pictureOfTheDay(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var reqParams transformers.PictureParams

	reqParams.Date, err = time.Parse("2006-01-02", values.Get("date"))
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	picture, err := c.picturesService.PictureOfTheDay(reqParams.Date)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithIMG(w, http.StatusOK, picture.File())
}

func respondWithIMG(w http.ResponseWriter, code int, file []byte) {
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(code)

	_, err := w.Write(file)
	if err != nil {
		log.Errorf(err.Error())
	}
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
