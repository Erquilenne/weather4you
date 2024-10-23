package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"weather4you/internal/city"
	"weather4you/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type cityRepo struct {
	db *sqlx.DB
}

// City repository constructor
func NewCityRepository(db *sqlx.DB) city.Repository {
	return &cityRepo{db: db}
}

func (d *cityRepo) Create(ctx context.Context, city *models.CityDB) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.Create")
	defer span.Finish()
	var id int64
	d.db.QueryRowContext(ctx, saveCity, city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	for _, prediction := range city.Predictions {
		_, err := d.db.Exec(savePrediction, id, prediction.Temp, prediction.Date, prediction.Info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *cityRepo) GetCitiesList(ctx context.Context) ([]*models.CityLight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCitiesList")
	defer span.Finish()
	var cities []*models.CityLight

	rows, err := d.db.QueryContext(ctx, getCityNameList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.CityLight
		err := rows.Scan(&city.Name)
		if err != nil {
			return nil, err
		}
		cities = append(cities, &city)
	}

	return cities, nil
}

func (d *cityRepo) GetCitiesLightListWithPredictions(ctx context.Context) ([]*models.CityLight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCitiesLightListWithPredictions")
	defer span.Finish()
	var cities []*models.CityLight
	rows, err := d.db.QueryContext(ctx, getCitiesLightListWithPredictions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.CityLight
		err := rows.Scan(&city.Name)
		if err != nil {
			return nil, err
		}
		cities = append(cities, &city)
	}

	return cities, nil
}

func (d *cityRepo) GetCitiesListWithPredictions(ctx context.Context) ([]*models.CityDB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCitiesListWithPredictions")
	defer span.Finish()
	var cities []*models.CityDB

	// err := d.db.SelectContext(ctx, &cities, GetCitiesListWithPredictions)
	rows, err := d.db.QueryContext(ctx, GetCitiesListWithPredictions)
	if err != nil {
		return nil, err
	}

	// for cityDB := range citiesDB {
	// 	prediction := []models.PredictionDB{}
	// 	prediction := json.Unmarshal(cityDB.Predictions, []prediction)
	// 	city := models.City{
	// 		Name:    cityDB.Name,
	// 		Country: cityDB.Country,
	// 		Lat:     cityDB.Lat,
	// 		Lon:     cityDB.Lon,
	// 	}
	// }
	fmt.Println(cities)
	defer rows.Close()

	// cityMap := make(map[string]models.CityDB) // Map to store cities by name

	for rows.Next() {
		// err := rows.Scan(&cityName, &country, &lat, &lon, &temp, &date, &info)
		var city models.CityDB
		var predictions []models.PredictionDB
		if err != nil {
			return nil, err
		}

		err := rows.Scan(&city.Name, &city.Country, &city.Lat, &city.Lon, &city.Predictions)
		if err != nil {
			fmt.Println("Scan error: ", err)
			return nil, err
		}
		cities = append(cities, &city)
		// city.Predictions = append(city.Predictions, predictions...)
		fmt.Println(predictions, "PREDISCTIONSSS")
	}

	// if city, ok := cityMap[city.Name]; ok {
	// 	city.Predictions = append(city.Predictions, predictions)
	// } else {
	// 	city := models.CityDB{
	// 		Name:    cityName,
	// 		Country: country,
	// 		Lat:     lat,
	// 		Lon:     lon,
	// 		Predictions: []models.PredictionDB{
	// 			{
	// 				Temp: temp,
	// 				Date: date,
	// 				Info: string(info),
	// 			},
	// 		},
	// 	}
	// 	cityMap[cityName] = city
	// }
	// 	cityMap[city.Name] = city
	// }

	// for _, city := range cityMap {
	// 	cities = append(cities, &city)
	// }

	return cities, nil
}

func (d *cityRepo) GetCityWithPrediction(ctx context.Context, name string, date time.Time) (*models.CityWithPrediction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCityWithPrediction")
	defer span.Finish()
	row := d.db.QueryRowContext(ctx, getCityWithPrediction, name, date)

	var city models.CityWithPrediction
	var prediction models.PredictionDB

	err := row.Scan(&city.Name, &city.Country, &city.Lat, &city.Lon, &prediction.Temp, &prediction.Date, &prediction.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	city.Prediction = prediction

	return &city, nil
}

func (d *cityRepo) Save(city models.CityDB) error {
	var id int64
	// prettyPrint(city.Name, city.Country, city.Lat, city.Lon)
	d.db.QueryRow(saveCity, city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	for _, prediction := range city.Predictions {
		_, err := d.db.Exec(savePrediction, id, time.Unix(prediction.Date, 0), prediction.Temp, prediction.Info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *cityRepo) Exists(cityName string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(cityExists, cityName).Scan(&exists)
	return exists, err
}
