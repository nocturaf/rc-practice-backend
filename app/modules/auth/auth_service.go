package auth

//ini importss
// commit rafli
//ss
import (
	"fmt"
	"github.com/dewasa98/rc-practice-backend/app/helpers"
	"github.com/dewasa98/rc-practice-backend/app/models"
	"github.com/dewasa98/rc-practice-backend/config"
)

func HashUserPassword(user *models.User) string {
	hashedPassword, _ := helpers.HashPassword(user.Password)
	return string(hashedPassword)
}

func InsertUser(user *models.User) (error, int) {

	conn, err := config.ConnectDB()
	if err != nil {
		fmt.Println(err)
	}

	query := `INSERT INTO "users"("first_name", "last_name", "email", "password") VALUES($1, $2, $3, $4) RETURNING id`
	statement, err := conn.Prepare(query)
	if err != nil {
		return err, 0
	}
	defer statement.Close()

	var userID int
	err = statement.QueryRow(user.FirstName, user.LastName, user.Email, HashUserPassword(user)).Scan(&userID)
	if err != nil {
		return err, 0
	}

	return nil, userID
}
