package cronjob

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/factory"
	"github.com/robfig/cron/v3"
	"log/slog"
	"time"
)

var serviceFactory *factory.ServiceFactory

func init() {
	location, _ := time.LoadLocation("Turkey")
	c := cron.New(cron.WithLocation(location))

	serviceFactory = factory.NewServiceFactory(database.GetDatabaseService(), cache.GetCacheService(), broker.GetBrokerService())

	dailyStart(c)
	dailyEnd(c)
	validateCache(c)
	checkAll(c)
	checkOutOfThresholds(c)

	c.Start()
}

func dailyStart(task *cron.Cron) {

	spec := "15 00 * * *"
	ds := serviceFactory.NewDailyService()

	entryID, err := task.AddFunc(spec, func() {
		ds.CreateDaily(true)
	})

	printLog(entryID, err, "dailyStart cron ID : ")

}

func dailyEnd(task *cron.Cron) {

	spec := "15 23 * * *"

	ds := serviceFactory.NewDailyService()

	entryID, err := task.AddFunc(spec, func() {
		ds.CreateDaily(false)
	})

	printLog(entryID, err, "dailyEnd cron ID : ")
}

func rsiCheck(task *cron.Cron) {
	spec := "30 10,15 * * *"

	entryId, err := task.AddFunc(spec, func() {
		// todo:
		//	service.RSIGraph(&rabbit.Publisher{})
	})
	printLog(entryId, err, "rsiCheck cron ID : ")
}

func validateCache(task *cron.Cron) {
	spec := "30 04 * * *"
	validateService := serviceFactory.NewValidateService()
	entryID, err := task.AddFunc(spec, func() {
		validateService.Validate()
	})
	printLog(entryID, err, "validateCache cron ID : ")
}

func checkAll(task *cron.Cron) {
	spec := "30 05 * * *"
	task.AddFunc(spec, func() {
		averageService := serviceFactory.NewAverageService()
		averageService.CheckAll(7, 25, 99)
	})
}

func checkOutOfThresholds(task *cron.Cron) {
	spec := "50 12 * * *"

	_, _ = task.AddFunc(spec, func() {
		bollingerService := serviceFactory.NewBollingerService()
		bollingerService.CheckThresholds()
	})
}

func printLog(entryID cron.EntryID, err error, message string) {
	if err != nil {
		slog.Error("printLog", "err", err, "message", message)
	}
	slog.Info("printLog", "entryID", entryID, "message", message)
}
