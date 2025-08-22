package pkg

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Проверка метода
func AllowMethod(w http.ResponseWriter, r *http.Request, allowed string) bool {
	if r.Method != allowed {
		WriteJSONResponse(w, http.StatusMethodNotAllowed, errors.New("Метод не разрешён!"))
		return false
	}
	return true
}

// Универсальный JSON-ответ
func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Формируем базовый ответ
	response := make(map[string]interface{})
	response["code"] = statusCode
	response["success"] = statusCode >= 200 && statusCode < 300

	switch v := payload.(type) {
	case string:
		// Просто сообщение
		response["message"] = v
	case error:
		// Ошибка
		response["error"] = v.Error()
	default:
		// Структура/массив/объект
		response["data"] = v
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"success":false,"error":"Internal JSON error"}`, http.StatusInternalServerError)
	}
}

// Получения id
func GetID(r *http.Request) (int64, error) {
	vars := mux.Vars(r) // получаем переменные пути
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		return 0, errors.New("параметр 'id' отсутствует")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("некорректный параметр 'id'")
	}

	return id, nil
}

//func GetID(r *http.Request) (int64, error) {
//	ids, ok := r.URL.Query()["id"]
//	if !ok || len(ids) == 0 || ids[0] == "" {
//		return 0, errors.New("параметр 'id' отсутствует")
//	}
//
//	id, err := strconv.ParseInt(ids[0], 10, 64)
//	if err != nil {
//		return 0, errors.New("некорректный параметр 'id'")
//	}
//
//	return id, nil
//}
