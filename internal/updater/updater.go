package updater

import (
	"context"
	"weather4you/config"
	"weather4you/internal/city"
	"weather4you/internal/city/repository"
	"weather4you/internal/weatherapi"
	"weather4you/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type Updater struct {
	cfg        *config.Config
	repository city.Repository
	logger     logger.Logger
}

func NewUpdater(cfg *config.Config, db *sqlx.DB, logger logger.Logger) *Updater {
	repository := repository.NewCityRepository(db)
	return &Updater{cfg: cfg, repository: repository, logger: logger}
}

func (u *Updater) Update() {
	ctx := context.Background()
	cityList, err := u.repository.GetUpdateList(ctx)
	if err != nil {
		u.logger.Fatalf("GetCitiesListWithPredictions error: %s", err)
	}
	if len(cityList) == 0 {
		u.logger.Fatal("No cities found. Try to 'make fillup' first")
		return
	}
	finder := weatherapi.NewCityFinder(u.cfg, u.logger)
	for _, city := range cityList {
		u.logger.Infof("Updating city: %s", city.Name)
		predictions := finder.FindPredictions(city.Lat, city.Lon)
		for _, prediction := range predictions {
			err = u.repository.SavePrediction(ctx, city.Id, &prediction)
			if err != nil {
				u.logger.Fatalf("UpdateCity error: %s", err)
			}
		}
		u.logger.Infof("CityId: %d", city.Id)
		err = u.repository.DeleteOldPredictions(ctx, city.Id)
		if err != nil {
			u.logger.Fatalf("DeleteOldPredictions error: %s", err)
		}
		u.logger.Info("old predictions deleted")
		u.logger.Infof("City updated: %s", city.Name)
	}
}
