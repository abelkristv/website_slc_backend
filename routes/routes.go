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

func RegisterContactUsRoutes(router *mux.Router, contactUsHandler *handlers.ContactUsHandler) {
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/contacts", contactUsHandler.GetAllContacts).Methods("GET")
	secured.HandleFunc("/contacts/{id:[0-9]+}", contactUsHandler.GetContactById).Methods("GET")
	router.HandleFunc("/contacts", contactUsHandler.CreateContact).Methods("POST")
	secured.HandleFunc("/contacts/{id:[0-9]+}", contactUsHandler.UpdateContact).Methods("PUT")
	secured.HandleFunc("/contacts/isread/{id}", contactUsHandler.UpdateIsRead).Methods("PATCH")
	secured.HandleFunc("/contacts/{id:[0-9]+}", contactUsHandler.DeleteContact).Methods("DELETE")
}

func RegisterAssistantSocialMediaRoutes(router *mux.Router, handler *handlers.AssistantSocialMediaHandler) {
	// router.HandleFunc("/assistant_social_media/{id:[0-9]+}", handler.GetAssistantSocialMediaByID).Methods("GET")
	router.HandleFunc("/assistant_social_media/assistant/{assistantId:[0-9]+}", handler.GetAssistantSocialMediaByAssistantID).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/assistant_social_media", handler.CreateOrUpdateAssistantSocialMedia).Methods("POST")
	// secured.HandleFunc("/assistant_social_media/{id:[0-9]+}", handler.UpdateAssistantSocialMedia).Methods("PUT")
	// secured.HandleFunc("/assistant_social_media/{id:[0-9]+}", handler.DeleteAssistantSocialMedia).Methods("DELETE")
}

func RegisterAwardRoutes(router *mux.Router, handler *handlers.AwardHandler) {
	router.HandleFunc("/awards", handler.CreateAward).Methods("POST")
	router.HandleFunc("/awards/{id:[0-9]+}", handler.GetAwardByID).Methods("GET")
	router.HandleFunc("/awards/{id:[0-9]+}", handler.UpdateAward).Methods("PUT")
	router.HandleFunc("/awards/{id:[0-9]+}", handler.DeleteAward).Methods("DELETE")
	router.HandleFunc("/awards", handler.GetAllAwards).Methods("GET")
}

func RegisterAssistantAwardRoutes(router *mux.Router, handler *handlers.AssistantAwardHandler) {

	router.HandleFunc("/assistant_awards/{id:[0-9]+}", handler.GetAssistantAwardByID).Methods("GET")
	router.HandleFunc("/assistant_awards/assistant/{assistantId:[0-9]+}", handler.GetAssistantAwardsByAssistantID).Methods("GET")
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/assistant_awards", handler.CreateAssistantAward).Methods("POST")
	secured.HandleFunc("/assistant_awards/{id:[0-9]+}", handler.UpdateAssistantAward).Methods("PUT")
	secured.HandleFunc("/assistant_awards/{id:[0-9]+}", handler.DeleteAssistantAward).Methods("DELETE")
}

func RegisterNewsRoutes(router *mux.Router, newsHandler *handlers.NewsHandler) {
	router.HandleFunc("/news", newsHandler.GetAllNews).Methods("GET")
	router.HandleFunc("/news/{id:[0-9]+}", newsHandler.GetNewsByID).Methods("GET")

	secured := router.PathPrefix("/").Subrouter()
	secured.Use(middleware.TokenValid)
	secured.HandleFunc("/news", newsHandler.CreateNews).Methods("POST")
	secured.HandleFunc("/news/{id:[0-9]+}", newsHandler.UpdateNews).Methods("PUT")
	secured.HandleFunc("/news/{id:[0-9]+}", newsHandler.DeleteNews).Methods("DELETE")
}
