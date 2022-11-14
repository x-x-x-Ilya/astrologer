package models

import "time"

type Pictures []Picture

type Picture struct {
	date PictureDateI
	file []byte
}

type PictureDateI interface {
	String() string
	AsTime() time.Time
}

type PictureDate struct {
	time.Time
}

func (p PictureDate) String() string {
	return p.Format("2006-01-02")
}

func (p PictureDate) AsTime() time.Time {
	return p.Time
}

func NewPicture(date time.Time, file []byte) Picture {
	return Picture{
		date: PictureDate{date},
		file: file,
	}
}

func (p Picture) AsTime() time.Time {
	return p.date.AsTime()
}

func (p Picture) Date() string {
	return p.date.String()
}

func (p Picture) File() []byte {
	return p.file
}
