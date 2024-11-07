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
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	secured.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	secured.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	secured.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
	secured.HandleFunc("/currentuser", userHandler.GetCurrentUser).Methods("GET")
	secured.HandleFunc("/logout", userHandler.Logout).Methods("POST")
	secured.HandleFunc("/change-password", userHandler.ChangePassword).Methods("PUT")
}

func RegisterAssistantRoutes(router *mux.Router, assistantHandler *handlers.AssistantHandler) {
	router.HandleFunc("/assistants", assistantHandler.GetAllAssistants).Methods("GET")
	router.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.GetAssistantById).Methods("GET")

	router.HandleFunc("/assistants/getgenerations", assistantHandler.GetAllGenerations).Methods("GET")
	// router.HandleFunc("/assistants/", assistantHandler.GetAssistantsByGeneration).Methods("GET")

	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/assistants", assistantHandler.CreateAssistant).Methods("POST")
	secured.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.UpdateAssistant).Methods("PUT")
	secured.HandleFunc("/assistants/{id:[0-9]+}", assistantHandler.DeleteAssistant).Methods("DELETE")
}

func RegisterEventRoutes(router *mux.Router, eventHandler *handlers.EventHandler) {
	router.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventById).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/events", eventHandler.CreateEvent).Methods("POST")
	secured.HandleFunc("/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods("PUT")
	secured.HandleFunc("/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods("DELETE")
}

func RegisterPeriodRoutes(router *mux.Router, periodHandler *handlers.PeriodHandler) {
	router.HandleFunc("/periods", periodHandler.GetAllPeriods).Methods("GET")
	router.HandleFunc("/periods/{id:[0-9]+}", periodHandler.GetPeriodById).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/periods", periodHandler.CreatePeriod).Methods("POST")
	secured.HandleFunc("/periods/{id:[0-9]+}", periodHandler.UpdatePeriod).Methods("PUT")
	secured.HandleFunc("/periods/{id:[0-9]+}", periodHandler.DeletePeriod).Methods("DELETE")
}

func RegisterTeachingHistoryRoutes(router *mux.Router, teachingHistoryHandler *handlers.TeachingHistoryHandler) {
	router.HandleFunc("/teaching-history", teachingHistoryHandler.GetTeachingHistoryByAssistantAndPeriod).Methods("GET")
	router.HandleFunc("/teaching-history/grouped", teachingHistoryHandler.GetTeachingHistoryGroupedByPeriod).Methods("GET")
}

func RegisterPositionRoutes(router *mux.Router, positionHandler *handlers.PositionHandler) {
	router.HandleFunc("/positions", positionHandler.GetAllPositions).Methods("GET")
	router.HandleFunc("/positions/{id}", positionHandler.GetPositionById).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/positions", positionHandler.CreatePosition).Methods("POST")
	secured.HandleFunc("/positions/{id}", positionHandler.UpdatePosition).Methods("PUT")
	secured.HandleFunc("/positions/{id}", positionHandler.DeletePosition).Methods("DELETE")
}
func RegisterAssistantPositionRoutes(router *mux.Router, handler *handlers.AssistantPositionHandler) {
	router.HandleFunc("/assistant_positions", handler.GetAllAssistantPositions).Methods("GET")
	router.HandleFunc("/assistant_positions/{id:[0-9]+}", handler.GetAssistantPositionById).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/assistant_positions", handler.CreatePositionByAssistant).Methods("POST")
	secured.HandleFunc("/assistant_positions/{id:[0-9]+}", handler.UpdateAssistantPosition).Methods("PUT")
	secured.HandleFunc("/assistant_positions/{id:[0-9]+}", handler.DeleteAssistantPosition).Methods("DELETE")
}

func RegisterContactUsRoutes(router *mux.Router, contactUsHandler *handlers.ContactUsHandler) {
	router.HandleFunc("/contacts", contactUsHandler.GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts/{id:[0-9]+}", contactUsHandler.GetContactById).Methods("GET")
	router.HandleFunc("/contacts", contactUsHandler.CreateContact).Methods("POST")
	router.HandleFunc("/contacts/{id:[0-9]+}", contactUsHandler.UpdateContact).Methods("PUT")
	router.HandleFunc("/contacts/{id:[0-9]+}", contactUsHandler.DeleteContact).Methods("DELETE")
}

func RegisterAssistantSocialMediaRoutes(router *mux.Router, handler *handlers.AssistantSocialMediaHandler) {
	router.HandleFunc("/assistant_social_media/{id:[0-9]+}", handler.GetAssistantSocialMediaByID).Methods("GET")
	router.HandleFunc("/assistant_social_media/assistant/{assistantId:[0-9]+}", handler.GetAssistantSocialMediaByAssistantID).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/assistant_social_media", handler.CreateOrUpdateAssistantSocialMedia).Methods("POST")
	// secured.HandleFunc("/assistant_social_media/{id:[0-9]+}", handler.UpdateAssistantSocialMedia).Methods("PUT")
	secured.HandleFunc("/assistant_social_media/{id:[0-9]+}", handler.DeleteAssistantSocialMedia).Methods("DELETE")
}
