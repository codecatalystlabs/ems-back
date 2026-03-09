package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"dispatch/internal/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	worker, err := bootstrap.NewWorker(ctx)
	if err != nil {
		log.Fatalf("bootstrap worker: %v", err)
	}

	if err := worker.Run(ctx); err != nil {
		log.Fatalf("run worker: %v", err)
	}
}
