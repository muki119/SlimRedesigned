package main

import (
	"fmt"
	"log/slog"
	"os"
	"v1/App"
	"v1/Utils"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	AuthenticationServiceApp := App.NewApp(Utils.Getenv("PORT"))
	AuthenticationServiceApp.Init()
	slog.Info(fmt.Sprintf("Starting server on port %s", AuthenticationServiceApp.Port), "Port", AuthenticationServiceApp.Port)
	defer AuthenticationServiceApp.Db.Close()
	err := AuthenticationServiceApp.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err.Error())
	}
	slog.Info("Server Closed")

}
