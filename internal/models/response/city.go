package models

type CityResponse struct {
	Name       string
	Lat        float64
	Lon        float64
	Country    string
	LocalNames map[string]string
}
