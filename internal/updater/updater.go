package updater

import (
	"context"
	"weather4you/config"
	"weather4you/internal/city"
	"weather4you/internal/city/repository"
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
	cityList, err := u.repository.GetCitiesListWithPredictions(context.Background())
	if err != nil {
		u.logger.Fatalf("GetCitiesListWithPredictions error: %s", err)
	}
	if len(cityList) == 0 {
		u.logger.Fatal("No cities found. Try to 'make fillup' first")
		return
	}
	for _, city := range cityList {
		u.logger.Infof("Updating city: %s", city.Name)
		//TODO
	}
}
