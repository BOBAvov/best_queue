package config

import (
	"github.com/spf13/viper"
	"log"
	"sso/models"
)

func LoadConfig() models.Config {
	if err := initConfig(); err != nil {
		log.Fatalf("init config err: %v", err.Error())
	}
	cfg := models.Config{
		Port: viper.GetString("port"),
		DB: models.DBConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
	}
	log.Println("Config loaded")
	return cfg
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
