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

	// test handler - for development

	// protected routes - must have a valid jwt token to access

	return mux
}
