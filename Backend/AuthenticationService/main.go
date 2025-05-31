package main

import (
	"fmt"
	"log"
	"net/http"
	"v1/Config"
	"v1/Models"
	"v1/Routes"
)

const port int = 2556

var httpServer = &http.Server{
	Addr:    fmt.Sprintf(":%d", port),
	Handler: http.DefaultServeMux, // Use default handler
}

func main() {
	Routes.Initial()

	http.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", Routes.AuthRouter))
	databaseError := Config.ConnectToDatabase()
	Models.InitialiseModels()
	if databaseError != nil {
		log.Fatalf("Error connecting to database: %v", databaseError)
	} else {
		log.Default().Println("Database Connected Successfully!")
	}
	log.Println("Starting server on port", port)
	serverError := httpServer.ListenAndServe()
	if serverError != nil {
		log.Fatalf("Error starting server: %v", serverError)
	} else {
		log.Default().Println("Server Started Successfully!")
		log.Default().Println("Listening on port", port)
	}

}
