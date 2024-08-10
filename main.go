package main

import (
	"log"
	"net/http"
	"os"
	"worldwide-coders/middleware"
	"worldwide-coders/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	r.Use(middleware.Authenticate)

	routes.RegisterAuthRoutes(r)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	handler := c.Handler(r)
	http.Handle("/", handler)
	// Retrieve the PORT environment variable, default to 9010 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "5112"
	}

	addr := ":" + port
	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
