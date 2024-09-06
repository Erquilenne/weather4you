package main

import (
	"log"
	"os"
	"weather4you/config"
	"weather4you/internal/fillup"
	"weather4you/pkg/db/postgres"
	"weather4you/pkg/logger"
	"weather4you/pkg/utils"
)

func main() {
	// TODO

	log.Println("Starting api server")

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
	appLogger.Infof("AppVersion: %s", "LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	db, err := postgres.NewPsqlDB(cfg)

	dbcities, err := db.GetCitiesList()
	if err != nil {
		log.Fatal("Error on getting cities:", err)
	}
	if len(dbcities) == 0 {
		cities := cfg.StartCities
		for _, city := range cities {
			err := fillup.FindAndSaveCity(city, db, cfg)
			if err != nil {
				log.Fatal("Error on saving city:", err)
			}
		}
	}
}
