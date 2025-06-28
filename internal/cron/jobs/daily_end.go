package jobs

import (
	"github.com/daddydemir/crypto/internal/port/daily"
	dc "github.com/daddydemir/crypto/internal/service/daily"
	"github.com/robfig/cron/v3"
)

type DailyEndJob struct {
	dailyService daily.DailyService
	creator      dc.DailyCreator
}

func NewDailyEndJob(s daily.DailyService, d dc.DailyCreator) *DailyEndJob {
	return &DailyEndJob{dailyService: s, creator: d}
}

func (j *DailyEndJob) Schedule(c *cron.Cron) error {
	_, err := c.AddFunc("15 23 * * *", func() {
		j.creator.CreateDaily(false)
	})
	return err
}
