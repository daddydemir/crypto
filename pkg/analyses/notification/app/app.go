package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/notification/domain"
	"github.com/daddydemir/crypto/pkg/analyses/notification/infra"
)

type App struct {
	infra *infra.Repository
}

func NewApp(infra *infra.Repository) *App {
	return &App{infra: infra}
}

func (a *App) GetAll(username string) []domain.Notification {
	return a.infra.GetAll(username)
}
