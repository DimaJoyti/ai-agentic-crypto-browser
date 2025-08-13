package common

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// ExchangeClient defines the unified interface for all cryptocurrency exchanges
type ExchangeClient interface {
	// Connection Management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	IsConnected() bool
	GetExchangeName() string

	// Market Data
	GetTicker(ctx context.Context, symbol string) (*TickerData, error)
	GetOrderBook(ctx context.Context, symbol string, limit int) (*OrderBookData, error)
	GetRecentTrades(ctx context.Context, symbol string, limit int) ([]*TradeData, error)
	GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*KlineData, error)

	// WebSocket Streaming
	SubscribeToTicker(ctx context.Context, symbol string) (<-chan *TickerData, error)
	SubscribeToOrderBook(ctx context.Context, symbol string) (<-chan *OrderBookData, error)
	SubscribeToTrades(ctx context.Context, symbol string) (<-chan *TradeData, error)
	SubscribeToUserData(ctx context.Context) (<-chan *UserDataUpdate, error)
	UnsubscribeAll(ctx context.Context) error

	// Order Management
	PlaceOrder(ctx context.Context, order *OrderRequest) (*OrderResponse, error)
	CancelOrder(ctx context.Context, symbol string, orderID string) (*OrderResponse, error)
	CancelAllOrders(ctx context.Context, symbol string) ([]*OrderResponse, error)
	GetOrder(ctx context.Context, symbol string, orderID string) (*OrderResponse, error)
	GetOpenOrders(ctx context.Context, symbol string) ([]*OrderResponse, error)
	GetOrderHistory(ctx context.Context, symbol string, limit int) ([]*OrderResponse, error)

	// Account Management
	GetAccountInfo(ctx context.Context) (*AccountInfo, error)
	GetBalances(ctx context.Context) ([]*Balance, error)
	GetTradingFees(ctx context.Context, symbol string) (*TradingFees, error)

	// Advanced Order Types
	PlaceStopLossOrder(ctx context.Context, req *StopLossOrderRequest) (*OrderResponse, error)
	PlaceTakeProfitOrder(ctx context.Context, req *TakeProfitOrderRequest) (*OrderResponse, error)
	PlaceIcebergOrder(ctx context.Context, req *IcebergOrderRequest) (*OrderResponse, error)
	PlaceTWAPOrder(ctx context.Context, req *TWAPOrderRequest) (*OrderResponse, error)

	// Risk Management
	GetPositionRisk(ctx context.Context, symbol string) (*PositionRisk, error)
	GetMaxOrderSize(ctx context.Context, symbol string) (decimal.Decimal, error)
	ValidateOrder(ctx context.Context, order *OrderRequest) error

	// Performance Metrics
	GetLatencyStats() *LatencyStats
	GetConnectionStats() *ConnectionStats
}

// TickerData represents 24hr ticker statistics
type TickerData struct {
	Symbol             string          `json:"symbol"`
	PriceChange        decimal.Decimal `json:"price_change"`
	PriceChangePercent decimal.Decimal `json:"price_change_percent"`
	WeightedAvgPrice   decimal.Decimal `json:"weighted_avg_price"`
	PrevClosePrice     decimal.Decimal `json:"prev_close_price"`
	LastPrice          decimal.Decimal `json:"last_price"`
	LastQty            decimal.Decimal `json:"last_qty"`
	BidPrice           decimal.Decimal `json:"bid_price"`
	BidQty             decimal.Decimal `json:"bid_qty"`
	AskPrice           decimal.Decimal `json:"ask_price"`
	AskQty             decimal.Decimal `json:"ask_qty"`
	OpenPrice          decimal.Decimal `json:"open_price"`
	HighPrice          decimal.Decimal `json:"high_price"`
	LowPrice           decimal.Decimal `json:"low_price"`
	Volume             decimal.Decimal `json:"volume"`
	QuoteVolume        decimal.Decimal `json:"quote_volume"`
	OpenTime           time.Time       `json:"open_time"`
	CloseTime          time.Time       `json:"close_time"`
	Count              int64           `json:"count"`
	Exchange           string          `json:"exchange"`
	Timestamp          time.Time       `json:"timestamp"`
}

// OrderBookData represents order book data
type OrderBookData struct {
	Symbol       string       `json:"symbol"`
	Bids         []PriceLevel `json:"bids"`
	Asks         []PriceLevel `json:"asks"`
	Timestamp    time.Time    `json:"timestamp"`
	Exchange     string       `json:"exchange"`
	LastUpdateID int64        `json:"last_update_id"`
}

// PriceLevel represents a price level in the order book
type PriceLevel struct {
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
}

// TradeData represents a trade
type TradeData struct {
	ID        string          `json:"id"`
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Side      string          `json:"side"`
	Timestamp time.Time       `json:"timestamp"`
	Exchange  string          `json:"exchange"`
	IsBuyer   bool            `json:"is_buyer"`
}

// KlineData represents candlestick data
type KlineData struct {
	Symbol      string          `json:"symbol"`
	OpenTime    time.Time       `json:"open_time"`
	CloseTime   time.Time       `json:"close_time"`
	Open        decimal.Decimal `json:"open"`
	High        decimal.Decimal `json:"high"`
	Low         decimal.Decimal `json:"low"`
	Close       decimal.Decimal `json:"close"`
	Volume      decimal.Decimal `json:"volume"`
	QuoteVolume decimal.Decimal `json:"quote_volume"`
	TradeCount  int64           `json:"trade_count"`
	Exchange    string          `json:"exchange"`
}

// OrderRequest represents an order request
type OrderRequest struct {
	Symbol        string                 `json:"symbol"`
	Side          OrderSide              `json:"side"`
	Type          OrderType              `json:"type"`
	Quantity      decimal.Decimal        `json:"quantity"`
	Price         decimal.Decimal        `json:"price,omitempty"`
	StopPrice     decimal.Decimal        `json:"stop_price,omitempty"`
	TimeInForce   TimeInForce            `json:"time_in_force"`
	ClientOrderID string                 `json:"client_order_id,omitempty"`
	IcebergQty    decimal.Decimal        `json:"iceberg_qty,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// OrderResponse represents an order response
type OrderResponse struct {
	OrderID         string          `json:"order_id"`
	ClientOrderID   string          `json:"client_order_id"`
	Symbol          string          `json:"symbol"`
	Side            OrderSide       `json:"side"`
	Type            OrderType       `json:"type"`
	Quantity        decimal.Decimal `json:"quantity"`
	Price           decimal.Decimal `json:"price"`
	StopPrice       decimal.Decimal `json:"stop_price,omitempty"`
	TimeInForce     TimeInForce     `json:"time_in_force"`
	Status          OrderStatus     `json:"status"`
	FilledQty       decimal.Decimal `json:"filled_qty"`
	RemainingQty    decimal.Decimal `json:"remaining_qty"`
	AvgFillPrice    decimal.Decimal `json:"avg_fill_price"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commission_asset"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Exchange        string          `json:"exchange"`
	LatencyMicros   int64           `json:"latency_micros"`
}

// UserDataUpdate represents user data stream updates
type UserDataUpdate struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Exchange  string      `json:"exchange"`
}

// AccountInfo represents account information
type AccountInfo struct {
	AccountType                 string          `json:"account_type"`
	CanTrade                    bool            `json:"can_trade"`
	CanWithdraw                 bool            `json:"can_withdraw"`
	CanDeposit                  bool            `json:"can_deposit"`
	UpdateTime                  time.Time       `json:"update_time"`
	TotalWalletBalance          decimal.Decimal `json:"total_wallet_balance"`
	TotalUnrealizedPnL          decimal.Decimal `json:"total_unrealized_pnl"`
	TotalMarginBalance          decimal.Decimal `json:"total_margin_balance"`
	TotalPositionInitialMargin  decimal.Decimal `json:"total_position_initial_margin"`
	TotalOpenOrderInitialMargin decimal.Decimal `json:"total_open_order_initial_margin"`
	Exchange                    string          `json:"exchange"`
}

// Balance represents account balance
type Balance struct {
	Asset  string          `json:"asset"`
	Free   decimal.Decimal `json:"free"`
	Locked decimal.Decimal `json:"locked"`
	Total  decimal.Decimal `json:"total"`
}

// TradingFees represents trading fees for a symbol
type TradingFees struct {
	Symbol          string          `json:"symbol"`
	MakerCommission decimal.Decimal `json:"maker_commission"`
	TakerCommission decimal.Decimal `json:"taker_commission"`
}

// Advanced Order Types

// StopLossOrderRequest represents a stop-loss order request
type StopLossOrderRequest struct {
	Symbol        string          `json:"symbol"`
	Side          OrderSide       `json:"side"`
	Quantity      decimal.Decimal `json:"quantity"`
	StopPrice     decimal.Decimal `json:"stop_price"`
	Price         decimal.Decimal `json:"price,omitempty"` // For stop-limit orders
	TimeInForce   TimeInForce     `json:"time_in_force"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
}

// TakeProfitOrderRequest represents a take-profit order request
type TakeProfitOrderRequest struct {
	Symbol        string          `json:"symbol"`
	Side          OrderSide       `json:"side"`
	Quantity      decimal.Decimal `json:"quantity"`
	StopPrice     decimal.Decimal `json:"stop_price"`
	Price         decimal.Decimal `json:"price,omitempty"` // For take-profit-limit orders
	TimeInForce   TimeInForce     `json:"time_in_force"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
}

// IcebergOrderRequest represents an iceberg order request
type IcebergOrderRequest struct {
	Symbol        string          `json:"symbol"`
	Side          OrderSide       `json:"side"`
	Quantity      decimal.Decimal `json:"quantity"`
	Price         decimal.Decimal `json:"price"`
	IcebergQty    decimal.Decimal `json:"iceberg_qty"`
	TimeInForce   TimeInForce     `json:"time_in_force"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
}

// TWAPOrderRequest represents a TWAP (Time-Weighted Average Price) order request
type TWAPOrderRequest struct {
	Symbol        string          `json:"symbol"`
	Side          OrderSide       `json:"side"`
	Quantity      decimal.Decimal `json:"quantity"`
	Duration      time.Duration   `json:"duration"`
	Interval      time.Duration   `json:"interval"`
	PriceLimit    decimal.Decimal `json:"price_limit,omitempty"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
}

// PositionRisk represents position risk information
type PositionRisk struct {
	Symbol           string          `json:"symbol"`
	PositionAmt      decimal.Decimal `json:"position_amt"`
	EntryPrice       decimal.Decimal `json:"entry_price"`
	MarkPrice        decimal.Decimal `json:"mark_price"`
	UnrealizedPnL    decimal.Decimal `json:"unrealized_pnl"`
	LiquidationPrice decimal.Decimal `json:"liquidation_price"`
	Leverage         decimal.Decimal `json:"leverage"`
	MaxNotionalValue decimal.Decimal `json:"max_notional_value"`
	MarginType       string          `json:"margin_type"`
	IsolatedMargin   decimal.Decimal `json:"isolated_margin"`
	IsAutoAddMargin  bool            `json:"is_auto_add_margin"`
	PositionSide     string          `json:"position_side"`
	UpdateTime       time.Time       `json:"update_time"`
}

// Performance Metrics

// LatencyStats represents latency statistics
type LatencyStats struct {
	AvgLatencyMicros int64     `json:"avg_latency_micros"`
	MinLatencyMicros int64     `json:"min_latency_micros"`
	MaxLatencyMicros int64     `json:"max_latency_micros"`
	P50LatencyMicros int64     `json:"p50_latency_micros"`
	P95LatencyMicros int64     `json:"p95_latency_micros"`
	P99LatencyMicros int64     `json:"p99_latency_micros"`
	SampleCount      int64     `json:"sample_count"`
	LastUpdated      time.Time `json:"last_updated"`
}

// ConnectionStats represents connection statistics
type ConnectionStats struct {
	IsConnected       bool      `json:"is_connected"`
	ConnectedSince    time.Time `json:"connected_since"`
	ReconnectCount    int64     `json:"reconnect_count"`
	LastReconnectTime time.Time `json:"last_reconnect_time"`
	MessagesSent      int64     `json:"messages_sent"`
	MessagesReceived  int64     `json:"messages_received"`
	BytesSent         int64     `json:"bytes_sent"`
	BytesReceived     int64     `json:"bytes_received"`
	ErrorCount        int64     `json:"error_count"`
	LastError         string    `json:"last_error"`
	LastErrorTime     time.Time `json:"last_error_time"`
}

// Enums

// OrderSide represents order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType represents order type
type OrderType string

const (
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
	OrderTypeIceberg         OrderType = "ICEBERG"
	OrderTypeTWAP            OrderType = "TWAP"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusPendingCancel   OrderStatus = "PENDING_CANCEL"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// TimeInForce represents order time in force
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Canceled
	TimeInForceIOC TimeInForce = "IOC" // Immediate Or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill Or Kill
	TimeInForceGTX TimeInForce = "GTX" // Good Till Crossing (Post Only)
)
