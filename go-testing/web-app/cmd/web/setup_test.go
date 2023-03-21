//setup_test is a reserved golang file name for setting up testing environment.

package main

import (
	"os"
	"testing"
)

// now app variable exists in entire scope of tests, so we don't need to declare in each _test.go
var app application

/*
* reserved golang testing function name
*	when we type $> go test .
*		will always be executed before other test functions run
*			golang tooling will look for setup_test.go, and TestMain function
*
**/
func TestMain(m *testing.M) {
	app.Session = getSession()

	os.Exit(m.Run())
}
