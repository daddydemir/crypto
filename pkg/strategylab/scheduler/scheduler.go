package scheduler

import (
	"github.com/daddydemir/crypto/pkg/strategylab/app"
	"github.com/robfig/cron/v3"
	"log/slog"
)

func Start(a *app.App) *cron.Cron {
	c := cron.New()
	_, e := c.AddFunc("15 5 * * *", func() { a.EvaluateRunning() })
	if e != nil {
		slog.Error("strategy evaluator schedule", "error", e)
	}
	c.Start()
	return c
}
