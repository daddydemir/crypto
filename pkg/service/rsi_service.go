package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/rsi"
	"log/slog"
	"net/http"
)

type rsiService struct {
	graphic graphs.Graph
	broker  broker.Broker
}

func NewRsiService(name string) *rsiService {
	return &rsiService{
		rsi.NewRsi(name),
		broker.GetBrokerService(),
	}
}

func (r *rsiService) Draw() func(w http.ResponseWriter, r *http.Request) {

	list := r.graphic.Calculate()
	if len(list) == 0 {
		slog.Error("Draw:graphic.Calculate", "error", "list is empty")
		return nil
	}

	draw := r.graphic.Draw(list)
	return draw
}

func (r *rsiService) CalculateIndex() {
	index := r.graphic.Index()

	if index < 30 || index > 70 {
		err := r.broker.SendMessage(fmt.Sprintf("RSI index %2.f", index))
		if err != nil {
			slog.Error("CalculateIndex:broker.SendMessage", "error", err)
		}
	}
}
