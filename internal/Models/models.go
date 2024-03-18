package Models

import (
	"database/sql"
	"time"
)

type Sex int

const (
	Male = iota
	Female
)

type Actor struct {
	Name      string    `json:"name"`
	Gender    Sex       `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type Film struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Year        int     `json:"year"`
	Rating      int8    `json:"rating"`
	Actors      []Actor `json:"actors"`
}

func GetAllFilms(db *sql.DB) ([]Film, error) {
	actors := make(map[int]Actor)
	actorQuery := "SELECT * FROM actors;"
	res, err := db.Query(actorQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var id int
		var name string
		var gender Sex
		var birthdate time.Time
		err = res.Scan(&id, &name, &gender, &birthdate)
		if err != nil {
			return nil, err
		}
	}

	// matches film_id with its actors
	match := make(map[int][]Actor)
	matchQuery := "SELECT * FROM actors_films;"
	res, err = db.Query(matchQuery)
	if err != nil {
		return nil, err
	}
	for res.Next() {

	}

	query := "SELECT * FROM films;"
}

type ActorWithFilms struct {
	Actor Actor  `json:"actor"`
	Films []Film `json:"films"`
}

func GetAllActorsWithFilms(db *sql.DB) ([]ActorWithFilms, error) {

}

func GetAllFilmsByFilmPart(db *sql.DB, filmNamePart string) ([]Film, error) {

}

func GetAllFilmsByActorPart(db *sql.DB, actorNamePart string) ([]Film, error) {

}
