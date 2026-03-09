package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	_ "dispatch/docs"
	"dispatch/internal/bootstrap"
)

// @title						Dispatch Backend API
// @version					1.0
// @description				Go modular monolith backend for emergency alert intake, ambulance dispatch, user management, and fleet readiness.
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.NewApp(ctx)
	if err != nil {
		log.Fatalf("bootstrap app: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
