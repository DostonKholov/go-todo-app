package model

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type TokenResponse struct {
	Token string `json:"token"`
}
