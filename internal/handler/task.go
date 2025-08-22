package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go.mood/internal/middleware"
	"go.mood/internal/model"
	"go.mood/pkg"
	"net/http"
	"time"
)

// Цвета для логов (ANSI escape codes)
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// log — универсальная функция для форматированных логов
func log(level string, color string, format string, a ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, a...)
	fmt.Printf("%s[%s] [%s] %s%s\n", color, timestamp, level, message, colorReset)
}

// Упрощённые вызовы логов
func logInfo(format string, a ...interface{})  { log("INFO", colorGreen, format, a...) }
func logWarn(format string, a ...interface{})  { log("WARN", colorYellow, format, a...) }
func logError(format string, a ...interface{}) { log("ERROR", colorRed, format, a...) }
func logDebug(format string, a ...interface{}) { log("DEBUG", colorCyan, format, a...) }

// GetAllTasksHandler — получает все задачи текущего пользователя
func (h *Handlers) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodGet); !ok {
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	logInfo("Пользователь %d запросил список всех своих задач", userID)

	tasks, err := h.service.TaskService.GetAllTasksByUserID(userID)
	if err != nil {
		logError("Ошибка при получении задач пользователя %d: %v", userID, err)
		pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("Ошибка сервера!"))
		return
	}
	logInfo("Пользователь %d получил %d задач", userID, len(tasks))
	pkg.WriteJSONResponse(w, http.StatusOK, tasks)
}

// GetTaskHandler — получает задачу по ID, только если она принадлежит текущему пользователю
func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodGet); !ok {
		return
	}
	id, err := pkg.GetID(r)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	logInfo("Пользователь %d пытается получить задачу %d", userID, id)

	task, err := h.service.TaskService.GetTaskByID(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logWarn("Задача %d не найдена для пользователя %d", id, userID)
			pkg.WriteJSONResponse(w, http.StatusNotFound, fmt.Errorf("Задача с id %d не найдена", id))
			return
		}
		// Проверяем, если ошибка — это "доступ запрещён", возвращаем 403
		if err.Error() == "доступ запрещён" {
			pkg.WriteJSONResponse(w, http.StatusForbidden, errors.New("доступ запрещён"))
			return
		}
		logError("Ошибка при получении задачи %d: %v", id, err)
		pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("Ошибка сервера!"))
		return
	}
	logInfo("Пользователь %d успешно получил задачу %d", userID, id)
	pkg.WriteJSONResponse(w, http.StatusOK, task)
}

// CreateTaskHandler — создаёт новую задачу для текущего пользователя
func (h *Handlers) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodPost); !ok {
		return
	}
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("Неверный формат данных!"))
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	logInfo("Пользователь %d создаёт задачу: %+v", userID, task)

	if err := h.service.TaskService.CreateTask(&task, userID); err != nil {
		logError("Ошибка создания задачи пользователем %d: %v", userID, err)
		pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("Ошибка создания задачи!"))
		return
	}
	logInfo("Пользователь %d успешно создал задачу ID %d", userID, task.Id)
	pkg.WriteJSONResponse(w, http.StatusCreated, task)
}

// DeleteTaskHandler — удаляет задачу, если она принадлежит текущему пользователю
func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodDelete); !ok {
		return
	}
	id, err := pkg.GetID(r)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	logInfo("Пользователь %d пытается удалить задачу %d", userID, id)

	if err := h.service.TaskService.DeleteTaskByIDWithCheck(id, userID); err != nil {
		logWarn("Пользователь %d не смог удалить задачу %d: %v", userID, id, err)
		pkg.WriteJSONResponse(w, http.StatusForbidden, err)
		return
	}
	logInfo("Пользователь %d успешно удалил задачу %d", userID, id)
	pkg.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Задача успешно удалена"})
}

// UpdateTaskHandler — обновляет задачу, если она принадлежит текущему пользователю
func (h *Handlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodPost); !ok {
		return
	}
	id, err := pkg.GetID(r)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	logInfo("Пользователь %d пытается обновить задачу %d данными %+v", userID, id, task)

	if err := h.service.TaskService.UpdateTaskByIDWithCheck(id, userID, &task); err != nil {
		logWarn("Пользователь %d не смог обновить задачу %d: %v", userID, id, err)
		pkg.WriteJSONResponse(w, http.StatusForbidden, err)
		return
	}
	logInfo("Пользователь %d успешно обновил задачу %d", userID, id)
	pkg.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Задача успешно обновлена"})
}

// UpdateTaskStatusHandler меняет статус задачи (true = сделано, false = не сделано)
func (h *Handlers) UpdateTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodPatch); !ok {
		return
	}
	id, err := pkg.GetID(r)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	var req struct {
		Status bool `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("Неверный формат данных!"))
		return
	}
	userID, _ := middleware.GetUserID(r.Context())
	if err := h.service.TaskService.UpdateTaskStatusByIDWithCheck(id, userID, req.Status); err != nil {
		pkg.WriteJSONResponse(w, http.StatusForbidden, err)
		return
	}
	pkg.WriteJSONResponse(w, http.StatusOK, map[string]any{"id": id, "status": req.Status})
}
