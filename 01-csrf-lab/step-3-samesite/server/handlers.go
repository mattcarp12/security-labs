// server/handlers.go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
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
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful. Cookie set!"))
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	// 1. Get the automatic session cookie
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Printf("🛡️ BLOCKED: Request to /transfer rejected. Reason: No session cookie (SameSite enforcement).")
		http.Error(w, "Unauthorized - No session cookie", http.StatusUnauthorized)
		return
	}

	// 2. NEW: Extract the manual CSRF token from the HTTP Headers
	// clientToken := r.Header.Get("X-CSRF-Token")
	// if clientToken == "" {
	// 	http.Error(w, "Forbidden - Missing CSRF Token", http.StatusForbidden)
	// 	return
	// }

	// 3. Fetch the user's ID AND their assigned CSRF token from the database
	var userID int
	var storedToken string
	err = db.QueryRow(
		"SELECT user_id, csrf_token FROM sessions WHERE id=$1",
		sessionCookie.Value,
	).Scan(&userID, &storedToken)

	if err != nil {
		log.Printf("Hello %d, your session ID is %s", userID, sessionCookie.Value)
		log.Printf("Error fetching session: %v", err)
		http.Error(w, "Unauthorized - Invalid session", http.StatusUnauthorized)
		return
	}

	// 4. NEW: The Critical Check. Does the header token match the database token?
	// if clientToken != storedToken {
	// 	http.Error(w, "Forbidden - Invalid CSRF Token", http.StatusForbidden)
	// 	return
	// }

	// 5. If everything matches, process the transfer
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

// Generates a cryptographically secure 32-byte string
func GenerateCSRF() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Endpoint: GET /csrf-token
func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	// 1. Identify the user via their session cookie
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Generate the new token
	token := GenerateCSRF()

	// 3. Save the token in the database attached to their session
	_, err = db.Exec(
		"UPDATE sessions SET csrf_token=$1 WHERE id=$2",
		token, sessionCookie.Value,
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 4. Send the token back to the frontend as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"csrf": token,
	})
}
