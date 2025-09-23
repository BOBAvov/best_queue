package main

import (
	"log"
	"sso/pkg/handler"
	"sso/pkg/repository"
	"sso/pkg/services"
)

func main() {
	repos := repository.NewAuthMemoryRepo()
	services := services.NewAuthService(repos)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
