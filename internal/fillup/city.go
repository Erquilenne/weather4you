package fillup

import (
	"fmt"
	"weather4you/internal/storage/pgsql"
	"weather4you/internal/weatherapi"
)

// TODO take city from api and save it to postgres
func SaveCity(city string, d *pgsql.Database) error {
	fmt.Println("Saving city:", city)
	findedCity, err := weatherapi.FindCity(city)
	if err != nil {
		return err
	}
	fmt.Println("Found city:", findedCity)
	err = d.SaveCity(findedCity)
	if err != nil {
		return err
	}
	return nil
}
