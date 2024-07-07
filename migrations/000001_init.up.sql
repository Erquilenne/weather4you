CREATE TABLE IF NOT EXISTS cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    lat FLOAT NOT NULL,
    lon FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    city_id INT NOT NULL,
    temp INT NOT NULL,
    date TIMESTAMP NOT NULL,
    info JSON NOT NULL,
    CONSTRAINT fk_city
        FOREIGN KEY (city_id)
        REFERENCES cities(id)
        ON DELETE CASCADE
);