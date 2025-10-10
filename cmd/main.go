// Package main содержит точку входа в HTTP SSO API приложение
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

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// main является точкой входа в приложение
func main() {
	// Загружаем конфигурацию приложения
	cfg := config.LoadConfig()
	fmt.Printf("Config: %+v\n", cfg)

	// Инициализируем подключение к базе данных PostgreSQL
	db, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("error initializing database: %s", err.Error())
	}
	fmt.Println("DB connected")

	// Создаем репозиторий для работы с базой данных
	authRepo := repository.NewRepository(db)

	// Инициализируем сервисы с бизнес-логикой
	authService := services.NewAuthService(authRepo)

	// Создаем HTTP обработчики
	handlers := handler.NewHandler(authService)

	// Инициализируем маршруты API
	router := handlers.InitRoutes()

	// Настраиваем канал для обработки сигналов завершения
	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	// Ожидаем сигнал завершения приложения
	<-killChan
	log.Println("Shutting down server...")

	// Закрываем подключение к базе данных
	if err := db.Close(); err != nil {
		log.Fatalf("error closing database: %s", err.Error())
	}

	log.Println("Server stopped gracefully")
}
