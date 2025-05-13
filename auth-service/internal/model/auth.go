package model

type User struct {
	ID           uint   `json:"id"`
	UserName     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}
