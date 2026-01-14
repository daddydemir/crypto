package app

import (
	"context"
	"github.com/daddydemir/crypto/pkg/analyses/alert/domain"
	"github.com/daddydemir/crypto/pkg/analyses/alert/infra"
)

type App struct {
	repo *infra.Repository
}

func NewApp(repo *infra.Repository) *App {
	return &App{repo: repo}
}

func (a *App) CreateAlert(ctx context.Context, coin string, price float32, isAbove bool) (*domain.Alert, error) {
	alert := domain.NewAlert(coin, price, isAbove)
	err := a.repo.Save(ctx, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (a *App) UpdateAlert(ctx context.Context, id uint, price float32, isAbove bool) error {
	alert, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	alert.Update(price, isAbove)
	return a.repo.Update(ctx, alert)
}

func (a *App) DeleteAlert(ctx context.Context, id uint) error {
	return a.repo.Delete(ctx, id)
}

func (a *App) ListAlerts(ctx context.Context) ([]domain.Alert, error) {
	return a.repo.List(ctx)
}
