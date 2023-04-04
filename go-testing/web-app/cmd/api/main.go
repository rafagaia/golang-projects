/*
* different main package from cmd/web
* This one is run as 'go run ./cmd/api', and the other as 'go run ./cmd/web from root dir.'
* Yet, they do share a common code base.
**/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"webapp/pkg/repository"
	"webapp/pkg/repository/dbrepo"
)

const PORT = 8181

type application struct {
	DSN       string
	DB        repository.DatabaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	var app application
	// app variable Domain, name, default domain, help text:
	flag.StringVar(&app.Domain, "domain", "mygolang.com", "Domain for application, e.g. enterprise.com")
	// DB Postgres Connect:
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "my_golang_jwt_signing_secret", "signing secret")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// close db pool of connections when main function exits.
	defer conn.Close()

	// wrap database connection in repository struct
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	log.Printf("Starting REST API on port: %d...\n", PORT)

	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
