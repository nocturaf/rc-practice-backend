package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"rc-practice-backend/app/helpers"
	"rc-practice-backend/app/models"
)

// GetUsersHandler returns all Users
// For auth jwt testing purpose only, ignore the file misplacement
// DELETE LATER
func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get cookie for authentication
	err := CheckCookie(w, r)
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-CheckCookie: %s\n", err)
		return
	}

	// Init User List 
	users, err := h.GetUsers()
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-GetUsers: %s\n", err)
		return
	}

	userBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-Marshal: %s\n", err)
		return
	}

	helpers.RenderJSON(w, userBytes, http.StatusOK)
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("auth_handler-Register-ReadAll: %s \n", err)
		messageByte := []byte(`message: "Failed to read body`)
		helpers.RenderJSON(w, messageByte, http.StatusBadRequest)
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("auth_handler-Register-Unmarshal: %s \n", err)
		return
	}

	err = h.InsertUser(user)
	if err != nil {
		fmt.Printf("auth_handler-Register-InsertUser: %s\n", err)
		helpers.RenderJSON(w, []byte(`{
			status: "failed",
			message: "Failed to register user"
		}`), http.StatusBadRequest)
		return
	}

	helpers.RenderJSON(w, []byte(`{
		status: "success",
		message: "Insert User Success"
	}`), http.StatusOK)
}

// Login handles user logging in
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("auth_handler-Login-ReadAll: %s\n", err)
		messageByte := []byte(`message: "Failed to read body`)
		helpers.RenderJSON(w, messageByte, http.StatusBadRequest)
		return
	}

	userCred := Credential{}

	err = json.Unmarshal(body, &userCred)
	if err != nil {
		fmt.Printf("auth_handler-Login-Unmarshal: %s\n", err)
		return
	}

	storedPassword, err := h.GetStoredPassword(userCred)
	if err != nil {
		fmt.Printf("auth_handler-Login-GetStoredPassword: %s\n", err)
		helpers.RenderJSON(w, []byte(`{
			status: "Failed",
			message: "no such user in db"
		}`), http.StatusBadRequest)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(userCred.Password)); err != nil {
		fmt.Printf("auth_handler-Login-ComparedHashAndPassword: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// JWT Token Below
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Email: userCred.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		fmt.Printf("auth_handler-Login-SignedString: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	helpers.RenderJSON(w, []byte(`{
		status: "success",
		message: "Login succeed"
	}`), http.StatusOK)
}