package service

import "github.com/daddydemir/crypto/pkg/model"

type DailyService struct {
	dailyRepo model.DailyRepository
}

func NewDailyService(dailyRepo model.DailyRepository) *DailyService {
	return &DailyService{dailyRepo: dailyRepo}
}

func (d *DailyService) FindByDateRange(start, end string) ([]model.DailyModel, error) {
	return d.dailyRepo.FindByDateRange(start, end)
}

func (d *DailyService) FindByIdAndDateRange(id, start, end string) ([]model.DailyModel, error) {
	return d.dailyRepo.FindByIdAndDateRange(id, start, end)
}

func (d *DailyService) FindTopSmallerByRate(start, end string) ([5]model.DailyModel, error) {
	return d.dailyRepo.FindTopSmallerByRate(start, end)
}

func (d *DailyService) FindTopBiggerByRate(start, end string) ([5]model.DailyModel, error) {
	return d.dailyRepo.FindTopBiggerByRate(start, end)
}
