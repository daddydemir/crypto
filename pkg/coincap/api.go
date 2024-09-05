package coincap

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	allCoins             = "https://api.coincap.io/v2/assets"
	priceHistoryWithId   = "https://api.coincap.io/v2/assets/%v/history?interval=d1"
	priceHistoryWithTime = "https://api.coincap.io/v2/assets/%v/history?interval=d1&start=%d&end=%d"
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

func HistoryWithTime(s string, start, end int64) []History {
	var data Data[History]
	url := fmt.Sprintf(priceHistoryWithTime, s, start/1_000_000, end/1_000_000)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("HistoryWithTime", "error", err)
		return nil
	}
	slog.Info("http get", "url", url, "statusCode", resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		slog.Error("HistoryWithTime json decode ", "error", err, "data", resp.Body)
		return nil
	}
	return data.Data
}
