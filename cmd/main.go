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

	jobRepo := jobs.NewRepository(db.DB)
	jobService := jobs.NewService(jobRepo)

	jobType := "email"
	payload := map[string]any{
		"to":      "user@example.com",
		"subject": "Welcome!",
		"body":    "Hello from TasQ!",
	}

	if err := jobService.EnqueueJob(jobType, payload); err != nil {
		log.Fatalf("Unable to submit job, %v", err)
	}

	log.Printf("✅ Job submitted successfully")
}
