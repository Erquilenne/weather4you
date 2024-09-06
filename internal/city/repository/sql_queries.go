package repository

const getCityWithPrediction = `
SELECT c.name, c.country, c.lat, c.lon, p.temp, p.date, p.info
FROM cities c
JOIN predictions p ON p.city_id = c.id
WHERE c.name = $1 AND p.date = $2
LIMIT 1
`

const getCitiesList = `
SELECT c.name, c.country, c.lat, c.lon, p.temp, p.date, p.info
FROM cities c
JOIN predictions p ON p.city_id = c.id
`

const getCitiesLightListWithPredictions = `
SELECT c.name
FROM cities c
WHERE EXISTS (
	SELECT 1
	FROM predictions p
	WHERE p.city_id = c.id
)
`

const getCityNameList = `
SELECT name
FROM cities
`

const addCity = `INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4)`

const addPrediction = `INSERT INTO predictions (city_id, date, temp, info) VALUES ($1, $2, $3, $4)`
