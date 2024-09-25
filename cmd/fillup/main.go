package main

import (
	"context"
	"log"
	"os"
	"weather4you/config"
	"weather4you/internal/city/repository"
	"weather4you/internal/fillup"
	"weather4you/pkg/db/postgres"
	"weather4you/pkg/logger"
	"weather4you/pkg/utils"

	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
)

func main() {
	log.Println("Starting filling up")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %s", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %s", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s", "LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}
	cityRepository := repository.NewCityRepository(db)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("fillup.main.GetCitiesList")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()
	dbcities, err := cityRepository.GetCitiesList(ctx)
	if err != nil {
		log.Fatal("Error on getting cities:", err)
	}
	if len(dbcities) == 0 {
		cities := cfg.StartCities
		for _, city := range cities {
			fillup.FindAndSaveCity(city, cityRepository, cfg, appLogger)
		}
	}
}
