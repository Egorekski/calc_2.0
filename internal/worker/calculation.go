package worker

import (
	"distributed_calculator/internal/constants"
	"distributed_calculator/internal/app/models"

	"go.uber.org/zap"
)

func (a *Agent) Calculate(task *models.Task) float64 {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			a.logger.Error(constants.ErrDivisionByZero,
				zap.String(constants.FieldTaskID, task.ID))
			panic(constants.ErrDivisionByZero)
		}
		return task.Arg1 / task.Arg2
	default:
		a.logger.Error(constants.ErrUnexpectedToken,
			zap.String(constants.FieldTaskID, task.ID),
			zap.String(constants.FieldOperation, task.Operation))
		panic(constants.ErrUnexpectedToken)
	}
}
