package orchestrator

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Agent struct {
	ID   string
	Addr string
}

type Task struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result,omitempty"`
	Status     string `json:"status"`
}

type Orchestrator struct {
	mu     sync.Mutex
	tasks  map[string]*Task
	agents []Agent
	next   int
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		tasks:  make(map[string]*Task),
		agents: []Agent{},
		next:   0,
	}
}

func (o *Orchestrator) RegisterAgent(agent Agent) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.agents = append(o.agents, agent)
	log.Printf("Registered agent: %s", agent.Addr)
}

func (o *Orchestrator) DispatchTask(task Task) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.agents) == 0 {
		log.Println("No available agents")
		return
	}

	agent := o.agents[o.next]
	o.next = (o.next + 1) % len(o.agents)

	task.Status = "processing"
	o.tasks[task.ID] = &task

	go func() {
		taskJSON, _ := json.Marshal(task)
		resp, err := http.Post("http://"+agent.Addr+"/task", "application/json", bytes.NewBuffer(taskJSON))
		if err != nil {
			log.Printf("Failed to send task to agent %s: %v", agent.Addr, err)
			return
		}
		defer resp.Body.Close()

		var result Task
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Println("Failed to decode response")
			return
		}

		o.mu.Lock()
		o.tasks[result.ID].Result = result.Result
		o.tasks[result.ID].Status = "completed"
		o.mu.Unlock()

		log.Printf("Task %s completed with result: %s", result.ID, result.Result)
	}()
}
