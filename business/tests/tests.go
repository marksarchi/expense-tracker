// Package tests contains supporting code for running tests
package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sarchimark/expense-tracker/business/auth"
	"github.com/sarchimark/expense-tracker/business/data/schema"
	"github.com/sarchimark/expense-tracker/business/data/user"
	"github.com/sarchimark/expense-tracker/foundation/database"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// Configuration for running tests.
var (
	dbImage = "postgres:latest"
	dbPort  = "5432"
	dbArgs  = []string{"-e", "POSTGRES_PASSWORD=postgres"}
	AdminID = "5cf37266-3473-4006-984f-9325122678b7"
	UserID  = "45b5fbd3-755f-4379-8f07-a58d4a30fa2f"
)

//DBContainer provides configuration for a container to run.
type DBContainer struct {
	Image string
	Port  string
	Args  []string
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, dbc DBContainer) (*log.Logger, *sqlx.DB, func()) {
	c := startContainer(t, dbc.Image, dbc.Port, dbc.Args...)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	// // Wait for the database to be ready. Wait 100ms longer between each attempt.
	// // Do not try more than 20 times.
	// var pingError error
	// maxAttempts := 20
	// for attempts := 1; attempts <= maxAttempts; attempts++ {
	// 	pingError = db.Ping()
	// 	if pingError == nil {
	// 		break
	// 	}
	// 	time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	// }

	// if pingError != nil {
	// 	dumpContainerLogs(t, c.ID)
	// 	stopContainer(t, c.ID)
	// 	t.Fatalf("Database never ready: %v", pingError)
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		stopContainer(t, c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c.ID)
	}

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	return log, db, teardown
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB   *sqlx.DB
	Log  *log.Logger
	Auth *auth.Auth
	//	KID      string
	Teardown func()

	t *testing.T
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, dbc DBContainer) *Test {
	log, db, teardown := NewUnit(t, dbc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}

	// Create RSA keys to enable authentication in our service.
	// privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // Build an authenticator using this key lookup function to retrieve
	// // the corresponding public key.
	// kidID := "4754d86b-7a6d-4df5-9c65-224741361492"
	// lookup := func(kid string) (*rsa.PublicKey, error) {
	// 	switch kid {
	// 	case kidID:
	// 		return &privateKey.PublicKey, nil
	// 	}
	// 	return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
	// }

	auth, err := auth.New("HS256")
	if err != nil {
		t.Fatal(err)
	}

	test := Test{
		//TraceID:  "00000000-0000-0000-0000-000000000000",
		DB:   db,
		Log:  log,
		Auth: auth,
		//KID:      kidID,
		t:        t,
		Teardown: teardown,
	}

	return &test
}

// Token generates an authenticated token for a user.
func (test *Test) Token(email, pass string) string {
	test.t.Log("Generating token for test ...")

	u := user.New(test.Log, test.DB)
	claims, err := u.Authenticate(email, pass)
	if err != nil {
		test.t.Fatal(err)
	}

	token, err := test.Auth.GenerateToken(claims)
	if err != nil {
		test.t.Fatal(err)
	}

	return token
}
