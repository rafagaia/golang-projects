//setup_test is a reserved golang file name for setting up testing environment.

package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
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
	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	/* Don't need this anymore, as we've implemented Test Repository to do testing without needing database
	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// don't close until current function exits
	defer conn.Close()
	*/

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
