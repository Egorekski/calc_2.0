package api

import (
	"encoding/json"
	"github.com/Egorekski/calc_2.0/internal/orchestrator"
	"net/http"
)

var orchestratorInstance = orchestrator.NewOrchestrator()

func RegisterAgentHandler(w http.ResponseWriter, r *http.Request) {
	var agent orchestrator.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	orchestratorInstance.RegisterAgent(agent)
	w.WriteHeader(http.StatusOK)
}

func SubmitTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task orchestrator.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	orchestratorInstance.DispatchTask(task)
	w.WriteHeader(http.StatusAccepted)
}
