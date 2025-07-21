package worker

import (
	"context"
	"log"
	"time"

	"github.com/CharlieAlbert/TasQ/internal/handlers"
	"github.com/CharlieAlbert/TasQ/internal/jobs"
)

type Worker struct {
	service *jobs.Service
}

func NewWorker(service *jobs.Service) *Worker {
	return &Worker{service: service}
}

func (w *Worker) StartPolling(interval time.Duration) {
	for {
		time.Sleep(interval)

		job, err := w.service.FetchNextPendingJob(context.Background())
		if err != nil {
			log.Printf("❌ Error fetching job: %v", err)
			continue
		}

		if job == nil {
			continue
		}

		log.Printf("🔄 Processing job ID %d of type %s\n: %s", job.ID, job.Type)

		payload := job.Payload
		if payload == nil {
			log.Printf("❌ Invalid job payload format\n")
			w.service.UpdateJobStatus(context.Background(), job.ID, "failed")
			continue
		}

		handlerFunc, ok := handlers.Registry[job.Type]
		if !ok {
			log.Printf("❌ No handler registered for job type: %s\n", job.Type)
			w.service.UpdateJobStatus(context.Background(), job.ID, "failed")
			continue
		}

		err = handlerFunc(payload)
		if err != nil {
			log.Printf("❌ Error processing job ID %d: %v\n", job.ID, err)
			w.service.UpdateJobStatus(context.Background(), job.ID, "failed")
		} else {
			log.Printf("✅ Job ID %d processed successfully\n", job.ID)
			w.service.UpdateJobStatus(context.Background(), job.ID, "completed")
		}
	}
}
