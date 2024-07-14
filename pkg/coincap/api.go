package coincap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	allCoins           = "https://api.coincap.io/v2/assets"
	priceHistoryWithId = "https://api.coincap.io/v2/assets/%v/history?interval=d1"
)

func ListCoins() []Coin {

	var data Data[Coin]

	resp, err := http.Get(allCoins)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	return data.Data
}

func HistoryWithId(s string) []History {

	var data Data[History]

	url := fmt.Sprintf(priceHistoryWithId, s)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	return data.Data
}
