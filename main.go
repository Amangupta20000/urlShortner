package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Amangupta20000/urlShortner/router"
)

func main() {
	fmt.Println("Project URL Shortner starting...")
	// Get the port from the environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Fallback to 8080 if not set
	}

	fmt.Println("MongoDb")
	r := router.Router()
	fmt.Printf("Server is getting Started on PORT %s...\n", port)

	// Listen and serve on the specified port
	log.Fatal(http.ListenAndServe(":"+port, r))
}
