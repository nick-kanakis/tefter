package model

//Account is used
type Account struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}
