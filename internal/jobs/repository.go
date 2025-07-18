package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SubmitJob(jobType string, payload map[string]any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `INSERT INTO jobs (type, payload)
	VALUES ($1, $2)`

	_, err = r.db.Exec(context.Background(), query, jobType, payloadJSON)
	if err != nil {
		return fmt.Errorf("failed to insert job: %w", err)
	}

	return nil
}

func (r *Repository) GetQueuedJobs(limit int) ([]Job, error) {
	query := `SELECT id, type, payload, status, created_at, updated_at
	FROM jobs
	WHERE status = 'queued'
	ORDER BY created_at ASC
	LIMIT $1`

	rows, err := r.db.Query(context.Background(), query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		var rawPayLoad []byte

		err := rows.Scan(&job.ID, &job.Type, &rawPayLoad, &job.Status, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(rawPayLoad, &job.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}