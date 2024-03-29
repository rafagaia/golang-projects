package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_app_getTokenFromHeaderAndVerify(t *testing.T) {
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPair(&testUser)

	var tests = []struct {
		name          string
		token         string
		errorExpected bool
		setHeader     bool
		issuer        string
	}{
		{"No-header", "", true, false, app.Domain},
		{"Valid-expired", fmt.Sprintf("Bearer %s", expiredToken), true, true, app.Domain},
		{"Valid-token", fmt.Sprintf("Bearer %s", tokens.Token), false, true, app.Domain},
		{"No-Bearer", fmt.Sprintf("Beeer %s", tokens.Token), true, true, app.Domain},
		{"Invalid-token", fmt.Sprintf("Bearer %s23", tokens.Token), true, true, app.Domain},
		{"Three-header-parts", fmt.Sprintf("Bearer %s thirdPart", tokens.Token), true, true, app.Domain},
		// "Wrong-issuer" case must be last in line to run, as generates new token pair with wrong issuer.
		{"Wrong-issuer", fmt.Sprintf("Bearer %s", tokens.Token), true, true, "wrongdomain.com"},
	}

	for _, e := range tests {
		if e.issuer != app.Domain {
			app.Domain = e.issuer
			// generate new token to simulate token not generated by us (wrong issuer)
			tokens, _ = app.generateTokenPair(&testUser)
		}
		req, _ := http.NewRequest("GET", "/", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}
		// create response recorder:
		rr := httptest.NewRecorder()
		_, _, err := app.getTokenFromHeaderAndVerify(rr, req)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: did not expect error, but got one: %s.", e.name, err.Error())
		}

		if err == nil && e.errorExpected {
			t.Errorf("%s: expected error, but did not get one.", e.name)
		}
		app.Domain = "mygolang.com"
	}
}
