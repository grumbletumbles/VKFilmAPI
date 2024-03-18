package Models

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"
)

func contains(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

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
	Id        int
	Name      string    `json:"name"`
	Gender    Sex       `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type Film struct {
	Id          int
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
			Id:        id,
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
			Id:          id,
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

func GetAllActorsWithFilms(db *sql.DB) ([]ActorWithFilms, error) {
	actors := make(map[int]ActorWithFilms)
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
		actors[id] = ActorWithFilms{
			Actor: Actor{
				Id:        id,
				Name:      name,
				Gender:    MakeSexFromString(gender),
				BirthDate: birthdate,
			},
			Films: nil,
		}
	}

	films, err := GetAllFilms(db)
	if err != nil {
		return nil, err
	}

	for _, film := range films {
		for _, a := range film.Actors {
			if entry, ok := actors[a.Id]; ok {
				entry.Films = append(entry.Films, film)
				actors[a.Id] = entry
			} else {
				log.Printf("cannot find actor %s in database for %s\n", a.Name, film.Name)
			}
		}
	}

	var result []ActorWithFilms
	for _, t := range actors {
		result = append(result, t)
	}

	return result, nil
}

/*
Note: the regex should be replaced with finding the substring since .* is
insanely ineffective and requires a lot of backtracking
*/

func GetAllFilmsByFilmPart(db *sql.DB, filmNamePart string) ([]Film, error) {
	films, err := GetAllFilms(db)
	if err != nil {
		return nil, err
	}
	var result []Film
	matchString := fmt.Sprintf(".*(%s).*", filmNamePart)
	re, err := regexp.Compile(matchString)
	if err != nil {
		return nil, err
	}

	for _, film := range films {
		if re.MatchString(film.Name) {
			result = append(result, film)
		}
	}

	return result, nil
}

func GetAllFilmsByActorPart(db *sql.DB, actorNamePart string) ([]Film, error) {
	films, err := GetAllFilms(db)
	if err != nil {
		return nil, err
	}
	var result []Film
	matchString := fmt.Sprintf(".*(%s).*", actorNamePart)
	re, err := regexp.Compile(matchString)
	if err != nil {
		return nil, err
	}

	for _, film := range films {
		found := false
		for _, actor := range film.Actors {
			if re.MatchString(actor.Name) {
				found = true
				break
			}
		}
		if found {
			result = append(result, film)
		}
	}

	return result, nil
}
