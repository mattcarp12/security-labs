package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// App holds our application state and dependencies
type App struct {
	OAuth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	http.HandleFunc("/", app.handleHome)
	http.HandleFunc("/login", app.handleLogin)
	http.HandleFunc("/callback", app.handleCallback)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// initializeApp sets up the OIDC provider and configuration
func initializeApp() (*App, error) {
	// Load the .env file. 
	// If it fails, we just log a message and continue, allowing the app 
	// to fall back to system environment variables (useful for production).
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: No .env file found, relying on system environment variables.")
	}

	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	if domain == "" || clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing Auth0 environment variables")
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC configuration: %w", err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &App{
		OAuth2Config: config,
		Verifier:     provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Log in with Auth0</a></body></html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w)
	url := a.OAuth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *App) handleCallback(w http.ResponseWriter, r *http.Request) {
	if err := verifyState(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := a.OAuth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	idToken, err := a.extractAndVerifyIDToken(ctx, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}

	data, _ := json.MarshalIndent(claims, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (a *App) extractAndVerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field found in the oauth2 token")
	}

	idToken, err := a.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID Token signature: %w", err)
	}

	return idToken, nil
}

func generateStateCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: time.Now().Add(1 * time.Hour)}
	http.SetCookie(w, &cookie)
	return state
}

func verifyState(r *http.Request) error {
	cookie, err := r.Cookie("oauthstate")
	if err != nil {
		return fmt.Errorf("missing state cookie")
	}
	if r.FormValue("state") != cookie.Value {
		return fmt.Errorf("invalid oauth state")
	}
	return nil
}