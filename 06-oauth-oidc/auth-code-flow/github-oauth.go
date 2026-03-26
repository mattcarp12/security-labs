package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig *oauth2.Config
	// In a real application, this state string should be generated randomly
	// for each request and stored in a session cookie to prevent CSRF attacks.
	// We'll exploit this exact vulnerability in Lab 4!
	oauthStateString = "pseudo-random-state-string"
)

func init() {
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"read:user"}, // Requesting read-only access to the user's profile
		Endpoint:     github.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Log in with GitHub</a></body></html>`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to GitHub's consent page
	url := oauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Verify the state parameter to protect against CSRF
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	// 2. Exchange the authorization code for an Access Token
	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Code exchange failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// 3. Use the Access Token to fetch the user's data from GitHub's API
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Display the raw JSON response to prove it worked
	w.Header().Set("Content-Type", "application/json")
	w.Write(userData)
}
