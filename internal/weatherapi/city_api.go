package weatherapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"weather4you/internal/config"
	db "weather4you/internal/models/db"
	response "weather4you/internal/models/response"
)

func FindCity(cityName string) (db.City, error) {

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration in city:", err)
	}
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, config.WeatherApiToken)

	resp, err := http.Get(url)
	if err != nil {
		return db.City{}, err
	}
	defer resp.Body.Close()

	var cities []response.CityResponse

	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return db.City{}, err
	}

	if len(cities) == 0 {
		return db.City{}, fmt.Errorf("City not found")
	}

	firstCity := cities[0]

	city := db.City{
		Name:    firstCity.Name,
		Lat:     firstCity.Lat,
		Lon:     firstCity.Lon,
		Country: firstCity.Country,
	}

	city.Predictions, err = GetPredictions(city.Lat, city.Lon)
	if err != nil {
		return db.City{}, err
	}

	return city, nil
}
