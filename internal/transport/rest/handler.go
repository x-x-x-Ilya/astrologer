package rest

import (
	"encoding/json"
	"github.com/x-x-x-Ilya/astrologer/internal/services"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/x-x-x-Ilya/astrologer/internal/transport/rest/transformers"
)

type PicturesController struct {
	picturesService services.PicturesServiceI
}

func (c *PicturesController) pictures(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	var reqParams transformers.PictureParams

	err = schema.NewDecoder().Decode(&reqParams, values)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	domainParams := reqParams.PicturesParametersToDomain()
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	pictures, err := c.picturesService.Pictures(domainParams)
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

func (c *PicturesController) pictureOfTheDay(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, errors.WithStack(err))
		return
	}

	var reqParams transformers.PictureParams

	err = schema.NewDecoder().Decode(&reqParams, values)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	domainParams := reqParams.PicturesParametersToDomain()
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	pictures, err := c.picturesService.PictureOfTheDay(domainParams)
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
