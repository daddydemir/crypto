package jobs

import (
	"github.com/daddydemir/crypto/internal/port/maintenance"
	"github.com/robfig/cron/v3"
)

type ValidateCacheJob struct {
	validateService maintenance.CacheValidator
}

func NewValidateCacheJob(s maintenance.CacheValidator) *ValidateCacheJob {
	return &ValidateCacheJob{s}
}

func (j *ValidateCacheJob) Schedule(c *cron.Cron) error {
	_, err := c.AddFunc("30 04 * * *", func() {
		j.validateService.Validate()
	})
	return err
}
