package routes

import (
	"worldwide-coders/controllers"

	"github.com/gorilla/mux"
)

var RegisterAuthRoutes = func(router *mux.Router) {
	router.HandleFunc("/create", controllers.Create).Methods("POST")
}
