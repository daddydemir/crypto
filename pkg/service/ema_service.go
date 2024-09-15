package service

import (
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
	"log/slog"
	"net/http"
)

type emaService struct {
	graphic graphs.Graph
}

func NewEmaService(name string, period int) *emaService {
	return &emaService{
		ma.NewEma(name, period),
	}
}

func (e *emaService) Draw() func(w http.ResponseWriter, r *http.Request) {
	list := e.graphic.Calculate()

	if len(list) == 0 {
		slog.Error("Draw:graphic.Calculate", "error", "list is empty")
		return nil
	}

	draw := e.graphic.Draw(list)
	return draw
}
