package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go.mood/internal/database/queries"
)

// Database — структура, содержащая все наборы запросов.
type Database struct {
	UserQueries *queries.UserQueries
	TaskQueries *queries.TaskQueries
}

// NewDatabase создает новый экземпляр Database.
func NewDatabase(conn *sql.DB) *Database {
	return &Database{
		UserQueries: queries.NewUserQueries(conn),
		TaskQueries: queries.NewTaskQueries(conn),
	}
}
