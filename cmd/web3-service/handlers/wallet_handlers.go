package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

func HandleConnectWallet(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
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
		var req web3.WalletConnectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resp, err := web3Service.ConnectWallet(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Wallet connect failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func HandleListWallets(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
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
		filter := web3.WalletListFilter{}
		if chain := r.URL.Query().Get("chain_id"); chain != "" {
			if v, err := strconv.Atoi(chain); err == nil {
				filter.ChainID = v
			}
		}
		if primary := r.URL.Query().Get("primary"); primary != "" {
			if primary == "true" {
				b := true
				filter.IsPrimary = &b
			} else if primary == "false" {
				b := false
				filter.IsPrimary = &b
			}
		}
		if page := r.URL.Query().Get("page"); page != "" {
			if v, err := strconv.Atoi(page); err == nil { filter.Page = v }
		}
		if ps := r.URL.Query().Get("page_size"); ps != "" {
			if v, err := strconv.Atoi(ps); err == nil { filter.PageSize = v }
		}
		wallets, pagination, err := web3Service.ListWallets(r.Context(), userID, filter)
		if err != nil {
			logger.Error(r.Context(), "List wallets failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"wallets": wallets, "pagination": pagination})
	}
}

func HandleGetBalance(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
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
		var req web3.BalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resp, err := web3Service.GetBalance(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Balance retrieval failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

