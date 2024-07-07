package models

import "time"

type Prediction struct {
	Temp int       `json:"temp"`
	Date time.Time `json:"date"`
	Info string    `json:"info"`
}
