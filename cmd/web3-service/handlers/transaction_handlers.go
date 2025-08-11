package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ai-agentic-browser/internal/web3"
	"github.com/ai-agentic-browser/pkg/middleware"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
)

func HandleCreateTransaction(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
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
		var req web3.TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resp, err := web3Service.CreateTransaction(r.Context(), userID, req)
		if err != nil {
			logger.Error(r.Context(), "Transaction creation failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func HandleListTransactions(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
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
		filter := web3.TransactionListFilter{}
		if chain := r.URL.Query().Get("chain_id"); chain != "" {
			if v, err := strconv.Atoi(chain); err == nil { filter.ChainID = v }
		}
		if status := r.URL.Query().Get("status"); status != "" {
			filter.Status = status
		}
		if page := r.URL.Query().Get("page"); page != "" {
			if v, err := strconv.Atoi(page); err == nil { filter.Page = v }
		}
		if ps := r.URL.Query().Get("page_size"); ps != "" {
			if v, err := strconv.Atoi(ps); err == nil { filter.PageSize = v }
		}
		transactions, pagination, err := web3Service.ListTransactions(r.Context(), userID, filter)
		if err != nil {
			logger.Error(r.Context(), "List transactions failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"transactions": transactions, "pagination": pagination})
	}
}

func HandleGetPrices(web3Service *web3.Service, logger *observability.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req web3.PriceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// If no body, use default
			req = web3.PriceRequest{Currency: "USD"}
		}
		// Support also query parameter ?token=ethereum&currency=EUR
		if token := r.URL.Query().Get("token"); strings.TrimSpace(token) != "" {
			req.Token = token
		}
		if cur := r.URL.Query().Get("currency"); strings.TrimSpace(cur) != "" {
			req.Currency = cur
		}
		resp, err := web3Service.GetPrices(r.Context(), req)
		if err != nil {
			logger.Error(r.Context(), "Price retrieval failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

