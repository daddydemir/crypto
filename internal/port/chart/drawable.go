package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"net/http"
)

type DrawableChart interface {
	Calculate() []model.ChartModel
	Draw(w http.ResponseWriter, r *http.Request)
}
