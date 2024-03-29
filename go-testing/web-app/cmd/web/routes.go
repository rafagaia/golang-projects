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

	/* no longer needs mux.Get("/user/profile", app.Profile):
	* a means to securing access to profile particular page.
	* whenever someone goes to /user/profile, following middleware is applied:
	**/
	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.auth) // middleware that determines whether or not user is logged in.
		mux.Get("/profile", app.Profile)
		mux.Post("/upload-profile-pic", app.UploadProfilePic)
	})

	// static assets (things that are not html. Js, css, images...)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
