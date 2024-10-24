package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abelkristv/slc_website/database"
	"github.com/abelkristv/slc_website/handlers"
	"github.com/abelkristv/slc_website/repositories"
	"github.com/abelkristv/slc_website/routes"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	db, err := database.InitializeDB()

	for _, arg := range os.Args {
		if arg == "--seed" {
			database.SeedDatabase(db)
		}
	}

	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	routes.RegisterUserRoutes(router, userHandler)

	log.Println("Starting server on :8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
