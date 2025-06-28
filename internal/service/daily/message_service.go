package daily

import (
	"context"
	"fmt"
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/daily"
	"github.com/daddydemir/crypto/pkg/broker"
	"log/slog"
	"time"
)

type DefaultDailyMessageService struct {
	dailyService daily.DailyService
	broker       port.Broker
}

func NewDefaultDailyMessageService(ds daily.DailyService, b broker.Broker) *DefaultDailyMessageService {
	return &DefaultDailyMessageService{
		dailyService: ds,
		broker:       b,
	}
}

func (s *DefaultDailyMessageService) SendDailySummaryMessage(ctx context.Context) error {
	startDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	bigger, err := s.dailyService.FindTopBiggerByRate(startDate, endDate)
	if err != nil {
		slog.Error("SendDailySummaryMessage:FindTopBiggerByRate", "error", err)
		return err
	}

	smaller, err := s.dailyService.FindTopSmallerByRate(startDate, endDate)
	if err != nil {
		slog.Error("SendDailySummaryMessage:FindTopSmallerByRate", "error", err)
		return err
	}

	biggerMsg := formatMessage("📈 En Yüksek Artış", bigger)
	smallerMsg := formatMessage("📉 En Büyük Düşüş", smaller)

	if err = s.broker.SendMessage(biggerMsg); err != nil {
		slog.Error("SendDailySummaryMessage:broker.SendMessage:bigger", "error", err)
	}

	if err = s.broker.SendMessage(smallerMsg); err != nil {
		slog.Error("SendDailySummaryMessage:broker.SendMessage:smaller", "error", err)
	}

	return nil
}

func formatMessage(title string, models [5]model.DailyModel) string {
	message := fmt.Sprintf("🔥 %s:\n", title)
	for _, d := range models {
		message += fmt.Sprintf("- (%s): %%%.2f | %.2f$\n", d.ExchangeId, d.Rate, d.Modulus)
	}
	return message
}
