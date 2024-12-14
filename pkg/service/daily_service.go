package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/model"
	"log/slog"
	"time"
)

type DailyService struct {
	dailyRepo model.DailyRepository
}

func NewDailyService(dailyRepo model.DailyRepository) *DailyService {
	return &DailyService{dailyRepo: dailyRepo}
}

func (d *DailyService) FindByDateRange(start, end string) ([]model.DailyModel, error) {
	return d.dailyRepo.FindByDateRange(start, end)
}

func (d *DailyService) FindByIdAndDateRange(id, start, end string) ([]model.DailyModel, error) {
	return d.dailyRepo.FindByIdAndDateRange(id, start, end)
}

func (d *DailyService) FindTopSmallerByRate(start, end string) ([5]model.DailyModel, error) {
	return d.dailyRepo.FindTopSmallerByRate(start, end)
}

func (d *DailyService) FindTopBiggerByRate(start, end string) ([5]model.DailyModel, error) {
	return d.dailyRepo.FindTopBiggerByRate(start, end)
}

func (d *DailyService) SaveAll(dailies []model.DailyModel) error {
	return d.dailyRepo.SaveAll(dailies)
}

func (d *DailyService) CreateDaily(morning bool) {
	dailies := GetDaily()

	if morning {
		err := d.SaveAll(dailies)
		if err != nil {
			slog.Error("CreateDaily.SaveAll", "error", err)
		}
		return
	} else {
		db := make(map[string]model.DailyModel)
		gecko := make(map[string]model.DailyModel)
		models, err := d.FindByDateRange(time.Now().AddDate(0, 0, -1).Format("2006-01-02"), time.Now().Format("2006-01-02"))
		if err != nil {
			slog.Error("CreateDaily.FindByDateRange", "error", err)
		}
		var save []model.DailyModel

		for i := 0; i < len(dailies); i++ {
			gecko = mergeMap(gecko, dailies[i].ToMap())
			db = mergeMap(db, models[i].ToMap())
		}

		current := model.DailyModel{}

		for key, value := range gecko {
			current = db[key]
			if current.ExchangeId == "" {
				continue
			}
			current.LastPrice = value.LastPrice
			if current.Min > value.Min {
				current.Min = value.Min
			}

			if current.Max < value.Max {
				current.Max = value.Max
			}

			current.Avg = (current.Min + current.Max) / 2
			current.Modulus = current.Max - current.Min
			current.Rate = current.Modulus * 100 / current.Avg

			save = append(save, current)
		}
		/*err = d.SaveAll(save)
		if err != nil {
			slog.Error("CreateDaily.SaveAll", "error", err)
		} */
		d.sendDailyMessage()
	}

}

func mergeMap(m1, m2 map[string]model.DailyModel) map[string]model.DailyModel {
	merged := make(map[string]model.DailyModel, len(m1)+len(m2))
	merged = m1
	for key, value := range m2 {
		merged[key] = value
	}

	return merged
}

func (d *DailyService) sendDailyMessage() {
	var m1, m2, rate, mod string

	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	bigger, err := d.FindTopBiggerByRate(endDate, startDate)
	if err != nil {
		slog.Error("FindTopBiggerByRate", "error", err, "startDate", startDate, "endDate", endDate)
	}
	smaller, err := d.FindTopSmallerByRate(endDate, startDate)
	if err != nil {
		slog.Error("FindTopSmallerByRate", "error", err, "startDate", startDate, "endDate", endDate)
	}

	for i := 0; i < 5; i++ {
		rate = fmt.Sprintf("%.2f", bigger[i].Rate)
		mod = fmt.Sprintf("%.1f", bigger[i].Modulus)
		m1 += "(" + bigger[i].ExchangeId + ")\t %" + rate + "\t | \t" + mod + "$ \n"

		rate = fmt.Sprintf("%.2f", smaller[i].Rate)
		mod = fmt.Sprintf("%v", smaller[i].Modulus)
		m2 += "(" + smaller[i].ExchangeId + ")\t %" + rate + "\t | \t" + mod + "$ \n"
	}

	brokerService := broker.GetBrokerService()
	err = brokerService.SendMessage(m1)
	if err != nil {
		slog.Error("SendMessage", "error", err)
	}
	err = brokerService.SendMessage(m2)
	if err != nil {
		slog.Error("SendMessage", "error", err)
	}
}
