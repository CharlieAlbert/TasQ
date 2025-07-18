package jobs

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) EnqueueJob(jobType string, payload map[string]any) error {
	return s.repo.SubmitJob(jobType, payload)
}

func (s *Service) FetchQueuedJobs(limit int) ([]Job, error) {
	return s.repo.GetQueuedJobs(limit)
}

func (s *Service) FetchNextPendingJob(ctx context.Context) (*Job, error) {
	return s.repo.FetchNextPendingJob(ctx)
}

func (s *Service) UpdateJobStatus(ctx context.Context, jobID int, status string) error {
	return s.repo.UpdateJobStatus(ctx, jobID, status)
}