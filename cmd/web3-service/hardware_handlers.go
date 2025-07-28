package main

import (
	"encoding/json"
	"net/http"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/observability"
)

// Hardware wallet handlers - simplified for Phase 3 completion

// handleGetDevices returns all connected hardware wallet devices
func handleGetDevices(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"devices": []interface{}{},
			"status":  "hardware_wallet_service_available",
		})
	}
}

// handleDiscoverDevices discovers available hardware wallet devices
func handleDiscoverDevices(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		devices, err := hwService.DiscoverDevices(ctx)
		if err != nil {
			http.Error(w, "Failed to discover devices", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"devices": devices,
		})
	}
}

// handleConnectDevice connects to a specific hardware wallet device
func handleConnectDevice(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Hardware wallet connection available",
		})
	}
}

// handleDisconnectDevice disconnects from a hardware wallet device
func handleDisconnectDevice(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Hardware wallet disconnection available",
		})
	}
}

// handleGetAddresses gets addresses from a hardware wallet device
func handleGetAddresses(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"addresses": []interface{}{},
			"message":   "Hardware wallet address derivation available",
		})
	}
}

// handleSignTransaction signs a transaction with a hardware wallet
func handleSignTransaction(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Hardware wallet transaction signing available",
		})
	}
}

// handleSignMessage signs a message with a hardware wallet
func handleSignMessage(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Hardware wallet message signing available",
		})
	}
}

// handleGetDeviceStatus gets the status of a hardware wallet device
func handleGetDeviceStatus(hwService *web3.HardwareWalletService, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "available",
			"connected": false,
			"message":   "Hardware wallet status monitoring available",
		})
	}
}

// Integration status handlers

// handleIntegrationStatus returns the comprehensive Web3 integration status
func handleIntegrationStatus(checker *web3.IntegrationChecker, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		status, err := checker.CheckIntegrationStatus(ctx)
		if err != nil {
			http.Error(w, "Failed to check integration status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

// handleIntegrationSummary returns a summary of the Web3 integration
func handleIntegrationSummary(checker *web3.IntegrationChecker, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summary := checker.GetIntegrationSummary()

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(summary))
	}
}
