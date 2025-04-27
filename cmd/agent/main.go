package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"distributed_calculator/internal/constants"
	"distributed_calculator/configs"
	"distributed_calculator/internal/logger"
	"distributed_calculator/internal/worker"

	"go.uber.org/zap"
)

func main() {

	opts := logger.DefaultOptions()
	opts.LogDir = "logs/agent"

	log, err := logger.New(opts)
	if err != nil {
		_, printErr := fmt.Fprintf(os.Stderr, constants.ErrFailedInitLogger, err)
		if printErr != nil {

			os.Exit(2)
		}
		os.Exit(1)
	}
	defer func() {
		if syncErr := log.Sync(); syncErr != nil {
			_, printErr := fmt.Fprintf(os.Stderr, constants.ErrFailedSyncLogger, syncErr)
			if printErr != nil {

				os.Exit(2)
			}
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := configs.NewWorkerConfig()
	if err != nil {
		log.Fatal(constants.ErrFailedInitConfig, zap.Error(err))
	}

	agent := worker.New(cfg, log)
	if err := agent.Start(); err != nil {
		log.Fatal(constants.ErrFailedStartAgent, zap.Error(err))
	}

	log.Info(constants.LogAgentStarted)

	<-ctx.Done()

	agent.Stop()
	log.Info(constants.LogAgentStoppedGrace)
}