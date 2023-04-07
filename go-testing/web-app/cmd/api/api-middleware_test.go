package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_app_enableCORS(t *testing.T) {
	// dummy handler:
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	var tableTests = []struct {
		name         string
		method       string
		expectHeader bool
	}{
		{"preflight", "OPTIONS", true},
		{"get", "GET", false},
	}

	for _, e := range tableTests {
		handlerToTest := app.enableCORS(nextHandler)

		req := httptest.NewRequest(e.method, "http://foo-testing", nil)
		rr := httptest.NewRecorder()

		handlerToTest.ServeHTTP(rr, req)

		if e.expectHeader && rr.Header().Get("Access-Control-Allow-Credentials") == "" {
			t.Errorf("%s: expected header, but didn't find it.", e.name)
		}

		if !e.expectHeader && rr.Header().Get("Access-Control-Allow-Credentials") != "" {
			t.Errorf("%s: expected no header, but got one.", e.name)
		}
	}
}

func Test_app_authRequired(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPair(&testUser)

	var tableTests = []struct {
		name             string
		token            string
		expectAuthorized bool
		setHeader        bool
	}{
		{name: "valid-token", token: fmt.Sprintf("Bearer %s", tokens.Token), expectAuthorized: true, setHeader: true},
		{name: "invalid-token", token: fmt.Sprintf("Bearer %s", expiredToken), expectAuthorized: false, setHeader: true},
		{name: "no-token", token: "", expectAuthorized: false, setHeader: false},
	}

	for _, e := range tableTests {
		req, _ := http.NewRequest("GET", "/", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}
		rr := httptest.NewRecorder()

		handlerToTest := app.authRequired(nextHandler)
		handlerToTest.ServeHTTP(rr, req)

		if e.expectAuthorized && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got unexpected code 401;", e.name)
		}

		if !e.expectAuthorized && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not get expected code 401;", e.name)
		}
	}
}
