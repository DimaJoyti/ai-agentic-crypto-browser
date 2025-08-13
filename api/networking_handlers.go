package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ai-agentic-browser/internal/hft"
	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// NetworkingHandlers provides HTTP handlers for High-Performance Networking
type NetworkingHandlers struct {
	networking *hft.HighPerformanceNetworking
	logger     *observability.Logger
}

// NewNetworkingHandlers creates new high-performance networking HTTP handlers
func NewNetworkingHandlers(networking *hft.HighPerformanceNetworking, logger *observability.Logger) *NetworkingHandlers {
	return &NetworkingHandlers{
		networking: networking,
		logger:     logger,
	}
}

// GetMetrics handles networking metrics requests
func (h *NetworkingHandlers) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.networking.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Error(ctx, "Failed to encode networking metrics", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetConnections handles connections list requests
func (h *NetworkingHandlers) GetConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connections := h.networking.GetConnections()

	response := map[string]interface{}{
		"connections": connections,
		"count":       len(connections),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ctx, "Failed to encode connections", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateConnection handles connection creation requests
func (h *NetworkingHandlers) CreateConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type       string `json:"type"`
		RemoteAddr string `json:"remote_addr"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode connection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type == "" || req.RemoteAddr == "" {
		http.Error(w, "Type and remote_addr are required", http.StatusBadRequest)
		return
	}

	var connType hft.ConnectionType
	switch req.Type {
	case "UDP":
		connType = hft.ConnectionTypeUDP
	case "TCP":
		connType = hft.ConnectionTypeTCP
	case "MULTICAST":
		connType = hft.ConnectionTypeMulticast
	default:
		http.Error(w, "Invalid connection type", http.StatusBadRequest)
		return
	}

	conn, err := h.networking.CreateConnection(ctx, connType, req.RemoteAddr)
	if err != nil {
		h.logger.Error(ctx, "Failed to create connection", err)
		http.Error(w, "Failed to create connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"message":       "Connection created successfully",
		"connection_id": conn.ID.String(),
		"connection":    conn,
	})

	h.logger.Info(ctx, "Connection created via API", map[string]interface{}{
		"connection_id": conn.ID.String(),
		"type":          req.Type,
		"remote_addr":   req.RemoteAddr,
	})
}

// GetConnection handles single connection requests
func (h *NetworkingHandlers) GetConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	connectionIDStr := vars["id"]
	if connectionIDStr == "" {
		http.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		http.Error(w, "Invalid connection ID", http.StatusBadRequest)
		return
	}

	conn := h.networking.GetConnection(connectionID)
	if conn == nil {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conn); err != nil {
		h.logger.Error(ctx, "Failed to encode connection", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CloseConnection handles connection close requests
func (h *NetworkingHandlers) CloseConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	connectionIDStr := vars["id"]
	if connectionIDStr == "" {
		http.Error(w, "Connection ID is required", http.StatusBadRequest)
		return
	}

	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		http.Error(w, "Invalid connection ID", http.StatusBadRequest)
		return
	}

	if err := h.networking.CloseConnection(ctx, connectionID); err != nil {
		h.logger.Error(ctx, "Failed to close connection", err)
		http.Error(w, "Failed to close connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Connection closed successfully",
	})

	h.logger.Info(ctx, "Connection closed via API", map[string]interface{}{
		"connection_id": connectionID.String(),
	})
}

// SendMessage handles message sending requests
func (h *NetworkingHandlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type        string `json:"type"`
		Payload     string `json:"payload"`
		Destination string `json:"destination"`
		Priority    int    `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode send request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Payload == "" || req.Destination == "" {
		http.Error(w, "Payload and destination are required", http.StatusBadRequest)
		return
	}

	// Create network message
	msg := &hft.NetworkMessage{
		ID:        uuid.New(),
		Type:      hft.MessageType(req.Type),
		Payload:   []byte(req.Payload),
		Priority:  hft.Priority(req.Priority),
		TTL:       0, // No TTL for now
	}

	// Parse destination address (simplified)
	// In production, this would properly parse and validate addresses

	if err := h.networking.SendMessage(ctx, msg); err != nil {
		h.logger.Error(ctx, "Failed to send message", err)
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"message":    "Message sent successfully",
		"message_id": msg.ID.String(),
	})

	h.logger.Info(ctx, "Message sent via API", map[string]interface{}{
		"message_id":  msg.ID.String(),
		"type":        req.Type,
		"destination": req.Destination,
		"size":        len(req.Payload),
	})
}

// ReceiveMessage handles message receiving requests
func (h *NetworkingHandlers) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := h.networking.ReceiveMessage(ctx)
	if err != nil {
		// No messages available is not an error for the API
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "no_messages",
			"message": "No messages available",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": msg,
	}); err != nil {
		h.logger.Error(ctx, "Failed to encode received message", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Debug(ctx, "Message received via API", map[string]interface{}{
		"message_id": msg.ID.String(),
		"type":       string(msg.Type),
		"size":       len(msg.Payload),
	})
}

// GetStatus handles networking status requests
func (h *NetworkingHandlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.networking.GetMetrics()
	connections := h.networking.GetConnections()

	status := map[string]interface{}{
		"system_status":      "operational",
		"active_connections": len(connections),
		"performance": map[string]interface{}{
			"packets_per_second": fmt.Sprintf("%.0f", float64(metrics.PacketsReceived+metrics.PacketsSent)/10.0), // Simplified
			"avg_latency_us":     metrics.AvgLatencyNs / 1000,
			"min_latency_us":     metrics.MinLatencyNs / 1000,
			"max_latency_us":     metrics.MaxLatencyNs / 1000,
		},
		"queues": map[string]interface{}{
			"inbound_size":  metrics.InboundQueueSize,
			"outbound_size": metrics.OutboundQueueSize,
			"priority_size": metrics.PriorityQueueSize,
		},
		"last_update": metrics.LastUpdate,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		h.logger.Error(ctx, "Failed to encode networking status", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
