package main

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

/*
*	1. takes a JSON payload from client, with username and password
*	2. validates user against database
*	3. if match found, generates and returns JSON with JWT Token
**/
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// read a JSON payload
	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	// look up the user by email address
	user, err := app.DB.GetUserByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	// generate tokens (jwt and refresh)
	tokenPairs, err := app.generateTokenPair(user)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
	}
	// send tokens to user
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)
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
