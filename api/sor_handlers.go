package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SORHandlers provides HTTP handlers for Smart Order Routing
type SORHandlers struct {
	sor    *hft.SmartOrderRouter
	logger *observability.Logger
}

// NewSORHandlers creates new SOR HTTP handlers
func NewSORHandlers(sor *hft.SmartOrderRouter, logger *observability.Logger) *SORHandlers {
	return &SORHandlers{
		sor:    sor,
		logger: logger,
	}
}

// RouteOrderRequest represents an order routing request
type RouteOrderRequest struct {
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`                    // "BUY" or "SELL"
	Quantity    string `json:"quantity"`                // Decimal string
	Type        string `json:"type"`                    // "MARKET", "LIMIT"
	Price       string `json:"price,omitempty"`         // Decimal string for limit orders
	TimeInForce string `json:"time_in_force,omitempty"` // "GTC", "IOC", "FOK"
	ClientID    string `json:"client_id,omitempty"`
	StrategyID  string `json:"strategy_id,omitempty"`
}

// RouteOrderResponse represents an order routing response
type RouteOrderResponse struct {
	OrderID       string                `json:"order_id"`
	Symbol        string                `json:"symbol"`
	TotalQuantity string                `json:"total_quantity"`
	Algorithm     string                `json:"algorithm"`
	Venues        []VenueAllocationResp `json:"venues"`
	ExpectedCost  string                `json:"expected_cost"`
	ExpectedTime  string                `json:"expected_time"`
	MarketImpact  float64               `json:"market_impact"`
	RiskScore     float64               `json:"risk_score"`
	Approved      bool                  `json:"approved"`
	ChildOrders   []ChildOrderResp      `json:"child_orders"`
	Timestamp     time.Time             `json:"timestamp"`
}

// VenueAllocationResp represents venue allocation in response
type VenueAllocationResp struct {
	VenueID       string `json:"venue_id"`
	Quantity      string `json:"quantity"`
	ExpectedPrice string `json:"expected_price"`
	Priority      int    `json:"priority"`
	OrderType     string `json:"order_type"`
	TimeInForce   string `json:"time_in_force"`
}

// ChildOrderResp represents child order in response
type ChildOrderResp struct {
	ID          string    `json:"id"`
	ParentID    string    `json:"parent_id"`
	VenueID     string    `json:"venue_id"`
	Symbol      string    `json:"symbol"`
	Side        string    `json:"side"`
	Quantity    string    `json:"quantity"`
	Price       string    `json:"price"`
	OrderType   string    `json:"order_type"`
	TimeInForce string    `json:"time_in_force"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RouteOrder handles order routing requests
func (h *SORHandlers) RouteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RouteOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode route order request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateRouteOrderRequest(&req); err != nil {
		h.logger.Error(ctx, "Invalid route order request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to internal order request
	orderReq, err := h.convertToOrderRequest(&req)
	if err != nil {
		h.logger.Error(ctx, "Failed to convert order request", err)
		http.Error(w, "Invalid order parameters", http.StatusBadRequest)
		return
	}

	// Route the order
	decision, err := h.sor.RouteOrder(ctx, orderReq)
	if err != nil {
		h.logger.Error(ctx, "Failed to route order", err)
		http.Error(w, "Order routing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	resp := h.convertToRouteOrderResponse(decision)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error(ctx, "Failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, "Order routed successfully", map[string]interface{}{
		"order_id":    decision.OrderID.String(),
		"symbol":      decision.Symbol,
		"venues_used": len(decision.Venues),
		"approved":    decision.Approved,
	})
}

// ExecuteOrder handles order execution requests
func (h *SORHandlers) ExecuteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var decision hft.RoutingDecision
	if err := json.NewDecoder(r.Body).Decode(&decision); err != nil {
		h.logger.Error(ctx, "Failed to decode execution request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute the order
	if err := h.sor.ExecuteOrder(ctx, &decision); err != nil {
		h.logger.Error(ctx, "Failed to execute order", err)
		http.Error(w, "Order execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"message":  "Order execution initiated",
		"order_id": decision.OrderID.String(),
	})

	h.logger.Info(ctx, "Order execution initiated", map[string]interface{}{
		"order_id": decision.OrderID.String(),
	})
}

// GetMetrics handles SOR metrics requests
func (h *SORHandlers) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.sor.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Error(ctx, "Failed to encode metrics", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetBestPrices handles best prices requests
func (h *SORHandlers) GetBestPrices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		symbol = "BTCUSDT" // Default symbol
	}

	bestPrices := h.sor.GetBestPrices(symbol)
	if bestPrices == nil {
		http.Error(w, "No market data available for symbol", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bestPrices); err != nil {
		h.logger.Error(ctx, "Failed to encode best prices", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetVenues handles venues information requests
func (h *SORHandlers) GetVenues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	venues := h.sor.GetVenues()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(venues); err != nil {
		h.logger.Error(ctx, "Failed to encode venues", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// validateRouteOrderRequest validates the route order request
func (h *SORHandlers) validateRouteOrderRequest(req *RouteOrderRequest) error {
	if req.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if req.Side != "BUY" && req.Side != "SELL" {
		return fmt.Errorf("side must be BUY or SELL")
	}
	if req.Quantity == "" {
		return fmt.Errorf("quantity is required")
	}
	if req.Type != "MARKET" && req.Type != "LIMIT" {
		return fmt.Errorf("type must be MARKET or LIMIT")
	}
	if req.Type == "LIMIT" && req.Price == "" {
		return fmt.Errorf("price is required for limit orders")
	}
	return nil
}

// convertToOrderRequest converts HTTP request to internal order request
func (h *SORHandlers) convertToOrderRequest(req *RouteOrderRequest) (*hft.OrderRequest, error) {
	quantity, err := decimal.NewFromString(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity: %w", err)
	}

	var price decimal.Decimal
	if req.Price != "" {
		price, err = decimal.NewFromString(req.Price)
		if err != nil {
			return nil, fmt.Errorf("invalid price: %w", err)
		}
	}

	var side hft.OrderSide
	if req.Side == "BUY" {
		side = hft.OrderSideBuy
	} else {
		side = hft.OrderSideSell
	}

	var orderType hft.OrderType
	if req.Type == "MARKET" {
		orderType = hft.OrderTypeMarket
	} else {
		orderType = hft.OrderTypeLimit
	}

	var timeInForce hft.TimeInForce = hft.TimeInForceGTC
	if req.TimeInForce != "" {
		switch req.TimeInForce {
		case "GTC":
			timeInForce = hft.TimeInForceGTC
		case "IOC":
			timeInForce = hft.TimeInForceIOC
		case "FOK":
			timeInForce = hft.TimeInForceFOK
		}
	}

	return &hft.OrderRequest{
		ID:          uuid.New(),
		Symbol:      req.Symbol,
		Side:        side,
		Quantity:    quantity,
		Type:        orderType,
		Price:       price,
		TimeInForce: timeInForce,
		ClientID:    req.ClientID,
		StrategyID:  req.StrategyID,
	}, nil
}

// convertToRouteOrderResponse converts internal routing decision to HTTP response
func (h *SORHandlers) convertToRouteOrderResponse(decision *hft.RoutingDecision) *RouteOrderResponse {
	venues := make([]VenueAllocationResp, len(decision.Venues))
	for i, venue := range decision.Venues {
		venues[i] = VenueAllocationResp{
			VenueID:       venue.VenueID,
			Quantity:      venue.Quantity.String(),
			ExpectedPrice: venue.ExpectedPrice.String(),
			Priority:      venue.Priority,
			OrderType:     string(venue.OrderType),
			TimeInForce:   string(venue.TimeInForce),
		}
	}

	childOrders := make([]ChildOrderResp, len(decision.ChildOrders))
	for i, child := range decision.ChildOrders {
		childOrders[i] = ChildOrderResp{
			ID:          child.ID.String(),
			ParentID:    child.ParentID.String(),
			VenueID:     child.VenueID,
			Symbol:      child.Symbol,
			Side:        string(child.Side),
			Quantity:    child.Quantity.String(),
			Price:       child.Price.String(),
			OrderType:   string(child.OrderType),
			TimeInForce: string(child.TimeInForce),
			Status:      string(child.Status),
			CreatedAt:   child.CreatedAt,
			UpdatedAt:   child.UpdatedAt,
		}
	}

	return &RouteOrderResponse{
		OrderID:       decision.OrderID.String(),
		Symbol:        decision.Symbol,
		TotalQuantity: decision.TotalQuantity.String(),
		Algorithm:     decision.Algorithm,
		Venues:        venues,
		ExpectedCost:  decision.ExpectedCost.String(),
		ExpectedTime:  decision.ExpectedTime.String(),
		MarketImpact:  decision.MarketImpact,
		RiskScore:     decision.RiskScore,
		Approved:      decision.Approved,
		ChildOrders:   childOrders,
		Timestamp:     decision.Timestamp,
	}
}
