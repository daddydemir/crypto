package jobs

import (
	"github.com/daddydemir/crypto/internal/service/movingaverage"
	"github.com/robfig/cron/v3"
)

type CheckAllMAJob struct {
	maService movingaverage.Service
}

func NewCheckAllMAJob(s movingaverage.Service) *CheckAllMAJob {
	return &CheckAllMAJob{s}
}

func (j *CheckAllMAJob) Schedule(c *cron.Cron) error {
	_, err := c.AddFunc("30 05 * * *", func() {
		j.maService.CheckAll(7, 25, 99)
	})
	return err
}
