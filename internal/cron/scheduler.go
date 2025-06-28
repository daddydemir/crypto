package cronjob

import (
	"github.com/robfig/cron/v3"
	"log/slog"
	"time"
)

type Job interface {
	Schedule(*cron.Cron) error
}

type Scheduler struct {
	cron    *cron.Cron
	jobs    []Job
	enabled bool
}

func NewScheduler(jobs ...Job) *Scheduler {
	return &Scheduler{
		cron:    cron.New(cron.WithLocation(turkey())),
		jobs:    jobs,
		enabled: true,
	}
}

func (s *Scheduler) Start() {
	if !s.enabled {
		slog.Warn("Scheduler disabled")
		return
	}

	for _, job := range s.jobs {
		if err := job.Schedule(s.cron); err != nil {
			slog.Error("Job Schedule failed", "error", err)
		}
	}

	s.cron.Start()
}

func turkey() *time.Location {
	loc, err := time.LoadLocation("Turkey")
	if err != nil {
		slog.Error("Failed to load Turkey timezone", "error", err)
		return time.Local
	}
	return loc
}
