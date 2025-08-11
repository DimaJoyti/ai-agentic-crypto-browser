package billing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// PerformanceFeeManager handles performance-based fee calculations
type PerformanceFeeManager struct {
	db     *sql.DB
	config *PerformanceFeeConfig
}

// PerformanceFeeConfig defines fee structure
type PerformanceFeeConfig struct {
	FeePercentage     decimal.Decimal `json:"fee_percentage"`     // e.g., 0.20 for 20%
	HighWaterMark     bool            `json:"high_water_mark"`    // Only charge on new highs
	MinimumProfit     decimal.Decimal `json:"minimum_profit"`     // Minimum profit to charge fees
	CalculationPeriod string          `json:"calculation_period"` // daily, weekly, monthly
}

// PerformanceFeeRecord tracks fee calculations
type PerformanceFeeRecord struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	StrategyID      string          `json:"strategy_id"`
	PeriodStart     time.Time       `json:"period_start"`
	PeriodEnd       time.Time       `json:"period_end"`
	StartingBalance decimal.Decimal `json:"starting_balance"`
	EndingBalance   decimal.Decimal `json:"ending_balance"`
	ProfitLoss      decimal.Decimal `json:"profit_loss"`
	FeeableProfit   decimal.Decimal `json:"feeable_profit"`
	FeeAmount       decimal.Decimal `json:"fee_amount"`
	FeePercentage   decimal.Decimal `json:"fee_percentage"`
	HighWaterMark   decimal.Decimal `json:"high_water_mark"`
	Status          string          `json:"status"` // pending, calculated, charged, paid
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// NewPerformanceFeeManager creates a new performance fee manager
func NewPerformanceFeeManager(db *sql.DB, config *PerformanceFeeConfig) *PerformanceFeeManager {
	return &PerformanceFeeManager{
		db:     db,
		config: config,
	}
}

// CalculatePerformanceFee calculates fees for a given period
func (pfm *PerformanceFeeManager) CalculatePerformanceFee(ctx context.Context, userID, strategyID string, periodStart, periodEnd time.Time) (*PerformanceFeeRecord, error) {
	// Get starting and ending balances
	startBalance, err := pfm.getPortfolioValue(ctx, userID, strategyID, periodStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get starting balance: %w", err)
	}

	endBalance, err := pfm.getPortfolioValue(ctx, userID, strategyID, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get ending balance: %w", err)
	}

	// Calculate profit/loss
	profitLoss := endBalance.Sub(startBalance)

	// Get high water mark if enabled
	var highWaterMark decimal.Decimal
	var feeableProfit decimal.Decimal

	if pfm.config.HighWaterMark {
		hwm, err := pfm.getHighWaterMark(ctx, userID, strategyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get high water mark: %w", err)
		}
		highWaterMark = hwm

		// Only charge fees on profits above high water mark
		if endBalance.GreaterThan(highWaterMark) {
			feeableProfit = endBalance.Sub(highWaterMark)
		} else {
			feeableProfit = decimal.Zero
		}
	} else {
		// Charge fees on all profits
		if profitLoss.GreaterThan(decimal.Zero) {
			feeableProfit = profitLoss
		} else {
			feeableProfit = decimal.Zero
		}
	}

	// Apply minimum profit threshold
	if feeableProfit.LessThan(pfm.config.MinimumProfit) {
		feeableProfit = decimal.Zero
	}

	// Calculate fee amount
	feeAmount := feeableProfit.Mul(pfm.config.FeePercentage)

	// Create performance fee record
	record := &PerformanceFeeRecord{
		ID:              generateID(),
		UserID:          userID,
		StrategyID:      strategyID,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		StartingBalance: startBalance,
		EndingBalance:   endBalance,
		ProfitLoss:      profitLoss,
		FeeableProfit:   feeableProfit,
		FeeAmount:       feeAmount,
		FeePercentage:   pfm.config.FeePercentage,
		HighWaterMark:   highWaterMark,
		Status:          "calculated",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save to database
	err = pfm.savePerformanceFeeRecord(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to save performance fee record: %w", err)
	}

	// Update high water mark if applicable
	if pfm.config.HighWaterMark && endBalance.GreaterThan(highWaterMark) {
		err = pfm.updateHighWaterMark(ctx, userID, strategyID, endBalance)
		if err != nil {
			return nil, fmt.Errorf("failed to update high water mark: %w", err)
		}
	}

	return record, nil
}

// ChargePerformanceFee processes the actual fee charge
func (pfm *PerformanceFeeManager) ChargePerformanceFee(ctx context.Context, recordID string) error {
	record, err := pfm.getPerformanceFeeRecord(ctx, recordID)
	if err != nil {
		return fmt.Errorf("failed to get performance fee record: %w", err)
	}

	if record.Status != "calculated" {
		return fmt.Errorf("performance fee record not in calculated status")
	}

	if record.FeeAmount.LessThanOrEqual(decimal.Zero) {
		// No fee to charge
		record.Status = "paid"
		return pfm.updatePerformanceFeeRecord(ctx, record)
	}

	// Deduct fee from user's portfolio
	err = pfm.deductFeeFromPortfolio(ctx, record.UserID, record.StrategyID, record.FeeAmount)
	if err != nil {
		return fmt.Errorf("failed to deduct fee from portfolio: %w", err)
	}

	// Update record status
	record.Status = "charged"
	record.UpdatedAt = time.Now()

	return pfm.updatePerformanceFeeRecord(ctx, record)
}

// GetPerformanceFeeHistory returns fee history for a user
func (pfm *PerformanceFeeManager) GetPerformanceFeeHistory(ctx context.Context, userID string, limit int) ([]*PerformanceFeeRecord, error) {
	query := `
		SELECT id, user_id, strategy_id, period_start, period_end, 
		       starting_balance, ending_balance, profit_loss, feeable_profit,
		       fee_amount, fee_percentage, high_water_mark, status, 
		       created_at, updated_at
		FROM performance_fees 
		WHERE user_id = $1 
		ORDER BY period_end DESC 
		LIMIT $2
	`

	rows, err := pfm.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*PerformanceFeeRecord
	for rows.Next() {
		record := &PerformanceFeeRecord{}
		err := rows.Scan(
			&record.ID, &record.UserID, &record.StrategyID,
			&record.PeriodStart, &record.PeriodEnd,
			&record.StartingBalance, &record.EndingBalance,
			&record.ProfitLoss, &record.FeeableProfit,
			&record.FeeAmount, &record.FeePercentage,
			&record.HighWaterMark, &record.Status,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// Helper functions
func (pfm *PerformanceFeeManager) getPortfolioValue(ctx context.Context, userID, strategyID string, timestamp time.Time) (decimal.Decimal, error) {
	// Implementation to get portfolio value at specific timestamp
	// This would integrate with your existing portfolio tracking system
	return decimal.Zero, nil
}

func (pfm *PerformanceFeeManager) getHighWaterMark(ctx context.Context, userID, strategyID string) (decimal.Decimal, error) {
	// Implementation to get current high water mark
	return decimal.Zero, nil
}

func (pfm *PerformanceFeeManager) updateHighWaterMark(ctx context.Context, userID, strategyID string, newMark decimal.Decimal) error {
	// Implementation to update high water mark
	return nil
}

func (pfm *PerformanceFeeManager) savePerformanceFeeRecord(ctx context.Context, record *PerformanceFeeRecord) error {
	// Implementation to save performance fee record to database
	return nil
}

func (pfm *PerformanceFeeManager) getPerformanceFeeRecord(ctx context.Context, recordID string) (*PerformanceFeeRecord, error) {
	// Implementation to get performance fee record from database
	return nil, nil
}

func (pfm *PerformanceFeeManager) updatePerformanceFeeRecord(ctx context.Context, record *PerformanceFeeRecord) error {
	// Implementation to update performance fee record
	return nil
}

func (pfm *PerformanceFeeManager) deductFeeFromPortfolio(ctx context.Context, userID, strategyID string, amount decimal.Decimal) error {
	// Implementation to deduct fee from user's portfolio
	return nil
}

func generateID() string {
	// Implementation to generate unique ID
	return fmt.Sprintf("pf_%d", time.Now().UnixNano())
}
