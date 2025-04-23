package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"goapi.railway.app/internal/api"
)

func main() {

	var cfg api.Config

	// Try to read environment variable for port (given by railway). Otherwise use default
	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 4000
	}

	// Set the port to run the API on
	cfg.Port = intPort

	// create the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// create the application
	app := &api.Application{
		Config: cfg,
		Logger: logger,
	}

	// create the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  45 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("server started", "addr", srv.Addr)

	// Start the server
	err = srv.ListenAndServe()
	logger.Error("Error on listen and serve: ", err.Error())
	os.Exit(1)

}
