package weatherapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather4you/internal/config"
	dbModels "weather4you/internal/models/db"
	responseModels "weather4you/internal/models/response"
)

func GetPredictions(lat, lon float64) ([]dbModels.Prediction, error) {

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration in predictions:", err)
	}
	// Find predictions with celsius units (units=metric)
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, config.WeatherApiToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response responseModels.OpenWeatherMapResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	predictions := make([]dbModels.Prediction, len(response.List))
	for i, item := range response.List {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		predictions[i] = dbModels.Prediction{
			Temp: int(item.Main.Temp),
			Date: time.Unix(item.Dt, 0),
			Info: string(itemJSON),
		}
	}

	return predictions, nil
}
