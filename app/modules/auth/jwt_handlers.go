package auth

import (
	"net/http"
	"time"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/dgrijalva/jwt-go"
)

// RefreshToken returns a new token with new expiration time
func (c *Claims) RefreshToken(w http.ResponseWriter) error {
	expirationTime := time.Now().Add(10 * time.Minute)
	c.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		fmt.Printf("user_handler-RefreshToken-SignedString: %s\n", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})

	return nil
}

// CheckCookie for user authentication
func CheckCookie(w http.ResponseWriter, r *http.Request) error{
	c, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			fmt.Printf("main-index-errnocookie: %s\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}
		fmt.Printf("main-index-cookie: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	tokenString := c.Value
	claims := &Claims{}

	err = godotenv.Load()
	if err != nil {
		fmt.Printf("main-index-Load: %s\n", err)
		return err
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		if err != jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}
		fmt.Printf("main-index-ParseWithClaims: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if !token.Valid {
		fmt.Printf("main-index-token.Valid: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 1 * time.Minute {
		claims.RefreshToken(w)	
	} 

	return nil
}

// SetCookie method to set cookie with JWT token
func (userCred Credential) SetCookie(w http.ResponseWriter, r *http.Request) error{
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Email: userCred.Email,
		Role: userCred.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		fmt.Printf("auth_handler-Login-SignedString: %s\n", err)
		return err
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	return nil
}
