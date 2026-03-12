// server/handlers.go
package main

import (
	"net/http"

	"github.com/google/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Parse the form data sent by the frontend
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var userID int

	// 1. Verify credentials against the database
	err := db.QueryRow(
		"SELECT id FROM users WHERE username=$1 AND password=$2",
		username, password,
	).Scan(&userID)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// 2. Generate a random Session ID (The "Hand Stamp")
	sessionID := uuid.New().String()

	// 3. Store the session in the database
	_, err = db.Exec(
		"INSERT INTO sessions(id, user_id) VALUES($1, $2)",
		sessionID, userID,
	)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// 4. Give the session ID to the user's browser as a Cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: sessionID,
		Path:  "/",
		// Notice what is missing here: No SameSite attributes, making it vulnerable!
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful. Cookie set!"))
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	// 1. Get the automatic session cookie
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized - No session cookie", http.StatusUnauthorized)
		return
	}

	// 2. Fetch the user's ID
	var userID int
	err = db.QueryRow(
		"SELECT user_id FROM sessions WHERE id=$1",
		sessionCookie.Value,
	).Scan(&userID)

	if err != nil {
		http.Error(w, "Unauthorized - Invalid session", http.StatusUnauthorized)
		return
	}

	// 3. If everything matches, process the transfer
	amount := r.FormValue("amount")
	if amount == "" { // Fallback if sent via JSON instead of form data
		amount = r.URL.Query().Get("amount")
	}

	_, err = db.Exec(
		"INSERT INTO transfers(user_id, amount) VALUES($1, $2)",
		userID, amount,
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Secure transfer complete!"))
}
