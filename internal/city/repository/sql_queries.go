package repository

const getCityWithPrediction string = `
SELECT c.name, c.country, c.lat, c.lon, p.temp, p.date, p.info
FROM cities c
JOIN predictions p ON p.city_id = c.id
WHERE c.name = $1 AND p.date = $2
LIMIT 1
`

const GetCitiesListWithPredictions string = `
SELECT 
  c.name, 
  c.country, 
  c.lat, 
  c.lon, 
  JSON_AGG(json_build_object('temp', p.temp, 'date', ROUND(EXTRACT(EPOCH FROM p.date)), 'info', p.info)) AS predictions
FROM 
  cities c
JOIN 
  predictions p ON p.city_id = c.id
GROUP BY 
  c.name, c.country, c.lat, c.lon
`

const getCitiesLightListWithPredictions string = `
SELECT c.name
FROM cities c
WHERE EXISTS (
	SELECT 1
	FROM predictions p
	WHERE p.city_id = c.id
)
`

const getCityNameList string = `
SELECT name
FROM cities
`

const saveCity = `
INSERT INTO cities (name, country, lat, lon)
VALUES ($1, $2, $3, $4)
RETURNING id
`

const savePrediction = `
INSERT INTO predictions (city_id, date, temp, info)
VALUES ($1, $2, $3, $4)
`

const cityExists = `SELECT EXISTS (SELECT 1 FROM cities WHERE name = $1)`
