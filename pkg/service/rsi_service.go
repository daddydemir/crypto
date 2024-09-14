package service

import (
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/rsi"
	"log/slog"
	"net/http"
)

type rsiService struct {
	graphic graphs.Graph
}

func NewRsiService(name string) *rsiService {
	return &rsiService{
		rsi.NewRsi(name),
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
