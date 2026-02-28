CREATE TABLE IF NOT EXISTS Computers (
    id SERIAL PRIMARY KEY,
    num TEXT(10) UNIQUE NOT NULL,
    price FLOAT NOT NULL,
    is_busy BOOLEAN DEFAULT FALSE,
    busy_since TIMESTAMP,
    busy_until TIMESTAMP
);