package auth

import(
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"remote-config/app/models"
	"remote-config/app/helpers"
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
		
		helpers.RenderJSON(res, []byte(`{"message" : "Internal Server Error"}`), 500)

	}
}