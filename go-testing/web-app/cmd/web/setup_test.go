//setup_test is a reserved golang file name for setting up testing environment.

package main

import (
	"log"
	"os"
	"testing"
	"webapp/pkg/db"
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

	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// don't close until current function exits
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	os.Exit(m.Run())
}
