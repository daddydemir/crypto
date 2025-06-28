package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"net/http"
)

type BollingerChart interface {
	CalculateBands() (middle, upper, lower, original []model.ChartModel)
	DrawBands(w http.ResponseWriter, r *http.Request)
}
