package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

func HandleDeFiInteraction(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		var req web3.DeFiProtocolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resp, err := web3Service.InteractWithDeFiProtocol(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "DeFi interaction failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func HandleGetProtocols(defiManager *web3.DeFiProtocolManager, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		protocols := defiManager.GetProtocols()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"protocols": protocols})
	}
}

func HandleGetProtocol(defiManager *web3.DeFiProtocolManager, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		protocolID := strings.TrimPrefix(r.URL.Path, "/web3/defi/protocols/")
		protocol, err := defiManager.GetProtocol(protocolID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(protocol)
	}
}

func HandleGetYieldOpportunities(defiManager *web3.DeFiProtocolManager, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This project doesn't currently expose a method; return protocols as placeholder for opportunities
		protocols := defiManager.GetProtocols()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"protocols": protocols})
	}
}
