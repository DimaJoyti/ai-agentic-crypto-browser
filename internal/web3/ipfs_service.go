package web3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/ai-agentic-browser/pkg/observability"
)

// IPFSService provides decentralized storage functionality
type IPFSService struct {
	client *shell.Shell
	logger *observability.Logger
	config IPFSConfig
}

// IPFSConfig holds IPFS configuration
type IPFSConfig struct {
	NodeURL     string        `json:"node_url"`
	Timeout     time.Duration `json:"timeout"`
	PinContent  bool          `json:"pin_content"`
	Gateway     string        `json:"gateway"`
	MaxFileSize int64         `json:"max_file_size"`
}

// IPFSObject represents an object stored in IPFS
type IPFSObject struct {
	Hash        string            `json:"hash"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	Pinned      bool              `json:"pinned"`
}

// IPFSUploadRequest represents a request to upload content to IPFS
type IPFSUploadRequest struct {
	Content     []byte            `json:"content"`
	ContentType string            `json:"content_type"`
	Filename    string            `json:"filename"`
	Metadata    map[string]string `json:"metadata"`
	Pin         bool              `json:"pin"`
}

// IPFSUploadResponse represents the response from uploading to IPFS
type IPFSUploadResponse struct {
	Object *IPFSObject `json:"object"`
	URL    string      `json:"url"`
}

// NewIPFSService creates a new IPFS service
func NewIPFSService(config IPFSConfig, logger *observability.Logger) *IPFSService {
	client := shell.NewShell(config.NodeURL)
	client.SetTimeout(config.Timeout)

	return &IPFSService{
		client: client,
		logger: logger,
		config: config,
	}
}

// Upload uploads content to IPFS
func (s *IPFSService) Upload(ctx context.Context, req IPFSUploadRequest) (*IPFSUploadResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.Upload")
	defer span.End()

	// Validate file size
	if int64(len(req.Content)) > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d bytes", s.config.MaxFileSize)
	}

	// Create reader from content
	reader := bytes.NewReader(req.Content)

	// Upload to IPFS
	hash, err := s.client.Add(reader, shell.Pin(req.Pin || s.config.PinContent))
	if err != nil {
		s.logger.Error(ctx, "Failed to upload to IPFS", err)
		return nil, fmt.Errorf("failed to upload to IPFS: %w", err)
	}

	// Create IPFS object
	object := &IPFSObject{
		Hash:        hash,
		Size:        int64(len(req.Content)),
		ContentType: req.ContentType,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		Pinned:      req.Pin || s.config.PinContent,
	}

	// Generate gateway URL
	url := s.generateGatewayURL(hash)

	response := &IPFSUploadResponse{
		Object: object,
		URL:    url,
	}

	s.logger.Info(ctx, "Content uploaded to IPFS", map[string]interface{}{
		"hash":         hash,
		"size":         object.Size,
		"content_type": req.ContentType,
		"pinned":       object.Pinned,
	})

	return response, nil
}

// Download downloads content from IPFS
func (s *IPFSService) Download(ctx context.Context, hash string) ([]byte, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.Download")
	defer span.End()

	// Validate hash format
	if !s.isValidIPFSHash(hash) {
		return nil, fmt.Errorf("invalid IPFS hash format: %s", hash)
	}

	// Download from IPFS
	reader, err := s.client.Cat(hash)
	if err != nil {
		s.logger.Error(ctx, "Failed to download from IPFS", err)
		return nil, fmt.Errorf("failed to download from IPFS: %w", err)
	}
	defer reader.Close()

	// Read content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read IPFS content: %w", err)
	}

	s.logger.Info(ctx, "Content downloaded from IPFS", map[string]interface{}{
		"hash": hash,
		"size": len(content),
	})

	return content, nil
}

// Pin pins content to ensure it stays available
func (s *IPFSService) Pin(ctx context.Context, hash string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.Pin")
	defer span.End()

	if !s.isValidIPFSHash(hash) {
		return fmt.Errorf("invalid IPFS hash format: %s", hash)
	}

	err := s.client.Pin(hash)
	if err != nil {
		s.logger.Error(ctx, "Failed to pin IPFS content", err)
		return fmt.Errorf("failed to pin IPFS content: %w", err)
	}

	s.logger.Info(ctx, "Content pinned to IPFS", map[string]interface{}{
		"hash": hash,
	})

	return nil
}

// Unpin unpins content from IPFS
func (s *IPFSService) Unpin(ctx context.Context, hash string) error {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.Unpin")
	defer span.End()

	if !s.isValidIPFSHash(hash) {
		return fmt.Errorf("invalid IPFS hash format: %s", hash)
	}

	err := s.client.Unpin(hash)
	if err != nil {
		s.logger.Error(ctx, "Failed to unpin IPFS content", err)
		return fmt.Errorf("failed to unpin IPFS content: %w", err)
	}

	s.logger.Info(ctx, "Content unpinned from IPFS", map[string]interface{}{
		"hash": hash,
	})

	return nil
}

// GetObjectInfo retrieves information about an IPFS object
func (s *IPFSService) GetObjectInfo(ctx context.Context, hash string) (*IPFSObject, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.GetObjectInfo")
	defer span.End()

	if !s.isValidIPFSHash(hash) {
		return nil, fmt.Errorf("invalid IPFS hash format: %s", hash)
	}

	// Get object stats
	stat, err := s.client.ObjectStat(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get IPFS object stats: %w", err)
	}

	// Check if pinned
	pins, err := s.client.Pins()
	if err != nil {
		s.logger.Warn(ctx, "Failed to check pin status", map[string]interface{}{
			"hash":  hash,
			"error": err.Error(),
		})
	}

	pinned := false
	if pins != nil {
		for pin := range pins {
			if pin == hash {
				pinned = true
				break
			}
		}
	}

	object := &IPFSObject{
		Hash:   hash,
		Size:   int64(stat.CumulativeSize),
		Pinned: pinned,
	}

	return object, nil
}

// UploadJSON uploads JSON data to IPFS
func (s *IPFSService) UploadJSON(ctx context.Context, data interface{}, metadata map[string]string) (*IPFSUploadResponse, error) {
	// Marshal JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req := IPFSUploadRequest{
		Content:     jsonData,
		ContentType: "application/json",
		Metadata:    metadata,
		Pin:         s.config.PinContent,
	}

	return s.Upload(ctx, req)
}

// DownloadJSON downloads and unmarshals JSON data from IPFS
func (s *IPFSService) DownloadJSON(ctx context.Context, hash string, target interface{}) error {
	content, err := s.Download(ctx, hash)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// ListPinnedObjects lists all pinned objects
func (s *IPFSService) ListPinnedObjects(ctx context.Context) ([]string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.ListPinnedObjects")
	defer span.End()

	pins, err := s.client.Pins()
	if err != nil {
		return nil, fmt.Errorf("failed to list pinned objects: %w", err)
	}

	var hashes []string
	for pin := range pins {
		hashes = append(hashes, pin)
	}

	return hashes, nil
}

// generateGatewayURL generates a gateway URL for accessing IPFS content
func (s *IPFSService) generateGatewayURL(hash string) string {
	if s.config.Gateway == "" {
		return fmt.Sprintf("https://ipfs.io/ipfs/%s", hash)
	}
	return fmt.Sprintf("%s/ipfs/%s", s.config.Gateway, hash)
}

// isValidIPFSHash validates IPFS hash format
func (s *IPFSService) isValidIPFSHash(hash string) bool {
	// Basic validation for IPFS hash
	if len(hash) < 46 {
		return false
	}
	
	// Check for common IPFS hash prefixes
	return strings.HasPrefix(hash, "Qm") || strings.HasPrefix(hash, "bafy") || strings.HasPrefix(hash, "bafk")
}

// GetNodeInfo returns information about the IPFS node
func (s *IPFSService) GetNodeInfo(ctx context.Context) (map[string]interface{}, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ipfs-service").Start(ctx, "ipfs.GetNodeInfo")
	defer span.End()

	id, err := s.client.ID()
	if err != nil {
		return nil, fmt.Errorf("failed to get node ID: %w", err)
	}

	version, _, err := s.client.Version()
	if err != nil {
		return nil, fmt.Errorf("failed to get node version: %w", err)
	}

	info := map[string]interface{}{
		"id":      id.ID,
		"version": version,
		"addresses": id.Addresses,
	}

	return info, nil
}
