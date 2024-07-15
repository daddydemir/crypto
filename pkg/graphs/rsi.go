package graphs

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/database/service"
	"github.com/daddydemir/crypto/pkg/model"
)

type RSI struct {
	model.RsiModel
}

func (r RSI) Calculate() {
	date := dao.Date{
		StartDate: "2024-04-01",
		EndDate:   "2024-06-30",
		Id:        "trx",
	}

	list := service.GetDailyWithId(date)
	println("length:", len(list))

	response := make([]model.RsiModel, 0)

	start, end := 0, 14

	for end < len(list) {
		rsiModel := calculateIndex(list[start:end])
		response = append(response, rsiModel)
		start++
		end++
	}

	for _, r := range response {
		fmt.Printf("date: %v index: %.3f \n", r.Date, r.Index)
	}

}

func calculateIndex(list []model.DailyModel) model.RsiModel {

	response := new(model.RsiModel)
	period := 14
	gain := make([]model.RsiDate, 0)
	loss := make([]model.RsiDate, 0)
	data := make([]float32, period)

	for i, d := range list {

		data = append(data, d.LastPrice)

		if i == 0 {
			continue
		}

		tmp := d.LastPrice - list[i-1].LastPrice

		if tmp > 0 {
			gain = append(gain, model.RsiDate{
				Date: d.Date,
				Val:  tmp,
			})
		} else {
			loss = append(loss, model.RsiDate{
				Date: d.Date,
				Val:  tmp,
			})
		}
	}

	gainAverage := sum(gain) / 14
	lossAverage := sum(loss) / 14 * -1

	index := 100 - (100 / (1 + (gainAverage / lossAverage)))

	response.Gain = gain
	response.Loss = loss
	response.Data = data
	response.Index = index
	response.Date = list[len(list)-1].Date

	return *response
}

// kapanis fiyatlarini getir
// fiyat degisimlerini hesaplayin (her bir gunun fiyatini bir onceki gunun fiyatindan cikarin)
// kazanc ve kayiplar olarak ayirin
// ortalama kazanc ve ortalama kayip olarak hesaplayin
// Relative Strength degerini hesaplayin = ortalama kazanc / ortalama kayip
// RSI formulunu uygulayalim = 100 - (100 / (1+RS))

func sum(arr []model.RsiDate) float32 {
	var t float32
	for _, d := range arr {
		t += d.Val
	}
	return t
}

func (r RSI) Draw() {

}
