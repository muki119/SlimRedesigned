package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	config "v1/Config"
	controllers "v1/Controllers"
	"v1/Helpers/token"
	middleware "v1/Middleware"
	models "v1/Models"
	routes "v1/Routes"
	userservices "v1/Services/user_services"
	"v1/Stream"
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
	EventBusConfig     *Stream.EventBusConfig
}
type App struct {
	HttpServer *http.Server
	Db         *pgxpool.Pool
	RedisCache *redis.Client
	Stream     *redis.Client
	EventBus   *Stream.StreamsEventBus
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
		EventBusConfig: &Stream.EventBusConfig{
			ConsumerName:  Utils.MustGetEnv("EVENT_BUS_CONSUMER_NAME"),
			ConsumerGroup: Utils.MustGetEnv("EVENT_BUS_CONSUMER_GROUP"),
			MaxCount:      int64(Utils.MustGetEnvInt("EVENT_BUS_MAX_COUNT")),
			Timeout:       time.Duration(Utils.MustGetEnvInt("EVENT_BUS_TIMEOUT")),
		},
	}
}

func (ac *AppConfig) NewApp() *App {
	// from the app configs , it should make a app instance
	return &App{
		Config: ac,
	}
}

func (app *App) initialise() error {
	// all system structs
	var err error
	app.Db, err = app.Config.UserDatabaseConfig.ConnectToDatabase()
	if err != nil {
		return err
	}
	app.RedisCache, err = app.Config.CacheConfig.Connect()
	if err != nil {
		return err
	}
	app.Stream, err = app.Config.StreamConfig.Connect()
	if err != nil {
		return err
	}
	tokenHelper, err := app.Config.TokenConfig.NewTokenHelper()
	if err != nil {
		return err
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

	app.Config.EventBusConfig.Connection = app.Stream        // adds the stream connection
	app.EventBus = app.Config.EventBusConfig.NewFromConfig() // creates New Event bus instance

	serverMux := http.NewServeMux()
	serverMux.Handle("/user/api/v1/", http.StripPrefix("/user/api/v1", userRoutes.GetUserRoutes()))

	// all listen and serve operations
	return nil
}

func (app *App) Start() error {
	// start app and start a listener thread for graaceful shutdown
	closedChan := make(chan bool)
	go func() {
		// listen for a sigint or something
		closeSignal := make(chan os.Signal, 1)                                     // make a signal
		signal.Notify(closeSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT) // signal is changed when theres a interruption
		<-closeSignal                                                              // wait until change in signal
		if err := app.HttpServer.Shutdown(context.Background()); err != nil {
			slog.Error(err.Error())
		}
		if err := app.EventBus.Close(); err != nil {
			slog.Error(err.Error())
		}
		close(closedChan)

	}()
	if err := app.initialise(); err != nil {
		return err
	}

	err := app.HttpServer.ListenAndServe()
	if err != nil {
		return err
	}
	eventBusErrorChan := app.EventBus.Listen() // starts the event bus for the streams -- only returns error
	select {
	case <-closedChan: // if the application is being closed
		fmt.Println("Application stopped")
		return nil
	case err := <-eventBusErrorChan: // if the eventbus Failed somewhere -- only returns error if it exists , no nulls
		return err
	}

}
