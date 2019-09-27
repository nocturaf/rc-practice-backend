package helpers

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

func RenderJSON(w http.ResponseWriter, data []byte, status int) {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// set HTTP respon type to JSON
	w.Header().Set("Content-Type", "application/json")

	// HTTP status (200 OK, 404 Not Found, 500 Internal Server Error, etc.)
	w.WriteHeader(status)

	// The actual data
	w.Write(data)
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}