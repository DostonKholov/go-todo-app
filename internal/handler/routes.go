package handler

import (
	"github.com/gorilla/mux"
	"go.mood/internal/database"
	"go.mood/internal/middleware"
	"go.mood/internal/service"
)

type Handlers struct {
	db      *database.Database
	service *service.Service
}

func NewHandler(s *service.Service, db *database.Database) *Handlers {
	return &Handlers{
		service: s,
		db:      db,
	}
}

// InitRoutes — инициализация всех маршрутов (роутов) приложения
func (h *Handlers) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// public: регистрация и логин
	router.HandleFunc("/gistreer", h.RegisterHandler)
	router.HandleFunc("/login", h.LoginHandler)

	// защищённые маршруты — нужно передать токен
	auth := router.PathPrefix("/").Subrouter()
	auth.Use(middleware.AuthMiddleware)

	// tasks
	auth.HandleFunc("/tasks", h.GetAllTasksHandler)                  // получить все задачи текущего пользователя
	auth.HandleFunc("/task", h.CreateTaskHandler)                    // создать новую задачу
	auth.HandleFunc("/tasks/{id}", h.GetTaskHandler)                 // получить конкретную задачу по ID
	auth.HandleFunc("/tasks/{id}", h.UpdateTaskHandler)              // обновить задачу
	auth.HandleFunc("/tasks/{id}", h.DeleteTaskHandler)              // удалить задачу
	auth.HandleFunc("/tasks/{id}/status", h.UpdateTaskStatusHandler) // изменить статус

	// admin-only routes (под /admin)
	admin := auth.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRole("admin"))

	// users
	admin.HandleFunc("/users", h.GetAllUsersHandler)
	admin.HandleFunc("/user/{id}", h.DeleteUserHandler)
	return router
}

//// InitRoutes — инициализация всех маршрутов (роутов) приложения
//func (h *Handlers) InitRoutes() *mux.Router {
//	router := mux.NewRouter()
//
//	// public: регистрация и логин
//	router.HandleFunc("/register", h.RegisterHandler)
//	router.HandleFunc("/login", h.LoginHandler)
//
//	// защищённые маршруты — нужно передать токен
//	auth := router.PathPrefix("/").Subrouter()
//	auth.Use(middleware.AuthMiddleware)
//
//	// tasks
//	auth.HandleFunc("/tasks", h.GetAllTasksHandler)
//	auth.HandleFunc("/task", h.GetTaskHandler)
//	auth.HandleFunc("/create_task", h.CreateTaskHandler)
//	auth.HandleFunc("/task", h.DeleteTaskHandler)
//	auth.HandleFunc("/update_task", h.UpdateTaskHandler)
//	router.HandleFunc("/tasks/{id}/status", h.UpdateTaskStatusHandler)
//
//	// admin-only routes (под /admin)
//	admin := auth.PathPrefix("/admin").Subrouter()
//	admin.Use(middleware.RequireRole("admin"))
//	admin.HandleFunc("/users", h.GetAllUsersHandler)
//	// удаление пользователя по пути /admin/delete_user/{id}
//	admin.HandleFunc("/delete_user/{id}", h.DeleteUserHandler)
//
//	return router
//}
