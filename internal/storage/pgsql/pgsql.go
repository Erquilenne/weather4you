package pgsql

import (
	"database/sql"
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

func (d *Database) SaveCity(city models.City) {
	_, err := d.db.Exec("INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4) RETURNING id", city.Name, city.Country, city.Lat, city.Lon)
	if err != nil {
		log.Fatal("Error saving city:", err)
	}
}
