package transformers

import (
	"github.com/x-x-x-Ilya/astrologer/internal/models"
	"time"
)

type PicturesParams struct {
	Limit  int64 `schema:"limit"`
	Offset int64 `schema:"offset"`
}

type PictureParams struct {
	Date time.Time `json:"date"`
}

type Picture struct {
	File []byte    `json:"file"`
	Date time.Time `json:"date"`
}

func ToRest(domain models.Picture) Picture {
	return Picture{
		domain.File(),
		domain.Date(),
	}
}

func ToRests(domains models.Pictures) []Picture {
	rests := make([]Picture, 0, len(domains))
	for _, picture := range domains {
		rests = append(rests, ToRest(picture))
	}

	return rests
}
