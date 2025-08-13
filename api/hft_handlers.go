package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ai-agentic-browser/internal/exchanges/common"
	"github.com/ai-agentic-browser/internal/hft"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// HFT Engine API handlers

// handleHFTStart starts the HFT engine
func (s *APIServer) handleHFTStart(w http.ResponseWriter, r *http.Request) {
	if s.hftEngine == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "HFT engine not available")
		return
	}

	err := s.hftEngine.Start(r.Context())
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]string{
		"status":  "started",
		"message": "HFT engine started successfully",
	})

	// Broadcast status update
	s.BroadcastMessage("hft_status", map[string]interface{}{
		"running": true,
		"message": "HFT engine started",
	})
}

// handleHFTStop stops the HFT engine
func (s *APIServer) handleHFTStop(w http.ResponseWriter, r *http.Request) {
	if s.hftEngine == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "HFT engine not available")
		return
	}

	err := s.hftEngine.Stop(r.Context())
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]string{
		"status":  "stopped",
		"message": "HFT engine stopped successfully",
	})

	// Broadcast status update
	s.BroadcastMessage("hft_status", map[string]interface{}{
		"running": false,
		"message": "HFT engine stopped",
	})
}

// handleHFTStatus returns the current status of the HFT engine
func (s *APIServer) handleHFTStatus(w http.ResponseWriter, r *http.Request) {
	if s.hftEngine == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "HFT engine not available")
		return
	}

	status := map[string]interface{}{
		"running": s.hftEngine.IsRunning(),
		"uptime":  s.hftEngine.GetUptime(),
		"version": "1.0.0",
	}

	s.sendJSON(w, r, http.StatusOK, status)
}

// handleHFTMetrics returns performance metrics for the HFT engine
func (s *APIServer) handleHFTMetrics(w http.ResponseWriter, r *http.Request) {
	if s.hftEngine == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "HFT engine not available")
		return
	}

	metrics := s.hftEngine.GetMetrics()
	s.sendJSON(w, r, http.StatusOK, metrics)
}

// handleHFTConfig handles HFT engine configuration
func (s *APIServer) handleHFTConfig(w http.ResponseWriter, r *http.Request) {
	if s.hftEngine == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "HFT engine not available")
		return
	}

	switch r.Method {
	case "GET":
		config := s.hftEngine.GetConfig()
		s.sendJSON(w, r, http.StatusOK, config)

	case "PUT":
		var config hft.HFTConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			s.sendError(w, r, http.StatusBadRequest, "Invalid configuration format")
			return
		}

		if err := s.hftEngine.UpdateConfig(config); err != nil {
			s.sendError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		s.sendJSON(w, r, http.StatusOK, map[string]string{
			"message": "Configuration updated successfully",
		})
	}
}

// Market Data API handlers

// handleMarketTickers returns ticker data for all symbols
func (s *APIServer) handleMarketTickers(w http.ResponseWriter, r *http.Request) {
	if s.binanceService == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Binance service not available")
		return
	}

	// Get query parameters
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		exchange = "binance"
	}

	// Mock ticker data for demonstration
	tickers := []map[string]interface{}{
		{
			"symbol":         "BTCUSDT",
			"price":          "45123.45",
			"change_24h":     "1234.56",
			"change_percent": "2.81",
			"volume":         "12345.67",
			"high_24h":       "46000.00",
			"low_24h":        "44000.00",
			"exchange":       exchange,
		},
		{
			"symbol":         "ETHUSDT",
			"price":          "3012.34",
			"change_24h":     "-45.67",
			"change_percent": "-1.49",
			"volume":         "23456.78",
			"high_24h":       "3100.00",
			"low_24h":        "2950.00",
			"exchange":       exchange,
		},
	}

	s.sendJSON(w, r, http.StatusOK, tickers)
}

// handleMarketTicker returns ticker data for a specific symbol
func (s *APIServer) handleMarketTicker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	if s.binanceService == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Binance service not available")
		return
	}

	ticker, err := s.binanceService.GetTicker(r.Context(), symbol)
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, ticker)
}

// handleMarketOrderbook returns orderbook data for a symbol
func (s *APIServer) handleMarketOrderbook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Mock orderbook data
	orderbook := map[string]interface{}{
		"symbol": symbol,
		"bids": [][]string{
			{"45120.00", "0.1234"},
			{"45119.50", "0.2345"},
			{"45119.00", "0.3456"},
		},
		"asks": [][]string{
			{"45125.00", "0.1111"},
			{"45125.50", "0.2222"},
			{"45126.00", "0.3333"},
		},
		"timestamp": "2024-01-15T10:30:00Z",
		"limit":     limit, // Use the limit parameter
	}

	s.sendJSON(w, r, http.StatusOK, orderbook)
}

// handleMarketTrades returns recent trades for a symbol
func (s *APIServer) handleMarketTrades(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Mock trades data
	trades := []map[string]interface{}{
		{
			"id":        "12345",
			"price":     "45123.45",
			"quantity":  "0.1234",
			"timestamp": "2024-01-15T10:30:00Z",
			"side":      "buy",
		},
		{
			"id":        "12346",
			"price":     "45122.00",
			"quantity":  "0.2345",
			"timestamp": "2024-01-15T10:29:58Z",
			"side":      "sell",
		},
	}

	// Limit results
	if len(trades) > limit {
		trades = trades[:limit]
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"symbol": symbol,
		"trades": trades,
		"count":  len(trades),
	})
}

// Trading API handlers

// handleTradingOrders handles order management
func (s *APIServer) handleTradingOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetOrders(w, r)
	case "POST":
		s.handleCreateOrder(w, r)
	}
}

// handleGetOrders returns all orders
func (s *APIServer) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	symbol := r.URL.Query().Get("symbol")
	status := r.URL.Query().Get("status")

	// Mock orders data
	orders := []map[string]interface{}{
		{
			"id":         "order_1",
			"symbol":     "BTCUSDT",
			"side":       "BUY",
			"type":       "LIMIT",
			"quantity":   "0.1",
			"price":      "45000.00",
			"filled_qty": "0.05",
			"status":     "PARTIAL_FILL",
			"created_at": "2024-01-15T10:30:00Z",
			"exchange":   "binance",
		},
		{
			"id":         "order_2",
			"symbol":     "ETHUSDT",
			"side":       "SELL",
			"type":       "MARKET",
			"quantity":   "1.0",
			"price":      "3000.00",
			"filled_qty": "1.0",
			"status":     "FILLED",
			"created_at": "2024-01-15T10:25:00Z",
			"exchange":   "binance",
		},
	}

	// Apply filters
	var filteredOrders []map[string]interface{}
	for _, order := range orders {
		if symbol != "" && order["symbol"] != symbol {
			continue
		}
		if status != "" && order["status"] != status {
			continue
		}
		filteredOrders = append(filteredOrders, order)
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"orders": filteredOrders,
		"count":  len(filteredOrders),
	})
}

// handleCreateOrder creates a new order
func (s *APIServer) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest struct {
		Symbol   string `json:"symbol"`
		Side     string `json:"side"`
		Type     string `json:"type"`
		Quantity string `json:"quantity"`
		Price    string `json:"price,omitempty"`
		Exchange string `json:"exchange,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid order format")
		return
	}

	// Validate required fields
	if orderRequest.Symbol == "" || orderRequest.Side == "" ||
		orderRequest.Type == "" || orderRequest.Quantity == "" {
		s.sendError(w, r, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Mock order creation
	order := map[string]interface{}{
		"id":         "order_" + strconv.FormatInt(1642248600, 10),
		"symbol":     orderRequest.Symbol,
		"side":       orderRequest.Side,
		"type":       orderRequest.Type,
		"quantity":   orderRequest.Quantity,
		"price":      orderRequest.Price,
		"filled_qty": "0",
		"status":     "NEW",
		"created_at": "2024-01-15T10:30:00Z",
		"exchange":   orderRequest.Exchange,
	}

	s.sendJSON(w, r, http.StatusCreated, order)

	// Broadcast order update
	s.BroadcastMessage("order_created", order)
}

// handleTradingOrder handles individual order operations
func (s *APIServer) handleTradingOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	switch r.Method {
	case "GET":
		// Mock order data
		order := map[string]interface{}{
			"id":         orderID,
			"symbol":     "BTCUSDT",
			"side":       "BUY",
			"type":       "LIMIT",
			"quantity":   "0.1",
			"price":      "45000.00",
			"filled_qty": "0.05",
			"status":     "PARTIAL_FILL",
			"created_at": "2024-01-15T10:30:00Z",
			"exchange":   "binance",
		}

		s.sendJSON(w, r, http.StatusOK, order)

	case "DELETE":
		// Mock order cancellation
		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"id":      orderID,
			"status":  "CANCELED",
			"message": "Order canceled successfully",
		})

		// Broadcast order update
		s.BroadcastMessage("order_canceled", map[string]interface{}{
			"id":     orderID,
			"status": "CANCELED",
		})
	}
}

// handleTradingPositions returns current positions
func (s *APIServer) handleTradingPositions(w http.ResponseWriter, r *http.Request) {
	// Mock positions data
	positions := []map[string]interface{}{
		{
			"symbol":         "BTCUSDT",
			"size":           "0.5",
			"avg_price":      "44500.00",
			"current_price":  "45123.45",
			"unrealized_pnl": "311.73",
			"realized_pnl":   "0.00",
			"commission":     "2.25",
			"open_time":      "2024-01-15T09:30:00Z",
			"exchange":       "binance",
		},
		{
			"symbol":         "ETHUSDT",
			"size":           "-2.0",
			"avg_price":      "3050.00",
			"current_price":  "3012.34",
			"unrealized_pnl": "75.32",
			"realized_pnl":   "0.00",
			"commission":     "6.10",
			"open_time":      "2024-01-15T09:45:00Z",
			"exchange":       "binance",
		},
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"positions": positions,
		"count":     len(positions),
	})
}

// handleTradingSignals returns recent trading signals
func (s *APIServer) handleTradingSignals(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Mock signals data
	signals := []map[string]interface{}{
		{
			"id":          "signal_1",
			"symbol":      "BTCUSDT",
			"side":        "BUY",
			"order_type":  "MARKET",
			"quantity":    "0.01",
			"price":       "45123.45",
			"confidence":  0.85,
			"strategy_id": "market_making_1",
			"timestamp":   "2024-01-15T10:30:00Z",
		},
		{
			"id":          "signal_2",
			"symbol":      "ETHUSDT",
			"side":        "SELL",
			"order_type":  "LIMIT",
			"quantity":    "0.5",
			"price":       "3015.00",
			"confidence":  0.72,
			"strategy_id": "arbitrage_1",
			"timestamp":   "2024-01-15T10:29:45Z",
		},
	}

	// Limit results
	if len(signals) > limit {
		signals = signals[:limit]
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"signals": signals,
		"count":   len(signals),
	})
}

// Exchange Management Handlers

// handleGetExchanges returns all active exchanges
func (s *APIServer) handleGetExchanges(w http.ResponseWriter, r *http.Request) {
	if s.exchangeManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Exchange manager not available")
		return
	}

	exchanges := s.exchangeManager.GetActiveExchanges()

	exchangeInfo := make([]map[string]interface{}, len(exchanges))
	for i, name := range exchanges {
		exchange, err := s.exchangeManager.GetExchange(name)
		if err != nil {
			continue
		}

		stats := exchange.GetConnectionStats()
		latencyStats := exchange.GetLatencyStats()

		exchangeInfo[i] = map[string]interface{}{
			"name":             name,
			"connected":        exchange.IsConnected(),
			"connection_stats": stats,
			"latency_stats":    latencyStats,
		}
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"exchanges": exchangeInfo,
		"count":     len(exchangeInfo),
	})
}

// handleGetTicker gets ticker data for a symbol
func (s *APIServer) handleGetTicker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := strings.ToUpper(vars["symbol"])
	exchangeName := r.URL.Query().Get("exchange")

	if symbol == "" {
		s.sendError(w, r, http.StatusBadRequest, "Symbol is required")
		return
	}

	if s.exchangeManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Exchange manager not available")
		return
	}

	var exchange common.ExchangeClient
	var err error

	if exchangeName != "" {
		exchange, err = s.exchangeManager.GetExchange(exchangeName)
	} else {
		exchange, err = s.exchangeManager.GetDefaultExchange()
	}

	if err != nil {
		s.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ticker, err := exchange.GetTicker(r.Context(), symbol)
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, ticker)
}

// handleGetOrderBook gets order book data for a symbol
func (s *APIServer) handleGetOrderBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := strings.ToUpper(vars["symbol"])
	exchangeName := r.URL.Query().Get("exchange")
	limitStr := r.URL.Query().Get("limit")

	if symbol == "" {
		s.sendError(w, r, http.StatusBadRequest, "Symbol is required")
		return
	}

	limit := 20 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if s.exchangeManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Exchange manager not available")
		return
	}

	var exchange common.ExchangeClient
	var err error

	if exchangeName != "" {
		exchange, err = s.exchangeManager.GetExchange(exchangeName)
	} else {
		exchange, err = s.exchangeManager.GetDefaultExchange()
	}

	if err != nil {
		s.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	orderBook, err := exchange.GetOrderBook(r.Context(), symbol, limit)
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, orderBook)
}

// handleGetBestPrice gets the best price across all exchanges
func (s *APIServer) handleGetBestPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := strings.ToUpper(vars["symbol"])
	sideStr := r.URL.Query().Get("side")

	if symbol == "" {
		s.sendError(w, r, http.StatusBadRequest, "Symbol is required")
		return
	}

	if sideStr == "" {
		s.sendError(w, r, http.StatusBadRequest, "Side is required (BUY or SELL)")
		return
	}

	var side common.OrderSide
	switch strings.ToUpper(sideStr) {
	case "BUY":
		side = common.OrderSideBuy
	case "SELL":
		side = common.OrderSideSell
	default:
		s.sendError(w, r, http.StatusBadRequest, "Invalid side, must be BUY or SELL")
		return
	}

	if s.exchangeManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Exchange manager not available")
		return
	}

	bestPrice, err := s.exchangeManager.GetBestPrice(r.Context(), symbol, side)
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, bestPrice)
}

// Order Management Handlers

// OrderRequest represents an order submission request
type OrderRequest struct {
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`
	Type        string `json:"type"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price,omitempty"`
	StopPrice   string `json:"stop_price,omitempty"`
	TimeInForce string `json:"time_in_force,omitempty"`
	Exchange    string `json:"exchange,omitempty"`
	Strategy    string `json:"strategy,omitempty"`
}

// handleSubmitOrder submits an order for execution
func (s *APIServer) handleSubmitOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if s.orderManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Order manager not available")
		return
	}

	// Validate required fields
	if req.Symbol == "" {
		s.sendError(w, r, http.StatusBadRequest, "Symbol is required")
		return
	}
	if req.Side == "" {
		s.sendError(w, r, http.StatusBadRequest, "Side is required")
		return
	}
	if req.Quantity == "" {
		s.sendError(w, r, http.StatusBadRequest, "Quantity is required")
		return
	}

	// Parse quantity
	quantity, err := decimal.NewFromString(req.Quantity)
	if err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid quantity")
		return
	}

	// Parse price if provided
	var price decimal.Decimal
	if req.Price != "" {
		price, err = decimal.NewFromString(req.Price)
		if err != nil {
			s.sendError(w, r, http.StatusBadRequest, "Invalid price")
			return
		}
	}

	// Parse stop price if provided
	var stopPrice decimal.Decimal
	if req.StopPrice != "" {
		stopPrice, err = decimal.NewFromString(req.StopPrice)
		if err != nil {
			s.sendError(w, r, http.StatusBadRequest, "Invalid stop price")
			return
		}
	}

	// Convert side
	var side common.OrderSide
	switch strings.ToUpper(req.Side) {
	case "BUY":
		side = common.OrderSideBuy
	case "SELL":
		side = common.OrderSideSell
	default:
		s.sendError(w, r, http.StatusBadRequest, "Invalid side, must be BUY or SELL")
		return
	}

	// Convert order type
	var orderType common.OrderType
	switch strings.ToUpper(req.Type) {
	case "MARKET":
		orderType = common.OrderTypeMarket
	case "LIMIT":
		orderType = common.OrderTypeLimit
	case "STOP_LOSS":
		orderType = common.OrderTypeStopLoss
	case "STOP_LOSS_LIMIT":
		orderType = common.OrderTypeStopLossLimit
	case "TAKE_PROFIT":
		orderType = common.OrderTypeTakeProfit
	case "TAKE_PROFIT_LIMIT":
		orderType = common.OrderTypeTakeProfitLimit
	default:
		orderType = common.OrderTypeLimit
	}

	// Convert time in force
	var timeInForce common.TimeInForce
	switch strings.ToUpper(req.TimeInForce) {
	case "GTC":
		timeInForce = common.TimeInForceGTC
	case "IOC":
		timeInForce = common.TimeInForceIOC
	case "FOK":
		timeInForce = common.TimeInForceFOK
	default:
		timeInForce = common.TimeInForceGTC
	}

	// Create order request
	orderReq := &common.OrderRequest{
		Symbol:      strings.ToUpper(req.Symbol),
		Side:        side,
		Type:        orderType,
		Quantity:    quantity,
		Price:       price,
		StopPrice:   stopPrice,
		TimeInForce: timeInForce,
		Metadata:    make(map[string]interface{}),
	}

	// Add strategy if specified
	if req.Strategy != "" {
		orderReq.Metadata["strategy"] = req.Strategy
	}

	// Submit order
	managedOrder, err := s.orderManager.SubmitOrder(r.Context(), orderReq)
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusCreated, managedOrder)

	// Broadcast order update
	s.BroadcastMessage("order_submitted", map[string]interface{}{
		"order_id": managedOrder.ID.String(),
		"symbol":   managedOrder.OriginalRequest.Symbol,
		"side":     string(managedOrder.OriginalRequest.Side),
		"quantity": managedOrder.OriginalRequest.Quantity.String(),
		"status":   string(managedOrder.Status),
	})
}

// handleGetOrder gets a managed order by ID
func (s *APIServer) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]

	if orderIDStr == "" {
		s.sendError(w, r, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid order ID")
		return
	}

	if s.orderManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Order manager not available")
		return
	}

	order, err := s.orderManager.GetOrder(orderID)
	if err != nil {
		s.sendError(w, r, http.StatusNotFound, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, order)
}

// handleCancelOrder cancels a managed order
func (s *APIServer) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]

	if orderIDStr == "" {
		s.sendError(w, r, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		s.sendError(w, r, http.StatusBadRequest, "Invalid order ID")
		return
	}

	if s.orderManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Order manager not available")
		return
	}

	err = s.orderManager.CancelOrder(r.Context(), orderID)
	if err != nil {
		s.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]string{
		"status":  "cancelled",
		"message": "Order cancelled successfully",
	})

	// Broadcast order update
	s.BroadcastMessage("order_cancelled", map[string]interface{}{
		"order_id": orderID.String(),
	})
}

// Smart Order Routing handlers

// handleSORRoute handles order routing requests
func (s *APIServer) handleSORRoute(w http.ResponseWriter, r *http.Request) {
	if s.smartOrderRouter == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Smart Order Router not available")
		return
	}

	// Create SOR handlers instance
	sorHandlers := NewSORHandlers(s.smartOrderRouter, s.logger)
	sorHandlers.RouteOrder(w, r)
}

// handleSORExecute handles order execution requests
func (s *APIServer) handleSORExecute(w http.ResponseWriter, r *http.Request) {
	if s.smartOrderRouter == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Smart Order Router not available")
		return
	}

	// Create SOR handlers instance
	sorHandlers := NewSORHandlers(s.smartOrderRouter, s.logger)
	sorHandlers.ExecuteOrder(w, r)
}

// handleSORMetrics handles SOR metrics requests
func (s *APIServer) handleSORMetrics(w http.ResponseWriter, r *http.Request) {
	if s.smartOrderRouter == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Smart Order Router not available")
		return
	}

	// Create SOR handlers instance
	sorHandlers := NewSORHandlers(s.smartOrderRouter, s.logger)
	sorHandlers.GetMetrics(w, r)
}

// handleSORBestPrices handles best prices requests
func (s *APIServer) handleSORBestPrices(w http.ResponseWriter, r *http.Request) {
	if s.smartOrderRouter == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Smart Order Router not available")
		return
	}

	// Create SOR handlers instance
	sorHandlers := NewSORHandlers(s.smartOrderRouter, s.logger)
	sorHandlers.GetBestPrices(w, r)
}

// handleSORVenues handles venues information requests
func (s *APIServer) handleSORVenues(w http.ResponseWriter, r *http.Request) {
	if s.smartOrderRouter == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Smart Order Router not available")
		return
	}

	// Create SOR handlers instance
	sorHandlers := NewSORHandlers(s.smartOrderRouter, s.logger)
	sorHandlers.GetVenues(w, r)
}

// Advanced Risk Management handlers

// handleRiskMetrics handles risk metrics requests
func (s *APIServer) handleRiskMetrics(w http.ResponseWriter, r *http.Request) {
	if s.advancedRiskManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Advanced Risk Manager not available")
		return
	}

	riskHandlers := NewRiskHandlers(s.advancedRiskManager, s.logger)
	riskHandlers.GetMetrics(w, r)
}

// handleRiskViolations handles risk violations requests
func (s *APIServer) handleRiskViolations(w http.ResponseWriter, r *http.Request) {
	if s.advancedRiskManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Advanced Risk Manager not available")
		return
	}

	riskHandlers := NewRiskHandlers(s.advancedRiskManager, s.logger)
	riskHandlers.GetViolations(w, r)
}

// handleRiskValidate handles order validation requests
func (s *APIServer) handleRiskValidate(w http.ResponseWriter, r *http.Request) {
	if s.advancedRiskManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Advanced Risk Manager not available")
		return
	}

	riskHandlers := NewRiskHandlers(s.advancedRiskManager, s.logger)
	riskHandlers.ValidateOrder(w, r)
}

// handleRiskEmergencyStop handles emergency stop requests
func (s *APIServer) handleRiskEmergencyStop(w http.ResponseWriter, r *http.Request) {
	if s.advancedRiskManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Advanced Risk Manager not available")
		return
	}

	riskHandlers := NewRiskHandlers(s.advancedRiskManager, s.logger)
	riskHandlers.EmergencyStop(w, r)
}

// handleRiskStatus handles risk status requests
func (s *APIServer) handleRiskStatus(w http.ResponseWriter, r *http.Request) {
	if s.advancedRiskManager == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Advanced Risk Manager not available")
		return
	}

	status := map[string]interface{}{
		"emergency_mode": s.advancedRiskManager.IsEmergencyMode(),
		"system_status":  "operational",
	}

	// Add current risk level assessment
	metrics := s.advancedRiskManager.GetMetrics()
	if metrics != nil {
		riskLevel := "LOW"
		if metrics.Concentration > 50.0 {
			riskLevel = "MEDIUM"
		}
		if metrics.MaxDrawdown > 10.0 {
			riskLevel = "HIGH"
		}
		if metrics.MaxDrawdown > 20.0 {
			riskLevel = "CRITICAL"
		}

		status["risk_level"] = riskLevel
		status["last_update"] = metrics.LastUpdate
	}

	s.sendJSON(w, r, http.StatusOK, status)
}

// Post-Trade Analytics handlers

// handleAnalyticsExecutionMetrics handles execution metrics requests
func (s *APIServer) handleAnalyticsExecutionMetrics(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)
	analyticsHandlers.GetExecutionMetrics(w, r)
}

// handleAnalyticsPerformance handles performance data requests
func (s *APIServer) handleAnalyticsPerformance(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)
	analyticsHandlers.GetPerformanceData(w, r)
}

// handleAnalyticsRealtime handles real-time metrics requests
func (s *APIServer) handleAnalyticsRealtime(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)
	analyticsHandlers.GetRealtimeMetrics(w, r)
}

// handleAnalyticsTrades handles trade history and recording requests
func (s *APIServer) handleAnalyticsTrades(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)

	switch r.Method {
	case "GET":
		analyticsHandlers.GetTradeHistory(w, r)
	case "POST":
		analyticsHandlers.RecordTrade(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleAnalyticsSummary handles analytics summary requests
func (s *APIServer) handleAnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)
	analyticsHandlers.GetAnalyticsSummary(w, r)
}

// handleAnalyticsSlippage handles slippage analysis requests
func (s *APIServer) handleAnalyticsSlippage(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)
	analyticsHandlers.GetSlippageAnalysis(w, r)
}

// handleAnalyticsMarketImpact handles market impact analysis requests
func (s *APIServer) handleAnalyticsMarketImpact(w http.ResponseWriter, r *http.Request) {
	if s.postTradeAnalytics == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Post-Trade Analytics not available")
		return
	}

	analyticsHandlers := NewAnalyticsHandlers(s.postTradeAnalytics, s.logger)
	analyticsHandlers.GetMarketImpactAnalysis(w, r)
}

// High-Performance Networking handlers

// handleNetworkingMetrics handles networking metrics requests
func (s *APIServer) handleNetworkingMetrics(w http.ResponseWriter, r *http.Request) {
	if s.hpNetworking == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "High-Performance Networking not available")
		return
	}

	networkingHandlers := NewNetworkingHandlers(s.hpNetworking, s.logger)
	networkingHandlers.GetMetrics(w, r)
}

// handleNetworkingConnections handles connections requests
func (s *APIServer) handleNetworkingConnections(w http.ResponseWriter, r *http.Request) {
	if s.hpNetworking == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "High-Performance Networking not available")
		return
	}

	networkingHandlers := NewNetworkingHandlers(s.hpNetworking, s.logger)

	switch r.Method {
	case "GET":
		networkingHandlers.GetConnections(w, r)
	case "POST":
		networkingHandlers.CreateConnection(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleNetworkingConnection handles single connection requests
func (s *APIServer) handleNetworkingConnection(w http.ResponseWriter, r *http.Request) {
	if s.hpNetworking == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "High-Performance Networking not available")
		return
	}

	networkingHandlers := NewNetworkingHandlers(s.hpNetworking, s.logger)

	switch r.Method {
	case "GET":
		networkingHandlers.GetConnection(w, r)
	case "DELETE":
		networkingHandlers.CloseConnection(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleNetworkingSend handles message sending requests
func (s *APIServer) handleNetworkingSend(w http.ResponseWriter, r *http.Request) {
	if s.hpNetworking == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "High-Performance Networking not available")
		return
	}

	networkingHandlers := NewNetworkingHandlers(s.hpNetworking, s.logger)
	networkingHandlers.SendMessage(w, r)
}

// handleNetworkingReceive handles message receiving requests
func (s *APIServer) handleNetworkingReceive(w http.ResponseWriter, r *http.Request) {
	if s.hpNetworking == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "High-Performance Networking not available")
		return
	}

	networkingHandlers := NewNetworkingHandlers(s.hpNetworking, s.logger)
	networkingHandlers.ReceiveMessage(w, r)
}

// handleNetworkingStatus handles networking status requests
func (s *APIServer) handleNetworkingStatus(w http.ResponseWriter, r *http.Request) {
	if s.hpNetworking == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "High-Performance Networking not available")
		return
	}

	networkingHandlers := NewNetworkingHandlers(s.hpNetworking, s.logger)
	networkingHandlers.GetStatus(w, r)
}

// Real-Time Dashboard handlers

// handleDashboardMetrics handles dashboard metrics requests
func (s *APIServer) handleDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.GetLiveMetrics(w, r)
}

// handleDashboardAlerts handles dashboard alerts requests
func (s *APIServer) handleDashboardAlerts(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.GetAlerts(w, r)
}

// handleDashboardAcknowledgeAlert handles alert acknowledgment requests
func (s *APIServer) handleDashboardAcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.AcknowledgeAlert(w, r)
}

// handleDashboardResolveAlert handles alert resolution requests
func (s *APIServer) handleDashboardResolveAlert(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.ResolveAlert(w, r)
}

// handleDashboardSessions handles session creation requests
func (s *APIServer) handleDashboardSessions(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.CreateSession(w, r)
}

// handleDashboardSession handles session management requests
func (s *APIServer) handleDashboardSession(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)

	switch r.Method {
	case "GET":
		dashboardHandlers.GetSession(w, r)
	case "PUT":
		dashboardHandlers.UpdateSession(w, r)
	case "DELETE":
		dashboardHandlers.CloseSession(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDashboardWidgets handles widgets requests
func (s *APIServer) handleDashboardWidgets(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.GetWidgets(w, r)
}

// handleDashboardLayouts handles layouts requests
func (s *APIServer) handleDashboardLayouts(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.GetLayouts(w, r)
}

// handleDashboardStatus handles dashboard status requests
func (s *APIServer) handleDashboardStatus(w http.ResponseWriter, r *http.Request) {
	if s.realtimeDashboard == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Real-Time Dashboard not available")
		return
	}

	dashboardHandlers := NewDashboardHandlers(s.realtimeDashboard, s.logger)
	dashboardHandlers.GetStatus(w, r)
}

// Testing & Simulation Framework handlers

// handleTestingTests handles test execution and listing requests
func (s *APIServer) handleTestingTests(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)

	switch r.Method {
	case "GET":
		testingHandlers.GetTestResults(w, r)
	case "POST":
		testingHandlers.RunTest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTestingTest handles single test requests
func (s *APIServer) handleTestingTest(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	// For now, redirect to test results
	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)
	testingHandlers.GetTestResults(w, r)
}

// handleTestingSimulations handles simulation execution and listing requests
func (s *APIServer) handleTestingSimulations(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)

	switch r.Method {
	case "GET":
		testingHandlers.GetSimulations(w, r)
	case "POST":
		testingHandlers.RunSimulation(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTestingSimulation handles single simulation requests
func (s *APIServer) handleTestingSimulation(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)
	testingHandlers.GetSimulation(w, r)
}

// handleTestingEnvironments handles test environment requests
func (s *APIServer) handleTestingEnvironments(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)

	switch r.Method {
	case "GET":
		testingHandlers.GetEnvironments(w, r)
	case "POST":
		testingHandlers.CreateEnvironment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTestingEnvironment handles single environment requests
func (s *APIServer) handleTestingEnvironment(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)
	testingHandlers.GetEnvironment(w, r)
}

// handleTestingResults handles test results requests
func (s *APIServer) handleTestingResults(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)
	testingHandlers.GetTestResults(w, r)
}

// handleTestingStatus handles testing framework status requests
func (s *APIServer) handleTestingStatus(w http.ResponseWriter, r *http.Request) {
	if s.testingFramework == nil {
		s.sendError(w, r, http.StatusServiceUnavailable, "Testing Framework not available")
		return
	}

	testingHandlers := NewHFTTestingHandlers(s.testingFramework, s.logger)
	testingHandlers.GetStatus(w, r)
}
