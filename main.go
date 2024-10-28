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
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	db, err := database.InitializeDB()

	for _, arg := range os.Args {
		if arg == "--seed" {
			database.SeedDatabase(db)
		} else if arg == "--clear" {
			database.ClearDatabase(db)
			return
		}
	}

	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Assistant handler
	assistantRepo := repositories.NewAssistantRepository(db)
	assistantService := services.NewAssistantService(assistantRepo)
	assistantHandler := handlers.NewAssistantHandler(assistantService)

	// Initialize Event handler
	eventRepo := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	// Set up CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	// Register routes for users, assistants, and events
	routes.RegisterUserRoutes(router, userHandler)
	routes.RegisterAssistantRoutes(router, assistantHandler)
	routes.RegisterEventRoutes(router, eventHandler)

	log.Println("Starting server on :8888")
	log.Fatal(http.ListenAndServe(":8888", handler))
}
