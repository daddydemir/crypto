package coincap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHistoryWithId(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Data[History]{
			Data: []History{
				{PriceUsd: 100.54654},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	url := server.URL + "/%s"
	priceHistoryWithId = url
	histories := HistoryWithId("1234")

	if len(histories) == 0 {
		t.Fatalf("expected non-empty history data")
	}

	expectedPrice := float32(100.54654)
	if histories[0].PriceUsd != expectedPrice {
		t.Errorf("expected price %.5f, got %.5f", expectedPrice, histories[0].PriceUsd)
	}
}
