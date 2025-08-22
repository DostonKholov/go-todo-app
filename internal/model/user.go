package model

import "time"

type User struct {
	Id           int       `gorm:"primary_key" json:"id"`
	Username     string    `json:"user_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `gorm:"default:'user'" json:"role"`
	CreateTime   time.Time `json:"create_time"`
}

type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
