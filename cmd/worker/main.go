package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"dispatch/internal/bootstrap"
	bloodworkers "dispatch/internal/modules/blood/workers"
	notificationworkers "dispatch/internal/modules/notifications/workers"
	"dispatch/internal/shared/constants"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	worker, err := bootstrap.NewWorker(ctx)
	if err != nil {
		log.Fatalf("bootstrap worker: %v", err)
	}

	broadcastConsumer := bloodworkers.NewBroadcastConsumer(
		worker.NotificationService,
		worker.BloodRecipientFinder,
		worker.Bus,
		worker.Logger,
	)

	notificationSendConsumer := notificationworkers.NewSendConsumer(
		worker.NotificationRepo,
		worker.NotificationSender,
		worker.Bus,
		worker.Logger,
	)

	go func() {
		if err := worker.Bus.Subscribe(
			ctx,
			[]string{constants.TopicBloodBroadcastRequested},
			worker.Group,
			broadcastConsumer.Handle,
		); err != nil && ctx.Err() == nil {
			log.Fatalf("subscribe blood broadcast: %v", err)
		}
	}()

	go func() {
		if err := worker.Bus.Subscribe(
			ctx,
			[]string{constants.TopicNotificationSendRequested},
			worker.Group,
			notificationSendConsumer.Handle,
		); err != nil && ctx.Err() == nil {
			log.Fatalf("subscribe notification send: %v", err)
		}
	}()

	if err := worker.Run(ctx); err != nil {
		log.Fatalf("run worker: %v", err)
	}
}
