package service

import "github.com/daddydemir/crypto/pkg/model"

type SignalService struct {
	signalRepo model.SignalRepository
}

func NewSignalService(signalRepo model.SignalRepository) *SignalService {
	return &SignalService{
		signalRepo: signalRepo,
	}
}

func (s *SignalService) SaveAll(signals []model.SignalModel) error {
	return s.signalRepo.SaveAll(signals)
}
