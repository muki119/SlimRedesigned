package App

import (
	"AuthenticationService/v1/Config"
	"AuthenticationService/v1/Controllers"
	"AuthenticationService/v1/Helpers/Token"
	"AuthenticationService/v1/Middleware"
	"AuthenticationService/v1/Models"
	"AuthenticationService/v1/Routes"
	"AuthenticationService/v1/Services"
	"AuthenticationService/v1/Utils"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5/pgxpool"
)

// App All connections used by the app
type App struct {
	Port       string
	Config     *appConfig
	Db         *pgxpool.Pool
	RedisDb    *redis.Client
	httpServer *http.Server
}

// The main configurations for the App
type appConfig struct {
	httpServer  *http.Server
	dbConfig    *Config.PGDatabase
	redisConfig *Config.RedisBlocklistConfig
	tokenConfig *Token.HelperTokenConfig
}

type appInterface interface {
	ListenAndServe() error
	Init()
	GenerateRoutes(*Controllers.Controllers, *Middleware.Middleware) *Routes.Routes
}

// NewApp Creates a default App instance
func NewApp(ServerPort string) *App {
	config := &appConfig{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf(":%s", ServerPort),
			ReadTimeout:    12 * time.Second,
			WriteTimeout:   12 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		dbConfig: &Config.PGDatabase{
			Host:    Utils.MustGetEnv("DB_HOST"),
			Port:    Utils.MustGetEnv("DB_PORT"),
			User:    Utils.MustGetEnv("DB_USER"),
			Name:    Utils.MustGetEnv("DB_NAME"),
			Timeout: Utils.MustGetEnv("DB_TIMEOUT"),
		},
		redisConfig: &Config.RedisBlocklistConfig{
			Addr:     Utils.MustGetEnv("REDIS_HOST"),
			Password: Utils.MustGetEnv("REDIS_PASSWORD"),
			DB:       Utils.MustGetEnvInt("REDIS_DB"),
		},
		tokenConfig: &Token.HelperTokenConfig{
			SecretKey:  Utils.MustGetEnv("JWT_SECRET_KEY"),
			PrivateKey: Utils.MustGetEnv("JWT_PRIVATE_KEY"),
		},
	}
	return &App{
		Port:   ServerPort,
		Config: config,
	}

}

func (a *App) GenerateRoutes(RouteControllers *Controllers.Controllers, RouteMiddleware *Middleware.Middleware) *Routes.Routes {
	return &Routes.Routes{
		RouteControllers: RouteControllers,
		Middleware:       RouteMiddleware,
	}
}

// Init Connects to database - does all the basic setup of repositories, services,
// controllers
func (a *App) Init() {
	var err error
	a.RedisDb = a.Config.redisConfig.ConnectToDatabase()               // connect to the redis db
	TokenHelpers := a.Config.tokenConfig.CreateTokenService(a.RedisDb) // creates the token services the controllers and middleware use

	a.Db, err = a.Config.dbConfig.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	dtoValidator := validator.New(validator.WithRequiredStructEnabled())
	userRepository := &Models.UserRepository{ // sets up the user repository
		Db: a.Db,
	}
	userServices := &Services.Services{UserRepository: userRepository} // then the user services

	RouteMiddleware := &Middleware.Middleware{ // create the middleware instance the routes will use
		TokenHelper: TokenHelpers,
	}
	RouteControllers := &Controllers.Controllers{ // and the controllers to be used
		UserServices: userServices,
		TokenHelpers: TokenHelpers,
		Validator:    dtoValidator,
	}
	userRoutes := a.GenerateRoutes(RouteControllers, RouteMiddleware)
	userRepository.InitialiseModels() // initialises the models for the database used.

	serverMux := http.NewServeMux()
	serverMux.Handle("/auth/", http.StripPrefix("/auth", RouteMiddleware.Cors(userRoutes.GetAuthRoutes()))) // mounts the routes to the multiplexer

	apiVersionMux := http.NewServeMux()
	apiVersionMux.Handle("/api/v1/", http.StripPrefix("/api/v1", serverMux))

	a.Config.httpServer.Handler = apiVersionMux
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
