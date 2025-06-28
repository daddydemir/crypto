package exchange

import (
	"log/slog"

	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
)

type ExchangeService struct {
	client port.CoingeckoClient
}

func NewExchangeService(client port.CoingeckoClient) *ExchangeService {
	return &ExchangeService{client: client}
}

func (e *ExchangeService) GetExchange() []model.ExchangeModel {
	exchanges, err := e.client.GetTopHundred()
	if err != nil {
		slog.Error("ExchangeService:GetTopHundred", "error", err)
		return nil
	}
	if len(exchanges) == 0 {
		slog.Warn("ExchangeService:GetTopHundred", "message", "empty list")
	}
	return exchanges
}
