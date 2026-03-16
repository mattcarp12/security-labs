package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Comment struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

func main() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil && db.Ping() == nil {
			fmt.Println("Connected to Database!")
			break
		}
		time.Sleep(2 * time.Second)
	}
	defer db.Close()

	http.HandleFunc("/login", Login)
	http.HandleFunc("/comments", CommentsHandler)

	fmt.Println("API running on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Simplified Login - Issues a standard session cookie
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	sessionID := uuid.New().String()

	// Vulnerability Setup: Notice HttpOnly is missing!
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true, // <-- THE ARMOR: JS can no longer read this cookie!
	})
	w.Write([]byte("Logged in"))
}

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// THE ARMOR: Content Security Policy
	// This tells the browser: "Only execute scripts from my own domain. Do NOT execute inline scripts (like onerror or <script> tags)."
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline';")
	// Note: We are allowing 'unsafe-inline' temporarily just so our index.html scripts still work, but ideally those would be moved to a separate .js file!

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "POST" {
		var c Comment
		json.NewDecoder(r.Body).Decode(&c)

		// THE ARMOR: Output Escaping
		// Convert < and > into safe HTML entities before saving to the database
		safeAuthor := html.EscapeString(c.Author)
		safeContent := html.EscapeString(c.Content)

		db.Exec("INSERT INTO comments (author, content) VALUES ($1, $2)", safeAuthor, safeContent)
		w.WriteHeader(http.StatusCreated)
		return
	}

	// GET: Fetches all comments exactly as they were stored
	rows, _ := db.Query("SELECT author, content FROM comments ORDER BY id ASC")
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		rows.Scan(&c.Author, &c.Content)
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
