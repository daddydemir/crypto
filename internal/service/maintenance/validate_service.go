package maintenance

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/internal/port"
	"log/slog"
	"time"
)

// Coin represents minimal coin data required
type Coin struct {
	Id string `json:"id"`
}

// CoinCache abstracts cache behavior needed by validation service.
type CoinCache interface {
	Get(key string) (any, error)
	SetList(key string, data any, ttl time.Duration) error
	GetList(key string, out any, start, end int64) error
}

type ValidateService struct {
	cache  CoinCache
	client port.CoinCapAPI
}

func NewValidateService(c CoinCache, client port.CoinCapAPI) *ValidateService {
	return &ValidateService{cache: c, client: client}
}

// Validate checks cached coin history and updates if outdated
func (v *ValidateService) Validate() {
	data, err := v.cache.Get("coinList")
	if err != nil {
		slog.Error("cache.Get failed", "err", err)
		return
	}

	raw, ok := data.(string)
	if !ok {
		slog.Error("coinList type assertion failed", "data", data)
		return
	}

	var coins []Coin
	if err = json.Unmarshal([]byte(raw), &coins); err != nil {
		slog.Error("coinList unmarshal failed", "err", err)
		return
	}

	for _, coin := range coins {
		if err = v.validateCoin(coin.Id); err != nil {
			slog.Error("coin validation failed", "coin", coin.Id, "error", err)
		} else {
			slog.Info("coin validation success", "coin", coin.Id)
		}
		time.Sleep(20 * time.Second)
	}
}

func (v *ValidateService) validateCoin(id string) error {
	var recent []struct {
		Date time.Time
	}
	err := v.cache.GetList(id, &recent, -1, -1)
	if err != nil {
		return fmt.Errorf("cache.GetList failed: %w", err)
	}
	if len(recent) == 0 {
		return fmt.Errorf("no cached history for %s", id)
	}

	start := recent[0].Date.Add(24 * time.Hour).UnixNano()
	end := time.Now().UnixNano()

	err, history := v.client.HistoryWithTime(id, start, end)
	if err != nil {
		return fmt.Errorf("HistoryWithTime failed: %w", err)
	}
	if len(history) == 0 {
		return fmt.Errorf("empty history from API")
	}

	err = v.cache.SetList(id, history, 0)
	if err != nil {
		return fmt.Errorf("cache.SetList failed: %w", err)
	}
	return nil
}
