package solana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TensorClient handles interactions with Tensor marketplace
type TensorClient struct {
	service    *Service
	baseURL    string
	httpClient *http.Client
}

// TensorCollection represents a collection from Tensor API
type TensorCollection struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Symbol      string          `json:"symbol"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	FloorPrice  decimal.Decimal `json:"floorPrice"`
	Volume24h   decimal.Decimal `json:"volume24h"`
	Volume7d    decimal.Decimal `json:"volume7d"`
	VolumeTotal decimal.Decimal `json:"volumeTotal"`
	MarketCap   decimal.Decimal `json:"marketCap"`
	Supply      int             `json:"supply"`
	NumListed   int             `json:"numListed"`
	NumOwners   int             `json:"numOwners"`
	IsVerified  bool            `json:"isVerified"`
}

// TensorNFT represents an NFT from Tensor API
type TensorNFT struct {
	MintAddress string            `json:"mintAddress"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Collection  string            `json:"collection"`
	Attributes  []TensorAttribute `json:"attributes"`
	RarityRank  *int              `json:"rarityRank,omitempty"`
	RarityScore *decimal.Decimal  `json:"rarityScore,omitempty"`
	ListPrice   *decimal.Decimal  `json:"listPrice,omitempty"`
	LastSale    *TensorSale       `json:"lastSale,omitempty"`
	Owner       string            `json:"owner"`
}

// TensorAttribute represents an NFT attribute with rarity
type TensorAttribute struct {
	TraitType string      `json:"traitType"`
	Value     interface{} `json:"value"`
	Rarity    float64     `json:"rarity"`
	Count     int         `json:"count"`
}

// TensorSale represents a sale record
type TensorSale struct {
	Price     decimal.Decimal `json:"price"`
	Timestamp time.Time       `json:"timestamp"`
	Buyer     string          `json:"buyer"`
	Seller    string          `json:"seller"`
	TxSig     string          `json:"txSig"`
}

// TensorListingRequest represents a listing request for Tensor
type TensorListingRequest struct {
	MintAddress string          `json:"mintAddress"`
	Price       decimal.Decimal `json:"price"`
	Seller      string          `json:"seller"`
	ExpiresAt   *time.Time      `json:"expiresAt,omitempty"`
}

// TensorBuyRequest represents a buy request for Tensor
type TensorBuyRequest struct {
	MintAddress string          `json:"mintAddress"`
	MaxPrice    decimal.Decimal `json:"maxPrice"`
	Buyer       string          `json:"buyer"`
}

// TensorRarityData represents rarity information
type TensorRarityData struct {
	Rank  int             `json:"rank"`
	Score decimal.Decimal `json:"score"`
	Total int             `json:"total"`
}

// NewTensorClient creates a new Tensor client
func NewTensorClient(service *Service) *TensorClient {
	return &TensorClient{
		service: service,
		baseURL: "https://api.tensor.trade/api/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListNFT lists an NFT for sale on Tensor
func (t *TensorClient) ListNFT(ctx context.Context, req ListNFTRequest) (*NFTListing, error) {
	tensorReq := TensorListingRequest{
		MintAddress: req.NFTMint.String(),
		Price:       req.Price,
		Seller:      req.Seller.String(),
		ExpiresAt:   req.ExpiresAt,
	}

	// Make listing request to Tensor API
	url := fmt.Sprintf("%s/listings", t.baseURL)
	reqBody, err := json.Marshal(tensorReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Tensor listing API error: %s", string(body))
	}

	// Create listing record
	listing := &NFTListing{
		ID:          uuid.New(),
		NFTMint:     req.NFTMint,
		Seller:      req.Seller,
		Price:       req.Price,
		Currency:    req.Currency,
		Marketplace: MarketplaceTensor,
		ListedAt:    time.Now(),
		ExpiresAt:   req.ExpiresAt,
		IsActive:    true,
	}

	t.service.logger.Info(ctx, "NFT listed on Tensor", map[string]interface{}{
		"nft_mint":   req.NFTMint.String(),
		"price":      req.Price.String(),
		"listing_id": listing.ID.String(),
	})

	return listing, nil
}

// BuyNFT purchases an NFT from Tensor
func (t *TensorClient) BuyNFT(ctx context.Context, req BuyNFTRequest) (*NFTSale, error) {
	tensorReq := TensorBuyRequest{
		MintAddress: req.NFTMint.String(),
		MaxPrice:    req.MaxPrice,
		Buyer:       req.Buyer.String(),
	}

	// Make buy request to Tensor API
	url := fmt.Sprintf("%s/buy", t.baseURL)
	reqBody, err := json.Marshal(tensorReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Tensor buy API error: %s", string(body))
	}

	// Create sale record
	sale := &NFTSale{
		ID:          uuid.New(),
		NFTMint:     req.NFTMint,
		Buyer:       req.Buyer,
		Price:       req.MaxPrice,
		Currency:    "SOL",
		Marketplace: MarketplaceTensor,
		Signature:   solana.Signature{}, // Would be actual transaction signature
		SoldAt:      time.Now(),
	}

	t.service.logger.Info(ctx, "NFT purchased from Tensor", map[string]interface{}{
		"nft_mint": req.NFTMint.String(),
		"buyer":    req.Buyer.String(),
		"price":    req.MaxPrice.String(),
		"sale_id":  sale.ID.String(),
	})

	return sale, nil
}

// GetFloorPrice gets the floor price for an NFT collection
func (t *TensorClient) GetFloorPrice(ctx context.Context, mintAddress solana.PublicKey) (decimal.Decimal, error) {
	// Get collection info for this NFT
	collectionID, err := t.getCollectionID(ctx, mintAddress)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get collection ID: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/stats", t.baseURL, collectionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("collection not found or API error")
	}

	var stats struct {
		FloorPrice decimal.Decimal `json:"floorPrice"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return stats.FloorPrice, nil
}

// GetRarityData gets rarity information for an NFT
func (t *TensorClient) GetRarityData(ctx context.Context, mintAddress solana.PublicKey) (*int, *decimal.Decimal, error) {

	url := fmt.Sprintf("%s/nfts/%s/rarity", t.baseURL, mintAddress.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("rarity data not found")
	}

	var rarityData TensorRarityData
	if err := json.NewDecoder(resp.Body).Decode(&rarityData); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	t.service.logger.Info(ctx, "Retrieved rarity data from Tensor", map[string]interface{}{
		"mint":  mintAddress.String(),
		"rank":  rarityData.Rank,
		"score": rarityData.Score.String(),
		"total": rarityData.Total,
	})

	return &rarityData.Rank, &rarityData.Score, nil
}

// GetCollections retrieves collections from Tensor
func (t *TensorClient) GetCollections(ctx context.Context, limit int, offset int) ([]*NFTCollection, error) {

	url := fmt.Sprintf("%s/collections?limit=%d&offset=%d&sortBy=volume24h", t.baseURL, limit, offset)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Tensor API error: %s", string(body))
	}

	var tensorCollections []TensorCollection
	if err := json.NewDecoder(resp.Body).Decode(&tensorCollections); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our format
	var collections []*NFTCollection
	for _, tCol := range tensorCollections {
		collection := &NFTCollection{
			ID:          uuid.New(),
			Name:        tCol.Name,
			Symbol:      tCol.Symbol,
			Description: tCol.Description,
			Image:       tCol.Image,
			TotalSupply: tCol.Supply,
			FloorPrice:  tCol.FloorPrice,
			Volume24h:   tCol.Volume24h,
			Volume7d:    tCol.Volume7d,
			VolumeTotal: tCol.VolumeTotal,
			Holders:     tCol.NumOwners,
			ListedCount: tCol.NumListed,
			MarketCap:   tCol.MarketCap,
			IsVerified:  tCol.IsVerified,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Parse collection address from ID
		if address, err := solana.PublicKeyFromBase58(tCol.ID); err == nil {
			collection.Address = address
		}

		collections = append(collections, collection)
	}

	t.service.logger.Info(ctx, "Retrieved collections from Tensor", map[string]interface{}{
		"count":  len(collections),
		"limit":  limit,
		"offset": offset,
	})

	return collections, nil
}

// GetNFTsByCollection retrieves NFTs from a specific collection
func (t *TensorClient) GetNFTsByCollection(ctx context.Context, collectionID string, limit int, offset int) ([]*TensorNFT, error) {

	url := fmt.Sprintf("%s/collections/%s/nfts?limit=%d&offset=%d&sortBy=rarityRank",
		t.baseURL, collectionID, limit, offset)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Tensor API error: %s", string(body))
	}

	var nfts []TensorNFT
	if err := json.NewDecoder(resp.Body).Decode(&nfts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to pointer slice
	var result []*TensorNFT
	for i := range nfts {
		result = append(result, &nfts[i])
	}

	t.service.logger.Info(ctx, "Retrieved NFTs from Tensor", map[string]interface{}{
		"collection": collectionID,
		"count":      len(result),
		"limit":      limit,
		"offset":     offset,
	})

	return result, nil
}

// Helper methods

func (t *TensorClient) getCollectionID(ctx context.Context, mintAddress solana.PublicKey) (string, error) {
	// In a real implementation, this would look up the collection ID from the NFT metadata
	// For now, return a placeholder
	return "example_collection_id", nil
}

// GetNFTAnalytics gets detailed analytics for an NFT
func (t *TensorClient) GetNFTAnalytics(ctx context.Context, mintAddress solana.PublicKey) (*NFTAnalytics, error) {

	url := fmt.Sprintf("%s/nfts/%s/analytics", t.baseURL, mintAddress.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analytics not found")
	}

	var analytics NFTAnalytics
	if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &analytics, nil
}

// NFTAnalytics represents detailed NFT analytics
type NFTAnalytics struct {
	MintAddress   string          `json:"mintAddress"`
	PriceHistory  []PricePoint    `json:"priceHistory"`
	VolumeHistory []VolumePoint   `json:"volumeHistory"`
	HolderHistory []HolderPoint   `json:"holderHistory"`
	TraitAnalysis []TraitAnalysis `json:"traitAnalysis"`
	MarketMetrics MarketMetrics   `json:"marketMetrics"`
}

// PricePoint represents a price data point
type PricePoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Price     decimal.Decimal `json:"price"`
	Volume    decimal.Decimal `json:"volume"`
}

// VolumePoint represents a volume data point
type VolumePoint struct {
	Timestamp time.Time       `json:"timestamp"`
	Volume    decimal.Decimal `json:"volume"`
	Sales     int             `json:"sales"`
}

// HolderPoint represents holder count data
type HolderPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Holders   int       `json:"holders"`
}

// TraitAnalysis represents trait rarity analysis
type TraitAnalysis struct {
	TraitType  string          `json:"traitType"`
	Value      string          `json:"value"`
	Rarity     float64         `json:"rarity"`
	FloorPrice decimal.Decimal `json:"floorPrice"`
}

// MarketMetrics represents market metrics
type MarketMetrics struct {
	FloorPrice      decimal.Decimal `json:"floorPrice"`
	CeilingPrice    decimal.Decimal `json:"ceilingPrice"`
	AveragePrice    decimal.Decimal `json:"averagePrice"`
	MedianPrice     decimal.Decimal `json:"medianPrice"`
	Volume24h       decimal.Decimal `json:"volume24h"`
	Sales24h        int             `json:"sales24h"`
	PriceChange24h  decimal.Decimal `json:"priceChange24h"`
	VolumeChange24h decimal.Decimal `json:"volumeChange24h"`
}
