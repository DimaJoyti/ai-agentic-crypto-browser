package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ai-agentic-browser/internal/mcp"
	"github.com/gorilla/mux"
)

// Firebase API handlers for MCP integration

// handleFirebaseAuth handles Firebase authentication operations
func (s *APIServer) handleFirebaseAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleCreateFirebaseUser(w, r)
	case http.MethodGet:
		s.handleGetFirebaseUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCreateFirebaseUser creates a new Firebase user
func (s *APIServer) handleCreateFirebaseUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get Firebase client from MCP integration
	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	user, err := firebaseClient.CreateUser(r.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusCreated, map[string]interface{}{
		"user":    user,
		"message": "User created successfully",
	})
}

// handleGetFirebaseUser retrieves a Firebase user
func (s *APIServer) handleGetFirebaseUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		http.Error(w, "User UID is required", http.StatusBadRequest)
		return
	}

	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	user, err := firebaseClient.GetUser(r.Context(), uid)
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusNotFound)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// handleFirebaseVerifyToken verifies a Firebase ID token
func (s *APIServer) handleFirebaseVerifyToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"id_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	token, err := firebaseClient.VerifyIDToken(r.Context(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"valid":  true,
		"claims": token.Claims,
		"uid":    token.UID,
	})
}

// handleFirestoreDocument handles Firestore document operations
func (s *APIServer) handleFirestoreDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	documentID := vars["documentId"]

	if collection == "" {
		http.Error(w, "Collection is required", http.StatusBadRequest)
		return
	}

	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.handleCreateFirestoreDocument(w, r, firebaseClient, collection, documentID)
	case http.MethodGet:
		s.handleGetFirestoreDocument(w, r, firebaseClient, collection, documentID)
	case http.MethodPut:
		s.handleUpdateFirestoreDocument(w, r, firebaseClient, collection, documentID)
	case http.MethodDelete:
		s.handleDeleteFirestoreDocument(w, r, firebaseClient, collection, documentID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCreateFirestoreDocument creates a new Firestore document
func (s *APIServer) handleCreateFirestoreDocument(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, collection, documentID string) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	document, err := firebaseClient.CreateDocument(r.Context(), collection, documentID, data)
	if err != nil {
		http.Error(w, "Failed to create document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusCreated, map[string]interface{}{
		"document": document,
		"message":  "Document created successfully",
	})
}

// handleGetFirestoreDocument retrieves a Firestore document
func (s *APIServer) handleGetFirestoreDocument(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, collection, documentID string) {
	if documentID == "" {
		// Query multiple documents
		limitStr := r.URL.Query().Get("limit")
		limit := 0
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		documents, err := firebaseClient.QueryDocuments(r.Context(), collection, limit)
		if err != nil {
			http.Error(w, "Failed to query documents: "+err.Error(), http.StatusInternalServerError)
			return
		}

		s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
			"documents": documents,
			"count":     len(documents),
		})
		return
	}

	document, err := firebaseClient.GetDocument(r.Context(), collection, documentID)
	if err != nil {
		http.Error(w, "Failed to get document: "+err.Error(), http.StatusNotFound)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"document": document,
	})
}

// handleUpdateFirestoreDocument updates a Firestore document
func (s *APIServer) handleUpdateFirestoreDocument(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, collection, documentID string) {
	if documentID == "" {
		http.Error(w, "Document ID is required for updates", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := firebaseClient.UpdateDocument(r.Context(), collection, documentID, data)
	if err != nil {
		http.Error(w, "Failed to update document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"message": "Document updated successfully",
	})
}

// handleDeleteFirestoreDocument deletes a Firestore document
func (s *APIServer) handleDeleteFirestoreDocument(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, collection, documentID string) {
	if documentID == "" {
		http.Error(w, "Document ID is required for deletion", http.StatusBadRequest)
		return
	}

	err := firebaseClient.DeleteDocument(r.Context(), collection, documentID)
	if err != nil {
		http.Error(w, "Failed to delete document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"message": "Document deleted successfully",
	})
}

// handleFirebaseRealtimeDB handles Realtime Database operations
func (s *APIServer) handleFirebaseRealtimeDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodPost, http.MethodPut:
		s.handleSetRealtimeData(w, r, firebaseClient, path)
	case http.MethodGet:
		s.handleGetRealtimeData(w, r, firebaseClient, path)
	case http.MethodPatch:
		s.handleUpdateRealtimeData(w, r, firebaseClient, path)
	case http.MethodDelete:
		s.handleDeleteRealtimeData(w, r, firebaseClient, path)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSetRealtimeData sets data in Realtime Database
func (s *APIServer) handleSetRealtimeData(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, path string) {
	var data interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := firebaseClient.SetRealtimeData(r.Context(), path, data)
	if err != nil {
		http.Error(w, "Failed to set data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"message": "Data set successfully",
	})
}

// handleGetRealtimeData gets data from Realtime Database
func (s *APIServer) handleGetRealtimeData(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, path string) {
	data, err := firebaseClient.GetRealtimeData(r.Context(), path)
	if err != nil {
		http.Error(w, "Failed to get data: "+err.Error(), http.StatusNotFound)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

// handleUpdateRealtimeData updates data in Realtime Database
func (s *APIServer) handleUpdateRealtimeData(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, path string) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := firebaseClient.UpdateRealtimeData(r.Context(), path, updates)
	if err != nil {
		http.Error(w, "Failed to update data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"message": "Data updated successfully",
	})
}

// handleDeleteRealtimeData deletes data from Realtime Database
func (s *APIServer) handleDeleteRealtimeData(w http.ResponseWriter, r *http.Request, firebaseClient *mcp.FirebaseClient, path string) {
	err := firebaseClient.DeleteRealtimeData(r.Context(), path)
	if err != nil {
		http.Error(w, "Failed to delete data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"message": "Data deleted successfully",
	})
}

// getFirebaseClient retrieves the Firebase client from MCP integration
func (s *APIServer) getFirebaseClient() *mcp.FirebaseClient {
	if s.mcpService == nil {
		return nil
	}

	// Access the Firebase client through the MCP integration service
	// This assumes the MCP service exposes the Firebase client
	// You may need to add a method to the MCP service to expose the Firebase client
	return s.mcpService.GetFirebaseClient()
}

// handleFirebaseStatus returns Firebase service status
func (s *APIServer) handleFirebaseStatus(w http.ResponseWriter, r *http.Request) {
	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	status := firebaseClient.GetStatus(r.Context())
	s.sendJSON(w, r, http.StatusOK, status)
}

// handleFirebaseBatchWrite performs batch operations on Firestore
func (s *APIServer) handleFirebaseBatchWrite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Operations []mcp.BatchOperation `json:"operations"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	firebaseClient := s.getFirebaseClient()
	if firebaseClient == nil {
		http.Error(w, "Firebase client not available", http.StatusServiceUnavailable)
		return
	}

	err := firebaseClient.BatchWrite(r.Context(), req.Operations)
	if err != nil {
		http.Error(w, "Failed to perform batch write: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, r, http.StatusOK, map[string]interface{}{
		"message":          "Batch write completed successfully",
		"operations_count": len(req.Operations),
	})
}
