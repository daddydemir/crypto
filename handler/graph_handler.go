package handler

import (
	"github.com/daddydemir/crypto/assets"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/bollingerBands"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

var client *coincap.CachedClient

func init() {
	client = coincap.NewCachedClient(*coincap.NewClient(), cache.GetCacheService())
}
type CoinData struct {
	Index    int
	Name     string
	Symbol   string
	PriceUsd float32
	Rsi      float32
	RsiClass string
	Id       string
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8;")

	tmpl := assets.GetTemplate("templates/coin.html")

	var coinData [100]CoinData

	err, coins := client.ListCoins()
	if err != nil {
		slog.Error( "mainHandler:ListCoins",  "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	go alertService.ControlAlerts(coins)
	rsi := graphs.RSI{}

	for i, coin := range coins {
		var class string
		var index float32 = 0

		if i <= 25 {
			index = rsi.Index(coin.Id)
		}

		if index == 0 {
			class = " "
		} else if index >= 70 {
			class = "bg-green-600 bg-opacity-50"
		} else if index <= 30 {
			class = "bg-red-600 bg-opacity-50"
		} else {
			class = "bg-yellow-400 bg-opacity-50"
		}

		coinData[i] = CoinData{
			Index:    i + 1,
			Name:     coin.Name,
			Symbol:   coin.Symbol,
			PriceUsd: coin.PriceUsd,
			Rsi:      index,
			RsiClass: class,
			Id:       coin.Id,
		}
	}

	err = tmpl.Execute(w, struct {
		Coins [100]CoinData
	}{
		Coins: coinData,
	})
	if err != nil {
		slog.Error("Execute", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func rsiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	graphicService := service.NewRsiService(coin)
	function := graphicService.Draw()
	function(w, r)
}

func smaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	graphicService := service.NewSmaService(coin, 10)
	function := graphicService.Draw()
	function(w, r)
}

func emaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	graphicService := service.NewEmaService(coin, 25)
	function := graphicService.Draw()
	function(w, r)
}

func maHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]
	newMa := ma.NewMa(coin, 0)

	draw := newMa.Draw(newMa.Calculate())
	draw(w, r)
}

func bollingerBandsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	bands := bollingerBands.NewBollingerBands(coin, 20)
	list := bands.Calculate()

	function := bands.Draw(list)

	function(w, r)

}
