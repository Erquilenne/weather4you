package repository

const getCityWithPrediction string = `
SELECT c.name, c.country, c.lat, c.lon, p.temp, ROUND(EXTRACT(EPOCH FROM p.date)) as date, p.info
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

const GetUpdateList string = `
SELECT c.name, c.lat, c.lon
FROM cities c
WHERE EXISTS (
  SELECT 1
  FROM predictions p
  WHERE p.city_id = c.id
)
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

const deleteOldPredictions string = `
        DELETE FROM predictions 
        WHERE city_id = $1 AND date < NOW()
    `

const getMaxDate string = `
        SELECT MAX(date)
        FROM predictions
        WHERE city_id = $1
    `

const getCityNameList string = `
SELECT name
FROM cities
`
const getCityId string = `
SELECT id
FROM cities
WHERE name = $1
`

const saveCity = `
INSERT INTO cities (name, country, lat, lon)
VALUES ($1, $2, $3, $4)
RETURNING id
`
const updateCity = `
UPDATE cities
SET name = $1, country = $2, lat = $3, lon = $4
WHERE id = $5
RETURNING id
`

const savePrediction = `
INSERT INTO predictions (city_id, date, temp, info)
VALUES ($1, $2, $3, $4)
`

const cityExists = `SELECT EXISTS (SELECT 1 FROM cities WHERE name = $1)`
