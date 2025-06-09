package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"user_rest_api/initializer"
	"user_rest_api/models"

	"golang.org/x/crypto/bcrypt"
)

// Create user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	var body struct {
		FirstName string
		LastName  string
		UserName  string
		Password  string
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// create user istance
	user := &models.User{FirstName: body.FirstName, LastName: body.LastName, UserName: body.UserName, Password: string(hash)}

	// decoding the request
	_ = json.NewDecoder(r.Body).Decode(&user)

	// inserting into database
	result, err := db.Model(user).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// returning user
	json.NewEncoder(w).Encode(user)

}
