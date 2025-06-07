package main

import (
	"log"
	"net/http"
	"user_rest_api/controllers"
	"user_rest_api/initializer"

	"github.com/gorilla/mux"
)

func main() {
	initializer.LoadEnvVariable()
	// db := initializer.ConnectToDb()
	// defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/users", controllers.CreateUser).Methods("Post")
	r.HandleFunc("/users", controllers.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", controllers.GetUserById).Methods("GET")

	// user := models.User{ID: 1, FirstName: "semre", LastName: "Bitaye", UserName: "semro", Password: "semreman"}
	// fmt.Println(user)

	log.Fatal(http.ListenAndServe(":3030", r))
}
