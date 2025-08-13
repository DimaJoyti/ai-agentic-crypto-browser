package mcp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"github.com/ai-agentic-browser/pkg/observability"
	"google.golang.org/api/option"
)

// FirebaseClient provides Firebase integration for MCP tools
type FirebaseClient struct {
	logger *observability.Logger
	config FirebaseConfig

	// Firebase services
	app       *firebase.App
	auth      *auth.Client
	firestore *firestore.Client
	database  *db.Client

	// State management
	isRunning bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// FirebaseConfig contains Firebase configuration
type FirebaseConfig struct {
	ProjectID          string         `json:"project_id"`
	CredentialsPath    string         `json:"credentials_path"`
	DatabaseURL        string         `json:"database_url"`
	StorageBucket      string         `json:"storage_bucket"`
	EnableAuth         bool           `json:"enable_auth"`
	EnableFirestore    bool           `json:"enable_firestore"`
	EnableRealtimeDB   bool           `json:"enable_realtime_db"`
	EnableStorage      bool           `json:"enable_storage"`
	EnableAnalytics    bool           `json:"enable_analytics"`
	EnableMessaging    bool           `json:"enable_messaging"`
	EnableRemoteConfig bool           `json:"enable_remote_config"`
	EnableDynamicLinks bool           `json:"enable_dynamic_links"`
	EnableMLKit        bool           `json:"enable_ml_kit"`
	EnablePerformance  bool           `json:"enable_performance"`
	EnableCrashlytics  bool           `json:"enable_crashlytics"`
	EnableAppCheck     bool           `json:"enable_app_check"`
	EnableExtensions   bool           `json:"enable_extensions"`
	EnableHosting      bool           `json:"enable_hosting"`
	EnableFunctions    bool           `json:"enable_functions"`
	EnableEmulators    bool           `json:"enable_emulators"`
	EmulatorConfig     EmulatorConfig `json:"emulator_config"`
}

// EmulatorConfig contains Firebase emulator configuration
type EmulatorConfig struct {
	AuthPort      int    `json:"auth_port"`
	FirestorePort int    `json:"firestore_port"`
	DatabasePort  int    `json:"database_port"`
	StoragePort   int    `json:"storage_port"`
	FunctionsPort int    `json:"functions_port"`
	HostingPort   int    `json:"hosting_port"`
	PubSubPort    int    `json:"pubsub_port"`
	UIPort        int    `json:"ui_port"`
	Host          string `json:"host"`
}

// FirebaseDocument represents a Firestore document
type FirebaseDocument struct {
	ID         string                 `json:"id"`
	Collection string                 `json:"collection"`
	Data       map[string]interface{} `json:"data"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// FirebaseUser represents a Firebase user
type FirebaseUser struct {
	UID           string                 `json:"uid"`
	Email         string                 `json:"email"`
	DisplayName   string                 `json:"display_name"`
	PhotoURL      string                 `json:"photo_url"`
	EmailVerified bool                   `json:"email_verified"`
	Disabled      bool                   `json:"disabled"`
	CustomClaims  map[string]interface{} `json:"custom_claims"`
	CreatedAt     time.Time              `json:"created_at"`
	LastSignIn    time.Time              `json:"last_sign_in"`
}

// NewFirebaseClient creates a new Firebase client
func NewFirebaseClient(logger *observability.Logger, config FirebaseConfig) *FirebaseClient {
	return &FirebaseClient{
		logger:   logger,
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// Start starts the Firebase client
func (fc *FirebaseClient) Start(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if fc.isRunning {
		return fmt.Errorf("firebase client is already running")
	}

	fc.logger.Info(ctx, "Starting Firebase client", map[string]interface{}{
		"project_id":       fc.config.ProjectID,
		"enable_auth":      fc.config.EnableAuth,
		"enable_firestore": fc.config.EnableFirestore,
		"enable_realtime":  fc.config.EnableRealtimeDB,
		"enable_emulators": fc.config.EnableEmulators,
	})

	// Initialize Firebase app
	if err := fc.initializeApp(ctx); err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Initialize services
	if err := fc.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize Firebase services: %w", err)
	}

	fc.isRunning = true

	fc.logger.Info(ctx, "Firebase client started successfully", nil)

	return nil
}

// Stop stops the Firebase client
func (fc *FirebaseClient) Stop(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if !fc.isRunning {
		return fmt.Errorf("firebase client is not running")
	}

	fc.logger.Info(ctx, "Stopping Firebase client", nil)

	close(fc.stopChan)
	fc.wg.Wait()

	// Close Firebase services
	if fc.firestore != nil {
		fc.firestore.Close()
	}

	fc.isRunning = false

	fc.logger.Info(ctx, "Firebase client stopped", nil)

	return nil
}

// IsHealthy checks if the Firebase client is healthy
func (fc *FirebaseClient) IsHealthy() bool {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	return fc.isRunning
}

// initializeApp initializes the Firebase app
func (fc *FirebaseClient) initializeApp(ctx context.Context) error {
	var opts []option.ClientOption

	// Add credentials if provided
	if fc.config.CredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(fc.config.CredentialsPath))
	}

	// Configure for emulators if enabled
	if fc.config.EnableEmulators {
		fc.logger.Info(ctx, "Configuring Firebase emulators", map[string]interface{}{
			"host":           fc.config.EmulatorConfig.Host,
			"auth_port":      fc.config.EmulatorConfig.AuthPort,
			"firestore_port": fc.config.EmulatorConfig.FirestorePort,
		})
	}

	// Initialize Firebase app
	config := &firebase.Config{
		ProjectID:     fc.config.ProjectID,
		DatabaseURL:   fc.config.DatabaseURL,
		StorageBucket: fc.config.StorageBucket,
	}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	fc.app = app

	return nil
}

// initializeServices initializes Firebase services
func (fc *FirebaseClient) initializeServices(ctx context.Context) error {
	// Initialize Auth service
	if fc.config.EnableAuth {
		authClient, err := fc.app.Auth(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Auth service: %w", err)
		}
		fc.auth = authClient
		fc.logger.Info(ctx, "Firebase Auth service initialized", nil)
	}

	// Initialize Firestore service
	if fc.config.EnableFirestore {
		firestoreClient, err := fc.app.Firestore(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Firestore service: %w", err)
		}
		fc.firestore = firestoreClient
		fc.logger.Info(ctx, "Firebase Firestore service initialized", nil)
	}

	// Initialize Realtime Database service
	if fc.config.EnableRealtimeDB {
		dbClient, err := fc.app.Database(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Realtime Database service: %w", err)
		}
		fc.database = dbClient
		fc.logger.Info(ctx, "Firebase Realtime Database service initialized", nil)
	}

	return nil
}

// Authentication methods

// CreateUser creates a new Firebase user
func (fc *FirebaseClient) CreateUser(ctx context.Context, email, password, displayName string) (*FirebaseUser, error) {
	if fc.auth == nil {
		return nil, fmt.Errorf("Firebase Auth is not enabled")
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(displayName).
		EmailVerified(false).
		Disabled(false)

	userRecord, err := fc.auth.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user := &FirebaseUser{
		UID:           userRecord.UID,
		Email:         userRecord.Email,
		DisplayName:   userRecord.DisplayName,
		PhotoURL:      userRecord.PhotoURL,
		EmailVerified: userRecord.EmailVerified,
		Disabled:      userRecord.Disabled,
		CustomClaims:  userRecord.CustomClaims,
		CreatedAt:     time.Unix(userRecord.UserMetadata.CreationTimestamp, 0),
		LastSignIn:    time.Unix(userRecord.UserMetadata.LastLogInTimestamp, 0),
	}

	fc.logger.Info(ctx, "Firebase user created", map[string]interface{}{
		"uid":   user.UID,
		"email": user.Email,
	})

	return user, nil
}

// GetUser retrieves a Firebase user by UID
func (fc *FirebaseClient) GetUser(ctx context.Context, uid string) (*FirebaseUser, error) {
	if fc.auth == nil {
		return nil, fmt.Errorf("Firebase Auth is not enabled")
	}

	userRecord, err := fc.auth.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := &FirebaseUser{
		UID:           userRecord.UID,
		Email:         userRecord.Email,
		DisplayName:   userRecord.DisplayName,
		PhotoURL:      userRecord.PhotoURL,
		EmailVerified: userRecord.EmailVerified,
		Disabled:      userRecord.Disabled,
		CustomClaims:  userRecord.CustomClaims,
		CreatedAt:     time.Unix(userRecord.UserMetadata.CreationTimestamp, 0),
		LastSignIn:    time.Unix(userRecord.UserMetadata.LastLogInTimestamp, 0),
	}

	return user, nil
}

// VerifyIDToken verifies a Firebase ID token
func (fc *FirebaseClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if fc.auth == nil {
		return nil, fmt.Errorf("Firebase Auth is not enabled")
	}

	token, err := fc.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	return token, nil
}

// SetCustomClaims sets custom claims for a user
func (fc *FirebaseClient) SetCustomClaims(ctx context.Context, uid string, claims map[string]interface{}) error {
	if fc.auth == nil {
		return fmt.Errorf("Firebase Auth is not enabled")
	}

	err := fc.auth.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		return fmt.Errorf("failed to set custom claims: %w", err)
	}

	fc.logger.Info(ctx, "Custom claims set for user", map[string]interface{}{
		"uid":    uid,
		"claims": claims,
	})

	return nil
}

// Firestore methods

// CreateDocument creates a new document in Firestore
func (fc *FirebaseClient) CreateDocument(ctx context.Context, collection, documentID string, data map[string]interface{}) (*FirebaseDocument, error) {
	if fc.firestore == nil {
		return nil, fmt.Errorf("Firestore is not enabled")
	}

	// Add timestamps
	now := time.Now()
	data["created_at"] = now
	data["updated_at"] = now

	var docRef *firestore.DocumentRef
	if documentID != "" {
		docRef = fc.firestore.Collection(collection).Doc(documentID)
	} else {
		docRef = fc.firestore.Collection(collection).NewDoc()
		documentID = docRef.ID
	}

	_, err := docRef.Set(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	document := &FirebaseDocument{
		ID:         documentID,
		Collection: collection,
		Data:       data,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	fc.logger.Info(ctx, "Firestore document created", map[string]interface{}{
		"collection":  collection,
		"document_id": documentID,
	})

	return document, nil
}

// GetDocument retrieves a document from Firestore
func (fc *FirebaseClient) GetDocument(ctx context.Context, collection, documentID string) (*FirebaseDocument, error) {
	if fc.firestore == nil {
		return nil, fmt.Errorf("Firestore is not enabled")
	}

	docRef := fc.firestore.Collection(collection).Doc(documentID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if !docSnap.Exists() {
		return nil, fmt.Errorf("document does not exist")
	}

	data := docSnap.Data()

	document := &FirebaseDocument{
		ID:         documentID,
		Collection: collection,
		Data:       data,
	}

	// Extract timestamps if they exist
	if createdAt, ok := data["created_at"].(time.Time); ok {
		document.CreatedAt = createdAt
	}
	if updatedAt, ok := data["updated_at"].(time.Time); ok {
		document.UpdatedAt = updatedAt
	}

	return document, nil
}

// UpdateDocument updates a document in Firestore
func (fc *FirebaseClient) UpdateDocument(ctx context.Context, collection, documentID string, data map[string]interface{}) error {
	if fc.firestore == nil {
		return fmt.Errorf("Firestore is not enabled")
	}

	// Add update timestamp
	data["updated_at"] = time.Now()

	docRef := fc.firestore.Collection(collection).Doc(documentID)
	_, err := docRef.Set(ctx, data, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	fc.logger.Info(ctx, "Firestore document updated", map[string]interface{}{
		"collection":  collection,
		"document_id": documentID,
	})

	return nil
}

// DeleteDocument deletes a document from Firestore
func (fc *FirebaseClient) DeleteDocument(ctx context.Context, collection, documentID string) error {
	if fc.firestore == nil {
		return fmt.Errorf("Firestore is not enabled")
	}

	docRef := fc.firestore.Collection(collection).Doc(documentID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	fc.logger.Info(ctx, "Firestore document deleted", map[string]interface{}{
		"collection":  collection,
		"document_id": documentID,
	})

	return nil
}

// QueryDocuments queries documents from Firestore
func (fc *FirebaseClient) QueryDocuments(ctx context.Context, collection string, limit int) ([]*FirebaseDocument, error) {
	if fc.firestore == nil {
		return nil, fmt.Errorf("Firestore is not enabled")
	}

	collectionRef := fc.firestore.Collection(collection)
	query := collectionRef.Query
	if limit > 0 {
		query = query.Limit(limit)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}

	var documents []*FirebaseDocument
	for _, doc := range docs {
		data := doc.Data()

		document := &FirebaseDocument{
			ID:         doc.Ref.ID,
			Collection: collection,
			Data:       data,
		}

		// Extract timestamps if they exist
		if createdAt, ok := data["created_at"].(time.Time); ok {
			document.CreatedAt = createdAt
		}
		if updatedAt, ok := data["updated_at"].(time.Time); ok {
			document.UpdatedAt = updatedAt
		}

		documents = append(documents, document)
	}

	return documents, nil
}

// Realtime Database methods

// SetRealtimeData sets data in Firebase Realtime Database
func (fc *FirebaseClient) SetRealtimeData(ctx context.Context, path string, data interface{}) error {
	if fc.database == nil {
		return fmt.Errorf("Realtime Database is not enabled")
	}

	ref := fc.database.NewRef(path)
	err := ref.Set(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to set realtime data: %w", err)
	}

	fc.logger.Info(ctx, "Realtime Database data set", map[string]interface{}{
		"path": path,
	})

	return nil
}

// GetRealtimeData gets data from Firebase Realtime Database
func (fc *FirebaseClient) GetRealtimeData(ctx context.Context, path string) (interface{}, error) {
	if fc.database == nil {
		return nil, fmt.Errorf("Realtime Database is not enabled")
	}

	ref := fc.database.NewRef(path)
	var data interface{}
	err := ref.Get(ctx, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to get realtime data: %w", err)
	}

	return data, nil
}

// UpdateRealtimeData updates data in Firebase Realtime Database
func (fc *FirebaseClient) UpdateRealtimeData(ctx context.Context, path string, updates map[string]interface{}) error {
	if fc.database == nil {
		return fmt.Errorf("Realtime Database is not enabled")
	}

	ref := fc.database.NewRef(path)
	err := ref.Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to update realtime data: %w", err)
	}

	fc.logger.Info(ctx, "Realtime Database data updated", map[string]interface{}{
		"path": path,
	})

	return nil
}

// DeleteRealtimeData deletes data from Firebase Realtime Database
func (fc *FirebaseClient) DeleteRealtimeData(ctx context.Context, path string) error {
	if fc.database == nil {
		return fmt.Errorf("Realtime Database is not enabled")
	}

	ref := fc.database.NewRef(path)
	err := ref.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete realtime data: %w", err)
	}

	fc.logger.Info(ctx, "Realtime Database data deleted", map[string]interface{}{
		"path": path,
	})

	return nil
}

// Utility methods

// GetConfig returns the Firebase configuration
func (fc *FirebaseClient) GetConfig() FirebaseConfig {
	return fc.config
}

// GetStatus returns the current status of Firebase services
func (fc *FirebaseClient) GetStatus(ctx context.Context) map[string]interface{} {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	status := map[string]interface{}{
		"running":    fc.isRunning,
		"project_id": fc.config.ProjectID,
		"services": map[string]bool{
			"auth":        fc.auth != nil,
			"firestore":   fc.firestore != nil,
			"realtime_db": fc.database != nil,
		},
		"emulators_enabled": fc.config.EnableEmulators,
	}

	return status
}

// BatchWrite performs a batch write operation in Firestore
func (fc *FirebaseClient) BatchWrite(ctx context.Context, operations []BatchOperation) error {
	if fc.firestore == nil {
		return fmt.Errorf("Firestore is not enabled")
	}

	batch := fc.firestore.Batch()

	for _, op := range operations {
		docRef := fc.firestore.Collection(op.Collection).Doc(op.DocumentID)

		switch op.Operation {
		case "create", "set":
			batch.Set(docRef, op.Data)
		case "update":
			batch.Update(docRef, op.Updates)
		case "delete":
			batch.Delete(docRef)
		default:
			return fmt.Errorf("unsupported batch operation: %s", op.Operation)
		}
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit batch write: %w", err)
	}

	fc.logger.Info(ctx, "Batch write completed", map[string]interface{}{
		"operations_count": len(operations),
	})

	return nil
}

// BatchOperation represents a single operation in a batch write
type BatchOperation struct {
	Operation  string                 `json:"operation"` // "create", "set", "update", "delete"
	Collection string                 `json:"collection"`
	DocumentID string                 `json:"document_id"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Updates    []firestore.Update     `json:"updates,omitempty"`
}

// ListenToDocument sets up a real-time listener for a Firestore document
func (fc *FirebaseClient) ListenToDocument(ctx context.Context, collection, documentID string, callback func(*FirebaseDocument, error)) error {
	if fc.firestore == nil {
		return fmt.Errorf("Firestore is not enabled")
	}

	docRef := fc.firestore.Collection(collection).Doc(documentID)

	fc.wg.Add(1)
	go func() {
		defer fc.wg.Done()

		snapshots := docRef.Snapshots(ctx)
		defer snapshots.Stop()

		for {
			select {
			case <-fc.stopChan:
				return
			default:
				snapshot, err := snapshots.Next()
				if err != nil {
					callback(nil, err)
					return
				}

				if snapshot.Exists() {
					data := snapshot.Data()
					document := &FirebaseDocument{
						ID:         documentID,
						Collection: collection,
						Data:       data,
					}

					// Extract timestamps if they exist
					if createdAt, ok := data["created_at"].(time.Time); ok {
						document.CreatedAt = createdAt
					}
					if updatedAt, ok := data["updated_at"].(time.Time); ok {
						document.UpdatedAt = updatedAt
					}

					callback(document, nil)
				} else {
					callback(nil, fmt.Errorf("document does not exist"))
				}
			}
		}
	}()

	fc.logger.Info(ctx, "Document listener started", map[string]interface{}{
		"collection":  collection,
		"document_id": documentID,
	})

	return nil
}
