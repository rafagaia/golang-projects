package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application
var expiredToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJhZG1pbiI6dHJ1ZSwiYXVkIjoibXlnb2xhbmcuY29tIiwiZXhwIjoxNjgwNDEyNTIwLCJp
c3MiOiJteWdvbGFuZy5jb20iLCJuYW1lIjoiQXJ5YSBTdGFyayIsInN1YiI6IjEifQ.
vUpTFS1FziDqArpK6Pb6F43YUr9hUqxifUjs8gUiJGU`

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "mygolang.com"
	app.JWTSecret = "my_golang_jwt_signing_secret"
	os.Exit(m.Run())
}
