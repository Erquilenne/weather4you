package models

import (
	"time"
	models "weather4you/internal/models/db"
)

type CityLight struct {
	Name string `json:"name"`
}

type CityShort struct {
	Name            string      `json:"name"`
	Country         string      `json:"country"`
	AverageTemp     int         `json:"average_temp"`
	PredictionDates []time.Time `json:"prediction_dates"`
}

type CityWithPrediction struct {
	Name       string            `json:"name"`
	Country    string            `json:"country"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Prediction models.Prediction `json:"prediction"`
}
