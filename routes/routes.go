package routes

import (
	"github.com/abelkristv/slc_website/handlers"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// secured := router.PathPrefix("/").Subrouter()
	// secured.Use(middleware.TokenValid)
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
}

func RegisterAssistantRoutes(router *mux.Router, assistantHandler *handlers.AssistantHandler) {
	router.HandleFunc("/assistants", assistantHandler.GetAllAssistants).Methods("GET")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.GetAssistantById).Methods("GET")
	router.HandleFunc("/assistants", assistantHandler.CreateAssistant).Methods("POST")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.UpdateAssistant).Methods("PUT")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.DeleteAssistant).Methods("DELETE")
	router.HandleFunc("/assistants/getgenerations", assistantHandler.GetAllGenerations).Methods("GET")
	// router.HandleFunc("/assistants/", assistantHandler.GetAssistantsByGeneration).Methods("GET")
}

func RegisterEventRoutes(router *mux.Router, eventHandler *handlers.EventHandler) {
	router.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventById).Methods("GET")
	router.HandleFunc("/events", eventHandler.CreateEvent).Methods("POST")
	router.HandleFunc("/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods("PUT")
	router.HandleFunc("/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods("DELETE")
}

func RegisterPeriodRoutes(router *mux.Router, periodHandler *handlers.PeriodHandler) {
	router.HandleFunc("/periods", periodHandler.GetAllPeriods).Methods("GET")
	router.HandleFunc("/periods/{id:[0-9]+}", periodHandler.GetPeriodById).Methods("GET")
	router.HandleFunc("/periods", periodHandler.CreatePeriod).Methods("POST")
	router.HandleFunc("/periods/{id:[0-9]+}", periodHandler.UpdatePeriod).Methods("PUT")
	router.HandleFunc("/periods/{id:[0-9]+}", periodHandler.DeletePeriod).Methods("DELETE")
}

func RegisterTeachingHistoryRoutes(router *mux.Router, teachingHistoryHandler *handlers.TeachingHistoryHandler) {
	router.HandleFunc("/teaching-history", teachingHistoryHandler.GetTeachingHistoryByAssistantAndPeriod).Methods("GET")
	router.HandleFunc("/teaching-history/grouped", teachingHistoryHandler.GetTeachingHistoryGroupedByPeriod).Methods("GET")
}

func RegisterPositionRoutes(router *mux.Router, positionHandler *handlers.PositionHandler) {
	router.HandleFunc("/positions", positionHandler.GetAllPositions).Methods("GET")
	router.HandleFunc("/positions/{id}", positionHandler.GetPositionById).Methods("GET")
	router.HandleFunc("/positions", positionHandler.CreatePosition).Methods("POST")
	router.HandleFunc("/positions/{id}", positionHandler.UpdatePosition).Methods("PUT")
	router.HandleFunc("/positions/{id}", positionHandler.DeletePosition).Methods("DELETE")
}
func RegisterAssistantPositionRoutes(router *mux.Router, handler *handlers.AssistantPositionHandler) {
	router.HandleFunc("/assistant_positions", handler.CreatePositionByAssistant).Methods("POST")
	router.HandleFunc("/assistant_positions", handler.GetAllAssistantPositions).Methods("GET")
	router.HandleFunc("/assistant_positions/{id:[0-9]+}", handler.GetAssistantPositionById).Methods("GET")
	router.HandleFunc("/assistant_positions/{id:[0-9]+}", handler.UpdateAssistantPosition).Methods("PUT")
	router.HandleFunc("/assistant_positions/{id:[0-9]+}", handler.DeleteAssistantPosition).Methods("DELETE")
}
