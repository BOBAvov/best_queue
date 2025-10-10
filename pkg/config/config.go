// Package config содержит функции для загрузки и управления конфигурацией приложения
package config

import (
	"log"
	"sso/models"

	"github.com/spf13/viper"
)

// LoadConfig загружает конфигурацию приложения из файла config.yml
func LoadConfig() models.Config {
	// Инициализируем конфигурацию
	if err := initConfig(); err != nil {
		log.Fatalf("init config err: %v", err.Error())
	}

	// Создаем структуру конфигурации из значений viper
	cfg := models.Config{
		Port: viper.GetString("port"), // Порт HTTP сервера
		DB: models.DBConfig{
			Host:     viper.GetString("db.host"),     // Хост базы данных
			Port:     viper.GetString("db.port"),     // Порт базы данных
			Username: viper.GetString("db.username"), // Имя пользователя БД
			Password: viper.GetString("db.password"), // Пароль пользователя БД
			DBName:   viper.GetString("db.dbname"),   // Имя базы данных
			SSLMode:  viper.GetString("db.sslmode"),  // Режим SSL подключения
		},
	}

	log.Println("Config loaded")
	return cfg
}

// initConfig инициализирует viper для чтения конфигурационного файла
func initConfig() error {
	// Устанавливаем путь к конфигурационным файлам
	viper.AddConfigPath("configs")
	// Устанавливаем имя конфигурационного файла (без расширения)
	viper.SetConfigName("config")
	// Читаем конфигурационный файл
	return viper.ReadInConfig()
}
