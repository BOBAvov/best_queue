package models

type Config struct {
	Port string
	DB   DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}
