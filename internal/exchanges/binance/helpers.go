package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/shopspring/decimal"
)

// RateLimiter methods

// Allow checks if a request is allowed under the rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastRefill) >= rl.refillRate {
		rl.tokens = rl.maxTokens
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// Client helper methods

// testConnectivity tests API connectivity
func (c *Client) testConnectivity(ctx context.Context) error {
	endpoint := "/api/v3/ping"
	_, err := c.makeRequest(ctx, "GET", endpoint, nil, false)
	return err
}

// makeRequest makes an HTTP request to Binance API
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values, signed bool) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}

	// Add timestamp for signed requests
	if signed {
		params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	// Build URL
	baseURL := c.config.BaseURL
	fullURL := baseURL + endpoint

	var body io.Reader
	var queryString string

	if method == "GET" || method == "DELETE" {
		if len(params) > 0 {
			queryString = params.Encode()
			fullURL += "?" + queryString
		}
	} else {
		queryString = params.Encode()
		body = strings.NewReader(queryString)
	}

	// Sign request if required
	if signed {
		signature := c.signRequest(queryString)
		if method == "GET" || method == "DELETE" {
			fullURL += "&signature=" + signature
		} else {
			queryString += "&signature=" + signature
			body = strings.NewReader(queryString)
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.config.APIKey != "" {
		req.Header.Set("X-MBX-APIKEY", c.config.APIKey)
	}

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.connectionStats.ErrorCount++
		c.connectionStats.LastError = err.Error()
		c.connectionStats.LastErrorTime = time.Now()
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Update connection stats
	c.connectionStats.MessagesSent++
	c.connectionStats.BytesSent += int64(len(queryString))

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Update connection stats
	c.connectionStats.MessagesReceived++
	c.connectionStats.BytesReceived += int64(len(respBody))

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// signRequest signs a request with HMAC SHA256
func (c *Client) signRequest(queryString string) string {
	mac := hmac.New(sha256.New, []byte(c.config.SecretKey))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

// buildOrderParams builds order parameters for API request
func (c *Client) buildOrderParams(order *common.OrderRequest) url.Values {
	params := url.Values{}
	params.Set("symbol", strings.ToUpper(order.Symbol))
	params.Set("side", string(order.Side))
	params.Set("type", c.convertOrderType(order.Type))
	params.Set("quantity", order.Quantity.String())

	if !order.Price.IsZero() {
		params.Set("price", order.Price.String())
	}

	if !order.StopPrice.IsZero() {
		params.Set("stopPrice", order.StopPrice.String())
	}

	if order.TimeInForce != "" {
		params.Set("timeInForce", string(order.TimeInForce))
	}

	if order.ClientOrderID != "" {
		params.Set("newClientOrderId", order.ClientOrderID)
	}

	if !order.IcebergQty.IsZero() {
		params.Set("icebergQty", order.IcebergQty.String())
	}

	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	return params
}

// convertOrderType converts common order type to Binance order type
func (c *Client) convertOrderType(orderType common.OrderType) string {
	switch orderType {
	case common.OrderTypeMarket:
		return "MARKET"
	case common.OrderTypeLimit:
		return "LIMIT"
	case common.OrderTypeStopLoss:
		return "STOP_LOSS"
	case common.OrderTypeStopLossLimit:
		return "STOP_LOSS_LIMIT"
	case common.OrderTypeTakeProfit:
		return "TAKE_PROFIT"
	case common.OrderTypeTakeProfitLimit:
		return "TAKE_PROFIT_LIMIT"
	case common.OrderTypeLimitMaker:
		return "LIMIT_MAKER"
	default:
		return "LIMIT"
	}
}

// convertOrderStatus converts Binance order status to common order status
func (c *Client) convertOrderStatus(status string) common.OrderStatus {
	switch status {
	case "NEW":
		return common.OrderStatusNew
	case "PARTIALLY_FILLED":
		return common.OrderStatusPartiallyFilled
	case "FILLED":
		return common.OrderStatusFilled
	case "CANCELED":
		return common.OrderStatusCanceled
	case "PENDING_CANCEL":
		return common.OrderStatusPendingCancel
	case "REJECTED":
		return common.OrderStatusRejected
	case "EXPIRED":
		return common.OrderStatusExpired
	default:
		return common.OrderStatusNew
	}
}

// convertOrderSide converts Binance order side to common order side
func (c *Client) convertOrderSide(side string) common.OrderSide {
	switch side {
	case "BUY":
		return common.OrderSideBuy
	case "SELL":
		return common.OrderSideSell
	default:
		return common.OrderSideBuy
	}
}

// Data conversion methods

// convertTickerData converts Binance ticker to common ticker data
func (c *Client) convertTickerData(ticker *BinanceTickerResponse) *common.TickerData {
	return &common.TickerData{
		Symbol:             ticker.Symbol,
		PriceChange:        ParseDecimal(ticker.PriceChange),
		PriceChangePercent: ParseDecimal(ticker.PriceChangePercent),
		WeightedAvgPrice:   ParseDecimal(ticker.WeightedAvgPrice),
		PrevClosePrice:     ParseDecimal(ticker.PrevClosePrice),
		LastPrice:          ParseDecimal(ticker.LastPrice),
		LastQty:            ParseDecimal(ticker.LastQty),
		BidPrice:           ParseDecimal(ticker.BidPrice),
		BidQty:             ParseDecimal(ticker.BidQty),
		AskPrice:           ParseDecimal(ticker.AskPrice),
		AskQty:             ParseDecimal(ticker.AskQty),
		OpenPrice:          ParseDecimal(ticker.OpenPrice),
		HighPrice:          ParseDecimal(ticker.HighPrice),
		LowPrice:           ParseDecimal(ticker.LowPrice),
		Volume:             ParseDecimal(ticker.Volume),
		QuoteVolume:        ParseDecimal(ticker.QuoteVolume),
		OpenTime:           ParseTime(ticker.OpenTime),
		CloseTime:          ParseTime(ticker.CloseTime),
		Count:              ticker.Count,
		Exchange:           "binance",
		Timestamp:          time.Now(),
	}
}

// convertOrderBookData converts Binance order book to common order book data
func (c *Client) convertOrderBookData(orderBook *BinanceOrderBookResponse, symbol string) *common.OrderBookData {
	bids := make([]common.PriceLevel, len(orderBook.Bids))
	for i, bid := range orderBook.Bids {
		if len(bid) >= 2 {
			bids[i] = common.PriceLevel{
				Price:    ParseDecimal(bid[0]),
				Quantity: ParseDecimal(bid[1]),
			}
		}
	}

	asks := make([]common.PriceLevel, len(orderBook.Asks))
	for i, ask := range orderBook.Asks {
		if len(ask) >= 2 {
			asks[i] = common.PriceLevel{
				Price:    ParseDecimal(ask[0]),
				Quantity: ParseDecimal(ask[1]),
			}
		}
	}

	return &common.OrderBookData{
		Symbol:       symbol,
		Bids:         bids,
		Asks:         asks,
		Timestamp:    time.Now(),
		Exchange:     "binance",
		LastUpdateID: orderBook.LastUpdateId,
	}
}

// convertTradeData converts Binance trade to common trade data
func (c *Client) convertTradeData(trade *BinanceTradeResponse, symbol string) *common.TradeData {
	side := "buy"
	if trade.IsBuyerMaker {
		side = "sell"
	}

	return &common.TradeData{
		ID:        strconv.FormatInt(trade.ID, 10),
		Symbol:    symbol,
		Price:     ParseDecimal(trade.Price),
		Quantity:  ParseDecimal(trade.Qty),
		Side:      side,
		Timestamp: ParseTime(trade.Time),
		Exchange:  "binance",
		IsBuyer:   !trade.IsBuyerMaker,
	}
}

// convertKlineData converts Binance kline to common kline data
func (c *Client) convertKlineData(kline []interface{}, symbol string) *common.KlineData {
	if len(kline) < 11 {
		return &common.KlineData{Symbol: symbol, Exchange: "binance"}
	}

	openTime, _ := kline[0].(float64)
	closeTime, _ := kline[6].(float64)
	tradeCount, _ := kline[8].(float64)

	return &common.KlineData{
		Symbol:      symbol,
		OpenTime:    ParseTime(int64(openTime)),
		CloseTime:   ParseTime(int64(closeTime)),
		Open:        ParseDecimal(kline[1].(string)),
		High:        ParseDecimal(kline[2].(string)),
		Low:         ParseDecimal(kline[3].(string)),
		Close:       ParseDecimal(kline[4].(string)),
		Volume:      ParseDecimal(kline[5].(string)),
		QuoteVolume: ParseDecimal(kline[7].(string)),
		TradeCount:  int64(tradeCount),
		Exchange:    "binance",
	}
}

// convertOrderResponse converts Binance order response to common order response
func (c *Client) convertOrderResponse(order *BinanceOrderResponse) *common.OrderResponse {
	return &common.OrderResponse{
		OrderID:         strconv.FormatInt(order.OrderId, 10),
		ClientOrderID:   order.ClientOrderId,
		Symbol:          order.Symbol,
		Side:            c.convertOrderSide(order.Side),
		Type:            c.convertCommonOrderType(order.Type),
		Quantity:        ParseDecimal(order.OrigQty),
		Price:           ParseDecimal(order.Price),
		StopPrice:       ParseDecimal(order.StopPrice),
		TimeInForce:     c.convertTimeInForce(order.TimeInForce),
		Status:          c.convertOrderStatus(order.Status),
		FilledQty:       ParseDecimal(order.ExecutedQty),
		RemainingQty:    ParseDecimal(order.OrigQty).Sub(ParseDecimal(order.ExecutedQty)),
		AvgFillPrice:    c.calculateAvgFillPrice(order),
		Commission:      decimal.Zero, // Would need separate API call
		CommissionAsset: "",
		CreatedAt:       ParseTime(order.Time),
		UpdatedAt:       ParseTime(order.UpdateTime),
		Exchange:        "binance",
	}
}

// convertCommonOrderType converts Binance order type to common order type
func (c *Client) convertCommonOrderType(orderType string) common.OrderType {
	switch orderType {
	case "MARKET":
		return common.OrderTypeMarket
	case "LIMIT":
		return common.OrderTypeLimit
	case "STOP_LOSS":
		return common.OrderTypeStopLoss
	case "STOP_LOSS_LIMIT":
		return common.OrderTypeStopLossLimit
	case "TAKE_PROFIT":
		return common.OrderTypeTakeProfit
	case "TAKE_PROFIT_LIMIT":
		return common.OrderTypeTakeProfitLimit
	case "LIMIT_MAKER":
		return common.OrderTypeLimitMaker
	default:
		return common.OrderTypeLimit
	}
}

// convertTimeInForce converts Binance time in force to common time in force
func (c *Client) convertTimeInForce(tif string) common.TimeInForce {
	switch tif {
	case "GTC":
		return common.TimeInForceGTC
	case "IOC":
		return common.TimeInForceIOC
	case "FOK":
		return common.TimeInForceFOK
	case "GTX":
		return common.TimeInForceGTX
	default:
		return common.TimeInForceGTC
	}
}

// calculateAvgFillPrice calculates average fill price
func (c *Client) calculateAvgFillPrice(order *BinanceOrderResponse) decimal.Decimal {
	executedQty := ParseDecimal(order.ExecutedQty)
	if executedQty.IsZero() {
		return decimal.Zero
	}

	cummulativeQuoteQty := ParseDecimal(order.CummulativeQuoteQty)
	if cummulativeQuoteQty.IsZero() {
		return ParseDecimal(order.Price)
	}

	return cummulativeQuoteQty.Div(executedQty)
}

// updateLatencyStats updates latency statistics
func (c *Client) updateLatencyStats(latencyMicros int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.latencyStats.SampleCount == 0 {
		c.latencyStats.MinLatencyMicros = latencyMicros
		c.latencyStats.MaxLatencyMicros = latencyMicros
		c.latencyStats.AvgLatencyMicros = latencyMicros
	} else {
		if latencyMicros < c.latencyStats.MinLatencyMicros {
			c.latencyStats.MinLatencyMicros = latencyMicros
		}
		if latencyMicros > c.latencyStats.MaxLatencyMicros {
			c.latencyStats.MaxLatencyMicros = latencyMicros
		}

		// Update running average
		totalLatency := c.latencyStats.AvgLatencyMicros * c.latencyStats.SampleCount
		c.latencyStats.AvgLatencyMicros = (totalLatency + latencyMicros) / (c.latencyStats.SampleCount + 1)
	}

	c.latencyStats.SampleCount++
	c.latencyStats.LastUpdated = time.Now()
}
