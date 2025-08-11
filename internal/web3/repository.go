package web3

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Pagination encapsulates paging parameters and results metadata
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// WalletListFilter defines filters for listing wallets
type WalletListFilter struct {
	ChainID   int  // 0 means all
	IsPrimary *bool
	Page      int
	PageSize  int
}

// TransactionListFilter defines filters for listing transactions
type TransactionListFilter struct {
	WalletID   uuid.UUID
	ChainID    int    // 0 means all
	Status     string // optional: pending|confirmed|failed
	FromTime   *time.Time
	ToTime     *time.Time
	Page       int
	PageSize   int
}

// WalletRepository abstracts wallet persistence
// Small, focused interface to enable mocking and clean separation
// of domain and data access.
type WalletRepository interface {
	Save(ctx context.Context, w *Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*Wallet, error)
	GetByAddress(ctx context.Context, userID uuid.UUID, address string, chainID int) (*Wallet, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, filter WalletListFilter) ([]*Wallet, Pagination, error)
	SetPrimary(ctx context.Context, userID uuid.UUID, walletID uuid.UUID) error
}

// TransactionRepository abstracts transaction persistence
type TransactionRepository interface {
	Save(ctx context.Context, t *Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	ListByUser(ctx context.Context, userID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error)
	ListByWallet(ctx context.Context, walletID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

