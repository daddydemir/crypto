package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	baseUrl     string
	redisClient *redis.Client
}

var assets = []string{"btcusdt", "ethusdt", "solusdt", "bnbusdt", "xrpusdt", "trxusdt", "dogeusdt", "adausdt",
	"bchusdt", "linkusdt", "xlmusdt", "hbarusdt", "ltcusdt", "avaxusdt", "suiusdt", "zecusdt", "tonusdt", "shibusdt",
	"wlfiusdt", "paxgusdt", "dotusdt", "uniusdt", "taousdt", "asterusdt", "skyusdt", "aaveusdt", "nearusdt", "pepeusdt",
	"icpusdt", "etcusdt", "ondousdt", "wldusdt", "polusdt", "atomusdt", "enausdt", "qntusdt", "algousdt", "morphousdt",
	"aptusdt", "filusdt", "renderusdt", "trumpusdt", "pumpusdt", "vetusdt", "zrousdt", "arbusdt", "nexousdt", "jupusdt",
	"dcrusdt", "stxusdt",
}

func NewClient(url string, redisClient *redis.Client) *Client {
	return &Client{
		baseUrl:     url,
		redisClient: redisClient,
	}
}

func (c *Client) Fetch() {

	baseUrl := c.baseUrl + strings.Join(mapAssetsToStreams(assets), "/")

	conn, _, err := websocket.DefaultDialer.Dial(baseUrl, nil)
	if err != nil {
		slog.Error("Failed to connect to websocket", "error", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message from websocket", "error", err)
			break
		}
		var ticker Ticker
		if err = json.Unmarshal(message, &ticker); err != nil {
			slog.Error("Failed to unmarshal message from websocket", "error", err)
			continue
		}

		ticker.Time = time.Now().Unix()
		ticker.Stream = ticker.Stream[:4]
		ticker.Data.Symbol = ticker.Data.Symbol[:len(ticker.Data.Symbol)-4]

		payload := fmt.Sprintf(`{"s": "%s", "p": "%s", "t": %d}`, ticker.Data.Symbol, ticker.Data.Price, ticker.Time)
		c.publish(payload)
	}
}

func mapAssetsToStreams(assets []string) []string {
	streams := make([]string, len(assets))
	for i, a := range assets {
		streams[i] = a + "@aggTrade"
	}
	return streams
}

type Ticker struct {
	Stream string `json:"stream"`
	Time   int64  `json:"time"`
	Data   struct {
		Symbol string `json:"s"`
		Price  string `json:"p"`
	} `json:"data"`
}

func (c *Client) publish(message string) {

	err := c.redisClient.Publish(context.Background(), "market:prices", message).Err()
	if err != nil {
		slog.Error("publish", "error", err)
	}
}
