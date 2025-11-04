package alert

import (
	"context"
	"github.com/daddydemir/crypto/pkg/domain/alert"
)

type AlertRepository interface {
	Save(ctx context.Context, a *alert.Alert) error
	Update(ctx context.Context, a *alert.Alert) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*alert.Alert, error)
	List(ctx context.Context) ([]alert.Alert, error)
}

type Service struct {
	Repo AlertRepository
}

func NewService(repo AlertRepository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) CreateAlert(ctx context.Context, coin string, price float32, isAbove bool) (*alert.Alert, error) {
	a := alert.NewAlert(coin, price, isAbove)
	err := s.Repo.Save(ctx, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) UpdateAlert(ctx context.Context, id uint, price float32, isAbove bool) error {
	a, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	a.Update(price, isAbove)
	return s.Repo.Update(ctx, a)
}

func (s *Service) DeleteAlert(ctx context.Context, id uint) error {
	return s.Repo.Delete(ctx, id)
}

func (s *Service) ListAlerts(ctx context.Context) ([]alert.Alert, error) {
	return s.Repo.List(ctx)
}
