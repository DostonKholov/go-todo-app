package handler

import (
	"errors"
	"go.mood/internal/middleware"
	"go.mood/pkg"
	"net/http"
)

// GetAllUsersHandler — получает всех пользователей.
func (h *Handlers) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodGet); !ok {
		return
	}

	users, err := h.service.UserService.GetAllUsers()
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	pkg.WriteJSONResponse(w, http.StatusOK, users)
}

// DeleteUserHandler — удаляет пользователя по ID, но с проверкой прав.
func (h *Handlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodDelete); !ok {
		return
	}

	idToDelete, err := pkg.GetID(r)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("неверный ID пользователя"))
		return
	}

	currentUserID, err := middleware.GetUserID(r.Context())
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("не удалось получить ID пользователя"))
		return
	}

	if err := h.service.UserService.DeleteUserByIDWithCheck(idToDelete, currentUserID); err != nil {
		pkg.WriteJSONResponse(w, http.StatusForbidden, err)
		return
	}

	pkg.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Пользователь успешно удален"})
}
