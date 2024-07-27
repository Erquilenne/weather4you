package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather4you/config"
	db "weather4you/internal/models/db"
	response "weather4you/internal/models/response"
	"weather4you/pkg/logger"
)

type Finder interface {
	FindCity(cityName string) db.City
	FindPredictions(lat float64, lon float64) []db.Prediction
}

type CityFinder struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewCityFinder(cfg *config.Config, logger logger.Logger) *CityFinder {
	return &CityFinder{cfg: cfg, logger: logger}
}

func (f *CityFinder) FindCity(cityName string) db.City {

	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, f.cfg.WeatherApiToken)

	resp, err := http.Get(url)
	if err != nil {
		f.logger.Fatalf("FindCity request error: %s", err)
	}
	defer resp.Body.Close()

	var cities []response.CityResponse

	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		f.logger.Fatalf("FindCity json decoder error: %s", err)
	}

	if len(cities) == 0 {
		f.logger.Warnf("FindCity: city not found: %s", cityName)
		return db.City{}
	}

	firstCity := cities[0]

	city := db.City{
		Name:    firstCity.Name,
		Lat:     firstCity.Lat,
		Lon:     firstCity.Lon,
		Country: firstCity.Country,
	}

	return city
}

func (f *CityFinder) FindPredictions(lat float64, lon float64) []db.Prediction {

	// Find predictions with celsius units (units=metric)
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, f.cfg.WeatherApiToken)

	resp, err := http.Get(url)
	if err != nil {
		f.logger.Fatalf("FindPredictions request error: %s", err)
	}
	defer resp.Body.Close()

	var response response.OpenWeatherMapResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		f.logger.Fatalf("FindPredictions json decoder error: %s", err)
	}

	predictions := make([]db.Prediction, len(response.List))
	for i, item := range response.List {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			f.logger.Fatalf("FindPredictions json marshal error: %s", err)
		}
		predictions[i] = db.Prediction{
			Temp: int(item.Main.Temp),
			Date: time.Unix(item.Dt, 0),
			Info: string(itemJSON),
		}
	}

	return predictions
}
