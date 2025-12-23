package application

import "github.com/daddydemir/crypto/pkg/atr/domain"

type PointService struct {
	repository domain.Repository
}

func NewPointService(repository domain.Repository) *PointService {
	return &PointService{
		repository: repository,
	}
}

func (s *PointService) GetPoints(symbol string) ([]domain.Point, error) {
	atrPoints, err := s.repository.GetPointsBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	points, err := domain.CalculateATR(atrPoints)
	return points, err
}
