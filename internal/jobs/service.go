package jobs

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