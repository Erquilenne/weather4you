package fillup

import (
	"weather4you/config"
	"weather4you/internal/models"
	"weather4you/internal/weatherapi"
	"weather4you/pkg/logger"
)

type Saver interface {
	Save(city models.CityDB) error
	Exists(cityName string) (bool, error)
}

func FindAndSaveCity(cityName string, d Saver, cfg *config.Config, logger logger.Logger) {
	exists, err := d.Exists(cityName)
	if err != nil {
		logger.Fatalf("Exists error: %s", err)
	}
	if exists {
		logger.Infof("City already exists: %s", cityName)
		return
	}
	finder := weatherapi.NewCityFinder(cfg, logger)
	city := findCity(cityName, finder)
	city.Predictions = FindPredictions(city.Lat, city.Lon, finder)

	if city.Predictions == nil {
		logger.Warnf("Predictions not found in city: %s", city.Name)
	}

	err = d.Save(city)
	if err != nil {
		logger.Fatalf("SaveCity error: %s", err)
	}
	logger.Infof("City saved: %s", city.Name)
}

func findCity(cityName string, finder weatherapi.Finder) models.CityDB {
	city := finder.FindCity(cityName)
	return city
}

func FindPredictions(lat float64, lon float64, finder weatherapi.Finder) []models.PredictionDB {
	predictions := finder.FindPredictions(lat, lon)
	return predictions
}
