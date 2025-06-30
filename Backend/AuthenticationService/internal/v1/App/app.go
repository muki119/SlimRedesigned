package App

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"v1/Config"
	"v1/Controllers"
	"v1/Models"
	"v1/Routes"
	"v1/Services"
	"v1/Utils"
)

type appConfig struct {
	httpServer *http.Server
	dbConfig   *Config.PGDatabase
}

type App struct {
	Port       string
	Config     *appConfig
	Db         *pgxpool.Pool
	httpServer *http.Server
}

type appInterface interface {
	ListenAndServe() error
	Init()
}

// NewApp Creates a new App instance
func NewApp(ServerPort string) *App {
	config := &appConfig{
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
	return &App{
		Port:   ServerPort,
		Config: config,
	}

}

// Init Connects to database - does all the basic setup of repositories, services,
// controllers
func (a *App) Init() {

	// open a new connection to db -- needs to be connected
	// make a userRepository -- needs init
	// make a user service with the repository -- just needs instantiation
	//make a controller with the userService -- needs insantiation
	// controller should be passes into initialise routes to make routes specific to the Controllers
	var err error
	a.Db, err = a.Config.dbConfig.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	userRepository := &Models.UserRepository{ // sets up the user repository
		Db: a.Db,
	}
	userServices := &Services.Services{UserRepository: userRepository} // then the user services

	RouteControllers := &Controllers.Controllers{ // and the controllers to be used
		UserServices: userServices,
	}
	userRepository.InitialiseModels() // initialises the models for the database used.

	serverMux := http.NewServeMux()
	serverMux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", Routes.InitialiseRoutes(RouteControllers))) // mounts the routes to the multiplexer

	a.Config.httpServer.Handler = serverMux
}

func (a *App) ListenAndServe() error {

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigint
		err := a.Config.httpServer.Shutdown(context.Background())
		if err != nil {
			slog.Error("Server Error", "error", err.Error())
		}
		close(idleConnsClosed)
	}()

	err := a.Config.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener:
		return err
	}
	<-idleConnsClosed
	return nil
}
