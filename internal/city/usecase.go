package city

import (
	"context"
	"weather4you/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, city *models.CityDB) error
	GetCitiesList(ctx context.Context) ([]*models.CityLight, error)
	GetCitiesLightListWithPredictions(ctx context.Context) ([]*models.CityLight, error)
	GetCitiesListWithPredictions(ctx context.Context) ([]*models.CityDB, error)
	GetCityWithPrediction(ctx context.Context, name string, timestamp string) (*models.CityWithPrediction, error)
}
