package weatherapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather4you/internal/config"
	"weather4you/internal/models"
)

type OpenWeatherMapResponse struct {
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
}

func GetPredictions(lat, lon float64) ([]models.Prediction, error) {

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s", lat, lon, config.WeatherApiToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response OpenWeatherMapResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	predictions := make([]models.Prediction, len(response.List))
	for i, item := range response.List {
		predictions[i] = models.Prediction{
			Temp: int(item.Main.Temp),
			Date: time.Unix(item.Dt, 0),
			Info: item.Weather[0].Description,
		}
	}
	// fmt.Println("predictions: ", predictions)

	return predictions, nil
}
