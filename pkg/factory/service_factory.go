package factory

import (
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/database/postgres"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
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

func (f *ServiceFactory) NewExchangeService() *service.ExchangeService {
	repository := postgres.NewPostgresExchangeRepository(f.db)
	return service.NewExchangeService(repository)
}

func (f *ServiceFactory) NewDailyService() *service.DailyService {
	repository := postgres.NewPostgresDailyRepository(f.db)
	return service.NewDailyService(repository)
}

func (f *ServiceFactory) NewAlertService() *service.AlertService {
	repository := postgres.NewPostgresAlertRepository(f.db)
	return service.NewAlertService(repository, f.broker)
}

func (f *ServiceFactory) NewCacheService() *service.CacheService {
	return service.NewCacheService(f.cache)
}

func (f *ServiceFactory) NewCoinCapClient() *coincap.Client {
	return coincap.NewClient()
}

func (f *ServiceFactory) NewCachedCoinCapClient() *coincap.CachedClient {
	return coincap.NewCachedClient(*f.NewCoinCapClient(), f.cache)
}

func (f *ServiceFactory) NewValidateService() *service.ValidateService {
	return service.NewValidateService(f.cache, f.NewCachedCoinCapClient())
}
