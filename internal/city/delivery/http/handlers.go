package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather4you/config"
	"weather4you/internal/city"
	"weather4you/internal/models"
	"weather4you/pkg/logger"

	"github.com/opentracing/opentracing-go"
)

// City handlers
type cityHandlers struct {
	cfg    *config.Config
	cityUC city.UseCase
	logger logger.Logger
}

// NewCityHandlers City handlers constructor
func NewCityHandlers(cfg *config.Config, cityUC city.UseCase, logger logger.Logger) city.Handlers {
	return &cityHandlers{cfg: cfg, cityUC: cityUC, logger: logger}
}

// Create godoc
// @Summary Get list
// @Description Get list of cities
// @Tags Cities
// @Accept json
// @Produce json
// @Success 201 {array} models.CityLight
// @Router /list/ [get]
func (h *cityHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("cityHandlers.GetList")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cities, err := h.cityUC.GetCitiesLightListWithPredictions(ctx)
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

// Create godoc
// @Summary Get predictions list
// @Description Get list of cities
// @Tags Cities
// @Accept json
// @Produce json
// @Success 201 {object} models.CityShort
// @Router /predictions/ [get]
func (h *cityHandlers) GetPredictionsList(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("cityHandlers.GetPredictionsList")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	citiesWithPredictions, err := h.cityUC.GetCitiesListWithPredictions(ctx)

	if err != nil {
		http.Error(w, "Error getting cities with predictions", http.StatusInternalServerError)
		return
	}

	var citiesShort []models.CityShort

	for _, city := range citiesWithPredictions {
		sumTemp := 0
		futurePredictions := 0
		for _, prediction := range city.Predictions {
			predictionDate := time.Unix(prediction.Date, 0)
			if predictionDate.After(time.Now()) {
				sumTemp += prediction.Temp
				futurePredictions++
			}
		}
		if futurePredictions > 0 {
			averageTemp := sumTemp / futurePredictions

			var futurePredictionDates []time.Time
			for _, prediction := range city.Predictions {
				predictionDate := time.Unix(prediction.Date, 0)
				if predictionDate.After(time.Now()) {
					futurePredictionDates = append(futurePredictionDates, predictionDate)
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

// Create godoc
// @Summary Get city predictions
// @Description Get full city info with prediction
// @Tags Cities
// @Accept json
// @Produce json
// @Success 201 {object} models.CityWithPrediction
// @Router /city/ [get]
func (h *cityHandlers) GetCityWithPrediction(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("cityHandlers.GetCityWithPrediction")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	name := r.URL.Query().Get("name")
	timestamp := r.URL.Query().Get("date")

	city, err := h.cityUC.GetCityWithPrediction(ctx, name, timestamp)
	if err != nil {
		h.logger.Error("Error getting city with prediction", err)
		http.Error(w, fmt.Sprintf("Error getting city with prediction - %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(city.Prediction.Info))
}
