package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_middleware(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false}, // defaut
		{"", "", "", true},  // couldn't get an address when we looked for it
		{"X-Forwarded-For", "192.4.2.1", "", false},
		{"", "", "hello:people", false}, // test for someone trying to spoof the address
	}

	// create a dummy handler that we'll use to check the context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// make sure that the value exists in the context
		val := req.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "not present")
		}

		// make sure we got a srtring back
		ip, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		// create the handler to test
		handlerToTest := app.addIPToContext(nextHandler)

		request := httptest.NewRequest("GET", "http://testing", nil)

		if e.emptyAddr {
			// test case where we call remoteAddr to get it from the request, and we get nothing back
			request.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			// add a header to request before actually execute handler to test
			request.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			// before we fire request to handler to test, set address
			request.RemoteAddr = e.addr
		}

		// call dummy handler to perform test. Requires ResponseWriter
		handlerToTest.ServeHTTP(httptest.NewRecorder(), request)
	}
}

func Test_application_ipFromContext(t *testing.T) {
	// get a context
	ctx := context.Background()

	// put something in the context
	ctx = context.WithValue(ctx, contextUserKey, "a value to context")

	// call the function
	ip := app.ipFromContext(ctx)

	// perform the test
	if !strings.EqualFold("a value to context", ip) {
		t.Error("wrong value returned from context")
	}
}
