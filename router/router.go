package router

import (
	"github.com/Amangupta20000/urlShortner/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/shorturl", controller.GenerateShortURL).Methods("POST")
	router.HandleFunc("/redirect/{shorturl}", controller.RedirectURLHandler).Methods("GET")

	return router
}
