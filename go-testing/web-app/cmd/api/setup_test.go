package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "mygolang.com"
	app.JWTSecret = "my_golang_jwt_signing_secret"
	os.Exit(m.Run())
}
