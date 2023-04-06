package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_app_authenticate(t *testing.T) {
	var tableTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"Valid user", `{"email":"admin@example.com","password":"secret"}`, http.StatusOK},
		{"Empty JSON", "{}", http.StatusBadRequest},
		{"Not JSON", "not a JSON in request body", http.StatusBadRequest},
		{"No email in JSON", `{"password":"secret"}`, http.StatusBadRequest},
		{"No password in JSON", `{"email":"admin@example.com"}`, http.StatusBadRequest},
		{"Empty email", `{"email":"","password":"secret"}`, http.StatusBadRequest},
		{"Empty password", `{"email":"admin@example.com","password":""}`, http.StatusBadRequest},
		{"Empty email and password", `{"email":"","password":""}`, http.StatusBadRequest},
		{"Two JSONs", `{"email":"em@example.com","password":"secret"}, {"email":"second@email.com"}`, http.StatusBadRequest},
		{"Disallowed fields in JSON", `{"email":"em@example.com","password":"aPassword", "field":"disallowed"}`, http.StatusBadRequest},
	}

	for _, e := range tableTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/auth", reader)
		resRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(resRecorder, req)

		if e.expectedStatusCode != resRecorder.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d.",
				e.name, e.expectedStatusCode, resRecorder.Code)
		}
	}
}
