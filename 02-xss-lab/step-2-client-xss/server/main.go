package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Vulnerable Search Endpoint (Reflected XSS)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// 1. Grab the search term from the URL
		query := r.URL.Query().Get("q")

		// 2. Reflect it directly back to the user without any escaping
		// If 'query' contains HTML/JS, the browser will execute it.
		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head><title>Search Results</title></head>
			<body style="font-family: sans-serif; padding: 20px;">
				<h2>Server-Side Search Engine</h2>
				<p>You searched for: <strong>%s</strong></p>
				<a href="http://localhost:3000">Go back to Home</a>
			</body>
			</html>
		`, query)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	fmt.Println("Reflected XSS Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
