package pgsql

import (
	"database/sql"
	"encoding/json"
	"log"

	"weather4you/internal/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
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

	log.Println("Migrations applied successfully!")
}

// func (d *Database) SaveCity(city models.City) {
// 	_, err := d.db.Exec("INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4) RETURNING id", city.Name, city.Country, city.Lat, city.Lon)
// 	if err != nil {
// 		log.Fatal("Error saving city:", err)
// 	}
// }

func (d *Database) SaveCity(city models.City) {
	var id int64
	d.db.QueryRow("INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4) RETURNING id", city.Name, city.Country, city.Lat, city.Lon).Scan(&id)
	for _, prediction := range city.Predictions {
		infoJSON, err := json.Marshal(prediction.Info)
		_, err = d.db.Exec("INSERT INTO predictions (city_id, temp, date, info) VALUES ($1, $2, $3, $4)", id, prediction.Temp, prediction.Date, infoJSON)
		if err != nil {
			log.Fatal("Error inserting prediction: ", err)
		}
	}
}

func (d *Database) GetCitiesList() ([]models.CityLight, error) {
	var cities []models.CityLight

	rows, err := d.db.Query("SELECT name FROM cities")
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
		cities = append(cities, city)
	}

	return cities, nil
}
