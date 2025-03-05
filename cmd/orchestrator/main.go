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

type Task struct {
	ID      string   `json:"id"`
	Expr    string   `json:"expr"`
	Status  string   `json:"status"`
	Result  float64  `json:"result"`
	Workers []string `json:"workers"`
}

type Orchestrator struct {
	mu      sync.Mutex
	tasks   map[string]*Task
	agents  []string
	taskMux sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		tasks:  make(map[string]*Task),
		agents: []string{"agent1", "agent2", "agent3"},
	}
}

func (o *Orchestrator) createTask(expr string) *Task {
	o.taskMux.Lock()
	defer o.taskMux.Unlock()

	taskID := fmt.Sprintf("%d", time.Now().UnixNano())
	task := &Task{
		ID:     taskID,
		Expr:   expr,
		Status: "pending",
	}

	task.Workers = o.assignWorkers(expr)

	o.tasks[taskID] = task
	return task
}

func (o *Orchestrator) assignWorkers(expr string) []string {
	numWorkers := len(o.agents)
	return o.agents[:numWorkers]
}

func (o *Orchestrator) computeTask(task *Task) {
	o.mu.Lock()
	defer o.mu.Unlock()

	task.Status = "in-progress"
	var wg sync.WaitGroup
	results := make([]float64, len(task.Workers))

	for i, agent := range task.Workers {
		wg.Add(1)
		go func(i int, agent string) {
			defer wg.Done()
			results[i] = math.Pow(2, float64(i))
			log.Printf("Agent %s computed part of task %s", agent, task.ID)
		}(i, agent)
	}

	wg.Wait()

	var totalResult float64
	for _, result := range results {
		totalResult += result
	}

	task.Result = totalResult
	task.Status = "completed"
}

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

	go o.computeTask(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(task)
}

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
