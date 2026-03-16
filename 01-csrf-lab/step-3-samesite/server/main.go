// server/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// Global database connection pool
var db *sql.DB

// Middleware to enable CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Explicitly allow the frontend origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// Explicitly allow cookies to be sent
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Allow our custom CSRF header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

		// Handle preflight requests (the browser asking for permission before sending the real request)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	var err error

	// Grab the database URL from the docker-compose environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Connect to PostgreSQL
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}
	defer db.Close()

	// Verify the connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	fmt.Println("Connected to Database successfully!")

	// Route definitions
	http.HandleFunc("/login", enableCORS(Login))
	http.HandleFunc("/transfer", enableCORS(Transfer))
	http.HandleFunc("/csrf-token", enableCORS(GetCSRFToken))

	// Start the server
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
