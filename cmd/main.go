package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"weather4you/internal/config"
	"weather4you/internal/fillup"
	"weather4you/internal/http-server/handlers"
	"weather4you/internal/storage/pgsql"
	"weather4you/pkg/logger"
	"weather4you/pkg/utils"
)

func main() {
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

	db, err := pgsql.NewDatabase(*cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", db.Stats())
	}
	defer db.Close()

	db.MakeMigrations()

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

	fmt.Println("Done!")

	handler := handlers.NewHandler(db)
	http.HandleFunc("/list/", handler.GetList)
	http.HandleFunc("/predictions/", handler.GetPredictionsList)
	http.HandleFunc("/prediction/", handler.GetCityWithPrediction)

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
