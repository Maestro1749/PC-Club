CREATE TABLE IF NOT EXISTS Users (
    id INT PRIMARY KEY,
    username TEXT(50) UNIQUE NOT NULL,
    fullname TEXT(100) NOT NULL,
    email TEXT(100) UNIQUE,
    phone_number TEXT(20) UNIQUE,
    passwd TEXT(50) NOT NULL,
    bithday DATE NOT NULL,
    balance float DEFAULT 0,
    registered TIMESTAMP NOT NULL,
    privilage TEXT DEFAULT 'Standart'
);