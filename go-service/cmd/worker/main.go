package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tubexxi/video-api/internal/dependencies"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cont, err := dependencies.NewContainer(ctxTimeout)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cont.Close()

	if cont.AsynqClient == nil {
		log.Fatalf("Asynq client is not initialized")
	}

	workerErrCh := make(chan error, 1)
	go func() {
		workerErrCh <- cont.AsynqClient.StartServer()
	}()

	select {
	case <-ctx.Done():
		stop()
		cont.Logger.Info("⚠️ Worker shutdown signal received")
		cont.AsynqClient.ShutdownServer()
		return
	case err := <-workerErrCh:
		if err != nil {
			cont.Logger.Error("Worker stopped with error", zap.Error(err))
			os.Exit(1)
		}
		cont.Logger.Info("Worker stopped")
	}
}
