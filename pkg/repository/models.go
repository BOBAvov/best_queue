package repository

type User struct {
	Id       int32  `db:"id"`
	Tg_name  string `db:"tg_name"`
	Password string `db:"password_hash"`
	IsAdmin  bool   `db:"is_admin"`
	Group    string `db:"group_name"`
}
