package transformers

import (
	"time"

	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type PicturesParams struct {
	Limit  int64 `schema:"limit"`
	Offset int64 `schema:"offset"`
}

type PictureParams struct {
	Date time.Time `json:"date"`
}

type Picture struct {
	Date string `json:"date"`
}

func ToRest(domain models.Picture) Picture {
	return Picture{
		domain.Date().Format("2006-01-02"),
	}
}

func ToRests(domains models.Pictures) []Picture {
	rests := make([]Picture, 0, len(domains))
	for _, picture := range domains {
		rests = append(rests, ToRest(picture))
	}

	return rests
}
