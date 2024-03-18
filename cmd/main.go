package main

import (
	"VKFilmAPI/internal/Db"
	"VKFilmAPI/internal/Models"
	"fmt"
	"log"
)

func main() {
	db, err := Db.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("database connected")

	err = Db.Prepare()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("database successfully configured")

	actors, err := Models.GetAllActorsWithFilms(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("ok")

	for _, actor := range actors {
		fmt.Println(actor)
	}
}
