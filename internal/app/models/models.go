package models

import (
	"time"
)

type ExpressionStatus string

const (
	StatusPending ExpressionStatus = "PENDING"
	StatusProgress ExpressionStatus = "IN_PROGRESS"
	StatusComplete ExpressionStatus = "COMPLETE"
	StatusError ExpressionStatus = "ERROR"
)

type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression,omitempty"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
	CreatedAt  time.Time        `json:"-"`
	UpdatedAt  time.Time        `json:"-"`
	Error      string           `json:"error,omitempty"`
}

type Task struct {
	ID               string
	ExpressionID     string
	Operation        string
	Arg1             float64
	Arg2             float64
	Result           *float64 // nil
	CreatedAt        time.Time
	DependsOnTaskIDs []string
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID string `json:"id"`
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type ExpressionResponse struct {
	Expression Expression `json:"expression"`
}

type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

type TaskResponse struct {
	Task Task `json:"task"`
}
