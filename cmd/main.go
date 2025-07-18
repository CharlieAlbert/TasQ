package main

import (
	"log"
	"time"

	"github.com/CharlieAlbert/TasQ/config"
	"github.com/CharlieAlbert/TasQ/internal/db"
	"github.com/CharlieAlbert/TasQ/internal/jobs"
	"github.com/CharlieAlbert/TasQ/internal/worker"
)

func main() {
	config.LoadEnv()
	db.Connect()
	defer db.Close()

	// Initialize repository and service
	jobRepo := jobs.NewRepository(db.DB)
	jobService := jobs.NewService(jobRepo)
	worker := worker.NewWorker(jobService)

	log.Println("🚀 Starting job worker...")
	go worker.StartPolling(5 * time.Second)

	select {}
}
