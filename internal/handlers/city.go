package handlers

import (
	"weather4you/internal/storage/pgsql"
	"weather4you/internal/weatherapi"
)

// TODO take city from api and save it to postgres
func SaveCity(city string, d *pgsql.Database) {
	findedCity, err := weatherapi.FindCity(city)
	if err != nil {
		return
	}
	d.SaveCity(findedCity)
}
