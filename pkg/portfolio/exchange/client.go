package exchange

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/pkg/portfolio/domain"
)

type Client struct {
	httpClient  *http.Client
	clockMu     sync.RWMutex
	clockOffset map[string]time.Duration
	clockSynced map[string]time.Time
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 20 * time.Second}, clockOffset: make(map[string]time.Duration), clockSynced: make(map[string]time.Time)}
}

type trade struct {
	ID              json.Number `json:"id"`
	TradeID         json.Number `json:"tradeId"`
	OrderID         json.Number `json:"orderId"`
	Price           string      `json:"price"`
	Qty             string      `json:"qty"`
	QuoteQty        string      `json:"quoteQty"`
	Commission      string      `json:"commission"`
	CommissionAsset string      `json:"commissionAsset"`
	Time            json.Number `json:"time"`
	IsBuyer         any         `json:"isBuyer"`
}
type trResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List []trade `json:"list"`
	} `json:"data"`
}

type btcturkTrade struct {
	ID                json.Number `json:"id"`
	Timestamp         json.Number `json:"timestamp"`
	Amount            string      `json:"amount"`
	PreciseAmount     float64     `json:"preciseAmount"`
	Fee               string      `json:"fee"`
	Tax               string      `json:"tax"`
	Price             string      `json:"price"`
	NumeratorSymbol   string      `json:"numeratorSymbol"`
	DenominatorSymbol string      `json:"denominatorSymbol"`
	OrderType         string      `json:"orderType"`
	OrderID           json.Number `json:"orderId"`
}

type btcturkResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Code    int            `json:"code"`
	Data    []btcturkTrade `json:"data"`
}

func (c *Client) Fetch(ctx context.Context, source, baseAsset, quoteAsset string, start, end time.Time) ([]domain.Transaction, error) {
	source = strings.ToUpper(source)
	baseAsset = strings.ToUpper(strings.TrimSpace(baseAsset))
	quoteAsset = strings.ToUpper(strings.TrimSpace(quoteAsset))
	if source != "BINANCE" && source != "BINANCE_TR" && source != "BTCTURK" {
		return nil, fmt.Errorf("unsupported exchange")
	}
	if !end.After(start) {
		return nil, fmt.Errorf("valid date range is required")
	}
	if source != "BTCTURK" && (baseAsset == "" || quoteAsset == "") {
		return nil, fmt.Errorf("pair is required")
	}
	if source == "BTCTURK" && (baseAsset == "") != (quoteAsset == "") {
		return nil, fmt.Errorf("BTCTURK için işlem çiftinin iki varlığı da girilmeli veya ikisi de boş bırakılmalı")
	}
	if end.Sub(start) > 90*24*time.Hour {
		return nil, fmt.Errorf("date range cannot exceed 90 days")
	}
	if source == "BINANCE" {
		if err := c.ensureBinanceClock(ctx); err != nil {
			return nil, err
		}
	}
	if source == "BTCTURK" {
		transactions, err := c.fetchBTCTurk(ctx, baseAsset, quoteAsset, start, end)
		if err != nil {
			return nil, err
		}
		return mergeFills(transactions), nil
	}
	result := make([]domain.Transaction, 0)
	for cursor := start; cursor.Before(end); {
		chunkEnd := cursor.Add(24*time.Hour - time.Millisecond)
		if chunkEnd.After(end) {
			chunkEnd = end
		}
		items, err := c.fetchChunk(ctx, source, baseAsset, quoteAsset, cursor, chunkEnd)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
		cursor = chunkEnd.Add(time.Millisecond)
	}
	return mergeFills(result), nil
}

func (c *Client) fetchBTCTurk(ctx context.Context, base, quote string, start, end time.Time) ([]domain.Transaction, error) {
	apiKey, secret := setting("BTCTURK_API_KEY"), setting("BTCTURK_SECRET_KEY")
	if apiKey == "" || secret == "" {
		return nil, fmt.Errorf("BTCTURK API credentials are not configured")
	}
	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("BTCTURK_SECRET_KEY geçerli bir Base64 değeri değil")
	}
	baseURL := "https://api.btcturk.com"
	if custom := setting("BTCTURK_BASE_URL"); custom != "" {
		baseURL = strings.TrimRight(custom, "/")
	}
	if err := c.ensureClock(ctx, "BTCTURK", baseURL+"/api/v2/server/time"); err != nil {
		return nil, err
	}
	params := url.Values{
		"startDate": {strconv.FormatInt(start.UnixMilli(), 10)},
		"endDate":   {strconv.FormatInt(end.UnixMilli(), 10)},
		"type":      {"buy", "sell"},
	}
	if base != "" && quote != "" {
		params.Set("pairSymbol", base+quote)
	}
	nonce := strconv.FormatInt(c.requestTime("BTCTURK").UnixMilli(), 10)
	mac := hmac.New(sha256.New, secretBytes)
	_, _ = mac.Write([]byte(apiKey + nonce))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/v1/users/transactions/trade?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-PCK", apiKey)
	req.Header.Set("X-Stamp", nonce)
	req.Header.Set("X-Signature", signature)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, exchangeError("BTCTURK", resp.StatusCode, body)
	}
	var response btcturkResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("BTCTURK yanıtı okunamadı: %w", err)
	}
	if !response.Success {
		return nil, fmt.Errorf("BTCTURK API hatası (%d): %s", response.Code, fallbackMessage(response.Message))
	}
	result := make([]domain.Transaction, 0, len(response.Data))
	for _, item := range response.Data {
		result = append(result, convertBTCTurk(item, base, quote))
	}
	return result, nil
}

func convertBTCTurk(item btcturkTrade, requestedBase, requestedQuote string) domain.Transaction {
	base := strings.ToUpper(item.NumeratorSymbol)
	quote := strings.ToUpper(item.DenominatorSymbol)
	if base == "" {
		base = requestedBase
	}
	if quote == "" {
		quote = requestedQuote
	}
	amount := math.Abs(item.PreciseAmount)
	if amount == 0 {
		amount, _ = strconv.ParseFloat(item.Amount, 64)
		amount = math.Abs(amount)
	}
	price, _ := strconv.ParseFloat(item.Price, 64)
	fee, _ := strconv.ParseFloat(item.Fee, 64)
	tax, _ := strconv.ParseFloat(item.Tax, 64)
	millis, _ := item.Timestamp.Int64()
	tradeTime := time.UnixMilli(millis)
	transaction := domain.Transaction{
		BaseAsset: base, QuoteAsset: quote, Platform: "BTCTURK", Source: "BTCTURK",
		ExternalTradeID: item.ID.String(), ExchangeOrderID: item.OrderID.String(),
		FeeAmount: math.Abs(fee) + math.Abs(tax), FeeAsset: quote,
		TradedAt: tradeTime, LastFillAt: tradeTime,
	}
	quoteAmount := amount * price
	if strings.EqualFold(item.OrderType, "buy") {
		transaction.TransactionType = domain.TransactionTypeBuy
		transaction.ReceivedAsset, transaction.ReceivedAmount = base, amount
		transaction.SpentAsset, transaction.SpentAmount = quote, quoteAmount
	} else {
		transaction.TransactionType = domain.TransactionTypeSell
		transaction.ReceivedAsset, transaction.ReceivedAmount = quote, quoteAmount
		transaction.SpentAsset, transaction.SpentAmount = base, amount
	}
	if quote == "USD" || quote == "USDT" || quote == "USDC" {
		transaction.USDValue = quoteAmount
	}
	return transaction
}

func (c *Client) fetchChunk(ctx context.Context, source, base, quote string, start, end time.Time) ([]domain.Transaction, error) {
	apiKey, secret, baseURL, endpoint, symbol := setting(source+"_API_KEY"), setting(source+"_SECRET_KEY"), "https://api.binance.com", "/api/v3/myTrades", base+quote
	if source == "BINANCE_TR" {
		baseURL, endpoint, symbol = "https://www.binance.tr", "/open/v1/orders/trades", base+"_"+quote
	}
	if custom := setting(source + "_BASE_URL"); custom != "" {
		baseURL = strings.TrimRight(custom, "/")
	}
	if apiKey == "" || secret == "" {
		return nil, fmt.Errorf("%s API credentials are not configured", source)
	}
	params := url.Values{"symbol": {symbol}, "startTime": {strconv.FormatInt(start.UnixMilli(), 10)}, "endTime": {strconv.FormatInt(end.UnixMilli(), 10)}, "limit": {"1000"}, "recvWindow": {"5000"}, "timestamp": {strconv.FormatInt(c.requestTime(source).UnixMilli(), 10)}}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(params.Encode()))
	params.Set("signature", hex.EncodeToString(mac.Sum(nil)))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MBX-APIKEY", apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, exchangeError(source, resp.StatusCode, body)
	}
	var raw []trade
	if source == "BINANCE_TR" {
		var wrapped trResponse
		if err = json.Unmarshal(body, &wrapped); err != nil {
			return nil, err
		}
		if wrapped.Code != 0 {
			return nil, fmt.Errorf("Binance TR hatası: %s", fallbackMessage(wrapped.Msg))
		}
		raw = wrapped.Data.List
	} else if err = json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	result := make([]domain.Transaction, 0, len(raw))
	for _, item := range raw {
		result = append(result, convert(item, source, base, quote))
	}
	return result, nil
}

func (c *Client) requestTime(source string) time.Time {
	c.clockMu.RLock()
	offset := c.clockOffset[source]
	c.clockMu.RUnlock()
	return time.Now().Add(offset)
}

func (c *Client) ensureBinanceClock(ctx context.Context) error {
	baseURL := "https://api.binance.com"
	if custom := setting("BINANCE_BASE_URL"); custom != "" {
		baseURL = strings.TrimRight(custom, "/")
	}
	return c.ensureClock(ctx, "BINANCE", baseURL+"/api/v3/time")
}

func (c *Client) ensureClock(ctx context.Context, source, endpoint string) error {
	c.clockMu.RLock()
	fresh := time.Since(c.clockSynced[source]) < 5*time.Minute
	c.clockMu.RUnlock()
	if fresh {
		return nil
	}
	requestStarted := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s saat bilgisi alınamadı: %w", source, err)
	}
	defer resp.Body.Close()
	var response struct {
		ServerTime int64 `json:"serverTime"`
	}
	if resp.StatusCode != http.StatusOK || json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&response) != nil || response.ServerTime == 0 {
		return fmt.Errorf("%s saat bilgisi alınamadı", source)
	}
	localMidpoint := requestStarted.Add(time.Since(requestStarted) / 2)
	c.clockMu.Lock()
	c.clockOffset[source] = time.UnixMilli(response.ServerTime).Sub(localMidpoint)
	c.clockSynced[source] = time.Now()
	c.clockMu.Unlock()
	return nil
}

func exchangeError(source string, status int, body []byte) error {
	var response struct {
		Code    any    `json:"code"`
		Msg     string `json:"msg"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(body, &response)
	message := response.Msg
	if message == "" {
		message = response.Message
	}
	if message == "" {
		message = strings.TrimSpace(string(body))
	}
	if message == "" {
		message = http.StatusText(status)
	}
	return fmt.Errorf("%s API hatası (%d): %s", strings.ReplaceAll(source, "_", " "), status, message)
}

func fallbackMessage(message string) string {
	if strings.TrimSpace(message) == "" {
		return "Borsa beklenmeyen bir hata döndürdü"
	}
	return message
}

func setting(key string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return strings.TrimSpace(config.Get(key))
}

func convert(item trade, source, base, quote string) domain.Transaction {
	qty, _ := strconv.ParseFloat(item.Qty, 64)
	quoteQty, _ := strconv.ParseFloat(item.QuoteQty, 64)
	fee, _ := strconv.ParseFloat(item.Commission, 64)
	millis, _ := item.Time.Int64()
	id := item.ID.String()
	if source == "BINANCE_TR" {
		id = item.TradeID.String()
	}
	isBuyer := item.IsBuyer == true || item.IsBuyer == float64(1)
	tradeTime := time.UnixMilli(millis)
	t := domain.Transaction{BaseAsset: base, QuoteAsset: quote, Platform: strings.ReplaceAll(source, "_", " "), Source: source, ExternalTradeID: id, ExchangeOrderID: item.OrderID.String(), FeeAmount: fee, FeeAsset: item.CommissionAsset, TradedAt: tradeTime, LastFillAt: tradeTime}
	if isBuyer {
		t.TransactionType = domain.TransactionTypeBuy
		t.ReceivedAsset = base
		t.ReceivedAmount = qty
		t.SpentAsset = quote
		t.SpentAmount = quoteQty
	} else {
		t.TransactionType = domain.TransactionTypeSell
		t.ReceivedAsset = quote
		t.ReceivedAmount = quoteQty
		t.SpentAsset = base
		t.SpentAmount = qty
	}
	if quote == "USD" || quote == "USDT" || quote == "USDC" {
		t.USDValue = quoteQty
	}
	return t
}

const fallbackMergeWindow = 5 * time.Minute

func mergeFills(transactions []domain.Transaction) []domain.Transaction {
	if len(transactions) < 2 {
		return transactions
	}
	sort.SliceStable(transactions, func(i, j int) bool { return transactions[i].TradedAt.Before(transactions[j].TradedAt) })
	merged := make([]domain.Transaction, 0, len(transactions))
	orderIndexes := make(map[string]int)
	for _, transaction := range transactions {
		if transaction.ExchangeOrderID != "" {
			orderKey := transaction.Source + ":" + transaction.ExchangeOrderID
			if index, ok := orderIndexes[orderKey]; ok {
				mergeInto(&merged[index], transaction)
				continue
			}
			transaction.ExternalTradeID = "order:" + transaction.ExchangeOrderID
			orderIndexes[orderKey] = len(merged)
			merged = append(merged, transaction)
			continue
		}

		if len(merged) > 0 && canFallbackMerge(merged[len(merged)-1], transaction) {
			mergeInto(&merged[len(merged)-1], transaction)
			continue
		}
		merged = append(merged, transaction)
	}
	return merged
}

func canFallbackMerge(current, next domain.Transaction) bool {
	lastFillAt := current.LastFillAt
	if lastFillAt.IsZero() {
		lastFillAt = current.TradedAt
	}
	return current.ExchangeOrderID == "" &&
		current.Source == next.Source &&
		current.BaseAsset == next.BaseAsset &&
		current.QuoteAsset == next.QuoteAsset &&
		current.TransactionType == next.TransactionType &&
		current.FeeAsset == next.FeeAsset &&
		next.TradedAt.Sub(lastFillAt) >= 0 &&
		next.TradedAt.Sub(lastFillAt) <= fallbackMergeWindow
}

func mergeInto(target *domain.Transaction, fill domain.Transaction) {
	target.ReceivedAmount += fill.ReceivedAmount
	target.SpentAmount += fill.SpentAmount
	target.USDValue += fill.USDValue
	if target.FeeAsset == fill.FeeAsset {
		target.FeeAmount += fill.FeeAmount
	}
	if fill.TradedAt.After(target.LastFillAt) {
		target.LastFillAt = fill.TradedAt
	}
}
