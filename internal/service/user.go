package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mood/internal/database"
	"go.mood/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserService — сервис для работы с пользователями.
type UserService struct {
	db        *database.Database
	jwtSecret []byte
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(db *database.Database, jwtSecret string) *UserService {
	return &UserService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// GetAllUsers получает всех пользователей.
func (s *UserService) GetAllUsers() ([]model.User, error) {
	users, err := s.db.UserQueries.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %w", err)
	}
	return users, nil
}

// DeleteUserByIDWithCheck удаляет пользователя, но с проверкой, не является ли он админом.
func (s *UserService) DeleteUserByIDWithCheck(idToDelete, userID int64) error {
	if idToDelete == userID {
		return errors.New("нельзя удалить собственный аккаунт")
	}

	userToDelete, err := s.db.UserQueries.GetUserByID(idToDelete)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя для удаления: %w", err)
	}

	if userToDelete.Role == "admin" {
		return errors.New("нельзя удалить другого администратора")
	}

	if err := s.db.UserQueries.DeleteUserByID(idToDelete); err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}
	return nil
}

// RegisterUser регистрирует нового пользователя.
func (s *UserService) RegisterUser(input *model.NewUser) (*model.User, error) {
	if input.Username == "" || input.Password == "" {
		return nil, errors.New("username и password обязательны")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	user := model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         "user",
	}

	if err := s.db.UserQueries.CreateUser(&user); err != nil {
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	user.PasswordHash = ""
	return &user, nil
}

// LoginUser аутентифицирует пользователя и возвращает JWT.
func (s *UserService) LoginUser(username, password string) (string, error) {
	user, err := s.db.UserQueries.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("неверные учетные данные")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("неверные учетные данные")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("не удалось подписать токен: %w", err)
	}

	return signedToken, nil
}
