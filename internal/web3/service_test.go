package web3

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

type mockWalletRepo struct {
	saveErr      error
	getByID      map[uuid.UUID]*Wallet
	getByAddress map[string]*Wallet
	countByUser  int
	listResult   []*Wallet
}

func (m *mockWalletRepo) Save(ctx context.Context, w *Wallet) error { return m.saveErr }
func (m *mockWalletRepo) GetByID(ctx context.Context, id uuid.UUID) (*Wallet, error) {
	if w, ok := m.getByID[id]; ok {
		return w, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *mockWalletRepo) GetByAddress(ctx context.Context, userID uuid.UUID, address string, chainID int) (*Wallet, error) {
	key := address
	if w, ok := m.getByAddress[key]; ok {
		return w, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *mockWalletRepo) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	return m.countByUser, nil
}
func (m *mockWalletRepo) ListByUser(ctx context.Context, userID uuid.UUID, filter WalletListFilter) ([]*Wallet, Pagination, error) {
	return m.listResult, Pagination{Page: 1, PageSize: len(m.listResult), TotalItems: len(m.listResult), TotalPages: 1}, nil
}
func (m *mockWalletRepo) SetPrimary(ctx context.Context, userID uuid.UUID, walletID uuid.UUID) error {
	return nil
}

type mockTxRepo struct {
	saveErr error
	list    []*Transaction
}

func (m *mockTxRepo) Save(ctx context.Context, t *Transaction) error { return m.saveErr }
func (m *mockTxRepo) GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	return nil, nil
}
func (m *mockTxRepo) ListByUser(ctx context.Context, userID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error) {
	return m.list, Pagination{Page: 1, PageSize: len(m.list), TotalItems: len(m.list), TotalPages: 1}, nil
}
func (m *mockTxRepo) ListByWallet(ctx context.Context, walletID uuid.UUID, filter TransactionListFilter) ([]*Transaction, Pagination, error) {
	return m.list, Pagination{Page: 1, PageSize: len(m.list), TotalItems: len(m.list), TotalPages: 1}, nil
}
func (m *mockTxRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error { return nil }

// construct service with mocks
func newServiceWithMocks() *Service {
	return &Service{
		config:     config.Web3Config{},
		logger:     observability.NewLogger(config.ObservabilityConfig{ServiceName: "web3-test", LogLevel: "error", LogFormat: "json"}),
		providers:  map[int]*ChainProvider{1: {ChainID: 1, RpcURL: "http://localhost"}},
		walletRepo: &mockWalletRepo{},
		txRepo:     &mockTxRepo{},
	}
}

func TestConnectWallet_SetsPrimaryOnFirst(t *testing.T) {
	s := newServiceWithMocks()
	mw := s.walletRepo.(*mockWalletRepo)
	mw.countByUser = 0

	userID := uuid.New()
	resp, err := s.ConnectWallet(context.Background(), userID, WalletConnectRequest{Address: "0xabc", ChainID: 1, WalletType: "metamask"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Wallet == nil || !resp.Wallet.IsPrimary {
		t.Fatalf("expected primary wallet on first connect")
	}
}

func TestListWallets(t *testing.T) {
	s := newServiceWithMocks()
	mw := s.walletRepo.(*mockWalletRepo)
	mw.listResult = []*Wallet{{ID: uuid.New(), UserID: uuid.New(), Address: "0xabc", ChainID: 1, WalletType: "metamask", CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	ws, pg, err := s.ListWallets(context.Background(), uuid.New(), WalletListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ws) != 1 || pg.TotalItems != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(ws))
	}
}

func TestCreateTransaction_Persists(t *testing.T) {
	s := newServiceWithMocks()
	mw := s.walletRepo.(*mockWalletRepo)
	walletID := uuid.New()
	userID := uuid.New()
	mw.getByID = map[uuid.UUID]*Wallet{walletID: {ID: walletID, UserID: userID, Address: "0xabc", ChainID: 1}}

	_, err := s.CreateTransaction(context.Background(), userID, TransactionRequest{WalletID: walletID, ToAddress: "0xdef", Value: big.NewInt(1)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTransactions(t *testing.T) {
	s := newServiceWithMocks()
	mt := s.txRepo.(*mockTxRepo)
	mt.list = []*Transaction{{ID: uuid.New(), UserID: uuid.New(), WalletID: uuid.New(), TxHash: "0x1", ChainID: 1, FromAddress: "0xabc", CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	list, pg, err := s.ListTransactions(context.Background(), uuid.New(), TransactionListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || pg.TotalItems != 1 {
		t.Fatalf("expected 1 transaction")
	}
}
