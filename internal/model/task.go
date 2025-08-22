package model

import "time"

type Task struct {
	Id        int       `gorm:"primary_key" json:"id"`
	UserId    int64     `json:"user_id"`
	Task      string    `json:"task"`       // Изменено с int на string
	Status    string    `json:"status"`     // Изменено с int на string
	CreatedAt time.Time `json:"created_at"` // Добавлено для created_at
	UpdatedAt time.Time `json:"updated_at"` // Добавлено для updated_at
}
