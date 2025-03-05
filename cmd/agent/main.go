package main

import (
	"github.com/Egorekski/calc_2.0/internal/worker"
	"log"
	"net/http"
	"os"
)

func main() {
	agentAddr := os.Getenv("AGENT_ADDR")
	if agentAddr == "" {
		agentAddr = "localhost:8081"
	}

	http.HandleFunc("/task", worker.HandleTask)

	log.Printf("Agent started at %s\n", agentAddr)
	log.Fatal(http.ListenAndServe(agentAddr, nil))
}
