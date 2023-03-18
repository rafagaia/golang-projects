package main

import (
	"html/template"
	"net/http"
)

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
	parsedTemplate, err := template.ParseFiles("./templates/" + tmplt)
	if err != nil {
		http.Error(write, "bad request", http.StatusBadRequest)
		return err
	}

	// execute the template, passing it data if any
	err = parsedTemplate.Execute(write, data)
	if err != nil {
		return err
	}
	return nil
}
