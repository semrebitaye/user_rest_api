package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"user_rest_api/initializer"
	"user_rest_api/models"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
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
	json.NewEncoder(w).Encode(result)

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

// get single user with pk
func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	// get id
	params := mux.Vars(r)
	userId := params["id"]
	userID, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal(err)
	}

	// get the user by the id
	user := &models.User{ID: userID}
	if err := db.Model(user).WherePK().Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// return product
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	// get id of url
	params := mux.Vars(r)
	userId := params["id"]

	// get the user based on the req id and delete it
	product := &models.User{}
	result, err := db.Model(product).Where("id = ?", userId).Delete()
	// or we can use
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// return result
	json.NewEncoder(w).Encode(result)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	// get the id of url
	params := mux.Vars(r)
	userId := params["id"]
	userID, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal(err)
	}

	// get the user and update based of the id of the req
	user := &models.User{ID: userID}
	_ = json.NewDecoder(r.Body).Decode(&user)

	_, err = db.Model(user).WherePK().Set("first_name = ?, last_name = ?, user_name = ?, password = ?", user.FirstName, user.LastName, user.UserName, user.Password).Update()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// connect to db
	db := initializer.ConnectToDb()
	defer db.Close()

	var body struct {
		Username string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	user.UserName = body.Username
	db.Model(&user).First()

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("secret")))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokenString)
}
