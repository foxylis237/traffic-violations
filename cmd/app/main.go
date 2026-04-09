package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/kvolis/tesgode/internal/adapters/broker"
	"github.com/kvolis/tesgode/internal/adapters/storage"
	"github.com/kvolis/tesgode/internal/config"
	"github.com/kvolis/tesgode/internal/service"
)

func main() {
	// Инициализируем конфиг с дефолтными значениями
	cfg := config.Default()

	// Graceful shutdown — контекст отменяется по Ctrl+C или SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Инициализируем адаптеры
	b := broker.New()
	if err := b.Connect(""); err != nil {
		log.Fatalf("broker connect: %v", err)
	}

	s := storage.New()
	if err := s.Connect(""); err != nil {
		log.Fatalf("storage connect: %v", err)
	}
	defer s.Close()

	// Запускаем сервис
	svc := service.New(b, s, cfg)
	if err := svc.Run(ctx); err != nil {
		log.Fatalf("service error: %v", err)
	}
}
