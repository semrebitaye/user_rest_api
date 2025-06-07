package main

import (
	"user_rest_api/initializer"
)

func main() {
	initializer.LoadEnvVariable()
	db := initializer.ConnectToDb()
	defer db.Close()

	// user := models.User{ID: 1, FirstName: "semre", LastName: "Bitaye", UserName: "semro", Password: "semreman"}
	// fmt.Println(user)
}
