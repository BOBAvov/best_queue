package main

import (
	"log"
	"sso/pkg/handler"
	"sso/pkg/repository"
	"sso/pkg/services"
)

func main() {
	// 1. Инициализация зависимостей (слои)
	repos := repository.NewAuthMemoryRepo()
	services := services.NewAuthService(repos)
	handlers := handler.NewHandler(services)

	// 2. Запуск сервера
	router := handlers.InitRoutes()
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}
