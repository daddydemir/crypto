package handler

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
	"github.com/gorilla/mux"
	"net/http"
)

func rsiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)

	coin := vars["coin"]
	fmt.Printf("Coin : %v \n", coin)

	rsi := graphs.RSI{}
	histories := rsi.Calculate(coin)
	draw := rsi.Draw(histories)
	draw(w, r)
}

func smaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)

	coin := vars["coin"]
	fmt.Printf("Coin : %v \n", coin)

	draw := ma.Draw(coin)
	draw(w, r)
}
