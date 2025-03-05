package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

// Task представляет задачу для вычисления арифметического выражения
type Task struct {
	ID      string  `json:"id"`
	Expr    string  `json:"expr"`
	Status  string  `json:"status"`
	Result  float64 `json:"result"`
	Workers []string `json:"workers"`
}

// Оркестратор управляет задачами и агентами
type Orchestrator struct {
	mu      sync.Mutex
	tasks   map[string]*Task
	agents  []string
	taskMux sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		tasks:  make(map[string]*Task),
		agents: []string{"agent1", "agent2", "agent3"}, // Пример агентов
	}
}

// createTask создаёт новую задачу для вычисления выражения
func (o *Orchestrator) createTask(expr string) *Task {
	o.taskMux.Lock()
	defer o.taskMux.Unlock()

	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	task := &Task{
		ID:     taskID,
		Expr:   expr,
		Status: "pending",
	}

	// Разделяем задачу на части для агентов
	task.Workers = o.assignWorkers(expr)

	o.tasks[taskID] = task
	return task
}

// assignWorkers распределяет задачу между агентами
func (o *Orchestrator) assignWorkers(expr string) []string {
	// Простая логика для распределения задач между агентами (например, разбиваем по частям)
	numWorkers := len(o.agents)
	return o.agents[:numWorkers]
}

// computeTask вычисляет задачу, вызывая агентов
func (o *Orchestrator) computeTask(task *Task) {
	o.mu.Lock()
	defer o.mu.Unlock()

	task.Status = "in-progress"
	// Имитируем вычисление, где каждый агент выполняет свою часть
	var wg sync.WaitGroup
	results := make([]float64, len(task.Workers))

	for i, agent := range task.Workers {
		wg.Add(1)
		go func(i int, agent string) {
			defer wg.Done()
			// Имитируем вычисления агента
			results[i] = math.Pow(2, float64(i)) // Пример вычисления (можно заменить на реальное вычисление)
			log.Printf("Agent %s computed part of task %s", agent, task.ID)
		}(i, agent)
	}

	wg.Wait()

	// Суммируем результаты
	var totalResult float64
	for _, result := range results {
		totalResult += result
	}

	task.Result = totalResult
	task.Status = "completed"
}

// handleComputeRequest обрабатывает запрос на вычисление выражения
func (o *Orchestrator) handleComputeRequest(w http.ResponseWriter, r *http.Request) {
	var request map[string]string
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	expr, ok := request["expr"]
	if !ok || expr == "" {
		http.Error(w, "Expression is required", http.StatusBadRequest)
		return
	}

	task := o.createTask(expr)

	// Асинхронно запускаем вычисление задачи
	go o.computeTask(task)

	// Отправляем ответ с ID задачи
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(task)
}

// handleTaskStatusRequest обрабатывает запрос на получение статуса задачи
func (o *Orchestrator) handleTaskStatusRequest(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	o.mu.Lock()
	task, exists := o.tasks[taskID]
	o.mu.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// handleResultRequest обрабатывает запрос на получение результата вычисления
func (o *Orchestrator) handleResultRequest(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	o.mu.Lock()
	task, exists := o.tasks[taskID]
	o.mu.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if task.Status != "completed" {
		http.Error(w, "Task is not completed yet", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func main() {
	orchestrator := NewOrchestrator()

	http.HandleFunc("/compute", orchestrator.handleComputeRequest)
	http.HandleFunc("/status", orchestrator.handleTaskStatusRequest)
	http.HandleFunc("/result", orchestrator.handleResultRequest)

	log.Println("Orchestrator started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
