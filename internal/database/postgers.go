package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
)

// NewConnectPostgres устанавливает соединение с PostgreSQL.
func NewConnectPostgres() *sql.DB {
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	user := viper.GetString("db.user")
	password := os.Getenv("DB_PASSWORD")
	dbname := viper.GetString("db.dbname")

	dbParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dbParams)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	log.Println("Connected to Postgres successfully")
	return db
}

//
//func NewConnectPostgres() *sql.DB {
//	host := viper.GetString("db.host")
//	port := viper.GetString("db.port")
//	user := viper.GetString("db.user")
//	password := os.Getenv("DB_PASSWORD")
//	dbname := viper.GetString("db.dbname")
//
//	dbParams := fmt.Sprintf(
//		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
//		host, port, user, password, dbname,
//	)
//
//	db, err := sql.Open("postgres", dbParams)
//	if err != nil {
//		log.Fatalf("Error opening database: %v", err)
//	}
//
//	if err = db.Ping(); err != nil {
//		log.Fatalf("Cannot connect to database: %v", err)
//	}
//
//	log.Println("Connected to Postgres successfully")
//	return db
//}
//
//func (d *Database) GetAllTasks() ([]model.Task, error) {
//	query := `SELECT id, user_id, task, status, created_at, updated_at FROM tasks`
//	rows, err := d.Connection.Query(query)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка получения списка задач: %v", err)
//	}
//	defer rows.Close()
//
//	var tasks []model.Task
//	for rows.Next() {
//		var task model.Task
//		if err := rows.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
//			return nil, fmt.Errorf("ошибка чтения задачи из строки: %v", err)
//		}
//		tasks = append(tasks, task)
//	}
//	return tasks, nil
//}
//
//func (d *Database) GetTasksByUserID(userID int64) ([]model.Task, error) {
//	query := `SELECT id, user_id, task, status, created_at, updated_at FROM tasks WHERE user_id = $1`
//	rows, err := d.Connection.Query(query, userID)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка получения задач пользователя: %v", err)
//	}
//	defer rows.Close()
//
//	var tasks []model.Task
//	for rows.Next() {
//		var task model.Task
//		if err := rows.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
//			return nil, fmt.Errorf("ошибка чтения задачи из строки: %v", err)
//		}
//		tasks = append(tasks, task)
//	}
//	return tasks, nil
//}
//
//func (d *Database) GetTaskByID(id int64) (model.Task, error) {
//	query := `SELECT id, user_id, task, status, created_at, updated_at FROM tasks WHERE id = $1`
//	row := d.Connection.QueryRow(query, id)
//	var task model.Task
//
//	if err := row.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
//		if err == sql.ErrNoRows {
//			return task, fmt.Errorf("задача с id %d не найдена", id)
//		}
//		return task, fmt.Errorf("ошибка получения задачи: %v", err)
//	}
//
//	return task, nil
//}
//
//func (d *Database) CreateTask(task *model.Task) error {
//	query := `INSERT INTO tasks (user_id, task) VALUES($1, $2) RETURNING id`
//
//	row := d.Connection.QueryRow(query, task.UserId, task.Task)
//
//	if err := row.Scan(&task.Id); err != nil {
//		return fmt.Errorf("ошибка создания задачи: %v", err)
//	}
//
//	return nil
//}
//
//// Новые методы: обновление/удаление с проверкой владельца
//func (d *Database) UpdateTaskByIDWithOwner(id int64, task *model.Task, ownerID int64) error {
//	query := `UPDATE tasks SET task = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND user_id = $3`
//
//	res, err := d.Connection.Exec(query, task.Task, id, ownerID)
//	if err != nil {
//		return fmt.Errorf("ошибка обновления: %v", err)
//	}
//	n, _ := res.RowsAffected()
//	if n == 0 {
//		return fmt.Errorf("задача не найдена или вы не являетесь владельцем")
//	}
//	return nil
//}
//
//func (d *Database) DeleteTaskByIDWithOwner(id int64, ownerID int64) error {
//	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
//
//	res, err := d.Connection.Exec(query, id, ownerID)
//	if err != nil {
//		return fmt.Errorf("ошибка удаления: %v", err)
//	}
//	n, _ := res.RowsAffected()
//	if n == 0 {
//		return fmt.Errorf("задача не найдена или вы не являетесь владельцем")
//	}
//	return nil
//}
//
//// Существующий метод, но если хотите — оставьте или используйте с админом
//func (d *Database) DeleteTaskByID(id int64) error {
//	query := `DELETE FROM tasks WHERE id = $1`
//
//	if _, err := d.Connection.Exec(query, id); err != nil {
//		return fmt.Errorf("ошибка удаления: %v", err)
//	}
//	return nil
//}
//
//// GetUserByUsername — новый метод для логина
//func (d *Database) GetUserByUsername(username string) (model.User, error) {
//	query := `SELECT id, user_name, email, password_hash, role, created_at FROM users WHERE user_name = $1`
//	row := d.Connection.QueryRow(query, username)
//
//	var user model.User
//	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
//		if err == sql.ErrNoRows {
//			return user, fmt.Errorf("пользователь не найден")
//		}
//		return user, fmt.Errorf("ошибка получения пользователя: %v", err)
//	}
//	return user, nil
//}
//
//// CreateUser — новый метод для регистрации
//func (d *Database) CreateUser(user *model.User) error {
//	query := `INSERT INTO users (user_name, email, password_hash, role) VALUES($1, $2, $3, $4) RETURNING id`
//	if err := d.Connection.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.Id); err != nil {
//		return fmt.Errorf("ошибка создания пользователя: %v", err)
//	}
//	return nil
//}
//
//func (d *Database) GetAllUsers() ([]model.User, error) {
//	query := `SELECT id, user_name, email, password_hash, role, created_at FROM users`
//	rows, err := d.Connection.Query(query)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка получения списка пользователей: %v", err)
//	}
//	defer rows.Close()
//
//	var users []model.User
//	for rows.Next() {
//		var user model.User
//		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
//			return nil, fmt.Errorf("ошибка : %v", err)
//		}
//		users = append(users, user)
//	}
//	return users, nil
//}
//
//func (d *Database) DeleteUserByID(id int64) error {
//	query := `DELETE FROM users WHERE id = $1`
//	if _, err := d.Connection.Exec(query, id); err != nil {
//		return fmt.Errorf("ошибка удаления: %v", err)
//	}
//	return nil
//}
//
//// UpdateTaskStatus обновляет статус задачи, если она принадлежит пользователю
//func (d *Database) UpdateTaskStatus(taskID int64, userID int64, status bool) error {
//	query := `UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`
//
//	res, err := d.Connection.Exec(query, status, taskID, userID)
//	if err != nil {
//		return err
//	}
//
//	rows, _ := res.RowsAffected()
//	if rows == 0 {
//		return errors.New("задача не найдена или вы не владелец")
//	}
//
//	return nil
//}
//
//// GetUserByID — получает пользователя по его ID.
//func (d *Database) GetUserByID(id int64) (model.User, error) {
//	query := `SELECT id, user_name, email, password_hash, role, created_at FROM users WHERE id = $1`
//	row := d.Connection.QueryRow(query, id)
//
//	var user model.User
//	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
//		if err == sql.ErrNoRows {
//			return user, fmt.Errorf("пользователь с id %d не найден", id)
//		}
//		return user, fmt.Errorf("ошибка получения пользователя: %v", err)
//	}
//	return user, nil
//}
//
////package database
////
////import (
////	"database/sql"
////	"fmt"
////	_ "github.com/lib/pq"
////	"github.com/spf13/viper"
////	"go.mood/internal/model"
////	"log"
////	"os"
////)
////
////func NewConnectPostgres() *sql.DB {
////	host := viper.GetString("db.host")
////	port := viper.GetString("db.port")
////	user := viper.GetString("db.user")
////	password := os.Getenv("DB_PASSWORD")
////	dbname := viper.GetString("db.dbname")
////
////	dbParams := fmt.Sprintf(
////		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
////		host, port, user, password, dbname,
////	)
////
////	db, err := sql.Open("postgres", dbParams)
////	if err != nil {
////		log.Fatalf("Error opening database: %v", err)
////	}
////
////	if err = db.Ping(); err != nil {
////		log.Fatalf("Cannot connect to database: %v", err)
////	}
////
////	log.Println("Connected to Postgres successfully")
////	return db
////}
////
////func (d *Database) GetAllTasks() ([]model.Task, error) {
////	query := `SELECT id, user_id, task, status, created_at, updated_at FROM tasks`
////	rows, err := d.Connection.Query(query)
////	if err != nil {
////		return nil, fmt.Errorf("ошибка получения списка задач: %v", err)
////	}
////	defer rows.Close()
////
////	var tasks []model.Task
////	for rows.Next() {
////		var task model.Task
////		if err := rows.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
////			return nil, fmt.Errorf("ошибка чтения задачи из строки: %v", err)
////		}
////		tasks = append(tasks, task)
////	}
////	return tasks, nil
////}
////
////func (d *Database) GetTaskByID(id int64) (model.Task, error) {
////	query := `SELECT id, user_id, task, status, created_at, updated_at FROM tasks WHERE id = $1`
////	row := d.Connection.QueryRow(query, id)
////	var task model.Task
////
////	if err := row.Scan(&task.Id, &task.UserId, &task.Task, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
////		if err == sql.ErrNoRows {
////			return task, fmt.Errorf("задача с id %d не найдена", id)
////		}
////		return task, fmt.Errorf("ошибка получения задачи: %v", err)
////	}
////
////	return task, nil
////}
////
////func (d *Database) CreateTask(task *model.Task) error {
////	query := `INSERT INTO tasks (user_id, task, status) VALUES($1, $2, $3) RETURNING id`
////
////	row := d.Connection.QueryRow(query, task.UserId, task.Task, task.Status)
////
////	if err := row.Scan(&task.Id); err != nil {
////		return fmt.Errorf("ошибка создания задачи: %v", err)
////	}
////
////	return nil
////}
////
////func (d *Database) DeleteTaskByID(id int64) error {
////	query := `DELETE FROM tasks WHERE id = $1`
////
////	if _, err := d.Connection.Exec(query, id); err != nil {
////		return fmt.Errorf("ошибка удаления: %v", err)
////	}
////	return nil
////}
////
////func (d *Database) UpdateTaskByID(id int64, task *model.Task) error {
////	query := `UPDATE tasks SET task = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
////
////	if _, err := d.Connection.Exec(query, task.Task, id); err != nil {
////		return fmt.Errorf("ошибка обнавления: %v", err)
////	}
////	return nil
////}
////
////func (d *Database) StatusTaskByID(id int64) {
////
////}
////
////func (d *Database) GetAllUsers() ([]model.User, error) {
////	query := `SELECT id, user_name, email, password_hash, role, created_at FROM users`
////	rows, err := d.Connection.Query(query)
////	if err != nil {
////		return nil, fmt.Errorf("ошибка получения списка пользвателей: %v", err)
////	}
////	defer rows.Close()
////
////	var users []model.User
////	for rows.Next() {
////		var user model.User
////		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreateTime); err != nil {
////			return nil, fmt.Errorf("ошибка : %v", err)
////		}
////		users = append(users, user)
////	}
////	return users, nil
////}
////
////func (d *Database) DeleteUserByID(id int64) error {
////	query := `DELETE FROM users WHERE id = $1`
////	if _, err := d.Connection.Exec(query, id); err != nil {
////		return fmt.Errorf("ошибка удаления: %v", err)
////	}
////	return nil
////}
