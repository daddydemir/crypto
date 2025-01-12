package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
)

type AlertService struct {
	alertRepo model.AlertRepository
	broker    broker.Broker
}

func NewAlertService(repo model.AlertRepository, b broker.Broker) *AlertService {
	return &AlertService{alertRepo: repo, broker: b}
}

func (a *AlertService) Save(alert model.Alert) error {
	return a.alertRepo.Save(alert)
}

func (a *AlertService) GetAll() ([]model.Alert, error) {
	return a.alertRepo.GetAll()
}

func (a *AlertService) Delete(id int) error {
	return a.alertRepo.Delete(id)
}

func (a *AlertService) ControlAlerts(coins []coincap.Coin) {
	alerts, err := a.GetAll()
	if err != nil {
		slog.Error("Failed to fetch alerts", "error", err)
		return
	}

	if len(alerts) == 0 {
		slog.Warn("No alerts to process")
		return
	}

	alertMap := make(map[string][]model.Alert)
	for _, alert := range alerts {
		alertMap[alert.Coin] = append(alertMap[alert.Coin], alert)
	}

	for _, coin := range coins {
		if coinAlerts, exists := alertMap[coin.Symbol]; exists {
			for _, alert := range coinAlerts {
				if coin.Symbol == alert.Coin {
					var message string
					if alert.IsAbove {
						if coin.PriceUsd >= alert.Price {
							message = fmt.Sprintf("%s %s'in üzerine çıktı: %s\n", alert.Coin, formatPrice(alert.Price), formatPrice(coin.PriceUsd))
						}
					} else {
						if coin.PriceUsd <= alert.Price {
							message = fmt.Sprintf("%s %s'in altına düştü: %s\n", alert.Coin, formatPrice(alert.Price), formatPrice(coin.PriceUsd))
						}
					}

					if message != "" {
						err = a.broker.SendMessage(message)
						if err != nil {
							slog.Error("SendMessage", "error", err)
							return
						}

						err = a.Delete(int(alert.ID))
						if err != nil {
							slog.Error("Delete", "error", err)
						}
					}
				}
			}
		}
	}
}

func formatPrice(value float32) string {
	if value < 1 {
		return fmt.Sprintf("$%.4f", value)
	}
	return fmt.Sprintf("$%.2f", value)
}
