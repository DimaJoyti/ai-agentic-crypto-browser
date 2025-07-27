package web3

import (
	"math/big"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGasOptimizer(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})

	// Create gas optimizer with mock clients
	_ = NewGasOptimizer(make(map[int]*ethclient.Client), logger)

	t.Run("ValidateGasStrategies", func(t *testing.T) {
		strategies := []GasStrategy{
			GasStrategyEconomical,
			GasStrategyStandard,
			GasStrategyFast,
			GasStrategyInstant,
		}

		for _, strategy := range strategies {
			assert.NotEmpty(t, string(strategy))
		}
	})

	t.Run("GasEstimateStructure", func(t *testing.T) {
		estimate := &GasEstimate{
			GasLimit:      21000,
			GasPrice:      big.NewInt(20000000000),     // 20 gwei
			EstimatedCost: big.NewInt(420000000000000), // 21000 * 20 gwei
			Strategy:      string(GasStrategyStandard),
			Confidence:    0.85,
			TimeToConfirm: 2 * time.Minute,
		}

		assert.Equal(t, uint64(21000), estimate.GasLimit)
		assert.Equal(t, "standard", estimate.Strategy)
		assert.Equal(t, 0.85, estimate.Confidence)
	})
}

func TestIPFSService(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})

	config := IPFSConfig{
		NodeURL:     "http://localhost:5001",
		Timeout:     30 * time.Second,
		PinContent:  true,
		Gateway:     "https://ipfs.io",
		MaxFileSize: 10 * 1024 * 1024, // 10MB
	}

	ipfsService := NewIPFSService(config, logger)

	t.Run("ValidateIPFSHash", func(t *testing.T) {
		validHashes := []string{
			"QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
			"bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi",
			"bafkreihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku",
		}

		for _, hash := range validHashes {
			assert.True(t, ipfsService.isValidIPFSHash(hash), "Hash should be valid: %s", hash)
		}

		invalidHashes := []string{
			"",
			"invalid",
			"Qm", // too short
			"notahash",
		}

		for _, hash := range invalidHashes {
			assert.False(t, ipfsService.isValidIPFSHash(hash), "Hash should be invalid: %s", hash)
		}
	})

	t.Run("GenerateGatewayURL", func(t *testing.T) {
		hash := "QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG"
		expectedURL := "https://ipfs.io/ipfs/" + hash
		actualURL := ipfsService.generateGatewayURL(hash)
		assert.Equal(t, expectedURL, actualURL)
	})

	t.Run("IPFSUploadRequest", func(t *testing.T) {
		req := IPFSUploadRequest{
			Content:     []byte("test content"),
			ContentType: "text/plain",
			Filename:    "test.txt",
			Metadata:    map[string]string{"author": "test"},
			Pin:         true,
		}

		assert.Equal(t, "text/plain", req.ContentType)
		assert.Equal(t, "test.txt", req.Filename)
		assert.True(t, req.Pin)
		assert.Equal(t, "test", req.Metadata["author"])
	})
}

func TestENSResolver(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})

	// Create ENS resolver with nil client for testing
	ensResolver := NewENSResolver(nil, logger)

	t.Run("ValidateENSNames", func(t *testing.T) {
		validNames := []string{
			"vitalik.eth",
			"test.eth",
			"my-domain.eth",
			"subdomain.test.eth",
			"example.xyz",
		}

		for _, name := range validNames {
			assert.True(t, ensResolver.isValidENSName(name), "Name should be valid: %s", name)
		}

		invalidNames := []string{
			"",
			"invalid",
			"test",      // no TLD
			"test.com",  // invalid TLD
			".eth",      // empty subdomain
			"test..eth", // double dot
			"test.eth.", // trailing dot
		}

		for _, name := range invalidNames {
			assert.False(t, ensResolver.isValidENSName(name), "Name should be invalid: %s", name)
		}
	})

	t.Run("ContentHashToURL", func(t *testing.T) {
		testCases := []struct {
			contentHash string
			expectedURL string
			shouldError bool
		}{
			{
				contentHash: "ipfs://QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
				expectedURL: "https://ipfs.io/ipfs/QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
				shouldError: false,
			},
			{
				contentHash: "QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
				expectedURL: "https://ipfs.io/ipfs/QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
				shouldError: false,
			},
			{
				contentHash: "https://example.com",
				expectedURL: "https://example.com",
				shouldError: false,
			},
			{
				contentHash: "",
				expectedURL: "",
				shouldError: true,
			},
		}

		for _, tc := range testCases {
			url, err := ensResolver.contentHashToURL(tc.contentHash)
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedURL, url)
			}
		}
	})

	t.Run("ENSRecord", func(t *testing.T) {
		record := &ENSRecord{
			Name:        "test.eth",
			Address:     common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b"),
			ContentHash: "QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
			TextRecords: map[string]string{
				"email":   "test@example.com",
				"website": "https://example.com",
			},
			ResolvedAt: time.Now(),
			TTL:        time.Hour,
		}

		assert.Equal(t, "test.eth", record.Name)
		assert.NotEqual(t, common.Address{}, record.Address)
		assert.Equal(t, "test@example.com", record.TextRecords["email"])
	})
}

func TestHardwareWalletService(t *testing.T) {
	logger := observability.NewLogger(config.ObservabilityConfig{})

	_ = NewHardwareWalletService(logger)

	t.Run("HardwareWalletTypes", func(t *testing.T) {
		types := []HardwareWalletType{
			HardwareWalletTypeLedger,
			HardwareWalletTypeTrezor,
			HardwareWalletTypeGridPlus,
		}

		for _, hwType := range types {
			assert.NotEmpty(t, string(hwType))
		}
	})

	t.Run("HardwareWallet", func(t *testing.T) {
		wallet := &HardwareWallet{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			DeviceType: HardwareWalletTypeLedger,
			DeviceID:   "device-123",
			Name:       "Ledger Nano S",
			Addresses: []HardwareAddress{
				{
					Address:        common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b"),
					DerivationPath: "m/44'/60'/0'/0/0",
					ChainID:        1,
					IsActive:       true,
				},
			},
			IsConnected: true,
			LastSeen:    time.Now(),
			CreatedAt:   time.Now(),
		}

		assert.Equal(t, HardwareWalletTypeLedger, wallet.DeviceType)
		assert.Equal(t, "device-123", wallet.DeviceID)
		assert.True(t, wallet.IsConnected)
		assert.Len(t, wallet.Addresses, 1)
	})

	t.Run("HardwareSignRequest", func(t *testing.T) {
		req := HardwareSignRequest{
			DeviceID:       "device-123",
			DerivationPath: "m/44'/60'/0'/0/0",
			ChainID:        1,
			Message:        []byte("test message"),
			Metadata:       map[string]interface{}{"type": "message"},
		}

		assert.Equal(t, "device-123", req.DeviceID)
		assert.Equal(t, "m/44'/60'/0'/0/0", req.DerivationPath)
		assert.Equal(t, 1, req.ChainID)
		assert.Equal(t, []byte("test message"), req.Message)
	})
}

func TestEnhancedTransactionRequest(t *testing.T) {
	t.Run("EnhancedTransactionRequest", func(t *testing.T) {
		req := EnhancedTransactionRequest{
			WalletID:    uuid.New(),
			ToAddress:   "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			Value:       big.NewInt(1000000000000000000), // 1 ETH
			Data:        "0x",
			GasStrategy: GasStrategyStandard,
			Metadata:    map[string]interface{}{"type": "transfer"},
			SimulateTx:  true,
		}

		assert.NotEqual(t, uuid.Nil, req.WalletID)
		assert.Equal(t, "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b", req.ToAddress)
		assert.Equal(t, GasStrategyStandard, req.GasStrategy)
		assert.True(t, req.SimulateTx)
	})

	t.Run("TransactionSimulation", func(t *testing.T) {
		simulation := TransactionSimulation{
			Success:       true,
			GasUsed:       21000,
			GasPrice:      big.NewInt(20000000000),
			EstimatedCost: big.NewInt(420000000000000),
			StateChanges:  map[string]interface{}{"balance_change": "-1000000000000000000"},
		}

		assert.True(t, simulation.Success)
		assert.Equal(t, uint64(21000), simulation.GasUsed)
		assert.Equal(t, "-1000000000000000000", simulation.StateChanges["balance_change"])
	})
}

func TestWeb3Config(t *testing.T) {
	t.Run("Web3Config", func(t *testing.T) {
		cfg := config.Web3Config{
			EthereumRPC:        "https://mainnet.infura.io/v3/test",
			PolygonRPC:         "https://polygon-mainnet.infura.io/v3/test",
			ArbitrumRPC:        "https://arbitrum-mainnet.infura.io/v3/test",
			OptimismRPC:        "https://optimism-mainnet.infura.io/v3/test",
			IPFSNodeURL:        "http://localhost:5001",
			IPFSGateway:        "https://ipfs.io",
			IPFSMaxFileSize:    10 * 1024 * 1024,
			GasOptimization:    true,
			HardwareWallets:    true,
			ENSResolution:      true,
			TransactionTimeout: 5 * time.Minute,
			MaxRetries:         3,
			RetryDelay:         2 * time.Second,
		}

		assert.NotEmpty(t, cfg.EthereumRPC)
		assert.NotEmpty(t, cfg.IPFSNodeURL)
		assert.True(t, cfg.GasOptimization)
		assert.True(t, cfg.HardwareWallets)
		assert.True(t, cfg.ENSResolution)
		assert.Equal(t, int64(10*1024*1024), cfg.IPFSMaxFileSize)
	})
}
