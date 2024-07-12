package main

import (
	"fmt"
	"log"
	"net/http"
	"weather4you/internal/config"
	"weather4you/internal/fillup"
	"weather4you/internal/http-server/handlers"
	"weather4you/internal/storage/pgsql"
)

func main() {
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)

	db, err := pgsql.NewDatabase(connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	defer db.Close()

	db.MakeMigrations()

	dbcities, err := db.GetCitiesList()
	if err != nil {
		log.Fatal("Error on getting cities:", err)
	}
	if len(dbcities) == 0 {
		cities := config.StartCities
		for _, city := range cities {
			err := fillup.SaveCity(city, db)
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
