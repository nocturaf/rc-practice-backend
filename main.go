package main

import (
	"log"
	"net/http"

	"rc-practice-backend/app/modules/auth"

	"github.com/gorilla/mux"
)

func main() {

	// Init router
	router := mux.NewRouter()

	// App routes
	router.HandleFunc("/api/users/register", auth.Register).Methods("POST")

	// listening
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8000", router))

}
