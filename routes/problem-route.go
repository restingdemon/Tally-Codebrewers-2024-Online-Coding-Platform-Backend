// routes/problem_routes.go
package routes

import (
	"github.com/gorilla/mux"
	"worldwide-coders/controllers"
)

func RegisterProblemRoutes(router *mux.Router) {
	router.HandleFunc("/problems/upload", controllers.CreateProblem).Methods("POST")
	router.HandleFunc("/problems/get", controllers.GetProblems).Methods("GET")
	router.HandleFunc("/problems/getnotvisible", controllers.GetNotVisibleProblems).Methods("GET")
	router.HandleFunc("/problems/update/{pid}", controllers.UpdateProblem).Methods("POST")
}
