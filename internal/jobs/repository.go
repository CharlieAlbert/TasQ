package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) FetchNextPendingJob(ctx context.Context) (*Job, error) {
	query := `SELECT id, type, payload, status, created_at, updated_at
	FROM jobs
	WHERE status = 'queued'
	ORDER BY created_at ASC
	LIMIT 1
	FOR UPDATE SKIP LOCKED`

	row := r.db.QueryRow(ctx, query)

	var job Job
	var payloadData []byte

	err := row.Scan(
		&job.ID,
		&job.Type,
		&payloadData,
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch pending job: %d", err)
	}

	err = json.Unmarshal(payloadData, &job.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %d/n", err)
	}

	return &job, nil
}

func (r *Repository) UpdateJobStatus(ctx context.Context, jobID int, status string) error {
	query := `UPDATE jobs
	SET status = $1, updated_at = NOW()
	WHERE id = $2`

	_, err := r.db.Exec(ctx, query, status, jobID)
	if err != nil {
		return fmt.Errorf("failed to update job status: %s/n", err)
	}

	return nil
}