package factory

import (
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/database/postgres"
	"github.com/daddydemir/crypto/pkg/service"
	"gorm.io/gorm"
)

type ServiceFactory struct {
	db     *gorm.DB
	cache  cache.Cache
	broker broker.Broker
}

func NewServiceFactory(db *gorm.DB, c cache.Cache, b broker.Broker) *ServiceFactory {
	return &ServiceFactory{
		db:     db,
		cache:  c,
		broker: b,
	}
}

func (f *ServiceFactory) NewAverageService() *service.MovingAverageService {
	repository := postgres.NewPostgresSignalRepository(f.db)
	signalService := service.NewSignalService(repository)
	return service.NewMovingAverageService(f.cache, f.broker, *signalService)
}

func (f *ServiceFactory) NewExchangeService() *service.ExchangeService {
	repository := postgres.NewPostgresExchangeRepository(f.db)
	return service.NewExchangeService(repository)
}
