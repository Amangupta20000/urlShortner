// File: url_controller.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CORSMiddleware adds the necessary CORS headers to the response
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the required CORS headers for every request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// If it's a preflight (OPTIONS) request, respond with 204 No Content
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Otherwise, call the next handler
		next.ServeHTTP(w, r)
	})
}

// ShortenURL is a handler that responds with a basic JSON message
func GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	// Handle preflight (OPTIONS) request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for actual requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if data.URL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	// Call createURL to generate shortened URL
	urlObject, err := createURL(data.URL)
	if err != nil {
		http.Error(w, "Failed to create shortened URL", http.StatusInternalServerError)
		log.Println("Error creating URL:", err)
		return
	}

	// Return the created shortened URL object as JSON
	err = json.NewEncoder(w).Encode(urlObject)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}
}

func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	// id := r.URL.Path[len("/redirect/")]
	// fmt.Println("id : ", id)

	params := mux.Vars(r)
	url, err := getURL(params["shorturl"])
	// url, err := getURL(id)

	if err != nil {
		http.Error(w, "Wrong URL", http.StatusNotFound)
	}

	// Check if the original URL is well-formed (starting with http:// or https://)
	if !(strings.HasPrefix(url.OriginalURL, "http://") || strings.HasPrefix(url.OriginalURL, "https://")) {
		// Prepend 'http://' if it doesn't have a scheme
		url.OriginalURL = "http://" + url.OriginalURL
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
