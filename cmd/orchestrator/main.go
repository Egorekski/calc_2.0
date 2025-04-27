package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed_calculator/internal/constants"
	"distributed_calculator/configs"
	"distributed_calculator/internal/logger"
	"distributed_calculator/internal/app"

	"go.uber.org/zap"
)

func main() {

	opts := logger.DefaultOptions()
	opts.LogDir = "logs/orchestrator"

	log, err := logger.New(opts)
	if err != nil {
		_, printErr := fmt.Fprintf(os.Stderr, constants.ErrFailedInitLogger, err)
		if printErr != nil {
			_, writeErr := fmt.Fprintln(os.Stderr, "Failed to write to stderr:", printErr)
			if writeErr != nil {
				os.Exit(2)
			}
			os.Exit(2)
		}
		os.Exit(1)
	}
	defer func() {
		if syncErr := log.Sync(); syncErr != nil {
			_, printErr := fmt.Fprintf(os.Stderr, constants.ErrFailedSyncLogger, syncErr)
			if printErr != nil {
				_, writeErr := fmt.Fprintln(os.Stderr, "Failed to write to stderr:", printErr)
				if writeErr != nil {
					os.Exit(2)
				}
				os.Exit(2)
			}
		}
	}()

	cfg, err := configs.NewServerConfig()
	if err != nil {
		log.Fatal(constants.ErrFailedInitConfig, zap.Error(err))
	}

	srv := server.New(cfg, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal(constants.ErrFailedStartServer, zap.Error(err))
		}
	}()

	log.Info(constants.LogOrchestratorStarted)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(constants.ErrServerShutdownFailed, zap.Error(err))
	}

	log.Info(constants.LogOrchestratorStoppedGrace)
}