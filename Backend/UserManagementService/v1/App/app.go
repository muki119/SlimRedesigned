package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
	config "v1/Config"
	controllers "v1/Controllers"
	"v1/Helpers/token"
	middleware "v1/Middleware"
	models "v1/Models"
	routes "v1/Routes"
	userservices "v1/Services/user_services"
	"v1/Utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AppConfig struct {
	Port               string
	UserDatabaseConfig *config.PGDatabase
	CacheConfig        *config.RedisCache
	StreamConfig       *config.RedisStream
	TokenConfig        *token.Config
}
type App struct {
	HttpServer *http.Server
	Db         *pgxpool.Pool
	RedisCache *redis.Client
	Stream     *redis.Client
	Config     *AppConfig
}

func DefaultAppConfig() *AppConfig {
	// the default application configuration
	return &AppConfig{
		Port: Utils.MustGetEnv("PORT"),
		UserDatabaseConfig: &config.PGDatabase{
			Host:     Utils.MustGetEnv("DB_HOST"),
			Port:     Utils.MustGetEnv("DB_PORT"),
			User:     Utils.MustGetEnv("DB_USER"),
			Name:     Utils.MustGetEnv("DB_NAME"),
			Password: Utils.MustGetEnv("DB_PASSWORD"),
		},
		CacheConfig: &config.RedisCache{
			Addr:     Utils.MustGetEnv("CACHE_ADDR"),
			Password: Utils.MustGetEnv("CACHE_PASSWORD"),
			DB:       Utils.MustGetEnvInt("CACHE_DB"),
		},
		StreamConfig: &config.RedisStream{
			Addr:     Utils.MustGetEnv("STREAM_ADDR"),
			Password: Utils.MustGetEnv("STREAM_PASSWORD"),
			DB:       Utils.MustGetEnvInt("STREAM_DB"),
		},
		TokenConfig: &token.Config{
			PublicKeyPath: Utils.MustGetEnv("PUBLIC_KEY_PATH"),
		},
	}
}

func (ac *AppConfig) NewApp() *App {
	// from the app configs , it should make a app instance
	return &App{
		Config: ac,
	}
}

func (app *App) initialise() {
	// all system structs
	var err error
	app.Db, err = app.Config.UserDatabaseConfig.ConnectToDatabase()
	if err != nil {
		// return error or slog
	}
	app.RedisCache, err = app.Config.CacheConfig.Connect()
	if err != nil {

	}
	app.Stream, err = app.Config.StreamConfig.Connect()
	if err != nil {

	}
	tokenHelper, err := app.Config.TokenConfig.NewTokenHelper()
	if err != nil {

	}

	UserRepository := &models.UserRepository{DatabaseConn: app.Db}

	userServices := &userservices.UserServices{UserRepository: UserRepository}
	userControllers := &controllers.UserControllers{UserServices: userServices}
	middlewareStruct := &middleware.Middleware{TokenHelper: tokenHelper}
	userRoutes := &routes.UserRoutes{UserControllers: userControllers, MiddleWare: middlewareStruct}
	app.HttpServer = &http.Server{
		Addr:           fmt.Sprintf(":%s", app.Config.Port),
		ReadTimeout:    12 * time.Second,
		WriteTimeout:   12 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/user/api/v1/", http.StripPrefix("/user/api/v1", userRoutes.GetUserRoutes()))

	// all listen and serve operations
}
func (app *App) Start() {
	// start app and start a listener thread for graaceful shutdown
	closedChan := make(chan bool)
	go func() {
		// listen for a sigint or something
		closeSignal := make(chan os.Signal, 1)            // make a signal
		signal.Notify(closeSignal, os.Interrupt, os.Kill) // signal is changed when theres a interruption
		<-closeSignal                                     // wait until change in signal
		if err := app.HttpServer.Shutdown(context.Background()); err != nil {
			slog.Error(err.Error())
		}
		close(closedChan)

	}()
	app.HttpServer.ListenAndServe()
	<-closedChan
	fmt.Println("Application stopped")
	os.Exit(0)
}
