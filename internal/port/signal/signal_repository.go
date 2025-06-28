package signal

import "github.com/daddydemir/crypto/internal/domain/model"

type SignalRepository interface {
	SaveAll([]model.SignalModel) error
}
