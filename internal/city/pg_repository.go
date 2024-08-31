//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package city

import (
	"context"
	"time"
	"weather4you/internal/models"
)

// type Repository interface {
// 	Create(ctx context.Context, city *models.City) (*models.News, error)
// 	Update(ctx context.Context, news *models.News) (*models.News, error)
// 	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error)
// 	Delete(ctx context.Context, newsID uuid.UUID) error
// 	GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error)
// 	SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error)
// }

type Repository interface {
	MakeMigrations()
	Create(ctx context.Context, city *models.CityDB) error
	GetCitiesList(ctx context.Context) ([]*models.CityLight, error)
	GetCitiesLightListWithPredictions(ctx context.Context) ([]*models.CityLight, error)
	GetCitiesListWithPredictions(ctx context.Context) ([]*models.CityDB, error)
	GetCityWithPrediction(ctx context.Context, name string, date time.Time) (*models.CityWithPrediction, error)
}
