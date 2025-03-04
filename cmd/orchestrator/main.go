package orchestration

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"github.com/gorilla/mux"
)

var (
	ErrEmptyExpression  = errors.New("expression is empty")
	ErrInvalidExpression = errors.New("invalid expression")
	ErrNotCorrInput     = errors.New("incorrect input")
)

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Orchestrator struct {
	mu           sync.Mutex
	expressions  map[string]*Expression
	tasks        []Task
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		expressions: make(map[string]*Expression),
		tasks:       []Task{},
	}
}

func (o *Orchestrator) AddExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	id := strconv.Itoa(rand.Intn(100000))
	o.mu.Lock()
	o.expressions[id] = &Expression{ID: id, Status: "pending"}
	// Разбиение выражения на задачи будет здесь
	o.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (o *Orchestrator) GetExpressions(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	exprs := make([]*Expression, 0, len(o.expressions))
	for _, expr := range o.expressions {
		exprs = append(exprs, expr)
	}
	o.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{ "expressions": exprs })
}

func (o *Orchestrator) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	o.mu.Lock()
	expr, exists := o.expressions[id]
	o.mu.Unlock()

	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{ "expression": expr })
}

func (o *Orchestrator) GetTask(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	if len(o.tasks) == 0 {
		o.mu.Unlock()
		http.Error(w, "No task available", http.StatusNotFound)
		return
	}
	task := o.tasks[0]
	o.tasks = o.tasks[1:]
	o.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{ "task": task })
}

func (o *Orchestrator) SubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	if expr, exists := o.expressions[result.ID]; exists {
		expr.Result = result.Result
		expr.Status = "completed"
	}
	o.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}
