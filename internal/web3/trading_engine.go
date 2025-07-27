package web3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TradingEngine provides autonomous trading capabilities
type TradingEngine struct {
	clients         map[int]*ethclient.Client
	logger          *observability.Logger
	riskAssessment  *RiskAssessmentService
	strategies      map[string]TradingStrategy
	activePositions map[string]*Position
	portfolios      map[uuid.UUID]*Portfolio
	config          TradingConfig
	isRunning       bool
	stopChan        chan struct{}
	mu              sync.RWMutex
}

// TradingConfig holds configuration for the trading engine
type TradingConfig struct {
	MaxPositionSize   decimal.Decimal `json:"max_position_size"`
	MaxDailyLoss      decimal.Decimal `json:"max_daily_loss"`
	RiskPerTrade      decimal.Decimal `json:"risk_per_trade"`
	MinLiquidity      decimal.Decimal `json:"min_liquidity"`
	SlippageTolerance decimal.Decimal `json:"slippage_tolerance"`
	GasMultiplier     decimal.Decimal `json:"gas_multiplier"`
	ExecutionInterval time.Duration   `json:"execution_interval"`
	RebalanceInterval time.Duration   `json:"rebalance_interval"`
	EnableStopLoss    bool            `json:"enable_stop_loss"`
	EnableTakeProfit  bool            `json:"enable_take_profit"`
	EmergencyStopLoss decimal.Decimal `json:"emergency_stop_loss"`
	AllowedTokens     []string        `json:"allowed_tokens"`
	BlacklistedTokens []string        `json:"blacklisted_tokens"`
}

// TradingStrategy interface for different trading strategies
type TradingStrategy interface {
	GetName() string
	GetDescription() string
	Analyze(ctx context.Context, market *MarketData) (*TradingSignal, error)
	ValidateSignal(ctx context.Context, signal *TradingSignal) error
	CalculatePositionSize(ctx context.Context, signal *TradingSignal, portfolio *Portfolio) (decimal.Decimal, error)
	GetRiskLevel() RiskLevel
	IsEnabled() bool
	GetParameters() map[string]interface{}
}

// TradingSignal represents a trading signal from a strategy
type TradingSignal struct {
	ID           uuid.UUID              `json:"id"`
	StrategyName string                 `json:"strategy_name"`
	Action       TradingAction          `json:"action"`
	TokenIn      string                 `json:"token_in"`
	TokenOut     string                 `json:"token_out"`
	AmountIn     decimal.Decimal        `json:"amount_in"`
	ExpectedOut  decimal.Decimal        `json:"expected_out"`
	Confidence   float64                `json:"confidence"`
	Urgency      SignalUrgency          `json:"urgency"`
	ValidUntil   time.Time              `json:"valid_until"`
	StopLoss     *decimal.Decimal       `json:"stop_loss,omitempty"`
	TakeProfit   *decimal.Decimal       `json:"take_profit,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
}

// TradingAction represents the type of trading action
type TradingAction string

const (
	ActionBuy     TradingAction = "buy"
	ActionSell    TradingAction = "sell"
	ActionHold    TradingAction = "hold"
	ActionSwap    TradingAction = "swap"
	ActionStake   TradingAction = "stake"
	ActionUnstake TradingAction = "unstake"
)

// SignalUrgency represents the urgency of a trading signal
type SignalUrgency string

const (
	UrgencyLow      SignalUrgency = "low"
	UrgencyMedium   SignalUrgency = "medium"
	UrgencyHigh     SignalUrgency = "high"
	UrgencyCritical SignalUrgency = "critical"
)

// Position represents an active trading position
type Position struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"user_id"`
	StrategyName  string                 `json:"strategy_name"`
	TokenAddress  string                 `json:"token_address"`
	TokenSymbol   string                 `json:"token_symbol"`
	Amount        decimal.Decimal        `json:"amount"`
	EntryPrice    decimal.Decimal        `json:"entry_price"`
	CurrentPrice  decimal.Decimal        `json:"current_price"`
	UnrealizedPnL decimal.Decimal        `json:"unrealized_pnl"`
	RealizedPnL   decimal.Decimal        `json:"realized_pnl"`
	StopLoss      *decimal.Decimal       `json:"stop_loss,omitempty"`
	TakeProfit    *decimal.Decimal       `json:"take_profit,omitempty"`
	Status        PositionStatus         `json:"status"`
	OpenedAt      time.Time              `json:"opened_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	ClosedAt      *time.Time             `json:"closed_at,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// PositionStatus represents the status of a position
type PositionStatus string

const (
	PositionStatusOpen    PositionStatus = "open"
	PositionStatusClosed  PositionStatus = "closed"
	PositionStatusPending PositionStatus = "pending"
)

// Portfolio represents a user's trading portfolio
type Portfolio struct {
	ID                uuid.UUID              `json:"id"`
	UserID            uuid.UUID              `json:"user_id"`
	Name              string                 `json:"name"`
	TotalValue        decimal.Decimal        `json:"total_value"`
	AvailableBalance  decimal.Decimal        `json:"available_balance"`
	InvestedAmount    decimal.Decimal        `json:"invested_amount"`
	TotalPnL          decimal.Decimal        `json:"total_pnl"`
	DailyPnL          decimal.Decimal        `json:"daily_pnl"`
	Holdings          map[string]*Holding    `json:"holdings"`
	ActivePositions   []uuid.UUID            `json:"active_positions"`
	TradingStrategies []string               `json:"trading_strategies"`
	RiskProfile       RiskProfile            `json:"risk_profile"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// Holding represents a token holding in a portfolio
type Holding struct {
	TokenAddress  string          `json:"token_address"`
	TokenSymbol   string          `json:"token_symbol"`
	Amount        decimal.Decimal `json:"amount"`
	AveragePrice  decimal.Decimal `json:"average_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	Value         decimal.Decimal `json:"value"`
	PnL           decimal.Decimal `json:"pnl"`
	PnLPercentage decimal.Decimal `json:"pnl_percentage"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// RiskProfile represents a user's risk tolerance
type RiskProfile struct {
	Level                string          `json:"level"`                  // conservative, moderate, aggressive
	MaxPositionSize      decimal.Decimal `json:"max_position_size"`      // % of portfolio
	MaxDailyLoss         decimal.Decimal `json:"max_daily_loss"`         // % of portfolio
	StopLossPercentage   decimal.Decimal `json:"stop_loss_percentage"`   // % below entry
	TakeProfitPercentage decimal.Decimal `json:"take_profit_percentage"` // % above entry
	AllowedStrategies    []string        `json:"allowed_strategies"`
}

// MarketData represents market data for analysis
type MarketData struct {
	TokenAddress   string                 `json:"token_address"`
	TokenSymbol    string                 `json:"token_symbol"`
	Price          decimal.Decimal        `json:"price"`
	Volume24h      decimal.Decimal        `json:"volume_24h"`
	MarketCap      decimal.Decimal        `json:"market_cap"`
	PriceChange24h decimal.Decimal        `json:"price_change_24h"`
	Liquidity      decimal.Decimal        `json:"liquidity"`
	Volatility     decimal.Decimal        `json:"volatility"`
	TechnicalData  *TechnicalIndicators   `json:"technical_data"`
	Timestamp      time.Time              `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// TechnicalIndicators represents technical analysis indicators
type TechnicalIndicators struct {
	RSI            decimal.Decimal `json:"rsi"`
	MACD           decimal.Decimal `json:"macd"`
	MACDSignal     decimal.Decimal `json:"macd_signal"`
	BollingerUpper decimal.Decimal `json:"bollinger_upper"`
	BollingerLower decimal.Decimal `json:"bollinger_lower"`
	SMA20          decimal.Decimal `json:"sma_20"`
	SMA50          decimal.Decimal `json:"sma_50"`
	EMA12          decimal.Decimal `json:"ema_12"`
	EMA26          decimal.Decimal `json:"ema_26"`
	Volume         decimal.Decimal `json:"volume"`
	VWAP           decimal.Decimal `json:"vwap"`
}

// NewTradingEngine creates a new trading engine
func NewTradingEngine(
	clients map[int]*ethclient.Client,
	logger *observability.Logger,
	riskAssessment *RiskAssessmentService,
) *TradingEngine {
	config := TradingConfig{
		MaxPositionSize:   decimal.NewFromFloat(0.1),   // 10% of portfolio
		MaxDailyLoss:      decimal.NewFromFloat(0.05),  // 5% daily loss limit
		RiskPerTrade:      decimal.NewFromFloat(0.02),  // 2% risk per trade
		MinLiquidity:      decimal.NewFromInt(10000),   // $10k minimum liquidity
		SlippageTolerance: decimal.NewFromFloat(0.005), // 0.5% slippage
		GasMultiplier:     decimal.NewFromFloat(1.2),   // 20% gas buffer
		ExecutionInterval: 30 * time.Second,
		RebalanceInterval: 1 * time.Hour,
		EnableStopLoss:    true,
		EnableTakeProfit:  true,
		EmergencyStopLoss: decimal.NewFromFloat(0.2), // 20% emergency stop
		AllowedTokens:     []string{"WETH", "USDC", "USDT", "DAI"},
		BlacklistedTokens: []string{},
	}

	engine := &TradingEngine{
		clients:         clients,
		logger:          logger,
		riskAssessment:  riskAssessment,
		strategies:      make(map[string]TradingStrategy),
		activePositions: make(map[string]*Position),
		portfolios:      make(map[uuid.UUID]*Portfolio),
		config:          config,
		stopChan:        make(chan struct{}),
	}

	// Initialize default strategies
	engine.initializeStrategies()

	return engine
}

// Start starts the trading engine
func (t *TradingEngine) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isRunning {
		return fmt.Errorf("trading engine is already running")
	}

	t.isRunning = true

	// Start trading loop
	go t.tradingLoop(ctx)

	// Start portfolio rebalancing loop
	go t.rebalancingLoop(ctx)

	t.logger.Info(ctx, "Trading engine started", map[string]interface{}{
		"strategies":         len(t.strategies),
		"active_positions":   len(t.activePositions),
		"portfolios":         len(t.portfolios),
		"execution_interval": t.config.ExecutionInterval.String(),
	})

	return nil
}

// Stop stops the trading engine
func (t *TradingEngine) Stop(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isRunning {
		return fmt.Errorf("trading engine is not running")
	}

	close(t.stopChan)
	t.isRunning = false

	t.logger.Info(ctx, "Trading engine stopped", nil)

	return nil
}

// CreatePortfolio creates a new trading portfolio
func (t *TradingEngine) CreatePortfolio(ctx context.Context, userID uuid.UUID, name string, initialBalance decimal.Decimal, riskProfile RiskProfile) (*Portfolio, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	portfolio := &Portfolio{
		ID:                uuid.New(),
		UserID:            userID,
		Name:              name,
		TotalValue:        initialBalance,
		AvailableBalance:  initialBalance,
		InvestedAmount:    decimal.Zero,
		TotalPnL:          decimal.Zero,
		DailyPnL:          decimal.Zero,
		Holdings:          make(map[string]*Holding),
		ActivePositions:   []uuid.UUID{},
		TradingStrategies: []string{},
		RiskProfile:       riskProfile,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Metadata:          make(map[string]interface{}),
	}

	t.portfolios[portfolio.ID] = portfolio

	t.logger.Info(ctx, "Portfolio created", map[string]interface{}{
		"portfolio_id":    portfolio.ID.String(),
		"user_id":         userID.String(),
		"name":            name,
		"initial_balance": initialBalance.String(),
		"risk_profile":    riskProfile.Level,
	})

	return portfolio, nil
}

// tradingLoop is the main trading execution loop
func (t *TradingEngine) tradingLoop(ctx context.Context) {
	ticker := time.NewTicker(t.config.ExecutionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stopChan:
			return
		case <-ticker.C:
			t.executeTrading(ctx)
		}
	}
}

// rebalancingLoop handles portfolio rebalancing
func (t *TradingEngine) rebalancingLoop(ctx context.Context) {
	ticker := time.NewTicker(t.config.RebalanceInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stopChan:
			return
		case <-ticker.C:
			t.rebalancePortfolios(ctx)
		}
	}
}

// executeTrading executes trading strategies
func (t *TradingEngine) executeTrading(ctx context.Context) {
	t.mu.RLock()
	portfolios := make([]*Portfolio, 0, len(t.portfolios))
	for _, portfolio := range t.portfolios {
		portfolios = append(portfolios, portfolio)
	}
	strategies := make([]TradingStrategy, 0, len(t.strategies))
	for _, strategy := range t.strategies {
		strategies = append(strategies, strategy)
	}
	t.mu.RUnlock()

	for _, portfolio := range portfolios {
		for _, strategy := range strategies {
			if !strategy.IsEnabled() {
				continue
			}

			// Check if strategy is allowed for this portfolio
			if !t.isStrategyAllowed(portfolio, strategy.GetName()) {
				continue
			}

			// Get market data for analysis
			marketData, err := t.getMarketData(ctx, portfolio)
			if err != nil {
				t.logger.Warn(ctx, "Failed to get market data", map[string]interface{}{
					"portfolio_id": portfolio.ID.String(),
					"strategy":     strategy.GetName(),
					"error":        err.Error(),
				})
				continue
			}

			// Analyze market and get trading signal
			signal, err := strategy.Analyze(ctx, marketData)
			if err != nil {
				t.logger.Warn(ctx, "Strategy analysis failed", map[string]interface{}{
					"portfolio_id": portfolio.ID.String(),
					"strategy":     strategy.GetName(),
					"error":        err.Error(),
				})
				continue
			}

			if signal == nil || signal.Action == ActionHold {
				continue
			}

			// Validate signal
			if err := strategy.ValidateSignal(ctx, signal); err != nil {
				t.logger.Warn(ctx, "Signal validation failed", map[string]interface{}{
					"signal_id": signal.ID.String(),
					"strategy":  strategy.GetName(),
					"error":     err.Error(),
				})
				continue
			}

			// Execute signal
			if err := t.executeSignal(ctx, portfolio, signal); err != nil {
				t.logger.Error(ctx, "Signal execution failed", err)
			}
		}
	}
}

// executeSignal executes a trading signal
func (t *TradingEngine) executeSignal(ctx context.Context, portfolio *Portfolio, signal *TradingSignal) error {
	// Perform risk assessment
	if err := t.assessSignalRisk(ctx, portfolio, signal); err != nil {
		return fmt.Errorf("signal risk assessment failed: %w", err)
	}

	// Check portfolio limits
	if err := t.checkPortfolioLimits(ctx, portfolio, signal); err != nil {
		return fmt.Errorf("portfolio limits exceeded: %w", err)
	}

	// Calculate position size
	strategy := t.strategies[signal.StrategyName]
	positionSize, err := strategy.CalculatePositionSize(ctx, signal, portfolio)
	if err != nil {
		return fmt.Errorf("failed to calculate position size: %w", err)
	}

	// Execute the trade
	position, err := t.executeTrade(ctx, portfolio, signal, positionSize)
	if err != nil {
		return fmt.Errorf("trade execution failed: %w", err)
	}

	// Update portfolio
	t.updatePortfolioAfterTrade(ctx, portfolio, position)

	t.logger.Info(ctx, "Signal executed successfully", map[string]interface{}{
		"signal_id":    signal.ID.String(),
		"position_id":  position.ID.String(),
		"portfolio_id": portfolio.ID.String(),
		"action":       string(signal.Action),
		"amount":       positionSize.String(),
	})

	return nil
}

// assessSignalRisk performs risk assessment on a trading signal
func (t *TradingEngine) assessSignalRisk(ctx context.Context, portfolio *Portfolio, signal *TradingSignal) error {
	// Create risk assessment request
	req := TransactionRiskRequest{
		FromAddress:     portfolio.ID.String(), // Use portfolio ID as identifier
		ToAddress:       signal.TokenOut,
		Value:           signal.AmountIn.BigInt(),
		ChainID:         1, // Default to Ethereum mainnet
		IncludeMLModels: true,
		Metadata: map[string]interface{}{
			"signal_id":     signal.ID.String(),
			"strategy_name": signal.StrategyName,
			"action":        string(signal.Action),
		},
	}

	assessment, err := t.riskAssessment.AssessTransactionRisk(ctx, req)
	if err != nil {
		return fmt.Errorf("risk assessment failed: %w", err)
	}

	// Check risk thresholds
	if assessment.RiskScore > 70 {
		return fmt.Errorf("signal risk score too high: %d", assessment.RiskScore)
	}

	if assessment.SafetyGrade == SafetyGradeF || assessment.SafetyGrade == SafetyGradeD {
		return fmt.Errorf("signal safety grade too low: %s", assessment.SafetyGrade)
	}

	// Add risk assessment to signal metadata
	signal.Metadata["risk_assessment"] = assessment

	return nil
}

// checkPortfolioLimits checks if the signal respects portfolio limits
func (t *TradingEngine) checkPortfolioLimits(ctx context.Context, portfolio *Portfolio, signal *TradingSignal) error {
	// Check daily loss limit
	if portfolio.DailyPnL.LessThan(t.config.MaxDailyLoss.Neg().Mul(portfolio.TotalValue)) {
		return fmt.Errorf("daily loss limit exceeded")
	}

	// Check position size limit
	maxPositionValue := t.config.MaxPositionSize.Mul(portfolio.TotalValue)
	if signal.AmountIn.GreaterThan(maxPositionValue) {
		return fmt.Errorf("position size exceeds limit")
	}

	// Check available balance
	if signal.AmountIn.GreaterThan(portfolio.AvailableBalance) {
		return fmt.Errorf("insufficient available balance")
	}

	return nil
}

// executeTrade executes the actual trade
func (t *TradingEngine) executeTrade(ctx context.Context, portfolio *Portfolio, signal *TradingSignal, positionSize decimal.Decimal) (*Position, error) {
	// Create position
	position := &Position{
		ID:            uuid.New(),
		UserID:        portfolio.UserID,
		StrategyName:  signal.StrategyName,
		TokenAddress:  signal.TokenOut,
		TokenSymbol:   signal.TokenOut, // Simplified - would resolve symbol
		Amount:        positionSize,
		EntryPrice:    signal.ExpectedOut.Div(signal.AmountIn), // Price per token
		CurrentPrice:  signal.ExpectedOut.Div(signal.AmountIn),
		UnrealizedPnL: decimal.Zero,
		RealizedPnL:   decimal.Zero,
		StopLoss:      signal.StopLoss,
		TakeProfit:    signal.TakeProfit,
		Status:        PositionStatusPending,
		OpenedAt:      time.Now(),
		UpdatedAt:     time.Now(),
		Metadata:      signal.Metadata,
	}

	// In a real implementation, this would interact with DEX contracts
	// For now, simulate successful execution
	position.Status = PositionStatusOpen

	// Store position
	t.mu.Lock()
	t.activePositions[position.ID.String()] = position
	portfolio.ActivePositions = append(portfolio.ActivePositions, position.ID)
	t.mu.Unlock()

	return position, nil
}

// updatePortfolioAfterTrade updates portfolio state after a trade
func (t *TradingEngine) updatePortfolioAfterTrade(ctx context.Context, portfolio *Portfolio, position *Position) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Update available balance
	portfolio.AvailableBalance = portfolio.AvailableBalance.Sub(position.Amount)
	portfolio.InvestedAmount = portfolio.InvestedAmount.Add(position.Amount)

	// Update holdings
	if holding, exists := portfolio.Holdings[position.TokenAddress]; exists {
		// Update existing holding
		totalAmount := holding.Amount.Add(position.Amount)
		totalValue := holding.Amount.Mul(holding.AveragePrice).Add(position.Amount.Mul(position.EntryPrice))
		holding.AveragePrice = totalValue.Div(totalAmount)
		holding.Amount = totalAmount
	} else {
		// Create new holding
		portfolio.Holdings[position.TokenAddress] = &Holding{
			TokenAddress:  position.TokenAddress,
			TokenSymbol:   position.TokenSymbol,
			Amount:        position.Amount,
			AveragePrice:  position.EntryPrice,
			CurrentPrice:  position.CurrentPrice,
			Value:         position.Amount.Mul(position.CurrentPrice),
			PnL:           decimal.Zero,
			PnLPercentage: decimal.Zero,
			LastUpdated:   time.Now(),
		}
	}

	portfolio.UpdatedAt = time.Now()
}

// initializeStrategies initializes default trading strategies
func (t *TradingEngine) initializeStrategies() {
	t.strategies["momentum"] = NewMomentumStrategy()
	t.strategies["mean_reversion"] = NewMeanReversionStrategy()
	t.strategies["arbitrage"] = NewArbitrageStrategy()
}

// rebalancePortfolios performs portfolio rebalancing
func (t *TradingEngine) rebalancePortfolios(ctx context.Context) {
	t.mu.RLock()
	portfolios := make([]*Portfolio, 0, len(t.portfolios))
	for _, portfolio := range t.portfolios {
		portfolios = append(portfolios, portfolio)
	}
	t.mu.RUnlock()

	for _, portfolio := range portfolios {
		// Check if portfolio needs rebalancing
		if t.shouldRebalancePortfolio(portfolio) {
			if err := t.rebalancePortfolio(ctx, portfolio); err != nil {
				t.logger.Error(ctx, "Portfolio rebalancing failed", err)
			}
		}
	}
}

// shouldRebalancePortfolio checks if a portfolio needs rebalancing
func (t *TradingEngine) shouldRebalancePortfolio(portfolio *Portfolio) bool {
	// Simple rebalancing logic - rebalance if daily P&L is significant
	threshold := portfolio.TotalValue.Mul(decimal.NewFromFloat(0.05)) // 5% threshold
	return portfolio.DailyPnL.Abs().GreaterThan(threshold)
}

// rebalancePortfolio performs portfolio rebalancing
func (t *TradingEngine) rebalancePortfolio(ctx context.Context, portfolio *Portfolio) error {
	t.logger.Info(ctx, "Rebalancing portfolio", map[string]interface{}{
		"portfolio_id": portfolio.ID.String(),
		"total_value":  portfolio.TotalValue.String(),
		"daily_pnl":    portfolio.DailyPnL.String(),
	})

	// Simple rebalancing: if we have significant losses, reduce position sizes
	if portfolio.DailyPnL.IsNegative() {
		for _, positionID := range portfolio.ActivePositions {
			position, exists := t.activePositions[positionID.String()]
			if !exists {
				continue
			}

			// Reduce position size by 10%
			reductionAmount := position.Amount.Mul(decimal.NewFromFloat(0.1))
			position.Amount = position.Amount.Sub(reductionAmount)
			portfolio.AvailableBalance = portfolio.AvailableBalance.Add(reductionAmount)
		}
	}

	return nil
}

// GetPortfolio returns a portfolio by ID
func (t *TradingEngine) GetPortfolio(portfolioID uuid.UUID) (*Portfolio, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	portfolio, exists := t.portfolios[portfolioID]
	if !exists {
		return nil, fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	return portfolio, nil
}

// isStrategyAllowed checks if a strategy is allowed for a portfolio
func (t *TradingEngine) isStrategyAllowed(portfolio *Portfolio, strategyName string) bool {
	if len(portfolio.TradingStrategies) == 0 {
		return true // No restrictions
	}

	for _, allowedStrategy := range portfolio.TradingStrategies {
		if allowedStrategy == strategyName {
			return true
		}
	}

	return false
}

// getMarketData gets market data for analysis (placeholder implementation)
func (t *TradingEngine) getMarketData(ctx context.Context, portfolio *Portfolio) (*MarketData, error) {
	// This would fetch real market data from exchanges/price feeds
	// For now, return mock data
	return &MarketData{
		TokenAddress:   "0xA0b86a33E6441e6e80D0c4C6C7556C974E1B2c20", // Mock token
		TokenSymbol:    "MOCK",
		Price:          decimal.NewFromFloat(100.0),
		Volume24h:      decimal.NewFromInt(1000000),
		MarketCap:      decimal.NewFromInt(100000000),
		PriceChange24h: decimal.NewFromFloat(5.0),
		Liquidity:      decimal.NewFromInt(5000000),
		Volatility:     decimal.NewFromFloat(0.15),
		TechnicalData: &TechnicalIndicators{
			RSI:            decimal.NewFromFloat(45.0),
			MACD:           decimal.NewFromFloat(2.5),
			MACDSignal:     decimal.NewFromFloat(2.0),
			BollingerUpper: decimal.NewFromFloat(105.0),
			BollingerLower: decimal.NewFromFloat(95.0),
			SMA20:          decimal.NewFromFloat(98.0),
			SMA50:          decimal.NewFromFloat(96.0),
			EMA12:          decimal.NewFromFloat(99.0),
			EMA26:          decimal.NewFromFloat(97.0),
			Volume:         decimal.NewFromInt(800000),
			VWAP:           decimal.NewFromFloat(99.5),
		},
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}, nil
}

// UpdatePortfolioValue updates portfolio value and P&L
func (t *TradingEngine) UpdatePortfolioValue(ctx context.Context, portfolioID uuid.UUID) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	portfolio, exists := t.portfolios[portfolioID]
	if !exists {
		return fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	// Calculate total value from holdings
	totalValue := portfolio.AvailableBalance
	for _, holding := range portfolio.Holdings {
		totalValue = totalValue.Add(holding.Value)
	}

	// Update P&L
	previousValue := portfolio.TotalValue
	portfolio.TotalValue = totalValue
	portfolio.DailyPnL = totalValue.Sub(previousValue)
	portfolio.TotalPnL = totalValue.Sub(portfolio.InvestedAmount)
	portfolio.UpdatedAt = time.Now()

	return nil
}

// GetActivePositions returns all active positions for a portfolio
func (t *TradingEngine) GetActivePositions(portfolioID uuid.UUID) ([]*Position, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	portfolio, exists := t.portfolios[portfolioID]
	if !exists {
		return nil, fmt.Errorf("portfolio not found: %s", portfolioID.String())
	}

	var positions []*Position
	for _, positionID := range portfolio.ActivePositions {
		if position, exists := t.activePositions[positionID.String()]; exists {
			positions = append(positions, position)
		}
	}

	return positions, nil
}

// ClosePosition closes a trading position
func (t *TradingEngine) ClosePosition(ctx context.Context, positionID uuid.UUID, reason string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	position, exists := t.activePositions[positionID.String()]
	if !exists {
		return fmt.Errorf("position not found: %s", positionID.String())
	}

	// Update position status
	position.Status = PositionStatusClosed
	now := time.Now()
	position.ClosedAt = &now
	position.UpdatedAt = now

	// Calculate realized P&L
	position.RealizedPnL = position.UnrealizedPnL

	// Remove from active positions
	delete(t.activePositions, positionID.String())

	// Update portfolio
	portfolio := t.portfolios[position.UserID]
	if portfolio != nil {
		// Remove from active positions list
		for i, activeID := range portfolio.ActivePositions {
			if activeID == positionID {
				portfolio.ActivePositions = append(portfolio.ActivePositions[:i], portfolio.ActivePositions[i+1:]...)
				break
			}
		}

		// Add back to available balance
		portfolio.AvailableBalance = portfolio.AvailableBalance.Add(position.Amount)
		portfolio.InvestedAmount = portfolio.InvestedAmount.Sub(position.Amount)
		portfolio.TotalPnL = portfolio.TotalPnL.Add(position.RealizedPnL)
	}

	t.logger.Info(ctx, "Position closed", map[string]interface{}{
		"position_id":  positionID.String(),
		"realized_pnl": position.RealizedPnL.String(),
		"reason":       reason,
	})

	return nil
}
