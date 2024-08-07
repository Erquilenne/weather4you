package postgres

import (
	"database/sql"
	"fmt"
	"time"
	"weather4you/config"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// Return new Postgresql db instance
func NewPsqlDB(c *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.DBName,
		c.Database.Password,
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
