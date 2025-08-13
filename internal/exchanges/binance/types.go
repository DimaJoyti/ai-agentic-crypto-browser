package binance

import (
	"time"

	"github.com/shopspring/decimal"
)

// Binance API Response Types

// BinanceTickerResponse represents Binance 24hr ticker response
type BinanceTickerResponse struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	Count              int64  `json:"count"`
}

// BinanceOrderBookResponse represents Binance order book response
type BinanceOrderBookResponse struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// BinanceTradeResponse represents Binance trade response
type BinanceTradeResponse struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

// BinanceOrderResponse represents Binance order response
type BinanceOrderResponse struct {
	Symbol                  string `json:"symbol"`
	OrderId                 int64  `json:"orderId"`
	OrderListId             int64  `json:"orderListId"`
	ClientOrderId           string `json:"clientOrderId"`
	TransactTime            int64  `json:"transactTime"`
	Price                   string `json:"price"`
	OrigQty                 string `json:"origQty"`
	ExecutedQty             string `json:"executedQty"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Status                  string `json:"status"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	Side                    string `json:"side"`
	StopPrice               string `json:"stopPrice,omitempty"`
	IcebergQty              string `json:"icebergQty,omitempty"`
	Time                    int64  `json:"time,omitempty"`
	UpdateTime              int64  `json:"updateTime,omitempty"`
	IsWorking               bool   `json:"isWorking,omitempty"`
	WorkingTime             int64  `json:"workingTime,omitempty"`
	OrigQuoteOrderQty       string `json:"origQuoteOrderQty,omitempty"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode,omitempty"`
}

// BinanceAccountResponse represents Binance account response
type BinanceAccountResponse struct {
	MakerCommission  int64                    `json:"makerCommission"`
	TakerCommission  int64                    `json:"takerCommission"`
	BuyerCommission  int64                    `json:"buyerCommission"`
	SellerCommission int64                    `json:"sellerCommission"`
	CanTrade         bool                     `json:"canTrade"`
	CanWithdraw      bool                     `json:"canWithdraw"`
	CanDeposit       bool                     `json:"canDeposit"`
	UpdateTime       int64                    `json:"updateTime"`
	AccountType      string                   `json:"accountType"`
	Balances         []BinanceBalanceResponse `json:"balances"`
	Permissions      []string                 `json:"permissions"`
}

// BinanceBalanceResponse represents Binance balance response
type BinanceBalanceResponse struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// BinanceExchangeInfoResponse represents Binance exchange info response
type BinanceExchangeInfoResponse struct {
	Timezone   string                    `json:"timezone"`
	ServerTime int64                     `json:"serverTime"`
	RateLimits []BinanceRateLimitInfo    `json:"rateLimits"`
	Symbols    []BinanceSymbolInfo       `json:"symbols"`
}

// BinanceRateLimitInfo represents rate limit information
type BinanceRateLimitInfo struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// BinanceSymbolInfo represents symbol information
type BinanceSymbolInfo struct {
	Symbol                     string                   `json:"symbol"`
	Status                     string                   `json:"status"`
	BaseAsset                  string                   `json:"baseAsset"`
	BaseAssetPrecision         int                      `json:"baseAssetPrecision"`
	QuoteAsset                 string                   `json:"quoteAsset"`
	QuoteAssetPrecision        int                      `json:"quoteAssetPrecision"`
	OrderTypes                 []string                 `json:"orderTypes"`
	IcebergAllowed             bool                     `json:"icebergAllowed"`
	OcoAllowed                 bool                     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool                     `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop          bool                     `json:"allowTrailingStop"`
	CancelReplaceAllowed       bool                     `json:"cancelReplaceAllowed"`
	IsSpotTradingAllowed       bool                     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                     `json:"isMarginTradingAllowed"`
	Filters                    []BinanceSymbolFilter    `json:"filters"`
	Permissions                []string                 `json:"permissions"`
}

// BinanceSymbolFilter represents symbol filter
type BinanceSymbolFilter struct {
	FilterType          string `json:"filterType"`
	MinPrice            string `json:"minPrice,omitempty"`
	MaxPrice            string `json:"maxPrice,omitempty"`
	TickSize            string `json:"tickSize,omitempty"`
	MultiplierUp        string `json:"multiplierUp,omitempty"`
	MultiplierDown      string `json:"multiplierDown,omitempty"`
	AvgPriceMins        int    `json:"avgPriceMins,omitempty"`
	MinQty              string `json:"minQty,omitempty"`
	MaxQty              string `json:"maxQty,omitempty"`
	StepSize            string `json:"stepSize,omitempty"`
	MinNotional         string `json:"minNotional,omitempty"`
	ApplyToMarket       bool   `json:"applyToMarket,omitempty"`
	Limit               int    `json:"limit,omitempty"`
	MaxNumOrders        int    `json:"maxNumOrders,omitempty"`
	MaxNumAlgoOrders    int    `json:"maxNumAlgoOrders,omitempty"`
	MaxNumIcebergOrders int    `json:"maxNumIcebergOrders,omitempty"`
	MaxPosition         string `json:"maxPosition,omitempty"`
}

// WebSocket Stream Types

// BinanceWSTickerData represents WebSocket ticker data
type BinanceWSTickerData struct {
	EventType          string `json:"e"`
	EventTime          int64  `json:"E"`
	Symbol             string `json:"s"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	WeightedAvgPrice   string `json:"w"`
	FirstTradePrice    string `json:"x"`
	LastPrice          string `json:"c"`
	LastQty            string `json:"Q"`
	BestBidPrice       string `json:"b"`
	BestBidQty         string `json:"B"`
	BestAskPrice       string `json:"a"`
	BestAskQty         string `json:"A"`
	OpenPrice          string `json:"o"`
	HighPrice          string `json:"h"`
	LowPrice           string `json:"l"`
	Volume             string `json:"v"`
	QuoteVolume        string `json:"q"`
	OpenTime           int64  `json:"O"`
	CloseTime          int64  `json:"C"`
	FirstTradeId       int64  `json:"F"`
	LastTradeId        int64  `json:"L"`
	TradeCount         int64  `json:"n"`
}

// BinanceWSDepthData represents WebSocket depth data
type BinanceWSDepthData struct {
	EventType        string     `json:"e"`
	EventTime        int64      `json:"E"`
	Symbol           string     `json:"s"`
	FirstUpdateId    int64      `json:"U"`
	FinalUpdateId    int64      `json:"u"`
	Bids             [][]string `json:"b"`
	Asks             [][]string `json:"a"`
}

// BinanceWSTradeData represents WebSocket trade data
type BinanceWSTradeData struct {
	EventType         string `json:"e"`
	EventTime         int64  `json:"E"`
	Symbol            string `json:"s"`
	TradeId           int64  `json:"t"`
	Price             string `json:"p"`
	Quantity          string `json:"q"`
	BuyerOrderId      int64  `json:"b"`
	SellerOrderId     int64  `json:"a"`
	TradeTime         int64  `json:"T"`
	IsBuyerMaker      bool   `json:"m"`
	Ignore            bool   `json:"M"`
}

// BinanceWSUserData represents WebSocket user data
type BinanceWSUserData struct {
	EventType string      `json:"e"`
	EventTime int64       `json:"E"`
	Data      interface{} `json:",inline"`
}

// BinanceWSExecutionReport represents order execution report
type BinanceWSExecutionReport struct {
	EventType                string `json:"e"`
	EventTime                int64  `json:"E"`
	Symbol                   string `json:"s"`
	ClientOrderId            string `json:"c"`
	Side                     string `json:"S"`
	OrderType                string `json:"o"`
	TimeInForce              string `json:"f"`
	OrderQuantity            string `json:"q"`
	OrderPrice               string `json:"p"`
	StopPrice                string `json:"P"`
	IcebergQuantity          string `json:"F"`
	OrderListId              int64  `json:"g"`
	OrigClientOrderId        string `json:"C"`
	CurrentExecutionType     string `json:"x"`
	CurrentOrderStatus       string `json:"X"`
	OrderRejectReason        string `json:"r"`
	OrderId                  int64  `json:"i"`
	LastExecutedQuantity     string `json:"l"`
	CumulativeFilledQuantity string `json:"z"`
	LastExecutedPrice        string `json:"L"`
	CommissionAmount         string `json:"n"`
	CommissionAsset          string `json:"N"`
	TransactionTime          int64  `json:"T"`
	TradeId                  int64  `json:"t"`
	Ignore1                  int64  `json:"I"`
	IsOrderOnBook            bool   `json:"w"`
	IsMaker                  bool   `json:"m"`
	Ignore2                  bool   `json:"M"`
	OrderCreationTime        int64  `json:"O"`
	CumulativeQuoteQty       string `json:"Z"`
	LastQuoteQty             string `json:"Y"`
	QuoteOrderQty            string `json:"Q"`
}

// Helper functions for decimal conversion

// ParseDecimal safely parses a string to decimal.Decimal
func ParseDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// ParseTime converts Unix timestamp to time.Time
func ParseTime(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000)
}
