package cronjob

import (
	"github.com/daddydemir/crypto/config/log"
	db "github.com/daddydemir/crypto/pkg/database/service"
	service "github.com/daddydemir/crypto/pkg/service"
	"github.com/robfig/cron/v3"
	"time"
)

func StartCronJob() {
	location, _ := time.LoadLocation("Turkey")
	c := cron.New(cron.WithLocation(location))

	dailyStart(c)
	dailyEnd(c)
	rsiCheck(c)

	c.Start()
}

func dailyStart(task *cron.Cron) {

	spec := "15 00 * * *"

	entryID, err := task.AddFunc(spec, func() {
		db.CreateDaily(true)
	})

	printLog(entryID, err, "dailyStart cron ID : ")

}

func dailyEnd(task *cron.Cron) {

	spec := "15 23 * * *"

	entryID, err := task.AddFunc(spec, func() {
		db.CreateDaily(false)
	})

	printLog(entryID, err, "dailyEnd cron ID : ")
}

func rsiCheck(task *cron.Cron) {
	spec := "30 10,15 * * *"

	entryId, err := task.AddFunc(spec, func() {
		service.RSIGraph()
	})
	printLog(entryId, err, "rsiCheck cron ID : ")
}

func printLog(entryID cron.EntryID, err error, message string) {
	if err != nil {
		log.Errorln("Error creating cron :  ", err)
	}
	log.Infoln(message, entryID)
}
