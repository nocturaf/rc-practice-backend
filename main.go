package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"rc-practice-backend/app/modules/auth"
	"rc-practice-backend/config"

	"github.com/gorilla/mux"
	gorillaHandler "github.com/gorilla/handlers"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	fmt.Println("Look here")

	handler := new(auth.Handler)

	err := config.ConnectDB(handler)
	if err != nil {
		fmt.Printf("main-ConnectDB: %s\n", err)
		return
	}
	defer handler.DB.Close()

	router := mux.NewRouter()
	headers := gorillaHandler.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := gorillaHandler.AllowedOrigins([]string{"*"})

	router.HandleFunc("/api", index).Methods("GET")
	router.HandleFunc("/api/users", handler.GetUsersHandler).Methods("GET")
	router.HandleFunc("/api/users/register", handler.Register).Methods("POST")
	router.HandleFunc("/api/users/login", handler.Login).Methods("POST")

	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(":8080", gorillaHandler.CORS(headers, methods, origins)(router)))
}
