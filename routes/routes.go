package routes

import (
	"github.com/abelkristv/slc_website/handlers"
	"github.com/abelkristv/slc_website/middleware"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	secured.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	secured.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	secured.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	secured.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
}

func RegisterAssistantRoutes(router *mux.Router, assistantHandler *handlers.AssistantHandler) {
	router.HandleFunc("/assistants", assistantHandler.GetAllAssistants).Methods("GET")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.GetAssistantById).Methods("GET")
	router.HandleFunc("/assistants", assistantHandler.CreateAssistant).Methods("POST")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.UpdateAssistant).Methods("PUT")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.DeleteAssistant).Methods("DELETE")
}
