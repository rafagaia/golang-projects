package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/cat", http.StatusNotFound},
	}

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	// when we're done with the test function, close the server:
	defer ts.Close()

	// range through test data
	for _, e := range theTests {
		// call built-in test server, request to specified url
		response, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if response.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, response.StatusCode)
		}
	}
}

func TestApp_Home(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"second visit", "hello, world!", "<small>From Session: hello, world!"},
	}

	for _, e := range tests {
		// create a request
		req, _ := http.NewRequest("GET", "/", nil)
		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context()) // make sure there's nothing in the Session before proceeding

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession) // "test" because its the expected key in the Home() Handler
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(res, req)

		// check status code
		if res.Code != http.StatusOK {
			t.Errorf("TestApp_Home returned wrong status code; expected 200 but got %d.", res.Code)
		}

		body, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: did not find %s in response body.", e.name, e.expectedHTML)
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}

/*
func TestApp_Home_Old(t *testing.T) {
	// create a request
	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(res, req)

	// check status code
	if res.Code != http.StatusOK {
		t.Errorf("TestApp_Home returned wrong status code; expected 200 but got %d.", res.Code)
	}

	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), `<small>From Session:`) {
		t.Error("Did not find correct text in HTML.")
	}
}*/
