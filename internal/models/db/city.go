package models

type City struct {
	Name        string       `json:"name"`
	Country     string       `json:"country"`
	Lat         float64      `json:"lat"`
	Lon         float64      `json:"lon"`
	Predictions []Prediction `json:"predictions"`
}
