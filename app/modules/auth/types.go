package auth

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
)

// Handler object to handle HTTP api Request
type Handler struct {
	DB *sql.DB
}

// Credential object for user email and password to login
type Credential struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// Claims object from JWT to store fields like expiry time
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}