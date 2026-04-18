package jobs

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pr0100pr0111/KV-redaction/internal/models"
)

type JobStore struct {
	jobs map[string]*models.ProcessingJob
	mu   sync.RWMutex
}

func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*models.ProcessingJob),
	}
}

func (s *JobStore) Create(fileName string) *models.ProcessingJob {
	jobID := uuid.New().String()
	job := &models.ProcessingJob{
		ID:        jobID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		InputFile: fileName,
		Progress:  0,
		Stage:     "upload",
	}
	s.mu.Lock()
	s.jobs[jobID] = job
	s.mu.Unlock()
	return job
}

func (s *JobStore) Get(id string) (*models.ProcessingJob, bool) {
	s.mu.RLock()
	job, exists := s.jobs[id]
	s.mu.RUnlock()
	return job, exists
}

func (s *JobStore) GetAll() []*models.ProcessingJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobList := make([]*models.ProcessingJob, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobList = append(jobList, job)
	}
	return jobList
}

func (s *JobStore) Update(id string, updateFunc func(job *models.ProcessingJob)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, exists := s.jobs[id]; exists {
		updateFunc(job)
		job.UpdatedAt = time.Now()
	}
}

func (s *JobStore) Delete(id string) {
	s.mu.Lock()
	delete(s.jobs, id)
	s.mu.Unlock()
}
