package repository

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}
type Authorization interface {
	CreateUser(user User) (int32, error)
	GetUserByID(id int32) (User, error)
}
