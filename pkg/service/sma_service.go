package service

import (
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
	"log/slog"
	"net/http"
)

type smaService struct {
	graphic graphs.Graph
}

func NewSmaService(name string, period int) *smaService {
	return &smaService{
		ma.NewSma(name, period),
	}
}

func (m *smaService) Draw() func(w http.ResponseWriter, r *http.Request) {
	list := m.graphic.Calculate()

	if len(list) == 0 {
		slog.Error("Draw:graphic.Calculate", "error", "list is empty")
		return nil
	}

	draw := m.graphic.Draw(list)
	return draw
}
