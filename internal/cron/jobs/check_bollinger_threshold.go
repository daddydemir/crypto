package jobs

import (
	"github.com/daddydemir/crypto/internal/service/bollinger"

	"github.com/robfig/cron/v3"
)

type CheckBollingerThresholdJob struct {
	bbService bollinger.Service
}

func NewCheckBollingerThresholdJob(s bollinger.Service) *CheckBollingerThresholdJob {
	return &CheckBollingerThresholdJob{s}
}

func (j *CheckBollingerThresholdJob) Schedule(c *cron.Cron) error {
	_, err := c.AddFunc("50 12 * * *", func() {
		j.bbService.CheckThresholds()
	})
	return err
}
