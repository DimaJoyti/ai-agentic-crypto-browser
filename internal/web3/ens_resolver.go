package web3

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/observability"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v3"
)

// ENSResolver provides Ethereum Name Service resolution functionality
type ENSResolver struct {
	client *ethclient.Client
	logger *observability.Logger
	cache  map[string]*ENSRecord
}

// ENSRecord represents an ENS record with metadata
type ENSRecord struct {
	Name        string            `json:"name"`
	Address     common.Address    `json:"address"`
	ContentHash string            `json:"content_hash"`
	TextRecords map[string]string `json:"text_records"`
	ResolvedAt  time.Time         `json:"resolved_at"`
	TTL         time.Duration     `json:"ttl"`
}

// ENSResolveRequest represents a request to resolve an ENS name
type ENSResolveRequest struct {
	Name           string   `json:"name"`
	ResolveAddress bool     `json:"resolve_address"`
	ResolveContent bool     `json:"resolve_content"`
	TextKeys       []string `json:"text_keys"`
}

// ENSResolveResponse represents the response from ENS resolution
type ENSResolveResponse struct {
	Record *ENSRecord `json:"record"`
	Cached bool       `json:"cached"`
}

// NewENSResolver creates a new ENS resolver
func NewENSResolver(client *ethclient.Client, logger *observability.Logger) *ENSResolver {
	return &ENSResolver{
		client: client,
		logger: logger,
		cache:  make(map[string]*ENSRecord),
	}
}

// Resolve resolves an ENS name to its associated records
func (r *ENSResolver) Resolve(ctx context.Context, req ENSResolveRequest) (*ENSResolveResponse, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ens-resolver").Start(ctx, "ens.Resolve")
	defer span.End()

	// Normalize ENS name
	name := strings.ToLower(strings.TrimSpace(req.Name))
	if !r.isValidENSName(name) {
		return nil, fmt.Errorf("invalid ENS name format: %s", name)
	}

	// Check cache first
	if cached := r.getCachedRecord(name); cached != nil {
		r.logger.Info(ctx, "ENS record found in cache", map[string]interface{}{
			"name": name,
		})
		return &ENSResolveResponse{
			Record: cached,
			Cached: true,
		}, nil
	}

	// Resolve from ENS
	record, err := r.resolveFromENS(ctx, name, req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve ENS name: %w", err)
	}

	// Cache the result
	r.cacheRecord(name, record)

	r.logger.Info(ctx, "ENS name resolved", map[string]interface{}{
		"name":    name,
		"address": record.Address.Hex(),
	})

	return &ENSResolveResponse{
		Record: record,
		Cached: false,
	}, nil
}

// ResolveAddress resolves an Ethereum address to its ENS name (reverse resolution)
func (r *ENSResolver) ResolveAddress(ctx context.Context, address common.Address) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ens-resolver").Start(ctx, "ens.ResolveAddress")
	defer span.End()

	// Perform reverse ENS lookup
	name, err := ens.ReverseResolve(r.client, address)
	if err != nil {
		return "", fmt.Errorf("failed to reverse resolve address: %w", err)
	}

	// Verify the forward resolution matches
	forwardAddress, err := ens.Resolve(r.client, name)
	if err != nil || forwardAddress != address {
		return "", fmt.Errorf("reverse resolution verification failed")
	}

	r.logger.Info(ctx, "Address reverse resolved", map[string]interface{}{
		"address": address.Hex(),
		"name":    name,
	})

	return name, nil
}

// resolveFromENS performs the actual ENS resolution
func (r *ENSResolver) resolveFromENS(ctx context.Context, name string, req ENSResolveRequest) (*ENSRecord, error) {
	record := &ENSRecord{
		Name:        name,
		TextRecords: make(map[string]string),
		ResolvedAt:  time.Now(),
		TTL:         time.Hour, // Default TTL
	}

	// Resolve address if requested
	if req.ResolveAddress {
		address, err := ens.Resolve(r.client, name)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve address: %w", err)
		}
		record.Address = address
	}

	// Resolve content hash if requested
	if req.ResolveContent {
		// Note: Content hash resolution would require specific ENS resolver contract calls
		// For now, we'll skip this implementation as it requires more complex contract interaction
		r.logger.Info(ctx, "Content hash resolution not implemented", map[string]interface{}{
			"name": name,
		})
	}

	// Resolve text records if requested
	if len(req.TextKeys) > 0 {
		// Note: Text record resolution would require specific ENS resolver contract calls
		// For now, we'll skip this implementation as it requires more complex contract interaction
		r.logger.Info(ctx, "Text record resolution not implemented", map[string]interface{}{
			"name": name,
			"keys": req.TextKeys,
		})
	}

	return record, nil
}

// GetContentURL resolves an ENS name to its content URL (IPFS, HTTP, etc.)
func (r *ENSResolver) GetContentURL(ctx context.Context, name string) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ens-resolver").Start(ctx, "ens.GetContentURL")
	defer span.End()

	// Note: Content hash resolution would require specific ENS resolver contract calls
	// For now, return a placeholder implementation
	r.logger.Info(ctx, "Content URL resolution not implemented", map[string]interface{}{
		"name": name,
	})

	return "", fmt.Errorf("content URL resolution not implemented for ENS name: %s", name)
}

// GetTextRecord retrieves a specific text record for an ENS name
func (r *ENSResolver) GetTextRecord(ctx context.Context, name, key string) (string, error) {
	ctx, span := observability.SpanFromContext(ctx).TracerProvider().Tracer("ens-resolver").Start(ctx, "ens.GetTextRecord")
	defer span.End()

	// Note: Text record resolution would require specific ENS resolver contract calls
	// For now, return a placeholder implementation
	r.logger.Info(ctx, "Text record resolution not implemented", map[string]interface{}{
		"name": name,
		"key":  key,
	})

	return "", fmt.Errorf("text record resolution not implemented for ENS name: %s, key: %s", name, key)
}

// IsENSName checks if a string is a valid ENS name
func (r *ENSResolver) IsENSName(name string) bool {
	return r.isValidENSName(name)
}

// isValidENSName validates ENS name format
func (r *ENSResolver) isValidENSName(name string) bool {
	if name == "" {
		return false
	}

	// Must end with .eth or other valid TLD
	validTLDs := []string{".eth", ".xyz", ".luxe", ".kred", ".art"}
	hasValidTLD := false
	for _, tld := range validTLDs {
		if strings.HasSuffix(name, tld) {
			hasValidTLD = true
			break
		}
	}

	if !hasValidTLD {
		return false
	}

	// Basic format validation
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return false
	}

	// Check for valid characters
	for _, part := range parts[:len(parts)-1] { // Exclude TLD
		if part == "" {
			return false
		}
		for _, char := range part {
			if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
				return false
			}
		}
	}

	return true
}

// getCachedRecord retrieves a cached ENS record
func (r *ENSResolver) getCachedRecord(name string) *ENSRecord {
	record, exists := r.cache[name]
	if !exists {
		return nil
	}

	// Check if cache entry is still valid
	if time.Since(record.ResolvedAt) > record.TTL {
		delete(r.cache, name)
		return nil
	}

	return record
}

// cacheRecord caches an ENS record
func (r *ENSResolver) cacheRecord(name string, record *ENSRecord) {
	r.cache[name] = record
}

// contentHashToURL converts a content hash to a URL
func (r *ENSResolver) contentHashToURL(contentHash string) (string, error) {
	if contentHash == "" {
		return "", fmt.Errorf("empty content hash")
	}

	// Handle IPFS content hashes
	if strings.HasPrefix(contentHash, "ipfs://") {
		hash := strings.TrimPrefix(contentHash, "ipfs://")
		return fmt.Sprintf("https://ipfs.io/ipfs/%s", hash), nil
	}

	// Handle HTTP/HTTPS URLs
	if strings.HasPrefix(contentHash, "http://") || strings.HasPrefix(contentHash, "https://") {
		return contentHash, nil
	}

	// Handle raw IPFS hashes
	if len(contentHash) == 46 && strings.HasPrefix(contentHash, "Qm") {
		return fmt.Sprintf("https://ipfs.io/ipfs/%s", contentHash), nil
	}

	// Handle other content hash formats
	if strings.HasPrefix(contentHash, "bafy") || strings.HasPrefix(contentHash, "bafk") {
		return fmt.Sprintf("https://ipfs.io/ipfs/%s", contentHash), nil
	}

	return "", fmt.Errorf("unsupported content hash format: %s", contentHash)
}

// ClearCache clears the ENS resolution cache
func (r *ENSResolver) ClearCache() {
	r.cache = make(map[string]*ENSRecord)
}

// GetCacheStats returns statistics about the cache
func (r *ENSResolver) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"entries": len(r.cache),
	}
}
