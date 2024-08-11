package routes

import (
	"github.com/gorilla/mux"
	"worldwide-coders/controllers"
)

func RegisterContestRoutes(router *mux.Router) {
	router.HandleFunc("/contests/create", controllers.CreateContest).Methods("POST")
	router.HandleFunc("/contests/get", controllers.GetContest).Methods("GET")
	router.HandleFunc("/contests/register/{contestId}", controllers.ContestRegister).Methods("POST")
	router.HandleFunc("/contests/get/registrations/{contestId}", controllers.GetAllRegistrations).Methods("GET")
	router.HandleFunc("/contests/check/registrations/{contestId}", controllers.CheckRegistration).Methods("GET")
	router.HandleFunc("/contests/leaderboard", controllers.GetLeaderboard).Methods("GET")
}
