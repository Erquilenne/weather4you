package usecase

import (
	"context"
	"time"
	"weather4you/config"
	city "weather4you/internal/city"
	"weather4you/internal/models"
	"weather4you/pkg/logger"

	"github.com/opentracing/opentracing-go"
)

// City UseCase
type cityUC struct {
	cfg      *config.Config
	cityRepo city.Repository
	logger   logger.Logger
}

// City UseCase constructor
func NewCityUseCase(cfg *config.Config, cityRepo city.Repository, logger logger.Logger) city.UseCase {
	return &cityUC{cfg: cfg, cityRepo: cityRepo, logger: logger}
}

func (u *cityUC) Create(ctx context.Context, city *models.CityDB) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityUC.Create")
	defer span.Finish()
	return u.cityRepo.Create(ctx, city)
}

func (u *cityUC) GetCitiesList(ctx context.Context) ([]*models.CityLight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityUC.GetCitiesList")
	defer span.Finish()
	return u.cityRepo.GetCitiesList(ctx)
}

func (u *cityUC) GetCitiesLightListWithPredictions(ctx context.Context) ([]*models.CityLight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityUC.GetCitiesLightListWithPredictions")
	defer span.Finish()
	return u.cityRepo.GetCitiesLightListWithPredictions(ctx)
}

func (u *cityUC) GetCitiesListWithPredictions(ctx context.Context) ([]*models.CityDB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityUC.GetCitiesListWithPredictions")
	defer span.Finish()
	return u.cityRepo.GetCitiesListWithPredictions(ctx)
}

func (u *cityUC) GetCityWithPrediction(ctx context.Context, name string, date time.Time) (*models.CityWithPrediction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityUC.GetCityWithPrediction")
	defer span.Finish()
	return u.cityRepo.GetCityWithPrediction(ctx, name, date)
}
