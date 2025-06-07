package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"user_rest_api/initializer"
	"user_rest_api/models"
)

// Create user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	// create user istance
	user := &models.User{}

	// decoding the request
	_ = json.NewDecoder(r.Body).Decode(&user)

	// inserting into database
	_, err := db.Model(user).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// returning user
	json.NewEncoder(w).Encode(user)

}

// get all users from the database
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	// creating user slice and select all from the database
	var users []models.User
	if err := db.Model(&users).Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// returning product
	json.NewEncoder(w).Encode(users)
}
