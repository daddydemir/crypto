package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/daddydemir/crypto/pkg/analyses/alert/domain"
	"github.com/daddydemir/crypto/pkg/analyses/alert/infra"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/redis/go-redis/v9"
)

type App struct {
	repo   *infra.Repository
	redis  *redis.Client
	broker broker.Broker
	mu     sync.RWMutex
	active map[string][]domain.Alert
}

func NewApp(repo *infra.Repository, redisClient *redis.Client, messageBroker broker.Broker) *App {
	return &App{repo: repo, redis: redisClient, broker: messageBroker, active: make(map[string][]domain.Alert)}
}

func (a *App) CreateAlert(ctx context.Context, username, coin string, price float32, isAbove bool) (*domain.Alert, error) {
	alert := domain.NewAlert(username, strings.ToUpper(strings.TrimSpace(coin)), price, isAbove)
	err := a.repo.Save(ctx, &alert)
	if err != nil {
		return nil, err
	}
	if err := a.reload(ctx); err != nil {
		return nil, err
	}
	return &alert, nil
}

func (a *App) UpdateAlert(ctx context.Context, username string, id uint, price float32, isAbove bool) (*domain.Alert, error) {
	alert, err := a.repo.FindByID(ctx, username, id)
	if err != nil {
		return nil, err
	}
	alert.Update(price, isAbove)
	if err := a.repo.Update(ctx, alert); err != nil {
		return nil, err
	}
	if err := a.reload(ctx); err != nil {
		return nil, err
	}
	return alert, nil
}

func (a *App) DeleteAlert(ctx context.Context, username string, id uint) error {
	if err := a.repo.Deactivate(ctx, username, id); err != nil {
		return err
	}
	return a.reload(ctx)
}

func (a *App) ListAlerts(ctx context.Context, username string) ([]domain.Alert, error) {
	return a.repo.List(ctx, username)
}

func (a *App) SetStatus(ctx context.Context, username string, id uint, active bool) (*domain.Alert, error) {
	if err := a.repo.SetStatus(ctx, username, id, active); err != nil {
		return nil, err
	}
	alert, err := a.repo.FindByID(ctx, username, id)
	if err != nil {
		return nil, err
	}
	if err := a.reload(ctx); err != nil {
		return nil, err
	}
	return alert, nil
}

type priceMessage struct {
	Symbol string `json:"s"`
	Price  string `json:"p"`
	Time   int64  `json:"t"`
}

func (a *App) StartPriceWorker() {
	if err := a.reload(context.Background()); err != nil {
		slog.Error("load active alerts", "error", err)
	}
	go func() {
		for {
			pubsub := a.redis.Subscribe(context.Background(), "market:prices")
			for message := range pubsub.Channel() {
				var price priceMessage
				if err := json.Unmarshal([]byte(message.Payload), &price); err != nil {
					slog.Error("alert price decode", "error", err)
					continue
				}
				value, err := strconv.ParseFloat(price.Price, 32)
				if err != nil {
					continue
				}
				a.evaluate(context.Background(), strings.ToUpper(price.Symbol), float32(value))
			}
			_ = pubsub.Close()
			slog.Warn("alert price subscription closed; reconnecting")
			time.Sleep(time.Second)
		}
	}()
}

func (a *App) evaluate(ctx context.Context, coin string, price float32) {
	a.mu.RLock()
	alerts := append([]domain.Alert(nil), a.active[coin]...)
	a.mu.RUnlock()
	for _, alert := range alerts {
		matched := alert.IsAbove && price >= alert.Price || !alert.IsAbove && price <= alert.Price
		if !matched {
			continue
		}
		at := time.Now()
		claimed, err := a.repo.MarkTriggered(ctx, alert.ID, price, at)
		if err != nil || !claimed {
			continue
		}
		a.remove(alert.ID)
		message := fmt.Sprintf("%s %s $%.4f hedefini gerçekleştirdi: $%.4f", alert.Coin, map[bool]string{true: "üzerine çıktı", false: "altına düştü"}[alert.IsAbove], alert.Price, price)
		if err := a.broker.SendMessage(message); err != nil {
			slog.Error("alert rabbit notification", "alert", alert.ID, "error", err)
		}
		notificationText := fmt.Sprintf("%s · $%.4f → $%.4f", alert.Coin, alert.Price, price)
		payload, _ := json.Marshal(map[string]any{"username": alert.Username, "tip": "PRICE_ALERT", "coin": notificationText, "gorsel": "", "olusturma_zamani": at.Unix()})
		if err := a.redis.ZAdd(ctx, "kripto:grafikler", redis.Z{Score: float64(at.Add(24 * time.Hour).Unix()), Member: payload}).Err(); err != nil {
			slog.Error("alert app notification", "alert", alert.ID, "error", err)
		}
		slog.Info("price alert triggered", "alert", alert.ID, "username", alert.Username, "coin", coin, "price", price)
	}
}

func (a *App) reload(ctx context.Context) error {
	alerts, err := a.repo.Active(ctx)
	if err != nil {
		return err
	}
	next := make(map[string][]domain.Alert)
	for _, alert := range alerts {
		coin := strings.ToUpper(alert.Coin)
		next[coin] = append(next[coin], alert)
	}
	a.mu.Lock()
	a.active = next
	a.mu.Unlock()
	return nil
}

func (a *App) remove(id uint) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for coin, alerts := range a.active {
		kept := alerts[:0]
		for _, alert := range alerts {
			if alert.ID != id {
				kept = append(kept, alert)
			}
		}
		if len(kept) == 0 {
			delete(a.active, coin)
		} else {
			a.active[coin] = kept
		}
	}
}
