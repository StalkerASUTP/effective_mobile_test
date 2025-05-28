package main

import (
	"fmt"
	"go/ef-mob-api/configs"
	"go/ef-mob-api/db"
	"go/ef-mob-api/person"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	//Repository
	personRepository := person.NewPersonRepository(db)

	//Handler
	person.NewPersonHandler(router, person.PersonHandlerDeps{
		PersonRepository: personRepository,
		Config:           conf,
	})
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
