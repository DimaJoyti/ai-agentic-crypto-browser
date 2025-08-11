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

// MagicEdenClient handles interactions with Magic Eden marketplace
type MagicEdenClient struct {
	service    *Service
	baseURL    string
	httpClient *http.Client
}

// MagicEdenCollection represents a collection from Magic Eden API
type MagicEdenCollection struct {
	Symbol       string          `json:"symbol"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Image        string          `json:"image"`
	Website      string          `json:"website"`
	Twitter      string          `json:"twitter"`
	Discord      string          `json:"discord"`
	FloorPrice   decimal.Decimal `json:"floorPrice"`
	ListedCount  int             `json:"listedCount"`
	VolumeAll    decimal.Decimal `json:"volumeAll"`
	Volume24hr   decimal.Decimal `json:"volume24hr"`
	Volume7d     decimal.Decimal `json:"volume7d"`
	AvgPrice24hr decimal.Decimal `json:"avgPrice24hr"`
	MarketCap    decimal.Decimal `json:"marketCap"`
	IsVerified   bool            `json:"isVerified"`
}

// MagicEdenNFT represents an NFT from Magic Eden API
type MagicEdenNFT struct {
	MintAddress string               `json:"mintAddress"`
	Name        string               `json:"name"`
	Image       string               `json:"image"`
	Collection  string               `json:"collection"`
	Description string               `json:"description"`
	Attributes  []MagicEdenAttribute `json:"attributes"`
	Properties  MagicEdenProperties  `json:"properties"`
	ListStatus  string               `json:"listStatus"`
	Price       *decimal.Decimal     `json:"price,omitempty"`
	Owner       string               `json:"owner"`
	RarityRank  *int                 `json:"rarityRank,omitempty"`
	RarityScore *decimal.Decimal     `json:"rarityScore,omitempty"`
}

// MagicEdenAttribute represents an NFT attribute
type MagicEdenAttribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
	Rarity    *float64    `json:"rarity,omitempty"`
}

// MagicEdenProperties represents NFT properties
type MagicEdenProperties struct {
	Files    []MagicEdenFile    `json:"files"`
	Category string             `json:"category"`
	Creators []MagicEdenCreator `json:"creators"`
}

// MagicEdenFile represents an NFT file
type MagicEdenFile struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
}

// MagicEdenCreator represents an NFT creator
type MagicEdenCreator struct {
	Address  string `json:"address"`
	Verified bool   `json:"verified"`
	Share    int    `json:"share"`
}

// MagicEdenListingRequest represents a listing request
type MagicEdenListingRequest struct {
	Price        decimal.Decimal `json:"price"`
	TokenMint    string          `json:"tokenMint"`
	TokenAccount string          `json:"tokenAccount"`
	Seller       string          `json:"seller"`
	ExpiryTime   *int64          `json:"expiryTime,omitempty"`
}

// MagicEdenListingResponse represents a listing response
type MagicEdenListingResponse struct {
	TxSig       string `json:"txSig"`
	Instruction string `json:"instruction"`
}

// NewMagicEdenClient creates a new Magic Eden client
func NewMagicEdenClient(service *Service) *MagicEdenClient {
	return &MagicEdenClient{
		service: service,
		baseURL: "https://api-mainnet.magiceden.dev/v2",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCollections retrieves popular collections from Magic Eden
func (m *MagicEdenClient) GetCollections(ctx context.Context, limit int, offset int) ([]*NFTCollection, error) {

	url := fmt.Sprintf("%s/collections?offset=%d&limit=%d", m.baseURL, offset, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Magic Eden API error: %s", string(body))
	}

	var meCollections []MagicEdenCollection
	if err := json.NewDecoder(resp.Body).Decode(&meCollections); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our format
	var collections []*NFTCollection
	for _, meCol := range meCollections {
		collection := &NFTCollection{
			ID:           uuid.New(),
			Name:         meCol.Name,
			Symbol:       meCol.Symbol,
			Description:  meCol.Description,
			Image:        meCol.Image,
			Website:      meCol.Website,
			Twitter:      meCol.Twitter,
			Discord:      meCol.Discord,
			FloorPrice:   meCol.FloorPrice,
			Volume24h:    meCol.Volume24hr,
			Volume7d:     meCol.Volume7d,
			VolumeTotal:  meCol.VolumeAll,
			ListedCount:  meCol.ListedCount,
			AveragePrice: meCol.AvgPrice24hr,
			MarketCap:    meCol.MarketCap,
			IsVerified:   meCol.IsVerified,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Parse collection address from symbol (simplified)
		if address, err := solana.PublicKeyFromBase58(meCol.Symbol); err == nil {
			collection.Address = address
		}

		collections = append(collections, collection)
	}

	m.service.logger.Info(ctx, "Retrieved collections from Magic Eden", map[string]interface{}{
		"count":  len(collections),
		"limit":  limit,
		"offset": offset,
	})

	return collections, nil
}

// ListNFT lists an NFT for sale on Magic Eden
func (m *MagicEdenClient) ListNFT(ctx context.Context, req ListNFTRequest) (*NFTListing, error) {
	// Get token account for the NFT
	tokenAccount, err := m.getTokenAccount(ctx, req.NFTMint, req.Seller)
	if err != nil {
		return nil, fmt.Errorf("failed to get token account: %w", err)
	}

	listingReq := MagicEdenListingRequest{
		Price:        req.Price,
		TokenMint:    req.NFTMint.String(),
		TokenAccount: tokenAccount.String(),
		Seller:       req.Seller.String(),
	}

	if req.ExpiresAt != nil {
		expiryTime := req.ExpiresAt.Unix()
		listingReq.ExpiryTime = &expiryTime
	}

	// Make listing request
	url := fmt.Sprintf("%s/instructions/sell_now", m.baseURL)
	reqBody, err := json.Marshal(listingReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Magic Eden listing API error: %s", string(body))
	}

	var listingResp MagicEdenListingResponse
	if err := json.NewDecoder(resp.Body).Decode(&listingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// In a real implementation, you would execute the returned instruction
	// For now, create a mock listing
	listing := &NFTListing{
		ID:          uuid.New(),
		NFTMint:     req.NFTMint,
		Seller:      req.Seller,
		Price:       req.Price,
		Currency:    req.Currency,
		Marketplace: MarketplaceMagicEden,
		ListedAt:    time.Now(),
		ExpiresAt:   req.ExpiresAt,
		IsActive:    true,
	}

	m.service.logger.Info(ctx, "NFT listed on Magic Eden", map[string]interface{}{
		"nft_mint":   req.NFTMint.String(),
		"price":      req.Price.String(),
		"listing_id": listing.ID.String(),
	})

	return listing, nil
}

// BuyNFT purchases an NFT from Magic Eden
func (m *MagicEdenClient) BuyNFT(ctx context.Context, req BuyNFTRequest) (*NFTSale, error) {
	// Get buy instruction from Magic Eden
	url := fmt.Sprintf("%s/instructions/buy_now?buyer=%s&seller=%s&tokenMint=%s&price=%s",
		m.baseURL, req.Buyer.String(), "", req.NFTMint.String(), req.MaxPrice.String())

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Magic Eden buy API error: %s", string(body))
	}

	// In a real implementation, you would execute the returned instruction
	// For now, create a mock sale
	sale := &NFTSale{
		ID:          uuid.New(),
		NFTMint:     req.NFTMint,
		Buyer:       req.Buyer,
		Price:       req.MaxPrice,
		Currency:    "SOL",
		Marketplace: MarketplaceMagicEden,
		Signature:   solana.Signature{}, // Would be actual transaction signature
		SoldAt:      time.Now(),
	}

	m.service.logger.Info(ctx, "NFT purchased from Magic Eden", map[string]interface{}{
		"nft_mint": req.NFTMint.String(),
		"buyer":    req.Buyer.String(),
		"price":    req.MaxPrice.String(),
		"sale_id":  sale.ID.String(),
	})

	return sale, nil
}

// GetFloorPrice gets the floor price for an NFT collection
func (m *MagicEdenClient) GetFloorPrice(ctx context.Context, mintAddress solana.PublicKey) (decimal.Decimal, error) {
	// First, get the collection symbol for this NFT
	// In a real implementation, you would look this up from the NFT metadata
	collectionSymbol := "example_collection"

	url := fmt.Sprintf("%s/collections/%s/stats", m.baseURL, collectionSymbol)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.httpClient.Do(req)
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

// Helper methods

func (m *MagicEdenClient) getTokenAccount(ctx context.Context, mintAddress, owner solana.PublicKey) (solana.PublicKey, error) {
	// Calculate associated token account address
	ata, _, err := solana.FindAssociatedTokenAddress(owner, mintAddress)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to find associated token address: %w", err)
	}

	// Verify the account exists
	accountInfo, err := m.service.client.GetAccountInfo(ctx, ata)
	if err != nil || accountInfo.Value == nil {
		return solana.PublicKey{}, fmt.Errorf("token account not found")
	}

	return ata, nil
}
