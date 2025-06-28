package alert

import "github.com/daddydemir/crypto/internal/domain/model"

type AlertRepository interface {
	Save(alert model.Alert) error
	GetAll() ([]model.Alert, error)
	Delete(id int) error
}
