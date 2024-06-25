package main

import (
	"context"
	"errors"
	"fmt"
	"httpserver/config"
	"httpserver/internal/database"
	"httpserver/internal/handlers"
	"httpserver/internal/middleware"
	"httpserver/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fileLog, err := openLogFile(cfg.Logging.File)

	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(fileLog, "[info]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	errorLog := log.New(fileLog, "[error]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	infoLog.Println("this is info")
	errorLog.Println("this is error")
	if err := database.InitDB(cfg); err != nil {
		log.Printf("Error initializing database: %v", err)
	} else if database.IsDBConnected() {
		defer database.CloseDB()
		log.Println("Database connection established")
	} else {
		log.Println("No database configuration provided, running without database")
	}

	r := router.NewRouter()

	r.AddRoute("GET", "/", middleware.Logging(handlers.GetRoot))

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
