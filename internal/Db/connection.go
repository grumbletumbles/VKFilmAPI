package Db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func Connect() (*sql.DB, error) {
	return sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@host.docker.internal/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME")))
}

func Prepare() error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
CREATE TABLE IF NOT EXISTS films (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	date DATE NOT NULL,
	rating INT8 NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sex') THEN
        CREATE TYPE sex AS ENUM (
			'male',
			'female'
		);
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS actors (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    gender sex NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS actors_films (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    film_id BIGINT REFERENCES films(id),
    actor_id BIGINT REFERENCES actors(id)
)
`
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
