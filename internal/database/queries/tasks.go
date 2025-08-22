package queries

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"go.mood/internal/model"
)

//go:embed sql/task/get_all.sql
var getAllTasksQuery string

//go:embed sql/task/get_by_user_id.sql
var getTasksByUserIDQuery string

//go:embed sql/task/get_by_id.sql
var getTaskByIDQuery string

//go:embed sql/task/create.sql
var createTaskQuery string

//go:embed sql/task/update_by_id_with_owner.sql
var updateTaskByIDWithOwnerQuery string

//go:embed sql/task/delete_by_id_with_owner.sql
var deleteTaskByIDWithOwnerQuery string

//go:embed sql/task/update_status.sql
var updateTaskStatusQuery string

// TaskQueries содержит методы для работы с задачами в БД.
type TaskQueries struct {
	db *sql.DB
}

// NewTaskQueries создает новый экземпляр TaskQueries.
func NewTaskQueries(db *sql.DB) *TaskQueries {
	return &TaskQueries{db: db}
}

// CreateTasksInBulk создает список задач в рамках одной транзакции.
// Если любая из операций создания задачи завершается ошибкой, вся транзакция
// будет отменена (откачена), и ни одна задача не будет сохранена в БД.
func (q *TaskQueries) CreateTasksInBulk(tasks []*model.Task) error {
	// 1. Начинаем транзакцию
	tx, err := q.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	// 2. Откладываем откат (Rollback) на случай ошибки.
	// Если транзакция будет зафиксирована (Commit) ниже, этот вызов будет проигнорирован.
	defer tx.Rollback()

	// 3. Выполняем каждую операцию в рамках транзакции
	for _, task := range tasks {
		row := tx.QueryRow(createTaskQuery, task.UserId, task.Task)
		if err := row.Scan(&task.Id); err != nil {
			return fmt.Errorf("ошибка создания задачи в транзакции: %w", err)
		}
	}

	// 4. Фиксируем транзакцию, если все операции успешны
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}

func (q *TaskQueries) GetAllTasks() ([]model.Task, error) {
	rows, err := q.db.Query(getAllTasksQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка задач: %v", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("ошибка чтения задачи из строки: %v", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (q *TaskQueries) GetTasksByUserID(userID int64) ([]model.Task, error) {
	rows, err := q.db.Query(getTasksByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задач пользователя: %v", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("ошибка чтения задачи из строки: %v", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (q *TaskQueries) GetTaskByID(id int64) (model.Task, error) {
	row := q.db.QueryRow(getTaskByIDQuery, id)
	var task model.Task

	if err := row.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("задача с id %d не найдена", id)
		}
		return task, fmt.Errorf("ошибка получения задачи: %v", err)
	}
	return task, nil
}

func (q *TaskQueries) CreateTask(task *model.Task) error {
	row := q.db.QueryRow(createTaskQuery, task.UserId, task.Task)
	if err := row.Scan(&task.Id); err != nil {
		return fmt.Errorf("ошибка создания задачи: %v", err)
	}
	return nil
}

func (q *TaskQueries) UpdateTaskByIDWithOwner(id int64, task *model.Task, ownerID int64) error {
	res, err := q.db.Exec(updateTaskByIDWithOwnerQuery, task.Task, id, ownerID)
	if err != nil {
		return fmt.Errorf("ошибка обновления: %v", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("задача не найдена или вы не являетесь владельцем")
	}
	return nil
}

func (q *TaskQueries) DeleteTaskByIDWithOwner(id int64, ownerID int64) error {
	res, err := q.db.Exec(deleteTaskByIDWithOwnerQuery, id, ownerID)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %v", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("задача не найдена или вы не являетесь владельцем")
	}
	return nil
}

func (q *TaskQueries) UpdateTaskStatus(taskID int64, userID int64, status bool) error {
	res, err := q.db.Exec(updateTaskStatusQuery, status, taskID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("задача не найдена или вы не владелец")
	}
	return nil
}
