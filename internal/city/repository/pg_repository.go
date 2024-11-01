package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (d *cityRepo) SaveCity(ctx context.Context, city *models.CityDB) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.Create")
	defer span.Finish()
	var id int64
	err := d.db.QueryRowContext(ctx, saveCity, city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *cityRepo) SavePrediction(ctx context.Context, cityId int64, prediction *models.PredictionDB) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.SavePrediction")
	defer span.Finish()
	_, err := d.db.ExecContext(ctx, savePrediction, cityId, prediction.Temp, prediction.Date, prediction.Info)
	return err
}
func (d *cityRepo) Update(ctx context.Context, id int64, city *models.CityDB) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.Create")
	defer span.Finish()
	d.db.QueryRowContext(ctx, updateCity, city.Name, city.Country, city.Lat, city.Lon, id)
	_, err := d.db.ExecContext(ctx, deleteOldPredictions, id)
	if err != nil {
		return err
	}
	var maxDate time.Time
	err = d.db.QueryRowContext(ctx, getMaxDate, id).Scan(&maxDate)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, prediction := range city.Predictions {
		predictionDate := time.Unix(prediction.Date, 0)
		if predictionDate.After(maxDate) {
			_, err := d.db.ExecContext(ctx, savePrediction, id, prediction.Temp, prediction.Date, prediction.Info)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// func GetMaxPredictionDate

func (d *cityRepo) GetMaxPredictionDate(ctx context.Context, id int64) (time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetMaxPredictionDate")
	defer span.Finish()
	var maxDate time.Time
	err := d.db.QueryRowContext(ctx, getMaxDate, id).Scan(&maxDate)
	return maxDate, err
}

func (d *cityRepo) DeleteOldPredictions(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.deleteOldPredictions")
	defer span.Finish()
	_, err := d.db.ExecContext(ctx, deleteOldPredictions, id)
	return err
}

func (d *cityRepo) GetUpdateList(ctx context.Context) ([]*models.CityToUpdate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetUpdateList")
	defer span.Finish()
	var cities []*models.CityToUpdate

	rows, err := d.db.QueryContext(ctx, GetUpdateList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.CityToUpdate
		err := rows.Scan(&city.Id, &city.Name, &city.Lat, &city.Lon)
		if err != nil {
			return nil, err
		}
		cities = append(cities, &city)
	}

	return cities, nil
}

func (d *cityRepo) GetCityId(ctx context.Context, cityName string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCityToUpdate")
	defer span.Finish()
	var id int64
	err := d.db.QueryRowContext(ctx, getCityId, cityName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err // Произошла ошибка при выполнении запроса
	}
	return id, nil
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

	rows, err := d.db.QueryContext(ctx, GetCitiesListWithPredictions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.CityDB
		var predictionsJSON []byte

		err := rows.Scan(&city.Name, &city.Country, &city.Lat, &city.Lon, &predictionsJSON)
		if err != nil {
			return nil, err
		}

		// Десериализация JSON в []models.PredictionDB
		err = json.Unmarshal(predictionsJSON, &city.Predictions)
		if err != nil {
			return nil, err
		}

		cities = append(cities, &city)
	}

	return cities, nil
}

func (d *cityRepo) GetCityWithPrediction(ctx context.Context, name string, date string) (*models.CityWithPrediction, error) {
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
