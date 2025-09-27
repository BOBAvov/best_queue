package repository

type User struct {
	Id       int32  `db:"id"`
	Username string `db:"username"`
	Tg_nick  string `db:"tg_nick"`
	GroupID  int    `db:"group_id"`
	Password string `db:"password_hash"`
	IsAdmin  bool   `db:"is_admin"`
	Group    string `db:"-"`
}
