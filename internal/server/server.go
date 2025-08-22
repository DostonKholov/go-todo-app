package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	server *http.Server
}

func (s *Server) ServerRun(handlers http.Handler, port string) error {
	// Устанавливаем более гибкие таймауты для сервера
	s.server = &http.Server{
		Addr:              "localhost:" + port,
		Handler:           handlers,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// Создаём канал для получения сигналов ОС (например, Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Запускаем сервер в отдельной горутине, чтобы не блокировать основной поток
	go func() {
		log.Println("Сервер запущен на localhost:" + port)
		// ListenAndServe вернёт ошибку, если сервер был закрыт
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Блокируем выполнение, пока не получим сигнал на завершение
	<-quit
	log.Println("Получен сигнал завершения. Завершаю работу сервера...")

	// Создаём контекст с таймаутом для "плавного" завершения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся плавно завершить работу сервера
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка при плавном завершении работы сервера: %w", err)
	}

	log.Println("Сервер успешно остановлен")
	return nil
}
