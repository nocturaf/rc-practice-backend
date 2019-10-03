package auth

import (
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"

	"rc-practice-backend/app/helpers"
	"rc-practice-backend/app/models"
)

// GetUsersHandler returns all Users(Access: SuperAdmin only)
// For auth jwt testing purpose only, ignore the file misplacement
// DELETE LATER
func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := []byte("")

	// Get cookie for authentication
	claims, err := CheckCookie(w, r)
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-CheckCookie: %s\n", err)
		status = http.StatusUnauthorized
		helpers.RenderJSON(w, message, status)
		return
	}

	// Check if user is superadmin,
	if err = claims.IsSuperAdmin(); err != nil {
		status = http.StatusUnauthorized
		message = []byte("Not Super Admin, Access Denied")
		helpers.RenderJSON(w, message, status)
		return
	}

	// Init User List
	users, err := h.GetUsers()
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-GetUsers: %s\n", err)
		status = http.StatusBadRequest
		helpers.RenderJSON(w, message, status)
		return
	}

	message, err = json.Marshal(users)
	if err != nil {
		fmt.Printf("auth_handler-GetUsersHandler-Marshal: %s\n", err)
		status = http.StatusBadRequest
		helpers.RenderJSON(w, message, status)
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
	// Create JSON Message Object
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Login User Success",
	}

	// Read Request Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("auth_handler-Login-ReadAll: %s\n", err.Error())
		message.Status = "Failed"
		message.Message = "Failed to read body"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// Unmarshall request body to Credential Object
	userCred := Credential{}
	err = json.Unmarshal(body, &userCred)
	if err != nil {
		fmt.Printf("auth_handler-Login-Unmarshal: %s\n", err)
		message.Status = "Failed"
		message.Message = "Failed to Unmarshal body to usercred"
		status = http.StatusInternalServerError
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// Get User password from DB
	storedPassword, role, err := h.GetStoredPassword(userCred)
	if err, ok := err.(*pq.Error); ok {
		fmt.Printf("auth_handler-Login-GetStoredPassword: %s\n", err)
		message.Status = "Failed"
		message.ErrorCode = fmt.Sprintf("%s", err.Code)
		message.Message = "User does not exist"	
		status = http.StatusUnauthorized
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}
	userCred.Role = role

	// Compare user password with hashed password stored in db
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(userCred.Password))
	if err != nil {
		fmt.Printf("auth_handler-Login-ComparedHashAndPassword: %s", err)
		message.Status =  "Failed"
		message.Message = "Username or password is wrong"
		status = http.StatusUnauthorized
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// JWT Token Below
	err = userCred.SetCookie(w, r)
	if err != nil {
		return
	}
	
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

// IsSuperAdmin as described
func (claims *Claims) IsSuperAdmin() error {
	if claims.Role != 2{
		return errors.New("Not Super Admin, Access Denied")
	}
	return nil
}