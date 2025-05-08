package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/rsi"
	"log/slog"
	"net/http"
)

type RsiService struct {
	graphic graphs.Graph
	broker  broker.Broker
}

func NewRsiService(name string) *RsiService {
	return &RsiService{
		rsi.NewRsi(name),
		broker.GetBrokerService(),
	}
}

func (r *RsiService) Draw() func(w http.ResponseWriter, r *http.Request) {

	list := r.graphic.Calculate()
	if len(list) == 0 {
		slog.Error("Draw:graphic.Calculate", "error", "list is empty")
		return nil
	}

	draw := r.graphic.Draw(list)
	return draw
}

func (r *RsiService) CalculateIndex() {
	index := r.graphic.Index()

	if index < 30 || index > 70 {
		err := r.broker.SendMessage(fmt.Sprintf("RSI index %2.f", index))
		if err != nil {
			slog.Error("CalculateIndex:broker.SendMessage", "error", err)
		}
	}
}
