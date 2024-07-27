package pgsql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"weather4you/internal/config"
	dbModels "weather4you/internal/models/db"
	reqModels "weather4you/internal/models/request"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Database struct {
	db  *sql.DB
	cfg config.Config
}

func NewDatabase(cfg config.Config) (*Database, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	database := &Database{db: db, cfg: cfg}

	return database, nil
}

func (d *Database) Stats() sql.DBStats {
	return d.db.Stats()
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) MakeMigrations() {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		log.Fatal("Error creating database driver instance:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Error creating migration instance:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations:", err)
	}

	log.Println("Migrations up to date")
}

func (d *Database) SaveCity(city dbModels.City) error {
	var id int64
	fmt.Println(city.Name)
	d.db.QueryRow("INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4) RETURNING id", city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	for _, prediction := range city.Predictions {
		_, err := d.db.Exec("INSERT INTO predictions (city_id, temp, date, info) VALUES ($1, $2, $3, $4)", id, prediction.Temp, prediction.Date, prediction.Info)
		if err != nil {
			log.Fatal("Error inserting prediction: ", err)
		}
	}
	return nil
}

func (d *Database) GetCitiesList() ([]reqModels.CityLight, error) {
	var cities []reqModels.CityLight

	rows, err := d.db.Query("SELECT name FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city reqModels.CityLight
		err := rows.Scan(&city.Name)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (d *Database) GetCitiesLightListWithPredictions() ([]reqModels.CityLight, error) {
	var cities []reqModels.CityLight

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
		var city reqModels.CityLight
		err := rows.Scan(&city.Name)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (d *Database) GetCitiesListWithPredictions() ([]dbModels.City, error) {
	var cities []dbModels.City

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

	cityMap := make(map[string]*dbModels.City) // Map to store cities by name

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
			prediction := dbModels.Prediction{
				Temp: temp,
				Date: date,
				Info: string(info),
			}
			city.Predictions = append(city.Predictions, prediction)
		} else {
			city := &dbModels.City{
				Name:    cityName,
				Country: country,
				Lat:     lat,
				Lon:     lon,
				Predictions: []dbModels.Prediction{
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
		cities = append(cities, *city)
	}

	return cities, nil
}

func (d *Database) GetCityWithPrediction(name string, date time.Time) (*reqModels.CityWithPrediction, error) {
	query := `
		SELECT c.name, c.country, c.lat, c.lon, p.temp, p.date, p.info
		FROM cities c
		JOIN predictions p ON p.city_id = c.id
		WHERE c.name = $1 AND p.date = $2
		LIMIT 1
	`
	row := d.db.QueryRow(query, name, date)

	var city reqModels.CityWithPrediction
	var prediction dbModels.Prediction

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
