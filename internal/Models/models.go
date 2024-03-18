package Models

import (
	"database/sql"
	"time"
)

type Sex int8

const (
	Male = iota
	Female
)

func MakeSexFromString(str string) Sex {
	if str == "male" {
		return Male
	} else {
		return Female
	}
}

type Actor struct {
	Name      string    `json:"name"`
	Gender    Sex       `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type Film struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Rating      int8      `json:"rating"`
	Actors      []Actor   `json:"actors"`
}

func GetAllFilms(db *sql.DB) ([]Film, error) {
	actors := make(map[int]Actor)
	actorQuery := "SELECT * FROM actors;"
	res1, err := db.Query(actorQuery)
	if err != nil {
		return nil, err
	}
	for res1.Next() {
		var id int
		var name string
		var gender string
		var birthdate time.Time
		err = res1.Scan(&id, &name, &gender, &birthdate)
		if err != nil {
			return nil, err
		}
		actors[id] = Actor{
			Name:      name,
			Gender:    MakeSexFromString(gender),
			BirthDate: birthdate,
		}
	}

	// matches film_id with its actors
	match := make(map[int][]Actor)
	matchQuery := "SELECT * FROM actors_films;"
	res2, err := db.Query(matchQuery)
	if err != nil {
		return nil, err
	}
	for res2.Next() {
		var id int
		var filmId int
		var actorId int
		err = res2.Scan(&id, &filmId, &actorId)
		if err != nil {
			return nil, err
		}
		match[filmId] = append(match[filmId], actors[actorId])
	}

	query := "SELECT * FROM films;"
	var films []Film
	res3, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	for res3.Next() {
		var id int
		var name string
		var description string
		var date time.Time
		var rating int8
		err = res3.Scan(&id, &name, &description, &date, &rating)
		if err != nil {
			return nil, err
		}
		films = append(films, Film{
			Name:        name,
			Description: description,
			Date:        date,
			Rating:      rating,
			Actors:      match[id],
		})
	}

	return films, nil
}

type ActorWithFilms struct {
	Actor Actor  `json:"actor"`
	Films []Film `json:"films"`
}

/*
func GetAllActorsWithFilms(db *sql.DB) ([]ActorWithFilms, error) {

}

func GetAllFilmsByFilmPart(db *sql.DB, filmNamePart string) ([]Film, error) {

}

func GetAllFilmsByActorPart(db *sql.DB, actorNamePart string) ([]Film, error) {

}
*/
