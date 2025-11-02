package service

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
	"time"
)

type ValidateService struct {
	cache.Cache
	client coincap.CoinCapClient
}

func NewValidateService(c cache.Cache, client coincap.CoinCapClient) *ValidateService {
	return &ValidateService{
		c,
		client,
	}
}

func (v *ValidateService) Validate() {
	cacheService := v.Cache
	client := v.client
	var coins []coincap.Coin
	data, err := cacheService.Get("coinList")
	if err != nil {
		slog.Error("Validate:cacheService.Get", "key", "coinList", "err", err)
		return
	}

	bytes, ok := data.(string)
	if !ok {
		slog.Error("Validate:data.(string)", "data", data)
		return
	}

	err = json.Unmarshal([]byte(bytes), &coins)
	if err != nil {
		slog.Error("Validate:json.Unmarshal", "bytes", bytes, "err", err)
		return
	}

	for _, i := range coins {
		array := make([]coincap.History, 0, 1)
		err = cacheService.GetList(i.Id, &array, -1, -1)
		if err != nil {
			slog.Error("Validate:cacheService.GetList", "coin", i.Id, "error", err)
			continue
		}

		if len(array) == 0 {
			slog.Info("Validate:len(array) == 0 - Fetching initial data", "coin", i.Id)

			err := v.fetchHistoricalDataByYear(i.Id, cacheService, client)
			if err != nil {
				slog.Error("Validate:fetchHistoricalDataByYear", "coin", i.Id, "error", err)
				continue
			}

		} else {
			_, histories := client.HistoryWithTime(i.Id, array[0].Date.Add(time.Hour*24).UnixNano(), time.Now().UnixNano())

			if len(histories) == 0 {
				slog.Error("Validate:client.HistoryWithTime", "coin", i.Id, "error", "no new data available")
				continue
			}

			err = cacheService.SetList(i.Id, histories, 0)
			if err != nil {
				slog.Error("Validate:cacheService.SetList", "coin", i.Id, "error", err)
				continue
			}

			slog.Info("Validation Success", "coin", i.Id, "new_records", len(histories))
			time.Sleep(time.Second * 20)
		}
	}
}

func (v *ValidateService) fetchHistoricalDataByYear(coinId string, cacheService cache.Cache, client coincap.CoinCapClient) error {
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	currentTime := time.Now()
	var allHistories []coincap.History

	for year := 2020; year <= currentTime.Year(); year++ {
		yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		var yearEnd time.Time

		if year == currentTime.Year() {
			yearEnd = currentTime
		} else {
			yearEnd = time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
		}

		if year == 2020 {
			yearStart = startDate
		}

		slog.Info("Fetching data for year", "coin", coinId, "year", year, "start", yearStart.Format("2006-01-02"), "end", yearEnd.Format("2006-01-02"))

		_, histories := client.HistoryWithTime(coinId, yearStart.UnixNano(), yearEnd.UnixNano())

		if len(histories) == 0 {
			slog.Warn("No data for year", "coin", coinId, "year", year)
		} else {
			allHistories = append(allHistories, histories...)
			slog.Info("Data fetched for year", "coin", coinId, "year", year, "count", len(histories))
		}

		time.Sleep(time.Second * 20)
	}

	if len(allHistories) == 0 {
		return fmt.Errorf("no historical data found for coin: %s", coinId)
	}

	err := cacheService.SetList(coinId, allHistories, 0)
	if err != nil {
		return fmt.Errorf("failed to cache historical data: %w", err)
	}

	slog.Info("All historical data cached successfully", "coin", coinId, "total_count", len(allHistories))
	return nil
}
