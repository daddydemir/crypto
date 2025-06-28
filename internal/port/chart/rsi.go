package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"net/http"
)

type RsiChart interface {
	Calculate() []model.RsiModel
	Draw(w http.ResponseWriter, r *http.Request)
	Index() float32
}
