package billing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// PerformanceFeeTracker tracks trading performance and calculates fees
type PerformanceFeeTracker struct {
	db *sql.DB
}

// NewPerformanceFeeTracker creates a new performance fee tracker
func NewPerformanceFeeTracker(db *sql.DB) *PerformanceFeeTracker {
	return &PerformanceFeeTracker{
		db: db,
	}
}

// TradeRecord represents a completed trade
type TradeRecord struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id"`
	Symbol         string                 `json:"symbol"`
	Side           string                 `json:"side"` // buy, sell
	Quantity       decimal.Decimal        `json:"quantity"`
	EntryPrice     decimal.Decimal        `json:"entry_price"`
	ExitPrice      decimal.Decimal        `json:"exit_price"`
	EntryTimestamp time.Time              `json:"entry_timestamp"`
	ExitTimestamp  time.Time              `json:"exit_timestamp"`
	PnL            decimal.Decimal        `json:"pnl"`             // Profit and Loss
	Fees           decimal.Decimal        `json:"fees"`            // Trading fees paid to exchange
	PerformanceFee decimal.Decimal        `json:"performance_fee"` // Our performance fee
	Status         string                 `json:"status"`          // completed, pending
	StrategyID     string                 `json:"strategy_id"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// PerformanceFeeConfig defines fee structure
type PerformanceFeeConfig struct {
	UserID               string          `json:"user_id"`
	FeePercentage        decimal.Decimal `json:"fee_percentage"`         // 2-20%
	HighWaterMark        decimal.Decimal `json:"high_water_mark"`        // Highest portfolio value
	MinimumFee           decimal.Decimal `json:"minimum_fee"`            // Minimum fee per trade
	MaximumFee           decimal.Decimal `json:"maximum_fee"`            // Maximum fee per trade
	FeeFrequency         string          `json:"fee_frequency"`          // per_trade, monthly, quarterly
	OnlyProfitableTrades bool            `json:"only_profitable_trades"` // Only charge on profitable trades
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// PerformanceSummary represents performance metrics for a user
type PerformanceSummary struct {
	UserID               string          `json:"user_id"`
	Period               string          `json:"period"`
	TotalTrades          int64           `json:"total_trades"`
	ProfitableTrades     int64           `json:"profitable_trades"`
	TotalPnL             decimal.Decimal `json:"total_pnl"`
	TotalPerformanceFees decimal.Decimal `json:"total_performance_fees"`
	WinRate              decimal.Decimal `json:"win_rate"`
	AverageReturn        decimal.Decimal `json:"average_return"`
	MaxDrawdown          decimal.Decimal `json:"max_drawdown"`
	SharpeRatio          decimal.Decimal `json:"sharpe_ratio"`
	CurrentHighWaterMark decimal.Decimal `json:"current_high_water_mark"`
	LastUpdated          time.Time       `json:"last_updated"`
}

// RecordTrade records a completed trade and calculates performance fee
func (pft *PerformanceFeeTracker) RecordTrade(ctx context.Context, trade *TradeRecord) error {
	// Get user's fee configuration
	feeConfig, err := pft.GetFeeConfig(ctx, trade.UserID)
	if err != nil {
		return fmt.Errorf("failed to get fee config: %v", err)
	}

	// Calculate performance fee
	performanceFee, err := pft.calculatePerformanceFee(ctx, trade, feeConfig)
	if err != nil {
		return fmt.Errorf("failed to calculate performance fee: %v", err)
	}

	trade.PerformanceFee = performanceFee
	trade.Status = "completed"

	// Store trade record
	err = pft.storeTrade(ctx, trade)
	if err != nil {
		return fmt.Errorf("failed to store trade: %v", err)
	}

	// Update high water mark if necessary
	err = pft.updateHighWaterMark(ctx, trade.UserID, trade.PnL)
	if err != nil {
		return fmt.Errorf("failed to update high water mark: %v", err)
	}

	// Update performance summary
	err = pft.updatePerformanceSummary(ctx, trade.UserID)
	if err != nil {
		return fmt.Errorf("failed to update performance summary: %v", err)
	}

	return nil
}

// calculatePerformanceFee calculates the performance fee for a trade
func (pft *PerformanceFeeTracker) calculatePerformanceFee(ctx context.Context, trade *TradeRecord, config *PerformanceFeeConfig) (decimal.Decimal, error) {
	// Only charge on profitable trades if configured
	if config.OnlyProfitableTrades && trade.PnL.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, nil
	}

	// Calculate base fee
	var feeBase decimal.Decimal
	if trade.PnL.GreaterThan(decimal.Zero) {
		feeBase = trade.PnL
	} else {
		// For losing trades, we don't charge performance fees
		return decimal.Zero, nil
	}

	// Apply fee percentage
	fee := feeBase.Mul(config.FeePercentage.Div(decimal.NewFromInt(100)))

	// Apply minimum and maximum fee limits
	if fee.LessThan(config.MinimumFee) {
		fee = config.MinimumFee
	}
	if config.MaximumFee.GreaterThan(decimal.Zero) && fee.GreaterThan(config.MaximumFee) {
		fee = config.MaximumFee
	}

	// High water mark logic - only charge fees if portfolio is above previous high
	currentPortfolioValue, err := pft.getCurrentPortfolioValue(ctx, trade.UserID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get current portfolio value: %v", err)
	}

	if currentPortfolioValue.LessThanOrEqual(config.HighWaterMark) {
		return decimal.Zero, nil // No fee if below high water mark
	}

	return fee, nil
}

// GetFeeConfig returns fee configuration for a user
func (pft *PerformanceFeeTracker) GetFeeConfig(ctx context.Context, userID string) (*PerformanceFeeConfig, error) {
	query := `
		SELECT user_id, fee_percentage, high_water_mark, minimum_fee, maximum_fee,
		       fee_frequency, only_profitable_trades, created_at, updated_at
		FROM performance_fee_configs 
		WHERE user_id = $1
	`

	config := &PerformanceFeeConfig{}
	err := pft.db.QueryRowContext(ctx, query, userID).Scan(
		&config.UserID, &config.FeePercentage, &config.HighWaterMark,
		&config.MinimumFee, &config.MaximumFee, &config.FeeFrequency,
		&config.OnlyProfitableTrades, &config.CreatedAt, &config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create default config for new user
		return pft.createDefaultFeeConfig(ctx, userID)
	}

	return config, err
}

// createDefaultFeeConfig creates default fee configuration for new users
func (pft *PerformanceFeeTracker) createDefaultFeeConfig(ctx context.Context, userID string) (*PerformanceFeeConfig, error) {
	config := &PerformanceFeeConfig{
		UserID:               userID,
		FeePercentage:        decimal.NewFromFloat(20.0),   // 20% default
		HighWaterMark:        decimal.NewFromFloat(10000),  // $10,000 starting value
		MinimumFee:           decimal.NewFromFloat(1.0),    // $1 minimum
		MaximumFee:           decimal.NewFromFloat(1000.0), // $1,000 maximum
		FeeFrequency:         "per_trade",
		OnlyProfitableTrades: true,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	query := `
		INSERT INTO performance_fee_configs (
			user_id, fee_percentage, high_water_mark, minimum_fee, maximum_fee,
			fee_frequency, only_profitable_trades, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := pft.db.ExecContext(ctx, query,
		config.UserID, config.FeePercentage, config.HighWaterMark,
		config.MinimumFee, config.MaximumFee, config.FeeFrequency,
		config.OnlyProfitableTrades, config.CreatedAt, config.UpdatedAt,
	)

	return config, err
}

// storeTrade stores a trade record in the database
func (pft *PerformanceFeeTracker) storeTrade(ctx context.Context, trade *TradeRecord) error {
	query := `
		INSERT INTO trade_records (
			id, user_id, symbol, side, quantity, entry_price, exit_price,
			entry_timestamp, exit_timestamp, pnl, fees, performance_fee,
			status, strategy_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := pft.db.ExecContext(ctx, query,
		trade.ID, trade.UserID, trade.Symbol, trade.Side, trade.Quantity,
		trade.EntryPrice, trade.ExitPrice, trade.EntryTimestamp, trade.ExitTimestamp,
		trade.PnL, trade.Fees, trade.PerformanceFee, trade.Status, trade.StrategyID,
	)

	return err
}

// updateHighWaterMark updates the high water mark for a user
func (pft *PerformanceFeeTracker) updateHighWaterMark(ctx context.Context, userID string, pnl decimal.Decimal) error {
	currentValue, err := pft.getCurrentPortfolioValue(ctx, userID)
	if err != nil {
		return err
	}

	query := `
		UPDATE performance_fee_configs 
		SET high_water_mark = GREATEST(high_water_mark, $1), updated_at = $2
		WHERE user_id = $3
	`

	_, err = pft.db.ExecContext(ctx, query, currentValue, time.Now(), userID)
	return err
}

// getCurrentPortfolioValue calculates current portfolio value for a user
func (pft *PerformanceFeeTracker) getCurrentPortfolioValue(ctx context.Context, userID string) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(pnl), 0) + 10000 as portfolio_value
		FROM trade_records 
		WHERE user_id = $1 AND status = 'completed'
	`

	var portfolioValue decimal.Decimal
	err := pft.db.QueryRowContext(ctx, query, userID).Scan(&portfolioValue)
	return portfolioValue, err
}

// updatePerformanceSummary updates performance metrics for a user
func (pft *PerformanceFeeTracker) updatePerformanceSummary(ctx context.Context, userID string) error {
	// This would calculate and update comprehensive performance metrics
	// Implementation would include win rate, Sharpe ratio, max drawdown, etc.

	query := `
		INSERT INTO performance_summaries (
			user_id, period, total_trades, profitable_trades, total_pnl,
			total_performance_fees, win_rate, last_updated
		)
		SELECT 
			$1 as user_id,
			'all_time' as period,
			COUNT(*) as total_trades,
			COUNT(CASE WHEN pnl > 0 THEN 1 END) as profitable_trades,
			COALESCE(SUM(pnl), 0) as total_pnl,
			COALESCE(SUM(performance_fee), 0) as total_performance_fees,
			CASE WHEN COUNT(*) > 0 THEN 
				COUNT(CASE WHEN pnl > 0 THEN 1 END)::decimal / COUNT(*)::decimal 
			ELSE 0 END as win_rate,
			NOW() as last_updated
		FROM trade_records 
		WHERE user_id = $1 AND status = 'completed'
		ON CONFLICT (user_id, period) 
		DO UPDATE SET
			total_trades = EXCLUDED.total_trades,
			profitable_trades = EXCLUDED.profitable_trades,
			total_pnl = EXCLUDED.total_pnl,
			total_performance_fees = EXCLUDED.total_performance_fees,
			win_rate = EXCLUDED.win_rate,
			last_updated = EXCLUDED.last_updated
	`

	_, err := pft.db.ExecContext(ctx, query, userID)
	return err
}

// GetPerformanceSummary returns performance summary for a user
func (pft *PerformanceFeeTracker) GetPerformanceSummary(ctx context.Context, userID, period string) (*PerformanceSummary, error) {
	query := `
		SELECT user_id, period, total_trades, profitable_trades, total_pnl,
		       total_performance_fees, win_rate, average_return, max_drawdown,
		       sharpe_ratio, current_high_water_mark, last_updated
		FROM performance_summaries
		WHERE user_id = $1 AND period = $2
	`

	summary := &PerformanceSummary{}
	err := pft.db.QueryRowContext(ctx, query, userID, period).Scan(
		&summary.UserID, &summary.Period, &summary.TotalTrades,
		&summary.ProfitableTrades, &summary.TotalPnL, &summary.TotalPerformanceFees,
		&summary.WinRate, &summary.AverageReturn, &summary.MaxDrawdown,
		&summary.SharpeRatio, &summary.CurrentHighWaterMark, &summary.LastUpdated,
	)

	return summary, err
}

// GenerateMonthlyPerformanceBill generates monthly performance fee bill
func (pft *PerformanceFeeTracker) GenerateMonthlyPerformanceBill(ctx context.Context, userID string, billingPeriod time.Time) (*PerformanceBill, error) {
	// Get trades for the billing period
	startDate := billingPeriod
	endDate := billingPeriod.AddDate(0, 1, 0)

	query := `
		SELECT COUNT(*),
		       COUNT(CASE WHEN pnl > 0 THEN 1 END),
		       COALESCE(SUM(pnl), 0),
		       COALESCE(SUM(performance_fee), 0)
		FROM trade_records
		WHERE user_id = $1 AND entry_timestamp >= $2 AND entry_timestamp < $3 AND status = 'completed'
	`

	var totalTrades, profitableTrades int64
	var totalPnL, totalFees decimal.Decimal

	err := pft.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(
		&totalTrades, &profitableTrades, &totalPnL, &totalFees,
	)
	if err != nil {
		return nil, err
	}

	// Get fee configuration
	config, err := pft.GetFeeConfig(ctx, userID)
	if err != nil {
		return nil, err
	}

	bill := &PerformanceBill{
		UserID:               userID,
		BillingPeriod:        billingPeriod,
		TotalTrades:          totalTrades,
		ProfitableTrades:     profitableTrades,
		TotalPnL:             totalPnL,
		TotalPerformanceFees: totalFees,
		FeeRate:              config.FeePercentage,
		HighWaterMarkStart:   config.HighWaterMark,
		HighWaterMarkEnd:     config.HighWaterMark, // Would be updated during the period
		Status:               "pending",
		GeneratedAt:          time.Now(),
		DueDate:              time.Now().AddDate(0, 0, 30),
	}

	return bill, nil
}

// PerformanceBill represents a performance fee bill
type PerformanceBill struct {
	ID                   string          `json:"id"`
	UserID               string          `json:"user_id"`
	BillingPeriod        time.Time       `json:"billing_period"`
	TotalTrades          int64           `json:"total_trades"`
	ProfitableTrades     int64           `json:"profitable_trades"`
	TotalPnL             decimal.Decimal `json:"total_pnl"`
	TotalPerformanceFees decimal.Decimal `json:"total_performance_fees"`
	FeeRate              decimal.Decimal `json:"fee_rate"`
	HighWaterMarkStart   decimal.Decimal `json:"high_water_mark_start"`
	HighWaterMarkEnd     decimal.Decimal `json:"high_water_mark_end"`
	Status               string          `json:"status"`
	StripeInvoiceID      string          `json:"stripe_invoice_id,omitempty"`
	GeneratedAt          time.Time       `json:"generated_at"`
	PaidAt               *time.Time      `json:"paid_at,omitempty"`
	DueDate              time.Time       `json:"due_date"`
}

// GetTradeHistory returns trade history for a user
func (pft *PerformanceFeeTracker) GetTradeHistory(ctx context.Context, userID string, limit int) ([]*TradeRecord, error) {
	query := `
		SELECT id, user_id, symbol, side, quantity, entry_price, exit_price,
		       entry_timestamp, exit_timestamp, pnl, fees, performance_fee,
		       status, strategy_id
		FROM trade_records
		WHERE user_id = $1
		ORDER BY entry_timestamp DESC
		LIMIT $2
	`

	rows, err := pft.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*TradeRecord
	for rows.Next() {
		trade := &TradeRecord{}
		err := rows.Scan(
			&trade.ID, &trade.UserID, &trade.Symbol, &trade.Side,
			&trade.Quantity, &trade.EntryPrice, &trade.ExitPrice,
			&trade.EntryTimestamp, &trade.ExitTimestamp, &trade.PnL,
			&trade.Fees, &trade.PerformanceFee, &trade.Status, &trade.StrategyID,
		)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
