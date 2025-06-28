package daily

import "github.com/daddydemir/crypto/internal/domain/model"

type DailyService interface {
	FindByDateRange(start, end string) ([]model.DailyModel, error)
	FindByIdAndDateRange(id, start, end string) ([]model.DailyModel, error)
	FindTopSmallerByRate(start, end string) ([5]model.DailyModel, error)
	FindTopBiggerByRate(start, end string) ([5]model.DailyModel, error)
	SaveAll([]model.DailyModel) error
}
