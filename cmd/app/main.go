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
	cfg := config.Default()

	// Create a context that is canceled on SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	b := broker.New()
	if err := b.Connect(""); err != nil {
		log.Fatalf("broker connect: %v", err)
	}

	s := storage.New()
	if err := s.Connect(""); err != nil {
		log.Fatalf("storage connect: %v", err)
	}
	defer s.Close()

	svc := service.New(b, s, cfg)
	if err := svc.Run(ctx); err != nil {
		log.Fatalf("service error: %v", err)
	}
}
