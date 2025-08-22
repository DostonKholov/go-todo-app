package service

import (
	"go.mood/internal/database"
)

// Service — структура, которая объединяет все сервисы приложения.
type Service struct {
	UserService
	TaskService
}

// NewService создает и инициализирует все сервисы.
func NewService(db *database.Database, jwtSecret string) *Service {
	return &Service{
		UserService: *NewUserService(db, jwtSecret),
		TaskService: *NewTaskService(db),
	}
}
