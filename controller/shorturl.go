// File: url_controller.go
package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ShortenURL is a handler that responds with a basic JSON message
func GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
	}

	// originalURL := r.URL.Query().Get("url")
	if data.URL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}
	urlObject, err := createURL(data.URL)
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
