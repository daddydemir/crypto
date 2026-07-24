package app

import "github.com/daddydemir/crypto/pkg/auth/basic/domain"

type App struct {
	repository domain.UserAuthRepository
}

func NewApp(repository domain.UserAuthRepository) *App {
	return &App{repository: repository}
}

func (a *App) CreateUser(user domain.User) error {
	return a.repository.CreateUser(user)
}

func (a *App) Login(username, password string) (*domain.LoginResponse, error) {
	return a.repository.Login(username, password)
}
