package usecase

import (
	"context"
	"strconv"
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
	id, err := u.cityRepo.SaveCity(ctx, city)
	if err != nil {
		return err
	}

	for _, prediction := range city.Predictions {
		err = u.cityRepo.SavePrediction(ctx, id, &prediction)
		if err != nil {
			return err
		}
	}
	return nil
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

func (u *cityUC) GetCityWithPrediction(ctx context.Context, name string, timestamp string) (*models.CityWithPrediction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityUC.GetCityWithPrediction")
	defer span.Finish()

	unixTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, err
	}

	// Преобразуем Unix timestamp в объект time.Time
	date := time.Unix(unixTimestamp, 0)

	// Преобразуем объект time.Time в строку формата 2006-01-02T15:04:05Z
	formattedDate := date.Format("2006-01-02 15:04:05")
	return u.cityRepo.GetCityWithPrediction(ctx, name, formattedDate)
}
