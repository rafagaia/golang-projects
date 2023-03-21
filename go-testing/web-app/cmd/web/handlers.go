package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

var pathToTemplates = "./templates/"

/*
*	- method with app receiver, as we'll need to share application data with handler
*		- as a handler, it needs something to write to
**/
func (app *application) Home(write http.ResponseWriter, request *http.Request) {
	_ = app.renderPage(write, request, "home.page.gohtml", &TemplateData{})
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
