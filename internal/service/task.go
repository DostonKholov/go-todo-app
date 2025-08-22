package service

import (
	"database/sql"
	"errors"
	"fmt"
	"go.mood/internal/database"
	"go.mood/internal/model"
)

// TaskService — сервис для работы с задачами.
type TaskService struct {
	db *database.Database
}

// NewTaskService создаёт новый экземпляр TaskService.
func NewTaskService(db *database.Database) *TaskService {
	return &TaskService{
		db: db,
	}
}

// GetAllTasksByUserID получает все задачи для конкретного пользователя.
func (s *TaskService) GetAllTasksByUserID(userID int64) ([]model.Task, error) {
	//  Вызов метода из нового объекта TaskQueries
	tasks, err := s.db.TaskQueries.GetTasksByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %w", err)
	}
	return tasks, nil
}

// GetTaskByID получает задачу по ID, проверяя, принадлежит ли она пользователю.
func (s *TaskService) GetTaskByID(taskID, userID int64) (*model.Task, error) {
	// Вызов метода из нового объекта TaskQueries
	task, err := s.db.TaskQueries.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("задача с id %d не найдена", taskID)
		}
		return nil, fmt.Errorf("ошибка при получении задачи: %w", err)
	}
	if task.UserId != userID {
		return nil, errors.New("доступ запрещён")
	}
	return &task, nil
}

// CreateTask создаёт новую задачу для пользователя.
func (s *TaskService) CreateTask(task *model.Task, userID int64) error {
	task.UserId = userID
	if err := s.db.TaskQueries.CreateTask(task); err != nil {
		return fmt.Errorf("ошибка при создании задачи: %w", err)
	}
	return nil
}

// DeleteTaskByIDWithCheck удаляет задачу, только если она принадлежит пользователю.
func (s *TaskService) DeleteTaskByIDWithCheck(taskID, userID int64) error {
	if err := s.db.TaskQueries.DeleteTaskByIDWithOwner(taskID, userID); err != nil {
		return fmt.Errorf("не удалось удалить задачу: %w", err)
	}
	return nil
}

// UpdateTaskByIDWithCheck обновляет задачу, только если она принадлежит пользователю.
func (s *TaskService) UpdateTaskByIDWithCheck(taskID, userID int64, updatedTask *model.Task) error {
	if err := s.db.TaskQueries.UpdateTaskByIDWithOwner(taskID, updatedTask, userID); err != nil {
		return fmt.Errorf("не удалось обновить задачу: %w", err)
	}
	return nil
}

// UpdateTaskStatusByIDWithCheck обновляет статус задачи, только если она принадлежит пользователю.
func (s *TaskService) UpdateTaskStatusByIDWithCheck(taskID, userID int64, status bool) error {
	if err := s.db.TaskQueries.UpdateTaskStatus(taskID, userID, status); err != nil {
		return fmt.Errorf("не удалось обновить статус задачи: %w", err)
	}
	return nil
}
