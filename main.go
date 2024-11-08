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

	// User setup
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Assistant setup
	assistantRepo := repositories.NewAssistantRepository(db)
	assistantService := services.NewAssistantService(assistantRepo)
	assistantHandler := handlers.NewAssistantHandler(assistantService)

	// Event setup
	eventRepo := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	periodRepo := repositories.NewPeriodRepository(db)
	periodService := services.NewPeriodService(periodRepo)
	periodHandler := handlers.NewPeriodHandler(periodService)

	teachingHistoryRepo := repositories.NewTeachingHistoryRepository(db)
	teachingHistoryService := services.NewTeachingHistoryService(teachingHistoryRepo)
	teachingHistoryHandler := handlers.NewTeachingHistoryHandler(teachingHistoryService)

	contactUsRepo := repositories.NewContactUsRepository(db)
	contactUsService := services.NewContactUsService(contactUsRepo)
	contactUsHandler := handlers.NewContactUsHandler(contactUsService)

	assistantSocialMediaRepo := repositories.NewAssistantSocialMediaRepository(db)
	assistantSocialMediaService := services.NewAssistantSocialMediaService(assistantSocialMediaRepo)
	assistantSocialMediaHandler := handlers.NewAssistantSocialMediaHandler(assistantSocialMediaService, *userService)

	awardRepo := repositories.NewAwardRepository(db)
	awardService := services.NewAwardService(awardRepo)
	awardHandler := handlers.NewAwardHandler(awardService)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	routes.RegisterUserRoutes(router, userHandler)
	routes.RegisterAssistantRoutes(router, assistantHandler)
	routes.RegisterEventRoutes(router, eventHandler)
	routes.RegisterPeriodRoutes(router, periodHandler)
	routes.RegisterTeachingHistoryRoutes(router, teachingHistoryHandler)

	routes.RegisterContactUsRoutes(router, contactUsHandler)
	routes.RegisterAssistantSocialMediaRoutes(router, assistantSocialMediaHandler)
	routes.RegisterAwardRoutes(router, awardHandler)

	log.Println("Starting server on :8888")
	log.Fatal(http.ListenAndServe(":8888", handler))

}
