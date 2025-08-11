package web3

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/database"
	"github.com/google/uuid"
)

// postgresWalletRepository implements WalletRepository using Postgres
type postgresWalletRepository struct {
	db *database.DB
}

func NewPostgresWalletRepository(db *database.DB) WalletRepository {
	return &postgresWalletRepository{db: db}
}

func (r *postgresWalletRepository) Save(ctx context.Context, w *Wallet) error {
	query := `
		INSERT INTO web3_wallets (id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, address, chain_id) DO UPDATE SET
		  wallet_type = EXCLUDED.wallet_type,
		  updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecWithMetrics(ctx, query, w.ID, w.UserID, w.Address, w.ChainID, w.WalletType, w.IsPrimary, w.CreatedAt, w.UpdatedAt)
	return err
}

func (r *postgresWalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*Wallet, error) {
	query := `SELECT id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at FROM web3_wallets WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	w := &Wallet{}
	if err := row.Scan(&w.ID, &w.UserID, &w.Address, &w.ChainID, &w.WalletType, &w.IsPrimary, &w.CreatedAt, &w.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("wallet not found: %w", err)
		}
		return nil, err
	}
	return w, nil
}

func (r *postgresWalletRepository) GetByAddress(ctx context.Context, userID uuid.UUID, address string, chainID int) (*Wallet, error) {
	query := `SELECT id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at FROM web3_wallets WHERE user_id = $1 AND address = $2 AND chain_id = $3`
	row := r.db.QueryRowContext(ctx, query, userID, strings.ToLower(address), chainID)
	w := &Wallet{}
	if err := row.Scan(&w.ID, &w.UserID, &w.Address, &w.ChainID, &w.WalletType, &w.IsPrimary, &w.CreatedAt, &w.UpdatedAt); err != nil {
		return nil, err
	}
	return w, nil
}

func (r *postgresWalletRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM web3_wallets WHERE user_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func paginate(page, pageSize int) (limit, offset int) {
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

func buildPagination(total, page, pageSize int) Pagination {
	if pageSize <= 0 {
		pageSize = 20
	}
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	return Pagination{Page: page, PageSize: pageSize, TotalItems: total, TotalPages: totalPages}
}

func (r *postgresWalletRepository) ListByUser(ctx context.Context, userID uuid.UUID, filter WalletListFilter) ([]*Wallet, Pagination, error) {
	var args []any
	where := []string{"user_id = $1"}
	args = append(args, userID)
	argPos := 2
	if filter.ChainID != 0 {
		where = append(where, fmt.Sprintf("chain_id = $%d", argPos))
		args = append(args, filter.ChainID)
		argPos++
	}
	if filter.IsPrimary != nil {
		where = append(where, fmt.Sprintf("is_primary = $%d", argPos))
		args = append(args, *filter.IsPrimary)
		argPos++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM web3_wallets WHERE %s", strings.Join(where, " AND "))
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, Pagination{}, err
	}

	limit, offset := paginate(filter.Page, filter.PageSize)
	listQuery := fmt.Sprintf(`
		SELECT id, user_id, address, chain_id, wallet_type, is_primary, created_at, updated_at
		FROM web3_wallets
		WHERE %s
		ORDER BY created_at DESC
		LIMIT %d OFFSET %d
	`, strings.Join(where, " AND "), limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, Pagination{}, err
	}
	defer rows.Close()

	var result []*Wallet
	for rows.Next() {
		w := &Wallet{}
		if err := rows.Scan(&w.ID, &w.UserID, &w.Address, &w.ChainID, &w.WalletType, &w.IsPrimary, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, Pagination{}, err
		}
		result = append(result, w)
	}
	return result, buildPagination(total, filter.Page, filter.PageSize), nil
}

func (r *postgresWalletRepository) SetPrimary(ctx context.Context, userID uuid.UUID, walletID uuid.UUID) error {
	// Set all user wallets to non-primary then set selected wallet to primary (transactional)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "UPDATE web3_wallets SET is_primary = false WHERE user_id = $1", userID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE web3_wallets SET is_primary = true WHERE id = $1 AND user_id = $2", walletID, userID); err != nil {
		return err
	}
	return tx.Commit()
}

// postgresTransactionRepository implements TransactionRepository using Postgres
type postgresTransactionRepository struct {
	db *database.DB
}

func NewPostgresTransactionRepository(db *database.DB) TransactionRepository {
	return &postgresTransactionRepository{db: db}
}

func (r *postgresTransactionRepository) Save(ctx context.Context, t *Transaction) error {
	metadataJSON, _ := jsonMarshalSafe(t.Metadata)
	query := `
		INSERT INTO web3_transactions (
		  id, user_id, wallet_id, tx_hash, chain_id, from_address, to_address, value, gas_used, gas_price,
		  status, block_number, transaction_type, metadata, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`
	_, err := r.db.ExecWithMetrics(ctx, query,
		t.ID, t.UserID, t.WalletID, t.TxHash, t.ChainID, t.FromAddress, t.ToAddress, t.Value,
		t.GasUsed, t.GasPrice, t.Status, t.BlockNumber, t.TransactionType, metadataJSON, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *postgresTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	query := `
		SELECT id, user_id, wallet_id, tx_hash, chain_id, from_address, to_address, value, gas_used, gas_price,
		       status, block_number, transaction_type, metadata, created_at, updated_at
		FROM web3_transactions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanTransaction(row)
}

func (r *postgresTransactionRepository) ListByUser(ctx context.Context, userID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error) {
	filter.WalletID = uuid.Nil // ensure by-user query
	return r.list(ctx, &listScope{UserID: &userID}, filter)
}

func (r *postgresTransactionRepository) ListByWallet(ctx context.Context, walletID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error) {
	return r.list(ctx, &listScope{WalletID: &walletID}, filter)
}

type listScope struct {
	UserID   *uuid.UUID
	WalletID *uuid.UUID
}

func (r *postgresTransactionRepository) list(ctx context.Context, scope *listScope, filter TransactionListFilter) ([]*Transaction, Pagination, error) {
	var where []string
	var args []any
	pos := 1

	if scope.UserID != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", pos))
		args = append(args, *scope.UserID)
		pos++
	}
	if scope.WalletID != nil {
		where = append(where, fmt.Sprintf("wallet_id = $%d", pos))
		args = append(args, *scope.WalletID)
		pos++
	}
	if filter.ChainID != 0 {
		where = append(where, fmt.Sprintf("chain_id = $%d", pos))
		args = append(args, filter.ChainID)
		pos++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", pos))
		args = append(args, filter.Status)
		pos++
	}
	if filter.FromTime != nil {
		where = append(where, fmt.Sprintf("created_at >= $%d", pos))
		args = append(args, *filter.FromTime)
		pos++
	}
	if filter.ToTime != nil {
		where = append(where, fmt.Sprintf("created_at <= $%d", pos))
		args = append(args, *filter.ToTime)
		pos++
	}

	if len(where) == 0 {
		where = append(where, "1=1")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM web3_transactions WHERE %s", strings.Join(where, " AND "))
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, Pagination{}, err
	}

	limit, offset := paginate(filter.Page, filter.PageSize)
	listQuery := fmt.Sprintf(`
		SELECT id, user_id, wallet_id, tx_hash, chain_id, from_address, to_address, value, gas_used, gas_price,
		       status, block_number, transaction_type, metadata, created_at, updated_at
		FROM web3_transactions
		WHERE %s
		ORDER BY created_at DESC
		LIMIT %d OFFSET %d`, strings.Join(where, " AND "), limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, Pagination{}, err
	}
	defer rows.Close()

	var result []*Transaction
	for rows.Next() {
		tx, err := scanTransaction(rows)
		if err != nil {
			return nil, Pagination{}, err
		}
		result = append(result, tx)
	}
	return result, buildPagination(total, filter.Page, filter.PageSize), nil
}

func (r *postgresTransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := "UPDATE web3_transactions SET status = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.ExecWithMetrics(ctx, query, status, time.Now(), id)
	return err
}

// Helpers
func scanTransaction(scanner interface{ Scan(dest ...any) error }) (*Transaction, error) {
	t := &Transaction{}
	var metadataRaw []byte
	if err := scanner.Scan(&t.ID, &t.UserID, &t.WalletID, &t.TxHash, &t.ChainID, &t.FromAddress, &t.ToAddress, &t.Value,
		&t.GasUsed, &t.GasPrice, &t.Status, &t.BlockNumber, &t.TransactionType, &metadataRaw, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found: %w", err)
		}
		return nil, err
	}
	if len(metadataRaw) > 0 {
		_ = jsonUnmarshalSafe(metadataRaw, &t.Metadata)
	}
	return t, nil
}

func jsonMarshalSafe(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("{}"), nil
	}
	return b, nil
}

func jsonUnmarshalSafe(b []byte, v any) error {
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, v)
}
