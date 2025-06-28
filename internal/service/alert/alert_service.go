package alert

import (
	"fmt"
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/alert"
	"log/slog"
)

type AlertService struct {
	repo   alert.AlertRepository
	broker port.Broker
}

func NewAlertService(repo alert.AlertRepository, broker port.Broker) *AlertService {
	return &AlertService{
		repo:   repo,
		broker: broker,
	}
}

func (a *AlertService) Save(alert model.Alert) error {
	return a.repo.Save(alert)
}

func (a *AlertService) GetAll() ([]model.Alert, error) {
	return a.repo.GetAll()
}

func (a *AlertService) Delete(id int) error {
	return a.repo.Delete(id)
}

func (a *AlertService) ControlAlerts(coins []model.Coin) {
	alerts, err := a.repo.GetAll()
	if err != nil {
		slog.Error("AlertService: GetAll", "error", err)
		return
	}
	if len(alerts) == 0 {
		slog.Info("AlertService: no alerts to evaluate")
		return
	}

	alertsBySymbol := groupAlertsBySymbol(alerts)

	for _, coin := range coins {
		if relatedAlerts, found := alertsBySymbol[coin.Symbol]; found {
			for _, alert := range relatedAlerts {
				if triggered(alert, coin.PriceUsd) {
					message := buildAlertMessage(alert, coin.PriceUsd)
					a.sendAndDelete(int(alert.ID), message)
				}
			}
		}
	}
}

func groupAlertsBySymbol(alerts []model.Alert) map[string][]model.Alert {
	grouped := make(map[string][]model.Alert)
	for _, alert := range alerts {
		grouped[alert.Coin] = append(grouped[alert.Coin], alert)
	}
	return grouped
}

func triggered(alert model.Alert, currentPrice float32) bool {
	if alert.IsAbove {
		return currentPrice >= alert.Price
	}
	return currentPrice <= alert.Price
}

func buildAlertMessage(alert model.Alert, currentPrice float32) string {
	if alert.IsAbove {
		return fmt.Sprintf("%s %s'in üzerine çıktı: %s", alert.Coin, formatPrice(alert.Price), formatPrice(currentPrice))
	}
	return fmt.Sprintf("%s %s'in altına düştü: %s", alert.Coin, formatPrice(alert.Price), formatPrice(currentPrice))
}

func (a *AlertService) sendAndDelete(id int, message string) {
	if err := a.broker.SendMessage(message); err != nil {
		slog.Error("AlertService: SendMessage", "error", err)
		return
	}
	if err := a.repo.Delete(id); err != nil {
		slog.Error("AlertService: Delete", "error", err)
	}
}

func formatPrice(value float32) string {
	if value < 1 {
		return fmt.Sprintf("$%.4f", value)
	}
	return fmt.Sprintf("$%.2f", value)
}
