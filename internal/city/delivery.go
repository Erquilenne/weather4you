package city

import (
	"net/http"
)

// City HTTP Handlers interface
type Handlers interface {
	GetPredictionsList(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)
	GetCityWithPrediction(w http.ResponseWriter, r *http.Request)
}
