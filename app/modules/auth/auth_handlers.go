package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rc-practice-backend/app/helpers"
	"rc-practice-backend/app/models"

	"github.com/lib/pq"
)

func Register(res http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {

		jsonMessage := []byte(`message: "Failed to read body"`)
		helpers.RenderJSON(res, jsonMessage, http.StatusBadRequest)
		return
	}

	// parsing request body
	user := models.NewUser()
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Println(err)
		return
	}

	// insert new user to db
	insert, insertedUserID := InsertUser(user)
	if insert == nil {

		user.ID = insertedUserID
		createdUser, _ := json.Marshal(user)
		helpers.RenderJSON(res, createdUser, http.StatusOK)

	} else {
		pqError := insert.(*pq.Error)
		// log.Println(pqError.Code)
		jsonMessage := []byte(`message: "Unknown Error"`)
		if pqError.Code == "23505" {
			jsonMessage = []byte(`message: "E-mail already exists"`)
		}
		helpers.RenderJSON(res, jsonMessage, 500)

	}
}
