package worker

import (
	"context"
	"net/http"
	"sync"
	"time"

	"distributed_calculator/internal/constants"
	"distributed_calculator/configs"
	"distributed_calculator/internal/logger"

	"go.uber.org/zap"
)

type Agent struct {
	config     *configs.WorkerConfig
	logger     *logger.Logger
	httpClient *http.Client
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// New создает нового агента.
func New(cfg *configs.WorkerConfig, log *logger.Logger) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	return &Agent{
		config: cfg,
		logger: log,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start запускает агента
func (a *Agent) Start() error {
	a.logger.Info("Starting agent",
		zap.Int(constants.FieldComputingPower, a.config.ComputingPower),
		zap.String(constants.FieldOrchestratorURL, a.config.OrchestratorURL))

	for i := 0; i < a.config.ComputingPower; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}

	return nil
}

// Stop останавливает.
func (a *Agent) Stop() {
	a.cancel()
	a.wg.Wait()
	a.logger.Info("Agent stopped")
}
