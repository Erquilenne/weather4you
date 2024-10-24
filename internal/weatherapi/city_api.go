package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather4you/config"
	"weather4you/internal/models"
	response "weather4you/internal/models/response"
	"weather4you/pkg/logger"
)

type Finder interface {
	FindCity(cityName string) models.CityDB
	FindPredictions(lat float64, lon float64) []models.PredictionDB
}

type CityFinder struct {
	client *http.Client
	cfg    *config.Config
	logger logger.Logger
}

func NewCityFinder(cfg *config.Config, logger logger.Logger) *CityFinder {
	client := &http.Client{
		Timeout: 5 * time.Second, // Устанавливаем таймаут 5 секунд
	}
	return &CityFinder{cfg: cfg, logger: logger, client: client}
}

func (f *CityFinder) FindCity(cityName string) models.CityDB {

	// url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, f.cfg.WeatherApiToken)
	url := fmt.Sprintf("http://localhost:8082/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, f.cfg.WeatherApiToken)

	resp, err := f.client.Get(url)
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
		return models.CityDB{}
	}

	firstCity := cities[0]

	city := models.CityDB{
		Name:    firstCity.Name,
		Lat:     firstCity.Lat,
		Lon:     firstCity.Lon,
		Country: firstCity.Country,
	}

	return city
}

func (f *CityFinder) FindPredictions(lat float64, lon float64) []models.PredictionDB {

	// Find predictions with celsius units (units=metric)
	// url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, f.cfg.WeatherApiToken)
	url := fmt.Sprintf("http://localhost:8082/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, f.cfg.WeatherApiToken)

	resp, err := f.client.Get(url)
	if err != nil {
		f.logger.Fatalf("FindPredictions request error: %s", err)
	}
	defer resp.Body.Close()

	var response response.OpenWeatherMapResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if len(response.List) == 0 {
		f.logger.Warnf("FindPredictions: city predictions not found: %f, %f", lat, lon)
		return nil
	}
	if err != nil {
		f.logger.Fatalf("FindPredictions json decoder error: %s", err)
	}

	predictions := make([]models.PredictionDB, len(response.List))
	for i, item := range response.List {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			f.logger.Fatalf("FindPredictions json marshal error: %s", err)
		}
		if err != nil {
			f.logger.Fatalf("FindPredictions time parse error: %s", err)
		}
		predictions[i] = models.PredictionDB{
			Temp: int(item.Main.Temp),
			Date: item.Dt,
			Info: itemJSON,
		}
	}

	return predictions
}
