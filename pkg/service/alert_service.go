package service

import "github.com/daddydemir/crypto/pkg/model"

type AlertService struct {
	alertRepo model.AlertRepository
}

func NewAlertService(repo model.AlertRepository) *AlertService {
	return &AlertService{alertRepo: repo}
}

func (a *AlertService) Save(alert model.Alert) error {
	return a.alertRepo.Save(alert)
}

func (a *AlertService) GetAll() ([]model.Alert, error) {
	return a.alertRepo.GetAll()
}
