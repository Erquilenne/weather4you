package repository

import (
	"context"
	"database/sql"
	"time"
	"weather4you/internal/city"
	"weather4you/internal/models"

	"github.com/opentracing/opentracing-go"
)

type cityRepo struct {
	db *sql.DB
}

// City repository constructor
func NewCityRepository(db *sql.DB) city.Repository {
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
	// rows, err := d.db.Query("SELECT name FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.CityLight
		err := rows.Scan(city.Name)
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

	rows, err := d.db.QueryContext(ctx, getCitiesList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cityMap := make(map[string]models.CityDB) // Map to store cities by name

	for rows.Next() {
		var cityName, country string
		var lat, lon float64
		var temp int
		var date time.Time
		var info []byte

		err := rows.Scan(&cityName, &country, &lat, &lon, &temp, &date, &info)
		if err != nil {
			return nil, err
		}

		if city, ok := cityMap[cityName]; ok {
			prediction := models.PredictionDB{
				Temp: temp,
				Date: date,
				Info: string(info),
			}
			city.Predictions = append(city.Predictions, prediction)
		} else {
			city := models.CityDB{
				Name:    cityName,
				Country: country,
				Lat:     lat,
				Lon:     lon,
				Predictions: []models.PredictionDB{
					{
						Temp: temp,
						Date: date,
						Info: string(info),
					},
				},
			}
			cityMap[cityName] = city
		}
	}

	for _, city := range cityMap {
		cities = append(cities, &city)
	}

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
	d.db.QueryRow(saveCity, city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	for _, prediction := range city.Predictions {
		_, err := d.db.Exec(savePrediction, id, prediction.Temp, prediction.Date, prediction.Info)
		if err != nil {
			return err
		}
	}
	return nil
}
