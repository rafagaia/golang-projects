//go:build integration

package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"webapp/pkg/data"
	"webapp/pkg/repository"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// package-level variable
// specify values we need for connection to test database
var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435" //different port from dev database, as this one's for test
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

// variables we'll need for docker when using ory/dockertest package
var resource *dockertest.Resource
var pool *dockertest.Pool

// pool of connections to database:
var testDB *sql.DB
var testRepo repository.DatabaseRepo //interface

func TestMain(m *testing.M) {
	fmt.Println("[starting_test] (users_postgres_test.go) TestMain.")
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to DockerTest; is it running? %s", err)
	}

	pool = p

	// opts(options): set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"}, //internal port for exposed from docker
		PortBindings: map[docker.Port][]docker.PortBinding{ // local machine
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (an instance of docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		if resource != nil {
			log.Println("[Pool_Run_Error] resource: ", resource)
			_ = pool.Purge(resource)
		}
		log.Fatalf("Could not start DockerTest resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("[Pool_Retry_Error] testDB:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		if resource != nil {
			log.Println("[Pool_Retry_Error] resource: ", resource)
			_ = pool.Purge(resource)
		}
		log.Fatalf("[Pool_Retry_Error] Could not connect to database: %s", err)
	}

	// populate the databse with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("Error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}

	// run tests
	code := m.Run()

	// clean up, so that next time tests run, it starts from a new state
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s.", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println("[CreateTables_Errored] ReadFile.")
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println("[CreateTables_Errored] testDB.Exec.")
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("Cant ping database.")
	}
}

func Test_PostgresDBRepo_InsertUser(t *testing.T) {
	testUser := data.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("InsertUser returned an error: %s.", err)
	}

	if id != 1 {
		t.Errorf("InsertUser returned wrong id; expected 1, but got: %d.", id)
	}
}

func Test_PostgresDBRepo_AllUsers(t *testing.T) {
	users, err := testRepo.AllUsers()
	if err != nil {
		t.Errorf("AllUsers returned an error: %s.", err)
	}

	if len(users) != 1 {
		t.Errorf("AllUsers returned wrong size; expected 1, but got: %d.", len(users))
	}

	testUser := data.User{
		FirstName: "Second",
		LastName:  "User",
		Email:     "secondUser@example.com",
		Password:  "secret",
		IsAdmin:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertUser(testUser)

	users, err = testRepo.AllUsers()
	if err != nil {
		t.Errorf("AllUsers returned an error: %s.", err)
	}

	if len(users) != 2 {
		t.Errorf("AllUsers returned wrong size; expected 2, but got: %d.", len(users))
	}
}

func Test_PostgresDBRepo_GetUserById(t *testing.T) {
	user, err := testRepo.GetUser(1)
	if err != nil {
		t.Errorf("Error getting user by id: %s.", err)
	}

	if user.Email != "admin@example.com" {
		t.Errorf("GetUser expected user email admin@example.com; but got: %s", user.Email)
	}

	_, err = testRepo.GetUser(3)
	if err == nil {
		t.Error("GetUserById: no error reported when getting non-existent user by ID.")
	}
}

func Test_PostgresDBRepo_GetUserByEmail(t *testing.T) {
	user, err := testRepo.GetUserByEmail("secondUser@example.com")
	if err != nil {
		t.Errorf("Error getting user by email: %s.", err)
	}

	if user.ID != 2 {
		t.Errorf("GetUser expected user ID 2; but got: %d", user.ID)
	}
}

func Test_PostgresDBRepo_UpdateUser(t *testing.T) {
	user, _ := testRepo.GetUser(2)
	user.FirstName = "Rafa"
	user.LastName = "Gaia"
	user.Email = "gaia@stark.com"

	err := testRepo.UpdateUser(*user)
	if err != nil {
		t.Errorf("UpdateUser with ID %d errored: %s.", 2, err)
	}

	user, _ = testRepo.GetUser(2)
	if user.FirstName != "Rafa" || user.LastName != "Gaia" || user.Email != "gaia@stark.com" {
		t.Errorf(`UpdateUser expected record with
			FirstName: Rafa, LastName: Gaia, Email: gaia@stark.com;
			but got: %s, %s, %s`, user.FirstName, user.LastName, user.Email)
	}
}

func Test_PostgresDBRepo_DeleteUser(t *testing.T) {
	err := testRepo.DeleteUser(2)
	if err != nil {
		t.Errorf("DeleteUser with ID: 2 errored: %s", err)
	}

	user, _ := testRepo.GetUser(2) // or _, err = testRepo.GetUser(2), then check if err == nil
	if user != nil {
		t.Error("DeleteUser expected to have deleted user with ID: 2; but did not.")
	}
}

func Test_PostgresDBRepo_ResetPassword(t *testing.T) {
	err := testRepo.ResetPassword(1, "new_password")
	if err != nil {
		t.Error("Error reseting user(ID=1)'s password: ", err)
	}

	user, _ := testRepo.GetUser(1)
	matches, err := user.PasswordMatches("new_password")
	if err != nil {
		t.Error("Error checking match in reset password: ", err)
	}

	if !matches {
		t.Error("ResetPassword expected to match 'new_password'; but does not")
	}
}

func Test_PostgresDBRepo_InsertUserImage(t *testing.T) {
	var image data.UserImage
	image.UserID = 1
	image.FileName = "test.jpg"
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	newID, err := testRepo.InsertUserImage(image)
	if err != nil {
		t.Error("Inserting user image failed: ", err)
	}

	if newID != 1 {
		t.Errorf("Got wrong id for image; expected 1, but got: %d.", newID)
	}

	image.UserID = 100 // ID that doesn't exist in DB
	_, err = testRepo.InsertUserImage(image)
	if err == nil {
		t.Error("Inserted a user image with a non-existent UserID.")
	}
}
