package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// dbPool is the database connection pool used for the tests that require database interaction
var dbPool *pgxpool.Pool
var ctx, cancel = context.WithCancel(context.Background())

func TestMain(m *testing.M) {
	dbConnStr := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to parse database connection string.\n")
	}

	for i := 1; i <= 5; i++ {
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil && dbPool != nil {
			break
		}
		log.Printf("Failed to initialise database connection pool")
		log.Printf(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if dbPool == nil {
		log.Fatalf("Failed to initialise database connection pool after 5 attempts")
	}
	defer dbPool.Close()

	// Run the tests
	code := m.Run()

	cancel()
	// Exit after running the tests
	os.Exit(code)
}

func TestUserSignUpFlow(t *testing.T) {
	// Prepare
	ts := httptest.NewServer(handlers.SignUp(dbPool))
	defer ts.Close()

	email := "testperson1@gmail.com"
	firstName := "per1"
	lastName := "son1"
	password := "password1"

	// Execute
	resp, err := ts.Client().PostForm(ts.URL, map[string][]string{
		"email":      {email},
		"first_name": {firstName},
		"last_name":  {lastName},
		"password":   {password},
	})
	if err != nil {
		t.Fatalf("Expected no error when sending POST request to /signup, got %v\n", err)
	}
	defer resp.Body.Close()

	// Verify
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v\n", resp.StatusCode)
	}

	// Tear down
	_, err = dbPool.Exec(ctx, "DELETE FROM users WHERE email = $1", email)
	if err != nil {
		t.Fatalf("Failed to delete user from database, %v\n", err)
	}
}

func TestUserLoginFlow(t *testing.T) {
	// Prepare
	email := "testperson1@gmail.com"
	password := "password1"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Insert test user into the database
	_, err := dbPool.Exec(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", email, string(hashedPassword))
	if err != nil {
		t.Fatalf("Failed to insert test user into database, %v\n", err)
	}

	ts := httptest.NewServer(handlers.Login(dbPool))
	defer ts.Close()

	// Execute
	resp, err := ts.Client().PostForm(ts.URL, map[string][]string{
		"email":    {email},
		"password": {password},
	})
	if err != nil {
		t.Fatalf("Expected no error when sending POST request to /login, got %v\n", err)
	}
	defer resp.Body.Close()

	// Verify
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v\n", resp.StatusCode)
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Expected no error when decoding response, got %v\n", err)
	}

	if _, ok := response["token"]; !ok {
		t.Errorf("Expected token in response, got %v\n", response)
	}

	// Tear down
	_, err = dbPool.Exec(ctx, "DELETE FROM users WHERE email = $1", email)
	if err != nil {
		t.Fatalf("Failed to delete user from database, %v\n", err)
	}
}
