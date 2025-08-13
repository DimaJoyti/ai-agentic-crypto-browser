package exchanges

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/shopspring/decimal"
)

// MarketDataService provides unified real-time market data across exchanges
type MarketDataService struct {
	logger      *observability.Logger
	manager     *Manager
	subscribers map[string]*SubscriberGroup
	config      MarketDataConfig

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// MarketDataConfig contains configuration for market data service
type MarketDataConfig struct {
	BufferSize        int           `json:"buffer_size"`
	EnableHeartbeat   bool          `json:"enable_heartbeat"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	MaxSubscribers    int           `json:"max_subscribers"`
	EnableAggregation bool          `json:"enable_aggregation"`
}

// SubscriberGroup manages subscribers for different data types
type SubscriberGroup struct {
	tickerSubs    map[string][]chan *AggregatedTicker
	orderBookSubs map[string][]chan *AggregatedOrderBook
	tradeSubs     map[string][]chan *AggregatedTrade
	mu            sync.RWMutex
}

// Aggregated data types

// AggregatedTicker represents ticker data aggregated across exchanges
type AggregatedTicker struct {
	Symbol           string                        `json:"symbol"`
	BestBid          *ExchangePrice                `json:"best_bid"`
	BestAsk          *ExchangePrice                `json:"best_ask"`
	WeightedAvgPrice *ExchangePrice                `json:"weighted_avg_price"`
	Volume24h        map[string]*common.TickerData `json:"volume_24h"`
	PriceChange24h   *ExchangePrice                `json:"price_change_24h"`
	Exchanges        map[string]*common.TickerData `json:"exchanges"`
	Timestamp        time.Time                     `json:"timestamp"`
}

// AggregatedTrade represents trade data aggregated across exchanges
type AggregatedTrade struct {
	Symbol    string                       `json:"symbol"`
	Trades    map[string]*common.TradeData `json:"trades"`
	BestPrice *ExchangePrice               `json:"best_price"`
	Volume    map[string]string            `json:"volume"`
	Timestamp time.Time                    `json:"timestamp"`
}

// ExchangePrice represents a price with exchange information
type ExchangePrice struct {
	Price    string `json:"price"`
	Exchange string `json:"exchange"`
	Volume   string `json:"volume,omitempty"`
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(logger *observability.Logger, manager *Manager, config MarketDataConfig) *MarketDataService {
	if config.BufferSize == 0 {
		config.BufferSize = 1000
	}
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 30 * time.Second
	}
	if config.MaxSubscribers == 0 {
		config.MaxSubscribers = 100
	}

	return &MarketDataService{
		logger:      logger,
		manager:     manager,
		subscribers: make(map[string]*SubscriberGroup),
		config:      config,
		stopChan:    make(chan struct{}),
	}
}

// Start begins the market data service
func (mds *MarketDataService) Start(ctx context.Context) error {
	mds.mu.Lock()
	defer mds.mu.Unlock()

	if mds.isRunning {
		return nil
	}

	mds.logger.Info(ctx, "Starting market data service", map[string]interface{}{
		"buffer_size":        mds.config.BufferSize,
		"enable_aggregation": mds.config.EnableAggregation,
		"max_subscribers":    mds.config.MaxSubscribers,
	})

	mds.isRunning = true

	// Start heartbeat if enabled
	if mds.config.EnableHeartbeat {
		mds.wg.Add(1)
		go mds.heartbeatMonitor(ctx)
	}

	mds.logger.Info(ctx, "Market data service started", nil)

	return nil
}

// Stop stops the market data service
func (mds *MarketDataService) Stop(ctx context.Context) error {
	mds.mu.Lock()
	defer mds.mu.Unlock()

	if !mds.isRunning {
		return nil
	}

	mds.logger.Info(ctx, "Stopping market data service", nil)

	close(mds.stopChan)
	mds.wg.Wait()

	// Close all subscriber channels
	for _, group := range mds.subscribers {
		group.mu.Lock()
		for _, subs := range group.tickerSubs {
			for _, ch := range subs {
				close(ch)
			}
		}
		for _, subs := range group.orderBookSubs {
			for _, ch := range subs {
				close(ch)
			}
		}
		for _, subs := range group.tradeSubs {
			for _, ch := range subs {
				close(ch)
			}
		}
		group.mu.Unlock()
	}

	mds.isRunning = false

	mds.logger.Info(ctx, "Market data service stopped", nil)

	return nil
}

// SubscribeToTicker subscribes to aggregated ticker updates for a symbol
func (mds *MarketDataService) SubscribeToTicker(ctx context.Context, symbol string) (<-chan *AggregatedTicker, error) {
	mds.mu.Lock()
	defer mds.mu.Unlock()

	if !mds.isRunning {
		return nil, fmt.Errorf("market data service is not running")
	}

	ch := make(chan *AggregatedTicker, mds.config.BufferSize)

	// Get or create subscriber group
	if mds.subscribers[symbol] == nil {
		mds.subscribers[symbol] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *AggregatedTicker),
			orderBookSubs: make(map[string][]chan *AggregatedOrderBook),
			tradeSubs:     make(map[string][]chan *AggregatedTrade),
		}
	}

	group := mds.subscribers[symbol]
	group.mu.Lock()
	if group.tickerSubs[symbol] == nil {
		group.tickerSubs[symbol] = make([]chan *AggregatedTicker, 0)
	}
	group.tickerSubs[symbol] = append(group.tickerSubs[symbol], ch)
	group.mu.Unlock()

	// Start aggregating ticker data for this symbol
	if err := mds.startTickerAggregation(ctx, symbol); err != nil {
		return nil, fmt.Errorf("failed to start ticker aggregation: %w", err)
	}

	mds.logger.Info(ctx, "Subscribed to ticker", map[string]interface{}{
		"symbol":      symbol,
		"subscribers": len(group.tickerSubs[symbol]),
	})

	return ch, nil
}

// SubscribeToOrderBook subscribes to aggregated order book updates for a symbol
func (mds *MarketDataService) SubscribeToOrderBook(ctx context.Context, symbol string) (<-chan *AggregatedOrderBook, error) {
	mds.mu.Lock()
	defer mds.mu.Unlock()

	if !mds.isRunning {
		return nil, fmt.Errorf("market data service is not running")
	}

	ch := make(chan *AggregatedOrderBook, mds.config.BufferSize)

	// Get or create subscriber group
	if mds.subscribers[symbol] == nil {
		mds.subscribers[symbol] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *AggregatedTicker),
			orderBookSubs: make(map[string][]chan *AggregatedOrderBook),
			tradeSubs:     make(map[string][]chan *AggregatedTrade),
		}
	}

	group := mds.subscribers[symbol]
	group.mu.Lock()
	if group.orderBookSubs[symbol] == nil {
		group.orderBookSubs[symbol] = make([]chan *AggregatedOrderBook, 0)
	}
	group.orderBookSubs[symbol] = append(group.orderBookSubs[symbol], ch)
	group.mu.Unlock()

	// Start aggregating order book data for this symbol
	if err := mds.startOrderBookAggregation(ctx, symbol); err != nil {
		return nil, fmt.Errorf("failed to start order book aggregation: %w", err)
	}

	return ch, nil
}

// SubscribeToTrades subscribes to aggregated trade updates for a symbol
func (mds *MarketDataService) SubscribeToTrades(ctx context.Context, symbol string) (<-chan *AggregatedTrade, error) {
	mds.mu.Lock()
	defer mds.mu.Unlock()

	if !mds.isRunning {
		return nil, fmt.Errorf("market data service is not running")
	}

	ch := make(chan *AggregatedTrade, mds.config.BufferSize)

	// Get or create subscriber group
	if mds.subscribers[symbol] == nil {
		mds.subscribers[symbol] = &SubscriberGroup{
			tickerSubs:    make(map[string][]chan *AggregatedTicker),
			orderBookSubs: make(map[string][]chan *AggregatedOrderBook),
			tradeSubs:     make(map[string][]chan *AggregatedTrade),
		}
	}

	group := mds.subscribers[symbol]
	group.mu.Lock()
	if group.tradeSubs[symbol] == nil {
		group.tradeSubs[symbol] = make([]chan *AggregatedTrade, 0)
	}
	group.tradeSubs[symbol] = append(group.tradeSubs[symbol], ch)
	group.mu.Unlock()

	// Start aggregating trade data for this symbol
	if err := mds.startTradeAggregation(ctx, symbol); err != nil {
		return nil, fmt.Errorf("failed to start trade aggregation: %w", err)
	}

	return ch, nil
}

// GetAggregatedTicker gets current aggregated ticker for a symbol
func (mds *MarketDataService) GetAggregatedTicker(ctx context.Context, symbol string) (*AggregatedTicker, error) {
	exchanges := mds.manager.GetAllExchanges()

	tickers := make(map[string]*common.TickerData)
	var bestBid, bestAsk *ExchangePrice

	for name, exchange := range exchanges {
		ticker, err := exchange.GetTicker(ctx, symbol)
		if err != nil {
			mds.logger.Error(ctx, "Failed to get ticker", err, map[string]interface{}{
				"exchange": name,
				"symbol":   symbol,
			})
			continue
		}

		tickers[name] = ticker

		// Find best bid and ask
		if bestBid == nil {
			bestBid = &ExchangePrice{
				Price:    ticker.BidPrice.String(),
				Exchange: name,
				Volume:   ticker.BidQty.String(),
			}
		} else {
			bestBidPrice, err := decimal.NewFromString(bestBid.Price)
			if err == nil && ticker.BidPrice.GreaterThan(bestBidPrice) {
				bestBid = &ExchangePrice{
					Price:    ticker.BidPrice.String(),
					Exchange: name,
					Volume:   ticker.BidQty.String(),
				}
			}
		}

		if bestAsk == nil {
			bestAsk = &ExchangePrice{
				Price:    ticker.AskPrice.String(),
				Exchange: name,
				Volume:   ticker.AskQty.String(),
			}
		} else {
			bestAskPrice, err := decimal.NewFromString(bestAsk.Price)
			if err == nil && ticker.AskPrice.LessThan(bestAskPrice) {
				bestAsk = &ExchangePrice{
					Price:    ticker.AskPrice.String(),
					Exchange: name,
					Volume:   ticker.AskQty.String(),
				}
			}
		}
	}

	return &AggregatedTicker{
		Symbol:    symbol,
		BestBid:   bestBid,
		BestAsk:   bestAsk,
		Exchanges: tickers,
		Timestamp: time.Now(),
	}, nil
}

// Private methods for aggregation

// startTickerAggregation starts ticker data aggregation for a symbol
func (mds *MarketDataService) startTickerAggregation(ctx context.Context, symbol string) error {
	exchanges := mds.manager.GetAllExchanges()

	for name, exchange := range exchanges {
		mds.wg.Add(1)
		go mds.aggregateTickerData(ctx, name, exchange, symbol)
	}

	return nil
}

// startOrderBookAggregation starts order book data aggregation for a symbol
func (mds *MarketDataService) startOrderBookAggregation(ctx context.Context, symbol string) error {
	exchanges := mds.manager.GetAllExchanges()

	for name, exchange := range exchanges {
		mds.wg.Add(1)
		go mds.aggregateOrderBookData(ctx, name, exchange, symbol)
	}

	return nil
}

// startTradeAggregation starts trade data aggregation for a symbol
func (mds *MarketDataService) startTradeAggregation(ctx context.Context, symbol string) error {
	exchanges := mds.manager.GetAllExchanges()

	for name, exchange := range exchanges {
		mds.wg.Add(1)
		go mds.aggregateTradeData(ctx, name, exchange, symbol)
	}

	return nil
}

// aggregateTickerData aggregates ticker data from an exchange
func (mds *MarketDataService) aggregateTickerData(ctx context.Context, exchangeName string, exchange common.ExchangeClient, symbol string) {
	defer mds.wg.Done()

	tickerCh, err := exchange.SubscribeToTicker(ctx, symbol)
	if err != nil {
		mds.logger.Error(ctx, "Failed to subscribe to ticker", err, map[string]interface{}{
			"exchange": exchangeName,
			"symbol":   symbol,
		})
		return
	}

	for {
		select {
		case <-mds.stopChan:
			return
		case ticker := <-tickerCh:
			if ticker != nil {
				mds.processTickerUpdate(ctx, exchangeName, ticker)
			}
		}
	}
}

// aggregateOrderBookData aggregates order book data from an exchange
func (mds *MarketDataService) aggregateOrderBookData(ctx context.Context, exchangeName string, exchange common.ExchangeClient, symbol string) {
	defer mds.wg.Done()

	orderBookCh, err := exchange.SubscribeToOrderBook(ctx, symbol)
	if err != nil {
		mds.logger.Error(ctx, "Failed to subscribe to order book", err, map[string]interface{}{
			"exchange": exchangeName,
			"symbol":   symbol,
		})
		return
	}

	for {
		select {
		case <-mds.stopChan:
			return
		case orderBook := <-orderBookCh:
			if orderBook != nil {
				mds.processOrderBookUpdate(ctx, exchangeName, orderBook)
			}
		}
	}
}

// aggregateTradeData aggregates trade data from an exchange
func (mds *MarketDataService) aggregateTradeData(ctx context.Context, exchangeName string, exchange common.ExchangeClient, symbol string) {
	defer mds.wg.Done()

	tradeCh, err := exchange.SubscribeToTrades(ctx, symbol)
	if err != nil {
		mds.logger.Error(ctx, "Failed to subscribe to trades", err, map[string]interface{}{
			"exchange": exchangeName,
			"symbol":   symbol,
		})
		return
	}

	for {
		select {
		case <-mds.stopChan:
			return
		case trade := <-tradeCh:
			if trade != nil {
				mds.processTradeUpdate(ctx, exchangeName, trade)
			}
		}
	}
}

// processTickerUpdate processes ticker updates and sends aggregated data
func (mds *MarketDataService) processTickerUpdate(ctx context.Context, exchangeName string, ticker *common.TickerData) {
	// This would implement the logic to aggregate ticker data and send to subscribers
	// For now, this is a placeholder
}

// processOrderBookUpdate processes order book updates and sends aggregated data
func (mds *MarketDataService) processOrderBookUpdate(ctx context.Context, exchangeName string, orderBook *common.OrderBookData) {
	// This would implement the logic to aggregate order book data and send to subscribers
	// For now, this is a placeholder
}

// processTradeUpdate processes trade updates and sends aggregated data
func (mds *MarketDataService) processTradeUpdate(ctx context.Context, exchangeName string, trade *common.TradeData) {
	// This would implement the logic to aggregate trade data and send to subscribers
	// For now, this is a placeholder
}

// heartbeatMonitor monitors connection health
func (mds *MarketDataService) heartbeatMonitor(ctx context.Context) {
	defer mds.wg.Done()

	ticker := time.NewTicker(mds.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-mds.stopChan:
			return
		case <-ticker.C:
			mds.checkConnectionHealth(ctx)
		}
	}
}

// checkConnectionHealth checks the health of all exchange connections
func (mds *MarketDataService) checkConnectionHealth(ctx context.Context) {
	exchanges := mds.manager.GetAllExchanges()

	for name, exchange := range exchanges {
		if !exchange.IsConnected() {
			mds.logger.Warn(ctx, "Exchange connection lost", map[string]interface{}{
				"exchange": name,
			})
		}
	}
}
