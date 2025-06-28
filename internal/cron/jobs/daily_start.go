package jobs

import (
	"github.com/daddydemir/crypto/internal/port/daily"
	dc "github.com/daddydemir/crypto/internal/service/daily"
	"github.com/robfig/cron/v3"
)

type DailyStartJob struct {
	dailyService daily.DailyService
	creator      dc.DailyCreator
}

func NewDailyStartJob(s daily.DailyService, d dc.DailyCreator) *DailyStartJob {
	return &DailyStartJob{dailyService: s, creator: d}
}

func (j *DailyStartJob) Schedule(c *cron.Cron) error {
	_, err := c.AddFunc("15 00 * * *", func() {
		j.creator.CreateDaily(true)
	})
	return err
}
