package main

import (
	"fmt"
	"sync"
	"time"
)

type JobFunc func() error

type Job struct {
	Name       string
	Interval   time.Duration
	Exec       JobFunc
	MaxRetries int

	mu         sync.Mutex
	retryCount int
	stopChan   chan struct{}
}

type Scheduler struct {
	jobs map[string]*Job
	mu   sync.Mutex
	wg   sync.WaitGroup
	quit chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: make(map[string]*Job),
		quit: make(chan struct{}),
	}
}

func (s *Scheduler) AddJob(name string, interval time.Duration, exec JobFunc, maxRetries int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.jobs[name]; exists {
		return fmt.Errorf("job %s already exists", name)
	}
	job := &Job{
		Name:       name,
		Interval:   interval,
		Exec:       exec,
		MaxRetries: maxRetries,
		stopChan:   make(chan struct{}),
	}
	s.jobs[name] = job
	s.wg.Add(1)
	go s.runJob(job)
	return nil
}

func (s *Scheduler) runJob(job *Job) {
	defer s.wg.Done()
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.executeJob(job)
		case <-job.stopChan:
			return
		case <-s.quit:
			return
		}
	}
}

func (s *Scheduler) executeJob(job *Job) {
	job.mu.Lock()
	defer job.mu.Unlock()
	err := job.Exec()
	if err != nil && job.retryCount < job.MaxRetries {
		job.retryCount++
		time.Sleep(job.Interval * time.Duration(job.retryCount))
		err = job.Exec()
	}
	if err == nil {
		job.retryCount = 0
	}
}

func (s *Scheduler) RemoveJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("job %s not found", name)
	}
	close(job.stopChan)
	delete(s.jobs, name)
	return nil
}

func (s *Scheduler) Stop() {
	close(s.quit)
	s.wg.Wait()
}
