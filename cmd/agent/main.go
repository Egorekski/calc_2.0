package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"time"
)

// Worker представляет собой агента, который обрабатывает задачи
type Worker struct {
	ID string
}

func (w *Worker) processTask(taskID string) {
	// Имитация выполнения задачи
	log.Printf("Worker %s started processing task %s", w.ID, taskID)
	// Здесь можно добавить реальную логику вычислений
	time.Sleep(2 * time.Second)
	log.Printf("Worker %s completed processing task %s", w.ID, taskID)
}

// handleTaskRequest обрабатывает запросы на выполнение задач
func (w *Worker) handleTaskRequest(wr http.ResponseWriter, req *http.Request) {
	// Чтение тела запроса
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(wr, "Failed to read body", http.StatusInternalServerError)
		return
	}
	taskID := string(body) // Для примера предполагаем, что запрос содержит только ID задачи

	// Выполнение задачи
	w.processTask(taskID)

	// Ответ оркестратору
	fmt.Fprintf(wr, "Task %s processed by worker %s", taskID, w.ID)
}

func main() {
	worker := &Worker{ID: "worker1"}

	http.HandleFunc("/task", worker.handleTaskRequest)

	log.Println("Worker started on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
