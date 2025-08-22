package queries

import (
	"database/sql"
	_ "embed"
	"fmt"
	"go.mood/internal/model"
)

//go:embed sql/user/get_by_username.sql
var getUserByUsernameQuery string

//go:embed sql/user/create.sql
var createUserQuery string

//go:embed sql/user/get_all.sql
var getAllUsersQuery string

//go:embed sql/user/delete_by_id.sql
var deleteUserByIDQuery string

//go:embed sql/user/get_by_id.sql
var getUserByIDQuery string

////go:embed sql/task/create.sql
//var createTaskQuery string

// UserQueries содержит методы для работы с пользователями в БД.
type UserQueries struct {
	db *sql.DB
}

// NewUserQueries создает новый экземпляр UserQueries.
func NewUserQueries(db *sql.DB) *UserQueries {
	return &UserQueries{db: db}
}

// CreateUserAndInitialTask создает пользователя и сразу же добавляет ему первую задачу,
// используя транзакцию. Это гарантирует, что либо обе операции будут выполнены,
// либо ни одна из них.
func (q *UserQueries) CreateUserAndInitialTask(user *model.User, task *model.Task) error {
	// Шаг 1: Начинаем транзакцию.
	tx, err := q.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	// Шаг 2: Откладываем откат. Если в коде ниже произойдет return по ошибке
	// и транзакция не будет зафиксирована (Commit), она автоматически откатится.
	defer tx.Rollback()

	// Первая операция в транзакции: создание пользователя.
	// Используем tx.QueryRow вместо q.db.QueryRow.
	if err := tx.QueryRow(createUserQuery, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.Id); err != nil {
		return fmt.Errorf("ошибка создания пользователя в транзакции: %w", err)
	}

	// Вторая операция в транзакции: создание задачи для этого пользователя.
	// Используем tx.Exec вместо q.db.Exec.
	// Мы используем user.Id, который только что получили от базы данных.
	if _, err := tx.Exec(createTaskQuery, user.Id, task.Task); err != nil {
		return fmt.Errorf("ошибка создания задачи в транзакции: %w", err)
	}

	// Шаг 3: Фиксируем транзакцию. Если все операции прошли успешно.
	// После этого все изменения будут сохранены в базе данных.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}

// GetUserByUsername получает пользователя по имени.
func (q *UserQueries) GetUserByUsername(username string) (model.User, error) {
	row := q.db.QueryRow(getUserByUsernameQuery, username)
	var user model.User
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("пользователь не найден")
		}
		return user, fmt.Errorf("ошибка получения пользователя: %v", err)
	}
	return user, nil
}

// CreateUser создает нового пользователя в БД.
func (q *UserQueries) CreateUser(user *model.User) error {
	if err := q.db.QueryRow(createUserQuery, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.Id); err != nil {
		return fmt.Errorf("ошибка создания пользователя: %v", err)
	}
	return nil
}

// GetAllUsers получает всех пользователей.
func (q *UserQueries) GetAllUsers() ([]model.User, error) {
	rows, err := q.db.Query(getAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка пользователей: %v", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
			return nil, fmt.Errorf("ошибка : %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// DeleteUserByID удаляет пользователя по ID.
func (q *UserQueries) DeleteUserByID(id int64) error {
	if _, err := q.db.Exec(deleteUserByIDQuery, id); err != nil {
		return fmt.Errorf("ошибка удаления: %v", err)
	}
	return nil
}

// GetUserByID — получает пользователя по его ID.
func (q *UserQueries) GetUserByID(id int64) (model.User, error) {
	row := q.db.QueryRow(getUserByIDQuery, id)

	var user model.User
	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("пользователь с id %d не найден", id)
		}
		return user, fmt.Errorf("ошибка получения пользователя: %v", err)
	}
	return user, nil
}
