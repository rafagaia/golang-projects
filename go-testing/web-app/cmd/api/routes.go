package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// register middleware
	mux.Use(middleware.Recoverer)
	// mux.Use(app.enableCORS)

	/*
	* authentication routes:
	* - auth handler
	* - refresh handler
	**/
	mux.Post("/auth", app.authenticate)
	mux.Post("/refresh-token", app.refresh)

	// test handler - for development
	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			Message string `json:"message"`
		}{
			Message: "Greetings, humans!",
		}
		_ = app.writeJSON(w, http.StatusOK, payload)
	})

	// protected routes - must have a valid jwt token to access

	// all nested routes here have "/users" to route url prepended
	mux.Route("/users", func(mux chi.Router) {
		// use auth middleware
		mux.Get("/", app.allUsers)
		mux.Get("/{userID}", app.getUser)
		mux.Delete("/{userID}", app.deleteUser)
		mux.Put("/", app.insertUser)
		mux.Patch("/", app.updateUser)
	})

	return mux
}
