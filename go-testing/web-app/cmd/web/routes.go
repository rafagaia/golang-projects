package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	/*
	*	- multiplexer to direct incoming http requests to appropriate handler functions
	*		- used to register middleware and routes as well as serve static assets
	**/
	mux := chi.NewRouter()

	// register middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(app.addIPToContext)
	mux.Use(app.Session.LoadAndSave) // persist session and load it on every request

	// register routes
	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)
	mux.Get("/user/profile", app.Profile)

	// static assets (things that are not html. Js, css, images...)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
