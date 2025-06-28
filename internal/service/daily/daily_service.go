package daily

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port/daily"
)

type DefaultDailyService struct {
	repo daily.DailyRepository
}

func NewDefaultDailyService(repo daily.DailyRepository) *DefaultDailyService {
	return &DefaultDailyService{repo: repo}
}

func (s *DefaultDailyService) FindByDateRange(start, end string) ([]model.DailyModel, error) {
	return s.repo.FindByDateRange(start, end)
}

func (s *DefaultDailyService) FindByIdAndDateRange(id, start, end string) ([]model.DailyModel, error) {
	return s.repo.FindByIdAndDateRange(id, start, end)
}

func (s *DefaultDailyService) FindTopSmallerByRate(start, end string) ([5]model.DailyModel, error) {
	return s.repo.FindTopSmallerByRate(start, end)
}

func (s *DefaultDailyService) FindTopBiggerByRate(start, end string) ([5]model.DailyModel, error) {
	return s.repo.FindTopBiggerByRate(start, end)
}

func (s *DefaultDailyService) SaveAll(dailies []model.DailyModel) error {
	return s.repo.SaveAll(dailies)
}
