package main

import (
	"net/http"
	"log"
	
	"remote-config/app/modules/auth"
	"github.com/gorilla/mux"
)

func main()  {

	// Init router
	router := mux.NewRouter()

	// App routes
	router.HandleFunc("/api/users/register", auth.Register).Methods("POST")

	
	// listening
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8000", router))
	
}