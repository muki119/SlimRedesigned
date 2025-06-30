package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"v1/Config"
	"v1/Controllers"
	"v1/Models"
	"v1/Routes"
	"v1/Services"
	"v1/Utils"
)

var ServerPort = Utils.Getenv("PORT")
var AppConfig = struct {
	httpServer *http.Server
	dbConfig   *Config.PGDatabase
}{
	httpServer: &http.Server{
		Addr: fmt.Sprintf(":%s", ServerPort),
	},
	dbConfig: &Config.PGDatabase{
		Host: Utils.Getenv("DB_HOST"),
		Port: Utils.Getenv("DB_PORT"),
		User: Utils.Getenv("DB_USER"),
		Name: Utils.Getenv("DB_NAME"),
	},
}

func main() {

	// open a new connection to db -- needs to be connected
	// make a userRepository -- needs init
	// make a user service with the repository -- just needs instantiation
	//make a controller with the userService -- needs insantiation
	// controller should be passes into initialise routes to make routes specific to the Controllers
	databaseConnection, err := AppConfig.dbConfig.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	userRepository := Models.UserRepository{
		Db: databaseConnection,
	}
	userServices := Services.Services{UserRepository: userRepository}
	RouteControllers := Controllers.Controllers{
		UserServices: userServices,
	}
	userRepository.InitialiseModels()
	serverMux := http.NewServeMux()
	serverMux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", Routes.InitialiseRoutes(&RouteControllers)))
	AppConfig.httpServer.Handler = serverMux
	log.Println("Starting server on port", ServerPort)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := AppConfig.httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	err = AppConfig.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	} else {
		log.Default().Println("Server Started Successfully!\n" + "listening on port: " + ServerPort)
	}
	<-idleConnsClosed

}
