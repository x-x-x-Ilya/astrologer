package models

import "time"

type Pictures struct {
	Items  []Picture
	Amount int64
}

type Picture struct {
	date time.Time
	file []byte
}
