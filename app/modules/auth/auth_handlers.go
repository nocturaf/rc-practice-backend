package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"

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
	status := http.StatusOK
	message := []byte("")

	err := CheckCookie(w, r)
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-CheckCookie: %s\n", err)
		status = http.StatusBadRequest
		return
	}

	// Init User List
	users, err := h.GetUsers()
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-GetUsers: %s\n", err)
		status = http.StatusBadRequest
		return
	}

	message, err = json.Marshal(users)
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-Marshal: %s\n", err)
		status = http.StatusBadRequest
		return
	}

	helpers.RenderJSON(w, message, status)
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var status = http.StatusOK
	var message = JSONMessage{
		Status:  "Success",
		Message: "User Registered",
	}

	// defer helpers.RenderJSON(w, helpers.MarshalJSON(message), status)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("auth_handler-Register-ReadAll: %s\n", err)
		message.Status =  "Failed"
		message.Message = fmt.Sprintf("ioutil.ReadAll request Body: %s", err.Error())
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("auth_handler-Register-Unmarshal: %s\n", err.Error())
		message.Status =  "Failed"
		message.Message = fmt.Sprintf("auth_handler-Register-Unmarshal: %s", err.Error())
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	err = h.InsertUser(user)
	if err, ok := err.(*pq.Error); ok {
		fmt.Printf("auth_handler-Register-InsertUser: %s\n", err.Error())
		message.Status =  "Failed"
		message.ErrorCode = fmt.Sprintf("%s", err.Code)
		message.Message = fmt.Sprintf("%s", err.Error())
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

// Login handles user logging in
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Login User Success",
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("auth_handler-Login-ReadAll: %s\n", err.Error())
		message.Status = "Failed"
		message.Message = "Failed to read body"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	userCred := Credential{}
	err = json.Unmarshal(body, &userCred)
	if err != nil {
		fmt.Printf("auth_handler-Login-Unmarshal: %s\n", err)
		message.Status = "Failed"
		message.Message = "Failed to Unmarshal body to usercred"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	storedPassword, err := h.GetStoredPassword(userCred)
	if err, ok := err.(*pq.Error); ok {
		fmt.Printf("auth_handler-Login-GetStoredPassword: %s\n", err)
		message.Status = "Failed"
		message.ErrorCode = fmt.Sprintf("%s", err.Code)
		message.Message = "User does not exist"	
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(userCred.Password)); err != nil {
		fmt.Printf("auth_handler-Login-ComparedHashAndPassword: %s", err)
		message.Status =  "Failed"
		message.Message = "Username or password is wrong"
		status = http.StatusUnauthorized
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)

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
		message.Status =  "Failed"
		message.Message = "Failed to signed server secret key"
		status = http.StatusInternalServerError
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	
	return
}
