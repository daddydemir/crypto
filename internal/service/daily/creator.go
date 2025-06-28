package daily

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"log/slog"
	"time"
)

type DailyCreator struct {
	service *DefaultDailyService
	client  port.CoingeckoClient
}

func NewDailyCreator(service *DefaultDailyService, client port.CoingeckoClient) *DailyCreator {
	return &DailyCreator{service: service, client: client}
}

func (c *DailyCreator) CreateDaily(morning bool) {
	dailies, err := c.client.GetTopHundredDaily()

	if morning {
		err = c.service.SaveAll(dailies)
		if err != nil {
			slog.Error("CreateDaily.SaveAll", "error", err)
		}
		return
	}

	gecko := make(map[string]model.DailyModel)
	db := make(map[string]model.DailyModel)

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	models, err := c.service.FindByDateRange(yesterday, today)
	if err != nil {
		slog.Error("CreateDaily.FindByDateRange", "error", err)
		return
	}

	var save []model.DailyModel
	for i := 0; i < len(dailies); i++ {
		gecko = mergeMap(gecko, dailies[i].ToMap())
		db = mergeMap(db, models[i].ToMap())
	}

	for key, value := range gecko {
		current := db[key]
		if current.ExchangeId == "" {
			continue
		}
		current.LastPrice = value.LastPrice
		current.Min = min(current.Min, value.Min)
		current.Max = max(current.Max, value.Max)
		current.Avg = (current.Min + current.Max) / 2
		current.Modulus = current.Max - current.Min
		current.Rate = current.Modulus * 100 / current.Avg
		save = append(save, current)
	}

	err = c.service.SaveAll(save)
	if err != nil {
		slog.Error("CreateDaily.SaveAll", "error", err)
	}
}

func mergeMap(m1, m2 map[string]model.DailyModel) map[string]model.DailyModel {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
