CREATE TABLE IF NOT EXISTS people (
    id serial PRIMARY KEY,
    name VARCHAR (50) NOT NULL,
    surname VARCHAR (50) NOT NULL,
    patronymic VARCHAR (50),
    age INT,
    gender VARCHAR (10),
    nationality VARCHAR (255)
);