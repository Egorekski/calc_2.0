package worker

import (
	"fmt"
	"time"

	"distributed_calculator/internal/constants"
	"go.uber.org/zap"
)

// worker  представляет собой горутину вычислений.
func (a *Agent) worker(id int) {
	defer a.wg.Done()

	a.logger.Info("Starting worker", zap.Int(constants.FieldWorkerID, id))

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Info("Worker stopped", zap.Int(constants.FieldWorkerID, id))
			return
		default:
			if err := a.processTask(id); err != nil {
				a.logger.Error("Error processing task", zap.Int(constants.FieldWorkerID, id), zap.Error(err))
			}
		}
	}
}

// processTask обрабатывает одну задачу.
func (a *Agent) processTask(workerID int) error {
	task, err := a.getTask()
	if err != nil {
		return fmt.Errorf(constants.ErrFormatWithWrap, "failed to get task", err)
	}

	if task == nil {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	a.logger.Info("Processing task",
		zap.Int(constants.FieldWorkerID, workerID),
		zap.String(constants.FieldTaskID, task.ID),
		zap.String(constants.FieldOperation, task.Operation))

	var operationTime time.Duration = 100 * time.Millisecond

	switch task.Operation {
	case "+":
		operationTime = 1000 * time.Millisecond
	case "-":
		operationTime = 1000 * time.Millisecond
	case "*":
		operationTime = 1000 * time.Millisecond
	case "/":
		operationTime = 1000 * time.Millisecond
	}

	time.Sleep(operationTime)

	result := a.Calculate(task)

	if err := a.sendResult(task.ID, result); err != nil {
		return fmt.Errorf(constants.ErrFormatWithWrap, constants.LogFailedSendResult, err)
	}

	return nil
}
