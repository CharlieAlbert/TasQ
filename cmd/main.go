package main

import (
	"log"

	"github.com/CharlieAlbert/TasQ/config"
	"github.com/CharlieAlbert/TasQ/internal/db"
	"github.com/CharlieAlbert/TasQ/internal/jobs"
)

func main() {
	config.LoadEnv()
	db.Connect()
	defer db.Close()

	// Initialize repository and service
	jobRepo := jobs.NewRepository(db.DB)
	jobService := jobs.NewService(jobRepo)

	payload := map[string]any{
		"user_id":  123,
		"task":     "send_email",
		"priority": "high",
	}

	// Enqueue job
	err := jobService.EnqueueJob("email_task", payload)
	if err != nil {
		log.Fatalf("❌ Failed to enqueue job: %v", err)
	}

	log.Println("✅ Job successfully enqueued")
}
