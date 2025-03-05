package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Task struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
	Status     string `json:"status"`
}

var (
	tasks = make(map[string]*Task)
	mu    sync.Mutex
)

// Обработчик для получения задачи от клиента
func SubmitTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	tasks[task.ID] = &task
	mu.Unlock()

	log.Printf("Received task: %s", task.Expression)
	w.WriteHeader(http.StatusAccepted)
}

// Обработчик для получения результата
func GetResultHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	mu.Lock()
	task, exists := tasks[taskID]
	mu.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}
