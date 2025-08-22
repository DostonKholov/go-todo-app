package main

import (
	"github.com/joho/godotenv"
	"go.mood/internal/database"
	"go.mood/internal/handler"
	"go.mood/internal/server"
	"go.mood/internal/service"
	"go.mood/pkg"
	"log"
	"os"
)

func main() {
	// 1. Загружаем конфиг из config.yaml
	if err := pkg.InitConfig(); err != nil {
		log.Fatal("Ошибка загрузки конфига:", err)
	}

	// 2. Загружаем переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	// Получаем JWT_SECRET из переменных окружения
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Переменная окружения JWT_SECRET не установлена")
	}

	// 3. Подключение к базе данных
	connection := database.NewConnectPostgres()
	if connection == nil {
		log.Fatal("Не удалось установить соединение с базой данных")
	}

	// 4. Создание объекта базы данных
	db := database.NewDatabase(connection)

	// 5. Создание сервисов (бизнес-логики)
	// Теперь передаём db И jwtSecret
	services := service.NewService(db, jwtSecret)

	// 6. Создание обработчиков
	// Теперь передаём db и services
	handler := handler.NewHandler(services, db)

	// 7. Создание и запуск сервера
	app := new(server.Server)
	if err := app.ServerRun(handler.InitRoutes(), "8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

//func main() {
//	// 1. Загружаем конфиг из config.yaml
//	if err := pkg.InitConfig(); err != nil {
//		log.Fatal("Ошибка загрузки конфига:", err)
//	}
//
//	// 2. Загружаем переменные окружения из .env
//	if err := godotenv.Load(); err != nil {
//		log.Fatal("Ошибка загрузки .env файла:", err)
//	}
//
//	// 3. Подключение к базе данных
//	connection := database.NewConnectPostgres()
//
//	// 4. Создание объекта базы данных
//	db := database.NewDatabase(connection)
//
//	// 5. Создание сервисов (бизнес-логики)
//	services := service.NewService(db)
//
//	// 6. Создание обработчиков
//	handler := handler.NewHandler(services, db)
//
//	// 7. Создание и запуск сервера
//	app := new(server.Server)
//	if err := app.ServerRun(handler.InitRoutes(), "8080"); err != nil {
//		log.Fatal("Ошибка запуска сервера:", err)
//	}
//}
