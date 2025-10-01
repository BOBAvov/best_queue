package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sso/pkg/config"
	"sso/pkg/handler"
	"sso/pkg/repository"
	"sso/pkg/services"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	fmt.Printf("Config: %+v\n", cfg)
	// Initialize repository
	db, err := repository.NewPostgresDB(cfg.DB)
	fmt.Println("DB connected")
	if err != nil {
		log.Fatalf("error initializing database: %s", err.Error())
	}
	authRepo := repository.NewRepository(db)
	// Initialize services
	authService := services.NewAuthService(authRepo)
	// Initialize handlers
	handlers := handler.NewHandler(authService)
	// Start the server
	router := handlers.InitRoutes()
	// Initialize TLS
	log.Printf("Starting server on %s", cfg.Port)

	kill_chan := make(chan os.Signal, 1)
	signal.Notify(kill_chan, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a separate goroutine
	go func() {
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	// Wait for termination signal
	<-kill_chan
	log.Println("Shutting down server...")
	if err := db.Close(); err != nil {
		log.Fatalf("error closing database: %s", err.Error())
	}

}
