package models

import (
	"encoding/json"
	"time"
)

// City base model
type CityDB struct {
	Name        string         `json:"name"`
	Country     string         `json:"country"`
	Lat         float64        `json:"lat"`
	Lon         float64        `json:"lon"`
	Predictions []PredictionDB `json:"predictions"`
}

// Predictions base model
type PredictionDB struct {
	Temp int             `json:"temp"`
	Date int64           `json:"date"`
	Info json.RawMessage `json:"info"`
}

// Model for cities list response
type CityLight struct {
	Name string `json:"name"`
}

// Model for city prediction list response
type CityShort struct {
	Name            string      `json:"name"`
	Country         string      `json:"country"`
	AverageTemp     int         `json:"average_temp"`
	PredictionDates []time.Time `json:"prediction_dates"`
}

// Model for full city response
type CityWithPrediction struct {
	Name       string       `json:"name"`
	Country    string       `json:"country"`
	Lat        float64      `json:"lat"`
	Lon        float64      `json:"lon"`
	Prediction PredictionDB `json:"prediction"`
}
