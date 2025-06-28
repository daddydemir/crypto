package signal

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/signal"
)

type DefaultSignalService struct {
	repo signal.SignalRepository
}

func NewDefaultSignalService(repo signal.SignalRepository) *DefaultSignalService {
	return &DefaultSignalService{repo: repo}
}

func (s *DefaultSignalService) SaveAll(signals []model.SignalModel) error {
	return s.repo.SaveAll(signals)
}
