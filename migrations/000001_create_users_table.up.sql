CREATE TABLE IF NOT EXISTS Users (
    id INT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    passwd VARCHAR(50) NOT NULL,
    bithday DATE NOT NULL,
    balance float DEFAULT 0,
    registered TIMESTAMP NOT NULL,
    privilage TEXT DEFAULT 'Standart'
);