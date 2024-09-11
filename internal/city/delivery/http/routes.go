package http

import (
	"github.com/gorilla/mux"

	"weather4you/internal/city"
	"weather4you/internal/middleware"
)

// Map news routes
func MapCityRoutes(newsGroup *mux.Router, h city.Handlers, mw *middleware.MiddlewareManager) {
	newsGroup.HandleFunc("/cities", h.GetList).Methods("GET")
	newsGroup.HandleFunc("/city", h.GetCityWithPrediction).Methods("GET")
	newsGroup.HandleFunc("/predictions", h.GetPredictionsList).Methods("GET")
}
