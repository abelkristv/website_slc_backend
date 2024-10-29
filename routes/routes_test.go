package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abelkristv/slc_website/database"
	"github.com/abelkristv/slc_website/handlers"
	"github.com/abelkristv/slc_website/repositories"
	"github.com/abelkristv/slc_website/routes"
	"github.com/abelkristv/slc_website/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// Initialize database connection for tests
	db, err := database.InitializeDB()
	if err != nil {
		panic(err) // Handle error properly in a real application
	}

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	assistantRepo := repositories.NewAssistantRepository(db)
	assistantService := services.NewAssistantService(assistantRepo)
	assistantHandler := handlers.NewAssistantHandler(assistantService)

	eventRepo := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)

	periodRepo := repositories.NewPeriodRepository(db)
	periodService := services.NewPeriodService(periodRepo)
	periodHandler := handlers.NewPeriodHandler(periodService)

	// Register routes
	routes.RegisterUserRoutes(router, userHandler)
	routes.RegisterAssistantRoutes(router, assistantHandler)
	routes.RegisterEventRoutes(router, eventHandler)
	routes.RegisterPeriodRoutes(router, periodHandler)

	return router
}

func TestUserRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		method     string
		url        string
		body       interface{}
		statusCode int
	}{
		{
			name:       "Login",
			method:     "POST",
			url:        "/login",
			body:       map[string]string{"username": "BL23-2", "password": "hehe"},
			statusCode: http.StatusOK, // Update as necessary based on your login handler
		},
		{
			name:       "Get All Users",
			method:     "GET",
			url:        "/users",
			statusCode: http.StatusOK,
		},
		// {
		// 	name:       "Create User",
		// 	method:     "POST",
		// 	url:        "/users",
		// 	body:       map[string]string{"username": "newuser", "password": "newpassword"},
		// 	statusCode: http.StatusCreated,
		// },
		{
			name:       "Get User By ID",
			method:     "GET",
			url:        "/users/1",
			statusCode: http.StatusOK,
		},
		// {
		// 	name:       "Update User",
		// 	method:     "PUT",
		// 	url:        "/users/1",
		// 	body:       map[string]string{"username": "updateduser"},
		// 	statusCode: http.StatusOK,
		// },
		{
			name:       "Delete User",
			method:     "DELETE",
			url:        "/users/1",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			if tt.body != nil {
				body, err := json.Marshal(tt.body)
				assert.NoError(t, err)
				reqBody = bytes.NewBuffer(body)
			} else {
				reqBody = nil
			}

			req := httptest.NewRequest(tt.method, tt.url, reqBody)
			if tt.method != "POST" && tt.method != "PUT" {
				// For GET and DELETE requests, you may need to set a valid token or mock it
				req.Header.Set("Authorization", "Bearer valid_token") // Adjust according to your middleware requirements
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestAssistantRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
	}{
		{"GetAllAssistants", "GET", "/assistants", http.StatusOK},
		{"GetAssistantByID", "GET", "/assistants/1", http.StatusOK},
		{"CreateAssistant", "POST", "/assistants", http.StatusCreated},
		{"UpdateAssistant", "PUT", "/assistants/1", http.StatusOK},
		{"DeleteAssistant", "DELETE", "/assistants/1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.method != "POST" && tt.method != "PUT" {
				req.Header.Set("Authorization", "Bearer valid_token") // Adjust for authorization
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestEventRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
	}{
		{"GetAllEvents", "GET", "/events", http.StatusOK},
		{"GetEventByID", "GET", "/events/1", http.StatusOK},
		{"CreateEvent", "POST", "/events", http.StatusCreated},
		{"UpdateEvent", "PUT", "/events/1", http.StatusOK},
		{"DeleteEvent", "DELETE", "/events/1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.method != "POST" && tt.method != "PUT" {
				req.Header.Set("Authorization", "Bearer valid_token") // Adjust for authorization
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestPeriodRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
	}{
		{"GetAllPeriods", "GET", "/periods", http.StatusOK},
		{"GetPeriodByID", "GET", "/periods/1", http.StatusOK},
		{"CreatePeriod", "POST", "/periods", http.StatusCreated},
		{"UpdatePeriod", "PUT", "/periods/1", http.StatusOK},
		{"DeletePeriod", "DELETE", "/periods/1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.method != "POST" && tt.method != "PUT" {
				req.Header.Set("Authorization", "Bearer valid_token") // Adjust for authorization
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}
