package main

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

func getSession() *scs.SessionManager {
	// create a new session manager
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	// avoid web browser issues
	session.Cookie.SameSite = http.SameSiteLaxMode
	// use encrypted cookies
	session.Cookie.Secure = true

	return session
}
