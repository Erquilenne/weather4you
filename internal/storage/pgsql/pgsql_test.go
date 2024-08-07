package pgsql

import (
	"time"
	dbModels "weather4you/internal/models/db"
	requestModels "weather4you/internal/models/request"
)

type DatabaseInterface interface {
	SaveCity(city dbModels.City) error
	GetCitiesList() ([]dbModels.City, error)
	GetCityWithPrediction(name string, date time.Time) (*requestModels.CityWithPrediction, error)
}

type FakeDatabase struct {
	SavedCity          dbModels.City
	CitiesList         []dbModels.City
	CityWithPrediction *requestModels.CityWithPrediction
}

func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{}
}

// SaveCity сохраняет город в фейковой базе данных.
func (d *FakeDatabase) SaveCity(city dbModels.City) error {
	d.SavedCity = city
	return nil
}

// GetCitiesList возвращает список городов из фейковой базы данных.
func (d *FakeDatabase) GetCitiesList() ([]dbModels.City, error) {
	return d.CitiesList, nil
}

// GetCityWithPrediction возвращает город с прогнозом из фейковой базы данных.
func (d *FakeDatabase) GetCityWithPrediction(name string, date time.Time) (*requestModels.CityWithPrediction, error) {
	return d.CityWithPrediction, nil
}
