package storage

import (
	"fmt"
	"sync"

	"distributed_calculator/internal/app/models"

	"go.uber.org/zap"
)

type Storage struct {
	expressions sync.Map
	tasks       sync.Map
	taskQueue   []models.Task // Slice to ensure FIFO order
	mu          sync.Mutex
	logger      *zap.Logger
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		taskQueue: make([]models.Task, 0),
		logger:    logger,
	}
}

func (s *Storage) GetTasksByDependency(taskID string) []*models.Task {
	var dependentTasks []*models.Task
	s.tasks.Range(func(_, value interface{}) bool {
		task := value.(*models.Task)
		for _, depID := range task.DependsOnTaskIDs {
			if depID == taskID {
				dependentTasks = append(dependentTasks, task)
				break
			}
		}
		return true
	})
	return dependentTasks
}

func (s *Storage) GetTaskResult(taskID string) (float64, error) {
	if value, ok := s.tasks.Load(taskID); ok {
		task := value.(*models.Task)
		if task.Result == nil {
			return 0, fmt.Errorf("task result not set: %s", taskID)
		}
		return *task.Result, nil
	}
	return 0, fmt.Errorf("task not found") 
}

func (s *Storage) GetTasksByExpressionID(expressionID string) []*models.Task {
	var tasks []*models.Task
	s.tasks.Range(func(_, value interface{}) bool {
		task := value.(*models.Task)
		if task.ExpressionID == expressionID {
			tasks = append(tasks, task)
		}
		return true
	})
	return tasks
}
