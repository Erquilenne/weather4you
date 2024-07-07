package weatherapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"weather4you/internal/config"
	"weather4you/internal/models"
)

type City struct {
	Name       string
	Lat        float64
	Lon        float64
	Country    string
	LocalNames map[string]string
}

func FindCity(cityName string) (models.City, error) {

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, config.WeatherApiToken)

	resp, err := http.Get(url)
	if err != nil {
		return models.City{}, err
	}
	defer resp.Body.Close()

	var cities []City

	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return models.City{}, err
	}

	if len(cities) == 0 {
		return models.City{}, fmt.Errorf("City not found")
	}

	firstCity := cities[0]

	city := models.City{
		Name:    firstCity.Name,
		Lat:     firstCity.Lat,
		Lon:     firstCity.Lon,
		Country: firstCity.Country,
	}

	city.Predictions, err = GetPredictions(city.Lat, city.Lon)
	if err != nil {
		return models.City{}, err
	}

	return city, nil
}
