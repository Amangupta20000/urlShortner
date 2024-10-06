// File: url_controller.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

// ShortenURL is a handler that responds with a basic JSON message
func GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}
	urlObject, err := createURL(originalURL)
	if err != nil {
		http.Error(w, "Failed to create shortened URL", http.StatusInternalServerError)
		log.Println("Error creating URL:", err)
		return
	}

	// Encode the entire inserted object as JSON in the response
	err = json.NewEncoder(w).Encode(urlObject)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}
}
