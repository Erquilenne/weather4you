package fillup

import (
	"weather4you/config"
	models "weather4you/internal/models/db"
	"weather4you/internal/weatherapi"
	"weather4you/pkg/logger"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=Saver
type Saver interface {
	SaveCity(city models.City) error
}

func FindAndSaveCity(cityName string, d Saver, cfg *config.Config, logger logger.Logger) {
	finder := weatherapi.NewCityFinder(cfg, logger)
	city := findCity(cityName, finder)
	city.Predictions = findPredictions(city.Lat, city.Lon, finder)

	if city.Predictions == nil {
		logger.Warnf("Predictions not found in city: %s", city.Name)
	}

	err := d.SaveCity(city)
	if err != nil {
		logger.Fatalf("SaveCity error: %s", err)
	}
	logger.Infof("City saved: %s", city.Name)
}

func findCity(cityName string, finder weatherapi.Finder) models.City {
	city := finder.FindCity(cityName)
	return city
}

func findPredictions(lat float64, lon float64, finder weatherapi.Finder) []models.Prediction {
	predictions := finder.FindPredictions(lat, lon)
	return predictions
}
