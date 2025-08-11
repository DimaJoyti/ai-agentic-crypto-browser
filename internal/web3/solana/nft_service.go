package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// NFTService handles Solana NFT operations
type NFTService struct {
	service         *Service
	magicEdenClient *MagicEdenClient
	tensorClient    *TensorClient
	metaplexClient  *MetaplexClient
}

// NFTMarketplace represents supported NFT marketplaces
type NFTMarketplace string

const (
	MarketplaceMagicEden NFTMarketplace = "magic_eden"
	MarketplaceTensor    NFTMarketplace = "tensor"
	MarketplaceOpenSea   NFTMarketplace = "opensea"
	MarketplaceMetaplex  NFTMarketplace = "metaplex"
)

// NFT represents a Solana NFT
type NFT struct {
	ID                 uuid.UUID         `json:"id"`
	MintAddress        solana.PublicKey  `json:"mint_address"`
	CollectionAddress  *solana.PublicKey `json:"collection_address,omitempty"`
	Name               string            `json:"name"`
	Symbol             string            `json:"symbol"`
	Description        string            `json:"description"`
	Image              string            `json:"image"`
	MetadataURI        string            `json:"metadata_uri"`
	Attributes         []NFTAttribute    `json:"attributes"`
	Creators           []NFTCreator      `json:"creators"`
	RarityRank         *int              `json:"rarity_rank,omitempty"`
	RarityScore        *decimal.Decimal  `json:"rarity_score,omitempty"`
	FloorPrice         decimal.Decimal   `json:"floor_price"`
	LastSalePrice      *decimal.Decimal  `json:"last_sale_price,omitempty"`
	EstimatedValue     decimal.Decimal   `json:"estimated_value"`
	Owner              solana.PublicKey  `json:"owner"`
	IsListed           bool              `json:"is_listed"`
	ListingPrice       *decimal.Decimal  `json:"listing_price,omitempty"`
	ListingMarketplace *NFTMarketplace   `json:"listing_marketplace,omitempty"`
	AcquiredAt         time.Time         `json:"acquired_at"`
	AcquiredPrice      *decimal.Decimal  `json:"acquired_price,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// NFTAttribute represents an NFT attribute/trait
type NFTAttribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
	Rarity    *float64    `json:"rarity,omitempty"`
}

// NFTCreator represents an NFT creator
type NFTCreator struct {
	Address  solana.PublicKey `json:"address"`
	Verified bool             `json:"verified"`
	Share    uint8            `json:"share"`
}

// NFTCollection represents an NFT collection
type NFTCollection struct {
	ID           uuid.UUID        `json:"id"`
	Address      solana.PublicKey `json:"address"`
	Name         string           `json:"name"`
	Symbol       string           `json:"symbol"`
	Description  string           `json:"description"`
	Image        string           `json:"image"`
	Website      string           `json:"website,omitempty"`
	Twitter      string           `json:"twitter,omitempty"`
	Discord      string           `json:"discord,omitempty"`
	TotalSupply  int              `json:"total_supply"`
	FloorPrice   decimal.Decimal  `json:"floor_price"`
	Volume24h    decimal.Decimal  `json:"volume_24h"`
	Volume7d     decimal.Decimal  `json:"volume_7d"`
	VolumeTotal  decimal.Decimal  `json:"volume_total"`
	Holders      int              `json:"holders"`
	ListedCount  int              `json:"listed_count"`
	AveragePrice decimal.Decimal  `json:"average_price"`
	MarketCap    decimal.Decimal  `json:"market_cap"`
	IsVerified   bool             `json:"is_verified"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// NFTListing represents an NFT listing
type NFTListing struct {
	ID          uuid.UUID        `json:"id"`
	NFTMint     solana.PublicKey `json:"nft_mint"`
	Seller      solana.PublicKey `json:"seller"`
	Price       decimal.Decimal  `json:"price"`
	Currency    string           `json:"currency"`
	Marketplace NFTMarketplace   `json:"marketplace"`
	ListedAt    time.Time        `json:"listed_at"`
	ExpiresAt   *time.Time       `json:"expires_at,omitempty"`
	IsActive    bool             `json:"is_active"`
}

// NFTSale represents an NFT sale transaction
type NFTSale struct {
	ID          uuid.UUID        `json:"id"`
	NFTMint     solana.PublicKey `json:"nft_mint"`
	Seller      solana.PublicKey `json:"seller"`
	Buyer       solana.PublicKey `json:"buyer"`
	Price       decimal.Decimal  `json:"price"`
	Currency    string           `json:"currency"`
	Marketplace NFTMarketplace   `json:"marketplace"`
	Signature   solana.Signature `json:"signature"`
	SoldAt      time.Time        `json:"sold_at"`
}

// ListNFTRequest represents a request to list an NFT
type ListNFTRequest struct {
	NFTMint     solana.PublicKey `json:"nft_mint"`
	Price       decimal.Decimal  `json:"price"`
	Currency    string           `json:"currency"`
	Marketplace NFTMarketplace   `json:"marketplace"`
	Seller      solana.PublicKey `json:"seller"`
	ExpiresAt   *time.Time       `json:"expires_at,omitempty"`
}

// BuyNFTRequest represents a request to buy an NFT
type BuyNFTRequest struct {
	NFTMint     solana.PublicKey `json:"nft_mint"`
	MaxPrice    decimal.Decimal  `json:"max_price"`
	Marketplace NFTMarketplace   `json:"marketplace"`
	Buyer       solana.PublicKey `json:"buyer"`
}

// NewNFTService creates a new NFT service
func NewNFTService(service *Service) *NFTService {
	return &NFTService{
		service:         service,
		magicEdenClient: NewMagicEdenClient(service),
		tensorClient:    NewTensorClient(service),
		metaplexClient:  NewMetaplexClient(service),
	}
}

// MetaplexClient interface for Metaplex operations
type MetaplexClient struct {
	service *Service
}

// NewMetaplexClient creates a new Metaplex client
func NewMetaplexClient(service *Service) *MetaplexClient {
	return &MetaplexClient{service: service}
}

// GetNFTMetadata gets NFT metadata using Metaplex standard
func (m *MetaplexClient) GetNFTMetadata(ctx context.Context, mintAddress solana.PublicKey) (*TokenMetadata, error) {
	return m.service.programMgr.GetTokenMetadata(ctx, mintAddress)
}

// GetUserNFTs retrieves all NFTs owned by a user
func (n *NFTService) GetUserNFTs(ctx context.Context, walletID uuid.UUID) ([]*NFT, error) {
	// Note: Using simplified logging since StartSpan is not available

	query := `
		SELECT id, mint_address, collection_address, name, symbol, description,
			   image_url, metadata_uri, attributes, creators, rarity_rank, rarity_score,
			   floor_price, last_sale_price, estimated_value, is_listed, listing_price,
			   listing_marketplace, acquired_at, acquired_price, created_at, updated_at
		FROM solana_nft_holdings 
		WHERE wallet_id = $1 
		ORDER BY acquired_at DESC
	`

	rows, err := n.service.db.QueryContext(ctx, query, walletID)
	if err != nil {
		n.service.logger.Error(ctx, "Failed to get user NFTs", err)
		return nil, fmt.Errorf("failed to get user NFTs: %w", err)
	}
	defer rows.Close()

	var nfts []*NFT
	for rows.Next() {
		var nft NFT
		var mintStr, collectionStr string
		var listingMarketplaceStr *string

		err := rows.Scan(
			&nft.ID,
			&mintStr,
			&collectionStr,
			&nft.Name,
			&nft.Symbol,
			&nft.Description,
			&nft.Image,
			&nft.MetadataURI,
			&nft.Attributes,
			&nft.Creators,
			&nft.RarityRank,
			&nft.RarityScore,
			&nft.FloorPrice,
			&nft.LastSalePrice,
			&nft.EstimatedValue,
			&nft.IsListed,
			&nft.ListingPrice,
			&listingMarketplaceStr,
			&nft.AcquiredAt,
			&nft.AcquiredPrice,
			&nft.CreatedAt,
			&nft.UpdatedAt,
		)
		if err != nil {
			n.service.logger.Error(ctx, "Failed to scan NFT row", err)
			continue
		}

		// Parse addresses
		nft.MintAddress, _ = solana.PublicKeyFromBase58(mintStr)
		if collectionStr != "" {
			collection, _ := solana.PublicKeyFromBase58(collectionStr)
			nft.CollectionAddress = &collection
		}

		if listingMarketplaceStr != nil {
			marketplace := NFTMarketplace(*listingMarketplaceStr)
			nft.ListingMarketplace = &marketplace
		}

		nfts = append(nfts, &nft)
	}

	return nfts, nil
}

// GetNFTMetadata retrieves metadata for a specific NFT
func (n *NFTService) GetNFTMetadata(ctx context.Context, mintAddress solana.PublicKey) (*NFT, error) {

	// Get metadata from Metaplex
	metadata, err := n.metaplexClient.GetNFTMetadata(ctx, mintAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT metadata: %w", err)
	}

	// Get market data from multiple sources
	floorPrice, err := n.getFloorPrice(ctx, mintAddress)
	if err != nil {
		n.service.logger.Warn(ctx, "Failed to get floor price", map[string]interface{}{
			"mint":  mintAddress.String(),
			"error": err.Error(),
		})
		floorPrice = decimal.Zero
	}

	// Get rarity data
	rarityRank, rarityScore, err := n.getRarityData(ctx, mintAddress)
	if err != nil {
		n.service.logger.Warn(ctx, "Failed to get rarity data", map[string]interface{}{
			"mint":  mintAddress.String(),
			"error": err.Error(),
		})
	}

	nft := &NFT{
		MintAddress:    mintAddress,
		Name:           metadata.Name,
		Symbol:         metadata.Symbol,
		Description:    metadata.Description,
		Image:          metadata.Image,
		MetadataURI:    "", // metadata.URI field not available
		Attributes:     convertAttributes(metadata.Attributes),
		Creators:       convertCreators(metadata.Creators),
		RarityRank:     rarityRank,
		RarityScore:    rarityScore,
		FloorPrice:     floorPrice,
		EstimatedValue: floorPrice, // Use floor price as estimated value
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if metadata.Collection != nil {
		nft.CollectionAddress = metadata.Collection
	}

	return nft, nil
}

// ListNFT lists an NFT for sale on a marketplace
func (n *NFTService) ListNFT(ctx context.Context, req ListNFTRequest) (*NFTListing, error) {

	var result *NFTListing
	var err error

	// List on the specified marketplace
	switch req.Marketplace {
	case MarketplaceMagicEden:
		result, err = n.magicEdenClient.ListNFT(ctx, req)
	case MarketplaceTensor:
		result, err = n.tensorClient.ListNFT(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported marketplace: %s", req.Marketplace)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list NFT: %w", err)
	}

	// Update NFT status in database
	err = n.updateNFTListingStatus(ctx, req.NFTMint, true, &req.Price, &req.Marketplace)
	if err != nil {
		n.service.logger.Error(ctx, "Failed to update NFT listing status", err)
		// Don't return error as the listing was successful
	}

	n.service.logger.Info(ctx, "NFT listed successfully", map[string]interface{}{
		"nft_mint":    req.NFTMint.String(),
		"price":       req.Price.String(),
		"marketplace": req.Marketplace,
		"listing_id":  result.ID.String(),
	})

	return result, nil
}

// BuyNFT purchases an NFT from a marketplace
func (n *NFTService) BuyNFT(ctx context.Context, req BuyNFTRequest) (*NFTSale, error) {
	var result *NFTSale
	var err error

	// Buy from the specified marketplace
	switch req.Marketplace {
	case MarketplaceMagicEden:
		result, err = n.magicEdenClient.BuyNFT(ctx, req)
	case MarketplaceTensor:
		result, err = n.tensorClient.BuyNFT(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported marketplace: %s", req.Marketplace)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to buy NFT: %w", err)
	}

	// Update NFT ownership and listing status
	err = n.updateNFTOwnership(ctx, req.NFTMint, req.Buyer, result.Price)
	if err != nil {
		n.service.logger.Error(ctx, "Failed to update NFT ownership", err)
		// Don't return error as the purchase was successful
	}

	n.service.logger.Info(ctx, "NFT purchased successfully", map[string]interface{}{
		"nft_mint":    req.NFTMint.String(),
		"buyer":       req.Buyer.String(),
		"price":       result.Price.String(),
		"marketplace": req.Marketplace,
		"signature":   result.Signature.String(),
	})

	return result, nil
}

// GetCollections retrieves popular NFT collections
func (n *NFTService) GetCollections(ctx context.Context, limit int, offset int) ([]*NFTCollection, error) {
	// Get collections from Magic Eden (primary source)
	collections, err := n.magicEdenClient.GetCollections(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	return collections, nil
}

// GetCollectionNFTs retrieves NFTs from a specific collection
func (n *NFTService) GetCollectionNFTs(ctx context.Context, collectionAddress solana.PublicKey, limit int, offset int) ([]*NFT, error) {
	query := `
		SELECT mint_address, name, symbol, description, image_url, metadata_uri,
			   attributes, creators, rarity_rank, rarity_score, floor_price,
			   last_sale_price, estimated_value, is_listed, listing_price
		FROM solana_nft_holdings
		WHERE collection_address = $1
		ORDER BY rarity_rank ASC NULLS LAST
		LIMIT $2 OFFSET $3
	`

	rows, err := n.service.db.QueryContext(ctx, query, collectionAddress.String(), limit, offset)
	if err != nil {
		n.service.logger.Error(ctx, "Failed to get collection NFTs", err)
		return nil, fmt.Errorf("failed to get collection NFTs: %w", err)
	}
	defer rows.Close()

	var nfts []*NFT
	for rows.Next() {
		var nft NFT
		var mintStr string

		err := rows.Scan(
			&mintStr,
			&nft.Name,
			&nft.Symbol,
			&nft.Description,
			&nft.Image,
			&nft.MetadataURI,
			&nft.Attributes,
			&nft.Creators,
			&nft.RarityRank,
			&nft.RarityScore,
			&nft.FloorPrice,
			&nft.LastSalePrice,
			&nft.EstimatedValue,
			&nft.IsListed,
			&nft.ListingPrice,
		)
		if err != nil {
			n.service.logger.Error(ctx, "Failed to scan collection NFT row", err)
			continue
		}

		nft.MintAddress, _ = solana.PublicKeyFromBase58(mintStr)
		nft.CollectionAddress = &collectionAddress

		nfts = append(nfts, &nft)
	}

	return nfts, nil
}

// Helper methods

func (n *NFTService) getFloorPrice(ctx context.Context, mintAddress solana.PublicKey) (decimal.Decimal, error) {
	// Try to get floor price from Magic Eden first
	floorPrice, err := n.magicEdenClient.GetFloorPrice(ctx, mintAddress)
	if err == nil {
		return floorPrice, nil
	}

	// Fallback to Tensor
	floorPrice, err = n.tensorClient.GetFloorPrice(ctx, mintAddress)
	if err == nil {
		return floorPrice, nil
	}

	return decimal.Zero, fmt.Errorf("failed to get floor price from any marketplace")
}

func (n *NFTService) getRarityData(ctx context.Context, mintAddress solana.PublicKey) (*int, *decimal.Decimal, error) {
	// Try to get rarity data from Tensor (known for rarity data)
	return n.tensorClient.GetRarityData(ctx, mintAddress)
}

func (n *NFTService) updateNFTListingStatus(ctx context.Context, mintAddress solana.PublicKey, isListed bool, price *decimal.Decimal, marketplace *NFTMarketplace) error {
	query := `
		UPDATE solana_nft_holdings 
		SET is_listed = $1, listing_price = $2, listing_marketplace = $3, updated_at = NOW()
		WHERE mint_address = $4
	`

	var marketplaceStr *string
	if marketplace != nil {
		str := string(*marketplace)
		marketplaceStr = &str
	}

	_, err := n.service.db.ExecContext(ctx, query, isListed, price, marketplaceStr, mintAddress.String())
	return err
}

func (n *NFTService) updateNFTOwnership(ctx context.Context, mintAddress, newOwner solana.PublicKey, salePrice decimal.Decimal) error {
	// Update listing status and last sale price
	query := `
		UPDATE solana_nft_holdings
		SET is_listed = false, listing_price = NULL, listing_marketplace = NULL,
			last_sale_price = $1, updated_at = NOW()
		WHERE mint_address = $2
	`

	_, err := n.service.db.ExecContext(ctx, query, salePrice, mintAddress.String())
	return err
}

func convertAttributes(metaplexAttrs []TokenAttribute) []NFTAttribute {
	var attrs []NFTAttribute
	for _, attr := range metaplexAttrs {
		attrs = append(attrs, NFTAttribute{
			TraitType: attr.TraitType,
			Value:     attr.Value,
		})
	}
	return attrs
}

func convertCreators(metaplexCreators []Creator) []NFTCreator {
	var creators []NFTCreator
	for _, creator := range metaplexCreators {
		creators = append(creators, NFTCreator{
			Address:  creator.Address,
			Verified: creator.Verified,
			Share:    creator.Share,
		})
	}
	return creators
}
