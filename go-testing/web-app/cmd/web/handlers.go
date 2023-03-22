package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
	"webapp/pkg/data"
)

var pathToTemplates = "./templates/"

/*
*	- method with app receiver, as we'll need to share application data with handler
*		- as a handler, it needs something to write to
**/
func (app *application) Home(write http.ResponseWriter, req *http.Request) {
	var td = make(map[string]any)

	if app.Session.Exists(req.Context(), "test") {
		msg := app.Session.GetString(req.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(
			req.Context(),
			"test",
			"Hit this page at "+time.Now().UTC().String())
	}
	_ = app.renderPage(write, req, "home.page.gohtml", &TemplateData{Data: td})
}

func (app *application) Profile(write http.ResponseWriter, req *http.Request) {
	_ = app.renderPage(write, req, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP    string
	Data  map[string]any // any = alias to interface
	Error string
	Flash string
	User  data.User
}

func (app *application) renderPage(write http.ResponseWriter, request *http.Request, tmplt string, td *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(
		path.Join(pathToTemplates, tmplt),
		path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(write, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(request.Context())

	td.Error = app.Session.PopString(request.Context(), "error")
	td.Flash = app.Session.PopString(request.Context(), "flash")

	// execute the template, passing it data if any
	err = parsedTemplate.Execute(write, td)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) Login(write http.ResponseWriter, req *http.Request) {
	// parse form data
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(write, "bad request", http.StatusBadRequest)
		return
	}

	// validate data
	form := NewForm(req.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		// redirect to the login page with error message
		app.Session.Put(req.Context(), "error", "Invalid login credentials")
		http.Redirect(write, req, "/", http.StatusSeeOther)
		return
	}

	email := req.Form.Get("email")
	password := req.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		// redirect to the login page with error message
		app.Session.Put(req.Context(), "error", "Invalid login!")
		http.Redirect(write, req, "/", http.StatusSeeOther)
		return
	}

	// got the User. Authenticate the user.
	// if not authenticated then redirect with error
	if !app.authenticate(req, user, password) {
		app.Session.Put(req.Context(), "error", "Invalid login!")
		http.Redirect(write, req, "/", http.StatusSeeOther)
		return
	}
	// prevent fixation attack
	_ = app.Session.RenewToken(req.Context())

	//store success message in session
	app.Session.Put(req.Context(), "flash", "Successfully logged in!")

	// redirect to some other page
	http.Redirect(write, req, "/user/profile", http.StatusSeeOther) // 303
}

func (app *application) authenticate(req *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}

	// user added to Session, that's how we know if logged in or not
	app.Session.Put(req.Context(), "user", user)

	return true
}
