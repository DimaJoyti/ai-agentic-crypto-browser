package web3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// HardwareWalletType represents different hardware wallet types
type HardwareWalletType string

const (
	HardwareWalletTypeLedger   HardwareWalletType = "ledger"
	HardwareWalletTypeTrezor   HardwareWalletType = "trezor"
	HardwareWalletTypeGridPlus HardwareWalletType = "gridplus"
)

// HardwareWallet represents a hardware wallet device
type HardwareWallet struct {
	ID          uuid.UUID          `json:"id"`
	UserID      uuid.UUID          `json:"user_id"`
	DeviceType  HardwareWalletType `json:"device_type"`
	DeviceID    string             `json:"device_id"`
	Name        string             `json:"name"`
	Addresses   []HardwareAddress  `json:"addresses"`
	IsConnected bool               `json:"is_connected"`
	LastSeen    time.Time          `json:"last_seen"`
	CreatedAt   time.Time          `json:"created_at"`
}

// HardwareAddress represents an address derived from a hardware wallet
type HardwareAddress struct {
	Address        common.Address `json:"address"`
	DerivationPath string         `json:"derivation_path"`
	Index          int            `json:"index"`
	ChainID        int            `json:"chain_id"`
	IsActive       bool           `json:"is_active"`
}

// HardwareWalletService manages hardware wallet interactions
type HardwareWalletService struct {
	logger     *observability.Logger
	devices    map[string]*HardwareWallet
	connectors map[HardwareWalletType]HardwareWalletConnector
}

// HardwareWalletConnector interface for different hardware wallet implementations
type HardwareWalletConnector interface {
	Connect(ctx context.Context, deviceID string) error
	Disconnect(ctx context.Context, deviceID string) error
	GetAddresses(ctx context.Context, deviceID string, chainID int, count int) ([]HardwareAddress, error)
	SignTransaction(ctx context.Context, deviceID string, tx *types.Transaction, derivationPath string) (*types.Transaction, error)
	SignMessage(ctx context.Context, deviceID string, message []byte, derivationPath string) ([]byte, error)
	GetDeviceInfo(ctx context.Context, deviceID string) (*HardwareDeviceInfo, error)
}

// HardwareDeviceInfo represents information about a hardware device
type HardwareDeviceInfo struct {
	DeviceID   string `json:"device_id"`
	Model      string `json:"model"`
	Version    string `json:"version"`
	IsLocked   bool   `json:"is_locked"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
}

// HardwareSignRequest represents a request to sign with hardware wallet
type HardwareSignRequest struct {
	DeviceID       string                 `json:"device_id"`
	Transaction    *types.Transaction     `json:"transaction,omitempty"`
	Message        []byte                 `json:"message,omitempty"`
	DerivationPath string                 `json:"derivation_path"`
	ChainID        int                    `json:"chain_id"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// HardwareSignResponse represents the response from hardware wallet signing
type HardwareSignResponse struct {
	SignedTransaction *types.Transaction `json:"signed_transaction,omitempty"`
	Signature         []byte             `json:"signature,omitempty"`
	DeviceConfirmed   bool               `json:"device_confirmed"`
	Error             string             `json:"error,omitempty"`
}

// NewHardwareWalletService creates a new hardware wallet service
func NewHardwareWalletService(logger *observability.Logger) *HardwareWalletService {
	service := &HardwareWalletService{
		logger:     logger,
		devices:    make(map[string]*HardwareWallet),
		connectors: make(map[HardwareWalletType]HardwareWalletConnector),
	}

	// Register hardware wallet connectors
	service.registerConnectors()

	return service
}

// registerConnectors registers different hardware wallet connectors
func (s *HardwareWalletService) registerConnectors() {
	// Register Ledger connector
	s.connectors[HardwareWalletTypeLedger] = NewLedgerConnector(s.logger)

	// Register Trezor connector
	s.connectors[HardwareWalletTypeTrezor] = NewTrezorConnector(s.logger)

	// Register GridPlus connector
	s.connectors[HardwareWalletTypeGridPlus] = NewGridPlusConnector(s.logger)
}

// DiscoverDevices discovers connected hardware wallets
func (s *HardwareWalletService) DiscoverDevices(ctx context.Context) ([]*HardwareWallet, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet-service").Start(ctx, "hardware.DiscoverDevices")
	defer span.End()

	var discoveredDevices []*HardwareWallet

	// Discover devices for each connector type
	for deviceType, connector := range s.connectors {
		devices, err := s.discoverDevicesForType(ctx, deviceType, connector)
		if err != nil {
			s.logger.Warn(ctx, "Failed to discover devices", map[string]interface{}{
				"device_type": string(deviceType),
				"error":       err.Error(),
			})
			continue
		}
		discoveredDevices = append(discoveredDevices, devices...)
	}

	s.logger.Info(ctx, "Hardware devices discovered", map[string]interface{}{
		"count": len(discoveredDevices),
	})

	return discoveredDevices, nil
}

// ConnectDevice connects to a specific hardware wallet device
func (s *HardwareWalletService) ConnectDevice(ctx context.Context, userID uuid.UUID, deviceType HardwareWalletType, deviceID string) (*HardwareWallet, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet-service").Start(ctx, "hardware.ConnectDevice")
	defer span.End()

	connector, exists := s.connectors[deviceType]
	if !exists {
		return nil, fmt.Errorf("unsupported hardware wallet type: %s", deviceType)
	}

	// Connect to the device
	err := connector.Connect(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to device: %w", err)
	}

	// Get device info
	deviceInfo, err := connector.GetDeviceInfo(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	// Create hardware wallet record
	wallet := &HardwareWallet{
		ID:          uuid.New(),
		UserID:      userID,
		DeviceType:  deviceType,
		DeviceID:    deviceID,
		Name:        fmt.Sprintf("%s %s", deviceInfo.Model, deviceInfo.DeviceID[:8]),
		Addresses:   []HardwareAddress{},
		IsConnected: true,
		LastSeen:    time.Now(),
		CreatedAt:   time.Now(),
	}

	// Store the device
	s.devices[deviceID] = wallet

	s.logger.Info(ctx, "Hardware wallet connected", map[string]interface{}{
		"device_id":   deviceID,
		"device_type": string(deviceType),
		"user_id":     userID.String(),
	})

	return wallet, nil
}

// GetAddresses retrieves addresses from a hardware wallet
func (s *HardwareWalletService) GetAddresses(ctx context.Context, deviceID string, chainID int, count int) ([]HardwareAddress, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet-service").Start(ctx, "hardware.GetAddresses")
	defer span.End()

	wallet, exists := s.devices[deviceID]
	if !exists {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	connector, exists := s.connectors[wallet.DeviceType]
	if !exists {
		return nil, fmt.Errorf("connector not found for device type: %s", wallet.DeviceType)
	}

	addresses, err := connector.GetAddresses(ctx, deviceID, chainID, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	// Update wallet addresses
	wallet.Addresses = append(wallet.Addresses, addresses...)
	wallet.LastSeen = time.Now()

	s.logger.Info(ctx, "Addresses retrieved from hardware wallet", map[string]interface{}{
		"device_id": deviceID,
		"chain_id":  chainID,
		"count":     len(addresses),
	})

	return addresses, nil
}

// SignTransaction signs a transaction using a hardware wallet
func (s *HardwareWalletService) SignTransaction(ctx context.Context, req HardwareSignRequest) (*HardwareSignResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet-service").Start(ctx, "hardware.SignTransaction")
	defer span.End()

	wallet, exists := s.devices[req.DeviceID]
	if !exists {
		return &HardwareSignResponse{
			Error: fmt.Sprintf("device not found: %s", req.DeviceID),
		}, nil
	}

	connector, exists := s.connectors[wallet.DeviceType]
	if !exists {
		return &HardwareSignResponse{
			Error: fmt.Sprintf("connector not found for device type: %s", wallet.DeviceType),
		}, nil
	}

	// Sign the transaction
	signedTx, err := connector.SignTransaction(ctx, req.DeviceID, req.Transaction, req.DerivationPath)
	if err != nil {
		return &HardwareSignResponse{
			Error: err.Error(),
		}, nil
	}

	// Update last seen
	wallet.LastSeen = time.Now()

	s.logger.Info(ctx, "Transaction signed with hardware wallet", map[string]interface{}{
		"device_id":       req.DeviceID,
		"derivation_path": req.DerivationPath,
		"chain_id":        req.ChainID,
	})

	return &HardwareSignResponse{
		SignedTransaction: signedTx,
		DeviceConfirmed:   true,
	}, nil
}

// SignMessage signs a message using a hardware wallet
func (s *HardwareWalletService) SignMessage(ctx context.Context, req HardwareSignRequest) (*HardwareSignResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet-service").Start(ctx, "hardware.SignMessage")
	defer span.End()

	wallet, exists := s.devices[req.DeviceID]
	if !exists {
		return &HardwareSignResponse{
			Error: fmt.Sprintf("device not found: %s", req.DeviceID),
		}, nil
	}

	connector, exists := s.connectors[wallet.DeviceType]
	if !exists {
		return &HardwareSignResponse{
			Error: fmt.Sprintf("connector not found for device type: %s", wallet.DeviceType),
		}, nil
	}

	// Sign the message
	signature, err := connector.SignMessage(ctx, req.DeviceID, req.Message, req.DerivationPath)
	if err != nil {
		return &HardwareSignResponse{
			Error: err.Error(),
		}, nil
	}

	// Update last seen
	wallet.LastSeen = time.Now()

	s.logger.Info(ctx, "Message signed with hardware wallet", map[string]interface{}{
		"device_id":       req.DeviceID,
		"derivation_path": req.DerivationPath,
		"message_length":  len(req.Message),
	})

	return &HardwareSignResponse{
		Signature:       signature,
		DeviceConfirmed: true,
	}, nil
}

// discoverDevicesForType discovers devices for a specific hardware wallet type
func (s *HardwareWalletService) discoverDevicesForType(ctx context.Context, deviceType HardwareWalletType, connector HardwareWalletConnector) ([]*HardwareWallet, error) {
	// This would implement device discovery logic specific to each hardware wallet type
	// For now, return empty list as this requires platform-specific USB/HID implementations
	return []*HardwareWallet{}, nil
}

// Helper functions for hardware wallet operations

// getChainCoinType returns the BIP44 coin type for a given chain ID
func getChainCoinType(chainID int) int {
	switch chainID {
	case 1, 11155111: // Ethereum mainnet and Sepolia
		return 60
	case 137: // Polygon
		return 966
	case 42161: // Arbitrum
		return 60 // Uses Ethereum coin type
	case 10: // Optimism
		return 60 // Uses Ethereum coin type
	default:
		return 60 // Default to Ethereum
	}
}

// generateMockAddress generates a mock address for testing purposes
func generateMockAddress(chainID int, index int) common.Address {
	// Generate a deterministic mock address based on chainID and index
	seed := fmt.Sprintf("hw_%d_%d", chainID, index)
	hash := crypto.Keccak256Hash([]byte(seed))
	return common.BytesToAddress(hash[:20])
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Hardware wallet connector implementations

// Ledger connector implementation
type LedgerConnector struct {
	logger  *observability.Logger
	devices map[string]*LedgerDevice
	mu      sync.RWMutex
}

type LedgerDevice struct {
	ID           string
	ProductName  string
	Manufacturer string
	SerialNumber string
	IsConnected  bool
	LastSeen     time.Time
	Addresses    map[string][]HardwareAddress // chainID -> addresses
}

func NewLedgerConnector(logger *observability.Logger) *LedgerConnector {
	return &LedgerConnector{
		logger:  logger,
		devices: make(map[string]*LedgerDevice),
	}
}

func (c *LedgerConnector) Connect(ctx context.Context, deviceID string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "ledger.Connect")
	defer span.End()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Simulate Ledger device detection and connection
	// In a real implementation, this would use the Ledger USB/HID protocol
	device := &LedgerDevice{
		ID:           deviceID,
		ProductName:  "Ledger Nano S Plus",
		Manufacturer: "Ledger",
		SerialNumber: fmt.Sprintf("LDG-%s", deviceID[:8]),
		IsConnected:  true,
		LastSeen:     time.Now(),
		Addresses:    make(map[string][]HardwareAddress),
	}

	c.devices[deviceID] = device

	c.logger.Info(ctx, "Ledger device connected", map[string]interface{}{
		"device_id":     deviceID,
		"product_name":  device.ProductName,
		"serial_number": device.SerialNumber,
	})

	return nil
}

func (c *LedgerConnector) Disconnect(ctx context.Context, deviceID string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "ledger.Disconnect")
	defer span.End()

	c.mu.Lock()
	defer c.mu.Unlock()

	device, exists := c.devices[deviceID]
	if !exists {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	device.IsConnected = false
	delete(c.devices, deviceID)

	c.logger.Info(ctx, "Ledger device disconnected", map[string]interface{}{
		"device_id": deviceID,
	})

	return nil
}

func (c *LedgerConnector) GetAddresses(ctx context.Context, deviceID string, chainID int, count int) ([]HardwareAddress, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "ledger.GetAddresses")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	// Generate mock addresses for demonstration
	// In a real implementation, this would derive addresses from the Ledger device
	addresses := make([]HardwareAddress, count)
	for i := 0; i < count; i++ {
		derivationPath := fmt.Sprintf("m/44'/%d'/0'/0/%d", getChainCoinType(chainID), i)
		address := generateMockAddress(chainID, i)

		addresses[i] = HardwareAddress{
			Address:        address,
			DerivationPath: derivationPath,
			Index:          i,
			ChainID:        chainID,
		}
	}

	// Cache addresses
	c.mu.Lock()
	chainKey := fmt.Sprintf("%d", chainID)
	device.Addresses[chainKey] = addresses
	c.mu.Unlock()

	c.logger.Info(ctx, "Generated Ledger addresses", map[string]interface{}{
		"device_id": deviceID,
		"chain_id":  chainID,
		"count":     count,
	})

	return addresses, nil
}

func (c *LedgerConnector) SignTransaction(ctx context.Context, deviceID string, tx *types.Transaction, derivationPath string) (*types.Transaction, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "ledger.SignTransaction")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	// Simulate transaction signing
	// In a real implementation, this would send the transaction to the Ledger device for signing
	c.logger.Info(ctx, "Signing transaction with Ledger", map[string]interface{}{
		"device_id":       deviceID,
		"derivation_path": derivationPath,
		"tx_hash":         tx.Hash().Hex(),
	})

	// Return the transaction as-is for simulation
	// In reality, this would return the signed transaction
	return tx, nil
}

func (c *LedgerConnector) SignMessage(ctx context.Context, deviceID string, message []byte, derivationPath string) ([]byte, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "ledger.SignMessage")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	// Simulate message signing
	c.logger.Info(ctx, "Signing message with Ledger", map[string]interface{}{
		"device_id":       deviceID,
		"derivation_path": derivationPath,
		"message_length":  len(message),
	})

	// Return mock signature
	signature := make([]byte, 65)
	copy(signature, message[:min(len(message), 32)])
	return signature, nil
}

func (c *LedgerConnector) GetDeviceInfo(ctx context.Context, deviceID string) (*HardwareDeviceInfo, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "ledger.GetDeviceInfo")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	return &HardwareDeviceInfo{
		DeviceID:   device.ID,
		Model:      device.ProductName,
		Version:    "2.1.0",
		IsLocked:   false,
		AppName:    "Ethereum",
		AppVersion: "1.9.0",
	}, nil
}

// Trezor connector implementation
type TrezorConnector struct {
	logger  *observability.Logger
	devices map[string]*TrezorDevice
	mu      sync.RWMutex
}

type TrezorDevice struct {
	ID           string
	ProductName  string
	Manufacturer string
	SerialNumber string
	IsConnected  bool
	LastSeen     time.Time
	Addresses    map[string][]HardwareAddress
}

func NewTrezorConnector(logger *observability.Logger) *TrezorConnector {
	return &TrezorConnector{
		logger:  logger,
		devices: make(map[string]*TrezorDevice),
	}
}

func (c *TrezorConnector) Connect(ctx context.Context, deviceID string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "trezor.Connect")
	defer span.End()

	c.mu.Lock()
	defer c.mu.Unlock()

	device := &TrezorDevice{
		ID:           deviceID,
		ProductName:  "Trezor Model T",
		Manufacturer: "SatoshiLabs",
		SerialNumber: fmt.Sprintf("TRZ-%s", deviceID[:8]),
		IsConnected:  true,
		LastSeen:     time.Now(),
		Addresses:    make(map[string][]HardwareAddress),
	}

	c.devices[deviceID] = device

	c.logger.Info(ctx, "Trezor device connected", map[string]interface{}{
		"device_id":     deviceID,
		"product_name":  device.ProductName,
		"serial_number": device.SerialNumber,
	})

	return nil
}

func (c *TrezorConnector) Disconnect(ctx context.Context, deviceID string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "trezor.Disconnect")
	defer span.End()

	c.mu.Lock()
	defer c.mu.Unlock()

	device, exists := c.devices[deviceID]
	if !exists {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	device.IsConnected = false
	delete(c.devices, deviceID)

	c.logger.Info(ctx, "Trezor device disconnected", map[string]interface{}{
		"device_id": deviceID,
	})

	return nil
}

func (c *TrezorConnector) GetAddresses(ctx context.Context, deviceID string, chainID int, count int) ([]HardwareAddress, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "trezor.GetAddresses")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	addresses := make([]HardwareAddress, count)
	for i := 0; i < count; i++ {
		derivationPath := fmt.Sprintf("m/44'/%d'/0'/0/%d", getChainCoinType(chainID), i)
		address := generateMockAddress(chainID, i+1000) // Offset for Trezor

		addresses[i] = HardwareAddress{
			Address:        address,
			DerivationPath: derivationPath,
			Index:          i,
			ChainID:        chainID,
			IsActive:       true,
		}
	}

	c.mu.Lock()
	chainKey := fmt.Sprintf("%d", chainID)
	device.Addresses[chainKey] = addresses
	c.mu.Unlock()

	c.logger.Info(ctx, "Generated Trezor addresses", map[string]interface{}{
		"device_id": deviceID,
		"chain_id":  chainID,
		"count":     count,
	})

	return addresses, nil
}

func (c *TrezorConnector) SignTransaction(ctx context.Context, deviceID string, tx *types.Transaction, derivationPath string) (*types.Transaction, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "trezor.SignTransaction")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	c.logger.Info(ctx, "Signing transaction with Trezor", map[string]interface{}{
		"device_id":       deviceID,
		"derivation_path": derivationPath,
		"tx_hash":         tx.Hash().Hex(),
	})

	return tx, nil
}

func (c *TrezorConnector) SignMessage(ctx context.Context, deviceID string, message []byte, derivationPath string) ([]byte, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "trezor.SignMessage")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	c.logger.Info(ctx, "Signing message with Trezor", map[string]interface{}{
		"device_id":       deviceID,
		"derivation_path": derivationPath,
		"message_length":  len(message),
	})

	signature := make([]byte, 65)
	copy(signature, message[:min(len(message), 32)])
	signature[64] = 1 // Different recovery ID for Trezor
	return signature, nil
}

func (c *TrezorConnector) GetDeviceInfo(ctx context.Context, deviceID string) (*HardwareDeviceInfo, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "trezor.GetDeviceInfo")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	return &HardwareDeviceInfo{
		DeviceID:   device.ID,
		Model:      device.ProductName,
		Version:    "2.5.3",
		IsLocked:   false,
		AppName:    "Ethereum",
		AppVersion: "1.11.1",
	}, nil
}

// GridPlus connector implementation
type GridPlusConnector struct {
	logger  *observability.Logger
	devices map[string]*GridPlusDevice
	mu      sync.RWMutex
}

type GridPlusDevice struct {
	ID           string
	ProductName  string
	Manufacturer string
	SerialNumber string
	IsConnected  bool
	LastSeen     time.Time
	Addresses    map[string][]HardwareAddress
}

func NewGridPlusConnector(logger *observability.Logger) *GridPlusConnector {
	return &GridPlusConnector{
		logger:  logger,
		devices: make(map[string]*GridPlusDevice),
	}
}

func (c *GridPlusConnector) Connect(ctx context.Context, deviceID string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "gridplus.Connect")
	defer span.End()

	c.mu.Lock()
	defer c.mu.Unlock()

	device := &GridPlusDevice{
		ID:           deviceID,
		ProductName:  "GridPlus Lattice1",
		Manufacturer: "GridPlus",
		SerialNumber: fmt.Sprintf("GP-%s", deviceID[:8]),
		IsConnected:  true,
		LastSeen:     time.Now(),
		Addresses:    make(map[string][]HardwareAddress),
	}

	c.devices[deviceID] = device

	c.logger.Info(ctx, "GridPlus device connected", map[string]interface{}{
		"device_id":     deviceID,
		"product_name":  device.ProductName,
		"serial_number": device.SerialNumber,
	})

	return nil
}

func (c *GridPlusConnector) Disconnect(ctx context.Context, deviceID string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "gridplus.Disconnect")
	defer span.End()

	c.mu.Lock()
	defer c.mu.Unlock()

	device, exists := c.devices[deviceID]
	if !exists {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	device.IsConnected = false
	delete(c.devices, deviceID)

	c.logger.Info(ctx, "GridPlus device disconnected", map[string]interface{}{
		"device_id": deviceID,
	})

	return nil
}

func (c *GridPlusConnector) GetAddresses(ctx context.Context, deviceID string, chainID int, count int) ([]HardwareAddress, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "gridplus.GetAddresses")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	addresses := make([]HardwareAddress, count)
	for i := 0; i < count; i++ {
		derivationPath := fmt.Sprintf("m/44'/%d'/0'/0/%d", getChainCoinType(chainID), i)
		address := generateMockAddress(chainID, i+2000) // Offset for GridPlus

		addresses[i] = HardwareAddress{
			Address:        address,
			DerivationPath: derivationPath,
			Index:          i,
			ChainID:        chainID,
			IsActive:       true,
		}
	}

	c.mu.Lock()
	chainKey := fmt.Sprintf("%d", chainID)
	device.Addresses[chainKey] = addresses
	c.mu.Unlock()

	c.logger.Info(ctx, "Generated GridPlus addresses", map[string]interface{}{
		"device_id": deviceID,
		"chain_id":  chainID,
		"count":     count,
	})

	return addresses, nil
}

func (c *GridPlusConnector) SignTransaction(ctx context.Context, deviceID string, tx *types.Transaction, derivationPath string) (*types.Transaction, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "gridplus.SignTransaction")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	c.logger.Info(ctx, "Signing transaction with GridPlus", map[string]interface{}{
		"device_id":       deviceID,
		"derivation_path": derivationPath,
		"tx_hash":         tx.Hash().Hex(),
	})

	return tx, nil
}

func (c *GridPlusConnector) SignMessage(ctx context.Context, deviceID string, message []byte, derivationPath string) ([]byte, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "gridplus.SignMessage")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists || !device.IsConnected {
		return nil, fmt.Errorf("device not connected: %s", deviceID)
	}

	c.logger.Info(ctx, "Signing message with GridPlus", map[string]interface{}{
		"device_id":       deviceID,
		"derivation_path": derivationPath,
		"message_length":  len(message),
	})

	signature := make([]byte, 65)
	copy(signature, message[:min(len(message), 32)])
	signature[64] = 2 // Different recovery ID for GridPlus
	return signature, nil
}

func (c *GridPlusConnector) GetDeviceInfo(ctx context.Context, deviceID string) (*HardwareDeviceInfo, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("hardware-wallet").Start(ctx, "gridplus.GetDeviceInfo")
	defer span.End()

	c.mu.RLock()
	device, exists := c.devices[deviceID]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	return &HardwareDeviceInfo{
		DeviceID:   device.ID,
		Model:      device.ProductName,
		Version:    "1.0.0",
		IsLocked:   false,
		AppName:    "Ethereum",
		AppVersion: "1.0.0",
	}, nil
}
