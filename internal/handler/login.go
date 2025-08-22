package handler

import (
	"encoding/json"
	"errors"
	"go.mood/internal/model"
	"go.mood/pkg"
	"net/http"
)

// RegisterHandler — регистрация нового пользователя
func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodPost); !ok {
		return
	}

	var input model.NewUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("Неверный формат данных"))
		return
	}

	// Вся логика перенесена в сервис
	user, err := h.service.UserService.RegisterUser(&input)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusInternalServerError, err)
		return
	}

	pkg.WriteJSONResponse(w, http.StatusCreated, user)
}

// LoginHandler — аутентификация и выдача JWT
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if ok := pkg.AllowMethod(w, r, http.MethodPost); !ok {
		return
	}

	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("неверный формат данных"))
		return
	}

	// Вся логика перенесена в сервис
	token, err := h.service.UserService.LoginUser(in.Username, in.Password)
	if err != nil {
		pkg.WriteJSONResponse(w, http.StatusUnauthorized, err)
		return
	}

	pkg.WriteJSONResponse(w, http.StatusOK, map[string]string{"token": token})
}

//
//// RegisterHandler — регистрация нового пользователя
//func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
//	if ok := pkg.AllowMethod(w, r, http.MethodPost); !ok {
//		return
//	}
//
//	var input model.NewUser
//	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
//		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("Неверный формат данных"))
//		return
//	}
//	if input.Username == "" || input.Password == "" {
//		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("username и password обязательны"))
//		return
//	}
//
//	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
//	if err != nil {
//		pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("ошибка хеширования пароля"))
//		return
//	}
//
//	user := model.User{
//		Username:     input.Username,
//		Email:        input.Email,
//		PasswordHash: string(hash),
//		Role:         "user",
//	}
//
//	if err := h.db.CreateUser(&user); err != nil {
//		pkg.WriteJSONResponse(w, http.StatusInternalServerError, err)
//		return
//	}
//
//	// Не возвращаем password_hash
//	user.PasswordHash = ""
//	pkg.WriteJSONResponse(w, http.StatusCreated, user)
//}
//
//// LoginHandler — аутентификация и выдача JWT
//func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
//	if ok := pkg.AllowMethod(w, r, http.MethodPost); !ok {
//		return
//	}
//
//	var in struct {
//		Username string `json:"username"`
//		Password string `json:"password"`
//	}
//
//	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
//		pkg.WriteJSONResponse(w, http.StatusBadRequest, errors.New("неверный формат данных"))
//		return
//	}
//
//	user, err := h.db.GetUserByUsername(in.Username)
//	if err != nil {
//		pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("неверные учетные данные"))
//		return
//	}
//
//	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
//		pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("неверные учетные данные"))
//		return
//	}
//
//	secret := os.Getenv("JWT_SECRET")
//	if secret == "" {
//		pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("server misconfigured"))
//		return
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"user_id": user.Id,
//		"role":    user.Role,
//		//"exp":     time.Now().Add(24 * time.Hour).Unix(),
//		"exp": time.Now().Add(1 * time.Hour).Unix(),
//	})
//
//	signed, err := token.SignedString([]byte(secret))
//	if err != nil {
//		pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("не удалось подписать токен"))
//		return
//	}
//
//	pkg.WriteJSONResponse(w, http.StatusOK, map[string]string{"token": signed})
//}
