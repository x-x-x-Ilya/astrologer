package transformers

import "time"

type PictureParams struct {
	Limit  int64 `schema:"limit"`
	Offset int64 `schema:"offset"`
}

type Picture struct {
	Date time.Time `json:"date"`
}

func ToRest(domain interface{}) Picture {
	return Picture{
		time.Now(),
	}
}

func (p *PictureParams) PicturesParametersToDomain() interface{} {
	return nil
}
