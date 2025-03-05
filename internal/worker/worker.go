package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Processing task %s: %s", task.ID, task.Expression)
	time.Sleep(2 * time.Second) // Имитация работы

	task.Result = "4" // Здесь будет парсинг выражения
	json.NewEncoder(w).Encode(task)
}
