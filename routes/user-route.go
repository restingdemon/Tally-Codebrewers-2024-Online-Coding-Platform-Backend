package routes

import (
	controller "worldwide-coders/controllers"

	"github.com/gorilla/mux"
)

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/users/get", controller.GetUserByEmail).Methods("GET")
	router.HandleFunc("/users/update/{email}", controller.UpdateUser).Methods("POST")
}

