package main

import (
	"fmt"
	"go/ef-mob-api/configs"
	_ "go/ef-mob-api/configs"
	"go/ef-mob-api/db"
	_ "go/ef-mob-api/db"
	_ "go/ef-mob-api/docs"
	"go/ef-mob-api/person"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Effective Mobile API
// @version         1.0
// @description     Сервис для обогащения ФИО информацией из открытых API (возраст, пол, национальность) и хранения данных в базе.
// @host            localhost:8081
// @BasePath        /
// @schemes         http

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repository
	personRepository := person.NewPersonRepository(db)

	// Handler
	person.NewPersonHandler(router, person.PersonHandlerDeps{
		PersonRepository: personRepository,
		Config:           conf,
	})

	router.Handle("/swagger/", httpSwagger.WrapHandler)

	return router
}

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}
	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
