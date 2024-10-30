//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package city

import (
	"context"
	"weather4you/internal/models"
)

type Repository interface {
	Create(ctx context.Context, city *models.CityDB) error
	GetCitiesList(ctx context.Context) ([]*models.CityLight, error)
	GetCitiesLightListWithPredictions(ctx context.Context) ([]*models.CityLight, error)
	GetCitiesListWithPredictions(ctx context.Context) ([]*models.CityDB, error)
	GetCityWithPrediction(ctx context.Context, name string, date string) (*models.CityWithPrediction, error)
	Save(city models.CityDB) error
	Exists(cityName string) (bool, error)
}
