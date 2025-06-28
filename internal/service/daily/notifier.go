package daily

import (
	"fmt"
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/daily"
	"time"
)

type DailyNotifier struct {
	Repo   daily.DailyRepository
	Output port.Broker
}

func (dn *DailyNotifier) SendSummary() {

	start, end := time.Now().AddDate(0, 0, -1).Format("2006-01-02"), time.Now().Format("2006-01-02")

	big, _ := dn.Repo.FindTopBiggerByRate(end, start)
	small, _ := dn.Repo.FindTopSmallerByRate(end, start)

	msg1 := formatTable(big)
	msg2 := formatTable(small)

	_ = dn.Output.SendMessage(msg1)
	_ = dn.Output.SendMessage(msg2)
}

func formatTable(items [5]model.DailyModel) string {
	var out string
	for _, item := range items {
		out += fmt.Sprintf("(%s)\t%%%.2f\t|\t%.2f$\n", item.ExchangeId, item.Rate, item.Modulus)
	}
	return out
}
