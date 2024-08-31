package repository

import (
	"context"
	"database/sql"
	"time"
	"weather4you/internal/city"
	"weather4you/internal/models"

	"github.com/opentracing/opentracing-go"
)

type newsRepo struct {
	db *sql.DB
}

// News repository constructor
func NewNewsRepository(db *sql.DB) city.Repository {
	return &newsRepo{db: db}
}

// func (d *newsRepo) MakeMigrations() {

// 	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
// 	if err != nil {
// 		log.Fatal("Error creating database driver instance:", err)
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://./migrations",
// 		"postgres", driver)
// 	if err != nil {
// 		log.Fatal("Error creating migration instance:", err)
// 	}

// 	err = m.Up()
// 	if err != nil && err != migrate.ErrNoChange {
// 		log.Fatal("Error applying migrations:", err)
// 	}

// 	log.Println("Migrations up to date")
// }

func (d *newsRepo) Create(ctx context.Context, city *models.CityDB) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.Create")
	defer span.Finish()
	var id int64
	d.db.QueryRow("INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4) RETURNING id", city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	for _, prediction := range city.Predictions {
		_, err := d.db.Exec("INSERT INTO predictions (city_id, temp, date, info) VALUES ($1, $2, $3, $4)", id, prediction.Temp, prediction.Date, prediction.Info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *newsRepo) GetCitiesList(ctx context.Context) ([]*models.CityLight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCitiesList")
	defer span.Finish()
	var cities []*models.CityLight

	rows, err := d.db.Query("SELECT name FROM cities")
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

func (d *newsRepo) GetCitiesLightListWithPredictions(ctx context.Context) ([]*models.CityLight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCitiesLightListWithPredictions")
	defer span.Finish()
	var cities []*models.CityLight

	query := `
		SELECT c.name
		FROM cities c
		WHERE EXISTS (
			SELECT 1
			FROM predictions p
			WHERE p.city_id = c.id
		)
	`
	rows, err := d.db.Query(query)
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

func (d *newsRepo) GetCitiesListWithPredictions(ctx context.Context) ([]*models.CityDB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCitiesListWithPredictions")
	defer span.Finish()
	var cities []*models.CityDB

	query := `
        SELECT c.name, c.country, c.lat, c.lon, p.temp, p.date, p.info
        FROM cities c
        JOIN predictions p ON p.city_id = c.id
    `
	rows, err := d.db.Query(query)
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

func (d *newsRepo) GetCityWithPrediction(ctx context.Context, name string, date time.Time) (*models.CityWithPrediction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cityRepo.GetCityWithPrediction")
	defer span.Finish()
	query := `
		SELECT c.name, c.country, c.lat, c.lon, p.temp, p.date, p.info
		FROM cities c
		JOIN predictions p ON p.city_id = c.id
		WHERE c.name = $1 AND p.date = $2
		LIMIT 1
	`
	row := d.db.QueryRow(query, name, date)

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
