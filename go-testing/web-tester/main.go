package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"webapp/pkg/data"
	"webapp/pkg/repository"
	"webapp/pkg/repository/dbrepo"

	"github.com/alexdwards/scs/v2"
)

type application struct {
	DSN		string
	DB		repository.DatabaseRepo
	Session	*scs.SessionManager
}

func main() {
	gob.Register(data.User{})

	// setup an app config
	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone...")
	flag.Parse()

	// ....
}