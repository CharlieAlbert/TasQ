package main

import (
	"context"
	"log"
	"time"

	"github.com/CharlieAlbert/TasQ/config"
	"github.com/CharlieAlbert/TasQ/internal/db"
	"github.com/CharlieAlbert/TasQ/internal/jobs"
)

func main() {
	config.LoadEnv()
	db.Connect()
	defer db.Close()

	repo := jobs.NewRepository(db.DB)
	service := jobs.NewService(repo)

	for {
		ctx := context.Background()

		job, err := service.FetchNextPendingJob(ctx)
		if err != nil {
			log.Println("No job found or error: ", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Picked job ID: %d, Type: %s\n", job.ID, job.Type)

		// Simulate job processing
		if err := service.UpdateJobStatus(ctx, job.ID, "in_progress"); err != nil {
			log.Printf("failed to update status to processing: %v", err)
			continue
		}

		// Mark as done
		if err := service.UpdateJobStatus(ctx, job.ID, "completed"); err != nil {
			log.Printf("failed to update status to done: %v", err)
		}
	}
}

func processJob(job *jobs.Job) error {
	log.Printf("Processing job ID %d with payload: %v/n", job.ID, job.Payload)
	time.Sleep(2 * time.Second)
	return nil
}