package auth

import (
	"fmt"

	"rc-practice-backend/app/models"

	"golang.org/x/crypto/bcrypt"
)

// GetUsers Query all users
func (h *Handler) GetUsers() ([]models.User, error) {
	query := "SELECT id, first_name, last_name, email, role FROM users;"
	rows, err := h.DB.Query(query)
	if err != nil {
		fmt.Printf("user_service-GetUsers-query: %s\n", err)
		return nil, err
	}

	var users []models.User

	for rows.Next(){
		user := models.User{}
		
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
		)
		if err != nil {
			fmt.Printf("user_service-GetUsers-Scan: %s \n",err)
		}
		
		users = append(users, user)
	}

	return users, nil
}

// InsertUser insert user object to database
func (h *Handler) InsertUser(user models.User) (error){

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	query := fmt.Sprintf("insert into users (first_name, last_name, email, password) values ('%s', '%s', '%s', '%s');", user.FirstName, user.LastName, user.Email, hashedPassword)

	_, err = h.DB.Exec(query)
	if err != nil {
		fmt.Printf("user_service-InsertUser-Exec: %s\n", err)
		return err
	}

	return nil
}


// GetStoredPassword returns user password stored in database, error if no such user exist
func (h *Handler) GetStoredPassword(cred Credential)(string, int, error){
	query := fmt.Sprintf("select password, role from users where email='%s';", cred.Email)
	rows := h.DB.QueryRow(query)
	
	var storedPassword string
	var role int
	
	err := rows.Scan(&storedPassword, &role)
	if err != nil {
		fmt.Printf("user_service-GetStoredPassword-Scan: %s \n", err)
		return "", 99, err
	}

	return storedPassword, role, nil
}