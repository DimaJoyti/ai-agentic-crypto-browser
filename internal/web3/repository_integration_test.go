package web3

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Integration tests for Postgres-backed repositories using testcontainers.
func TestRepositoryIntegration(t *testing.T) {
	ctx := context.Background()

	// Start Postgres container
	pgReq := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(60 * time.Second),
	}
	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: pgReq, Started: true})
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgC.Terminate(ctx) })

	host, err := pgC.Host(ctx)
	require.NoError(t, err)
	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)
	dsn := "postgres://postgres:postgres@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	db, err := database.NewPostgresDB(config.DatabaseConfig{URL: dsn, MaxOpenConns: 10, MaxIdleConns: 5, ConnMaxLifetime: time.Minute})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	// Create minimal schema required for repositories
	_, err = db.ExecWithMetrics(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	require.NoError(t, err)
	_, err = db.ExecWithMetrics(ctx, `CREATE TABLE IF NOT EXISTS web3_wallets (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL,
		address VARCHAR(42) NOT NULL,
		chain_id INTEGER NOT NULL,
		wallet_type VARCHAR(50) NOT NULL,
		is_primary BOOLEAN DEFAULT false,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(user_id, address, chain_id)
	);`)
	require.NoError(t, err)
	_, err = db.ExecWithMetrics(ctx, `CREATE TABLE IF NOT EXISTS web3_transactions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL,
		wallet_id UUID NOT NULL,
		tx_hash VARCHAR(66) NOT NULL,
		chain_id INTEGER NOT NULL,
		from_address VARCHAR(42) NOT NULL,
		to_address VARCHAR(42),
		value DECIMAL(36, 18),
		gas_used BIGINT,
		gas_price DECIMAL(36, 18),
		status VARCHAR(20) DEFAULT 'pending',
		block_number BIGINT,
		transaction_type VARCHAR(50),
		metadata JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`)
	require.NoError(t, err)

	walletRepo := NewPostgresWalletRepository(db)
	txRepo := NewPostgresTransactionRepository(db)

	userID := uuid.New()
	walletID := uuid.New()
	w := &Wallet{ID: walletID, UserID: userID, Address: "0xabc", ChainID: 1, WalletType: "metamask", IsPrimary: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, walletRepo.Save(ctx, w))

	// GetByID
	got, err := walletRepo.GetByID(ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, w.Address, got.Address)

	// GetByAddress
	got2, err := walletRepo.GetByAddress(ctx, userID, "0xabc", 1)
	require.NoError(t, err)
	require.Equal(t, walletID, got2.ID)

	// ListByUser
	list, pg, err := walletRepo.ListByUser(ctx, userID, WalletListFilter{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Equal(t, 1, len(list))
	require.Equal(t, 1, pg.TotalItems)

	// Upsert duplicate wallet with new wallet type
	w2 := &Wallet{ID: uuid.New(), UserID: userID, Address: "0xabc", ChainID: 1, WalletType: "rabby", IsPrimary: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, walletRepo.Save(ctx, w2))
	got3, err := walletRepo.GetByAddress(ctx, userID, "0xabc", 1)
	require.NoError(t, err)
	require.Equal(t, "rabby", got3.WalletType)

	// Transactions
	txID := uuid.New()
	tran := &Transaction{ID: txID, UserID: userID, WalletID: walletID, TxHash: "0xhash", ChainID: 1, FromAddress: "0xabc", ToAddress: "0xdef", Value: nil, Status: TxStatusPending, TransactionType: "transfer", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, txRepo.Save(ctx, tran))

	// ListByUser
	trs, pg2, err := txRepo.ListByUser(ctx, userID, TransactionListFilter{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Equal(t, 1, len(trs))
	require.Equal(t, 1, pg2.TotalItems)

	// UpdateStatus
	require.NoError(t, txRepo.UpdateStatus(ctx, txID, TxStatusConfirmed))
	gotTx, err := txRepo.GetByID(ctx, txID)
	require.NoError(t, err)
	require.Equal(t, TxStatusConfirmed, gotTx.Status)
}
