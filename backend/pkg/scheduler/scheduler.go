package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Job interface {
	Name() string
	Execute(ctx context.Context) error
	Schedule() string
	IsEnabled() bool
}

type JobStatus struct {
	LastRun     time.Time
	LastError   error
	IsRunning   bool
	SuccessRuns int64
	ErrorRuns   int64
}

type Scheduler struct {
	cron     *cron.Cron
	jobs     map[string]Job
	status   map[string]*JobStatus
	mu       sync.RWMutex
	logger   *log.Logger
	stopChan chan struct{}
}

func New(opts ...Option) *Scheduler {
	s := &Scheduler{
		cron:     cron.New(),
		jobs:     make(map[string]Job),
		status:   make(map[string]*JobStatus),
		logger:   log.New(os.Stdout, "[SCHEDULER] ", log.LstdFlags),
		stopChan: make(chan struct{}),
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Scheduler) executeJob(job Job) {
	jobStatus := s.status[job.Name()]

	if jobStatus.IsRunning {
		s.logger.Printf("Job %s is already running, skipping execution", job.Name())
		return
	}

	s.mu.Lock()
	jobStatus.IsRunning = true
	jobStatus.LastRun = time.Now()
	s.mu.Unlock()

	ctx := context.Background()
	err := job.Execute(ctx)

	s.mu.Lock()
	jobStatus.IsRunning = false
	if err != nil {
		jobStatus.LastError = err
		jobStatus.ErrorRuns++
		s.logger.Printf("Error executing job %s: %v", job.Name(), err)
	} else {
		jobStatus.LastError = nil
		jobStatus.SuccessRuns++
	}
	s.mu.Unlock()
}

func (s *Scheduler) RegisterJob(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.Name()]; exists {
		return fmt.Errorf("job %s already registered", job.Name())
	}

	if !job.IsEnabled() {
		s.logger.Printf("Job %s is disabled, skipping registration", job.Name())
		return nil
	}

	_, err := s.cron.AddFunc(job.Schedule(), func() {
		s.executeJob(job)
	})

	if err != nil {
		return fmt.Errorf("failed to add job %s: %w", job.Name(), err)
	}

	s.jobs[job.Name()] = job
	s.status[job.Name()] = &JobStatus{}
	s.logger.Printf("Registering job %s with schedule %s", job.Name(), job.Schedule())
	return nil
}

func (s *Scheduler) Start() {
	s.logger.Println("Starting scheduler...")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.logger.Println("Stopping scheduler...")

	s.cron.Stop()
	close(s.stopChan)
}

func (s *Scheduler) GetJobStatus(name string) (*JobStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.status[name]
	if !exists {
		return nil, fmt.Errorf("job %s not found", name)
	}
	return status, nil
}

func (s *Scheduler) Logger() *log.Logger {
	return s.logger
}
