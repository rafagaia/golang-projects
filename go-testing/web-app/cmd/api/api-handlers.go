package main

import "net/http"

/*
*	1. takes a JSON payload from client, with username and password
*	2. validates user against database
*	3. if match found, generates and returns JSON with JWT Token
**/
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {

}

// takes JSON from client with a valid JWT Token, generates a new one, and returns it back
func (app *application) refresh(w http.ResponseWriter, r *http.Request) {

}

func (app *application) allUsers(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertUser(w http.ResponseWriter, r *http.Request) {

}
