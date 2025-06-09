package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"v1/Config"
	"v1/Models"
	"v1/Routes"
)

var port string = os.Getenv("PORT")

var httpServer = &http.Server{
	Addr:    fmt.Sprintf(":%s", port),
	Handler: http.DefaultServeMux, // Use default handler
}

func main() {
	databaseError := Config.ConnectToDatabase()
	if databaseError != nil {
		log.Fatalf("Error connecting to database: %v", databaseError)
	} else {
		log.Default().Println("Database Connected Successfully!")
	}
	Models.InitialiseModels()
	Routes.InitialiseRoutes()

	http.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", Routes.AuthRouter))
	log.Println("Starting server on port", port)
	serverError := httpServer.ListenAndServe()
	if serverError != nil {
		log.Fatalf("Error starting server: %v", serverError)
	} else {
		log.Default().Println("Server Started Successfully!\n" + "listening on port: " + port)
	}

}
