package models

import "time"

type Pictures []Picture

type Picture struct {
	date time.Time
	file []byte
}

func NewPicture(date time.Time, file []byte) Picture {
	return Picture{
		date: date,
		file: file,
	}
}

func (p Picture) Date() time.Time {
	return p.date
}

func (p Picture) File() []byte {
	return p.file
}
