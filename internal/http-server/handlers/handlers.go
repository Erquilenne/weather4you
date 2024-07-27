package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	db "weather4you/internal/models/db"
	models "weather4you/internal/models/request"
	"weather4you/internal/storage/pgsql"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=DatabaseGetter
type DatabaseGetter interface {
	GetCitiesListWithPredictions() ([]db.City, error)
	GetCitiesLightListWithPredictions() ([]models.CityLight, error)
	GetCityWithPrediction(name string, date time.Time) (*models.CityWithPrediction, error)
}

type Handler struct {
	db DatabaseGetter
}

func NewHandler(db *pgsql.Database) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	cities, err := h.db.GetCitiesLightListWithPredictions()
	if err != nil {
		http.Error(w, "Error getting cities with predictions", http.StatusInternalServerError)
		return
	}

	citiesJSON, err := json.Marshal(cities)
	if err != nil {
		http.Error(w, "Error marshaling cities to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(citiesJSON)
}

func (h *Handler) GetPredictionsList(w http.ResponseWriter, r *http.Request) {
	citiesWithPredictions, err := h.db.GetCitiesListWithPredictions()

	if err != nil {
		http.Error(w, "Error getting cities with predictions", http.StatusInternalServerError)
		return
	}

	var citiesShort []models.CityShort

	for _, city := range citiesWithPredictions {
		sumTemp := 0
		futurePredictions := 0
		for _, prediction := range city.Predictions {
			if prediction.Date.After(time.Now()) {
				sumTemp += prediction.Temp
				futurePredictions++
			}
		}
		if futurePredictions > 0 {
			averageTemp := sumTemp / futurePredictions

			var futurePredictionDates []time.Time
			for _, prediction := range city.Predictions {
				if prediction.Date.After(time.Now()) {
					futurePredictionDates = append(futurePredictionDates, prediction.Date)
				}
			}

			cityShort := models.CityShort{
				Name:            city.Name,
				Country:         city.Country,
				AverageTemp:     averageTemp,
				PredictionDates: futurePredictionDates,
			}
			citiesShort = append(citiesShort, cityShort)
		}
	}

	citiesShortJSON, err := json.Marshal(citiesShort)
	if err != nil {
		http.Error(w, "Error marshaling cities to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(citiesShortJSON)
}

func (h *Handler) GetCityWithPrediction(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	dateStr := r.URL.Query().Get("date")
	fmt.Println(name, dateStr)

	date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	city, err := h.db.GetCityWithPrediction(name, date)
	if err != nil {
		http.Error(w, "Error getting city with prediction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(city.Prediction.Info))
}
