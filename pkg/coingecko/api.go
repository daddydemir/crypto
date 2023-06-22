package coingecko

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/adapter"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	topHundred string = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1&sparkline=false"
)

func GetTopHundred() []adapter.Adapter {

	req, err := http.NewRequest(http.MethodGet, topHundred, nil)
	if err != nil {
		log.Println("::GetTopHundred:: NewRequest err:{}", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("::GetTopHundred:: Do err:{}", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("::GetTopHundred:: ReadAll err:{}", err)
		return nil
	}
	var adapts []adapter.Adapter
	err = json.Unmarshal(body, &adapts)
	if err != nil {
		log.Println("::GetTopHundred:: Unmarshal err:{}", err)
	}
	return adapts
}
