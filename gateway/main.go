package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	application := &application{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	application.logger = logger
	server := &http.Server{
		Addr:     ":8080",
		Handler:  *application.getRoutes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	log.Printf("Gateway server start on port: %d", 8080)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Cannot start gateway server. Error: " + err.Error())
	}
}
