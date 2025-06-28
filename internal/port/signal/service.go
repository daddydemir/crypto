package signal

import "github.com/daddydemir/crypto/internal/domain/model"

type SignalService interface {
	SaveAll([]model.SignalModel) error
}
