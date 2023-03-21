package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
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

type TemplateData struct {
	IP   string
	Data map[string]any // any = alias to interface
}

func (app *application) renderPage(write http.ResponseWriter, request *http.Request, tmplt string, data *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(
		path.Join(pathToTemplates, tmplt),
		path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(write, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(request.Context())

	// execute the template, passing it data if any
	err = parsedTemplate.Execute(write, data)
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
		fmt.Fprint(write, "Failed Login Validation.")
		return
	}

	email := req.Form.Get("email")
	password := req.Form.Get("password")

	log.Println(email, password)

	fmt.Fprint(write, email)
}
