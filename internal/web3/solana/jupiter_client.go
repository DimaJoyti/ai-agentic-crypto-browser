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
	"github.com/shopspring/decimal"
)

// JupiterClient handles interactions with Jupiter DEX aggregator
type JupiterClient struct {
	service    *Service
	baseURL    string
	httpClient *http.Client
}

// JupiterQuoteRequest represents a Jupiter quote request
type JupiterQuoteRequest struct {
	InputMint           string `json:"inputMint"`
	OutputMint          string `json:"outputMint"`
	Amount              string `json:"amount"`
	SlippageBps         int    `json:"slippageBps"`
	OnlyDirectRoutes    bool   `json:"onlyDirectRoutes,omitempty"`
	AsLegacyTransaction bool   `json:"asLegacyTransaction,omitempty"`
	PlatformFeeBps      int    `json:"platformFeeBps,omitempty"`
	MaxAccounts         int    `json:"maxAccounts,omitempty"`
}

// JupiterQuoteResponse represents a Jupiter quote response
type JupiterQuoteResponse struct {
	InputMint            string              `json:"inputMint"`
	InAmount             string              `json:"inAmount"`
	OutputMint           string              `json:"outputMint"`
	OutAmount            string              `json:"outAmount"`
	OtherAmountThreshold string              `json:"otherAmountThreshold"`
	SwapMode             string              `json:"swapMode"`
	SlippageBps          int                 `json:"slippageBps"`
	PlatformFee          *JupiterPlatformFee `json:"platformFee,omitempty"`
	PriceImpactPct       string              `json:"priceImpactPct"`
	RoutePlan            []JupiterRoutePlan  `json:"routePlan"`
	ContextSlot          int64               `json:"contextSlot"`
	TimeTaken            float64             `json:"timeTaken"`
}

// JupiterPlatformFee represents platform fee information
type JupiterPlatformFee struct {
	Amount     string `json:"amount"`
	FeeBps     int    `json:"feeBps"`
	FeeAccount string `json:"feeAccount"`
}

// JupiterRoutePlan represents a step in the swap route
type JupiterRoutePlan struct {
	SwapInfo JupiterSwapInfo `json:"swapInfo"`
	Percent  int             `json:"percent"`
}

// JupiterSwapInfo represents swap information for a route step
type JupiterSwapInfo struct {
	AmmKey     string `json:"ammKey"`
	Label      string `json:"label"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	InAmount   string `json:"inAmount"`
	OutAmount  string `json:"outAmount"`
	FeeAmount  string `json:"feeAmount"`
	FeeMint    string `json:"feeMint"`
}

// JupiterSwapRequest represents a Jupiter swap execution request
type JupiterSwapRequest struct {
	QuoteResponse                 JupiterQuoteResponse `json:"quoteResponse"`
	UserPublicKey                 string               `json:"userPublicKey"`
	WrapAndUnwrapSol              bool                 `json:"wrapAndUnwrapSol"`
	UseSharedAccounts             bool                 `json:"useSharedAccounts,omitempty"`
	FeeAccount                    string               `json:"feeAccount,omitempty"`
	TrackingAccount               string               `json:"trackingAccount,omitempty"`
	ComputeUnitPriceMicroLamports int                  `json:"computeUnitPriceMicroLamports,omitempty"`
}

// JupiterSwapResponse represents a Jupiter swap execution response
type JupiterSwapResponse struct {
	SwapTransaction      string `json:"swapTransaction"`
	LastValidBlockHeight int64  `json:"lastValidBlockHeight"`
	PriorityFeeEstimate  struct {
		PriorityFeeEstimate int `json:"priorityFeeEstimate"`
	} `json:"priorityFeeEstimate"`
}

// NewJupiterClient creates a new Jupiter client
func NewJupiterClient(service *Service) *JupiterClient {
	return &JupiterClient{
		service: service,
		baseURL: "https://quote-api.jup.ag/v6",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetSwapQuote gets a swap quote from Jupiter
func (j *JupiterClient) GetSwapQuote(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	j.service.logger.Info(ctx, "Getting Jupiter swap quote", map[string]interface{}{
		"operation":   "GetSwapQuote",
		"input_mint":  req.InputMint.String(),
		"output_mint": req.OutputMint.String(),
		"amount":      req.Amount.String(),
	})

	// Convert amount to smallest unit (considering token decimals)
	// For simplicity, assuming 9 decimals for most tokens
	amount := req.Amount.Mul(decimal.NewFromInt(1e9)).String()

	quoteReq := JupiterQuoteRequest{
		InputMint:   req.InputMint.String(),
		OutputMint:  req.OutputMint.String(),
		Amount:      amount,
		SlippageBps: int(req.SlippageBps),
	}

	// Make request to Jupiter API
	url := fmt.Sprintf("%s/quote", j.baseURL)
	reqBody, err := json.Marshal(quoteReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal quote request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := j.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jupiter API error: %s", string(body))
	}

	var quoteResp JupiterQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return nil, fmt.Errorf("failed to decode quote response: %w", err)
	}

	// Convert response to SwapResult
	inputAmount, _ := decimal.NewFromString(quoteResp.InAmount)
	outputAmount, _ := decimal.NewFromString(quoteResp.OutAmount)
	priceImpact, _ := decimal.NewFromString(quoteResp.PriceImpactPct)

	// Convert amounts back from smallest unit
	inputAmount = inputAmount.Div(decimal.NewFromInt(1e9))
	outputAmount = outputAmount.Div(decimal.NewFromInt(1e9))

	// Build route information
	var routes []SwapRoute
	for _, routePlan := range quoteResp.RoutePlan {
		inAmount, _ := decimal.NewFromString(routePlan.SwapInfo.InAmount)
		outAmount, _ := decimal.NewFromString(routePlan.SwapInfo.OutAmount)
		feeAmount, _ := decimal.NewFromString(routePlan.SwapInfo.FeeAmount)

		inputMint, _ := solana.PublicKeyFromBase58(routePlan.SwapInfo.InputMint)
		outputMint, _ := solana.PublicKeyFromBase58(routePlan.SwapInfo.OutputMint)

		route := SwapRoute{
			Protocol:     ProtocolJupiter,
			InputMint:    inputMint,
			OutputMint:   outputMint,
			InputAmount:  inAmount.Div(decimal.NewFromInt(1e9)),
			OutputAmount: outAmount.Div(decimal.NewFromInt(1e9)),
			Fee:          feeAmount.Div(decimal.NewFromInt(1e9)),
		}
		routes = append(routes, route)
	}

	result := &SwapResult{
		InputAmount:  inputAmount,
		OutputAmount: outputAmount,
		PriceImpact:  priceImpact.Abs(),
		Route:        routes,
		Success:      true,
	}

	j.service.logger.Info(ctx, "Jupiter quote retrieved", map[string]interface{}{
		"input_mint":    req.InputMint.String(),
		"output_mint":   req.OutputMint.String(),
		"input_amount":  result.InputAmount.String(),
		"output_amount": result.OutputAmount.String(),
		"price_impact":  result.PriceImpact.String(),
		"route_steps":   len(routes),
	})

	return result, nil
}

// ExecuteSwap executes a swap through Jupiter
func (j *JupiterClient) ExecuteSwap(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// Log the operation start
	j.service.logger.Info(ctx, "Executing Jupiter swap", map[string]interface{}{
		"operation": "ExecuteSwap",
		"amount":    req.Amount.String(),
	})

	// First get a quote
	quote, err := j.GetSwapQuote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Get the raw quote response for swap execution
	quoteResp, err := j.getRawQuote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw quote: %w", err)
	}

	// Prepare swap request
	swapReq := JupiterSwapRequest{
		QuoteResponse:                 *quoteResp,
		UserPublicKey:                 req.UserPublicKey.String(),
		WrapAndUnwrapSol:              true,
		UseSharedAccounts:             true,
		ComputeUnitPriceMicroLamports: j.getComputeUnitPrice(req.SlippageBps),
	}

	// Make swap request to Jupiter API
	url := fmt.Sprintf("%s/swap", j.baseURL)
	reqBody, err := json.Marshal(swapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal swap request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := j.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jupiter swap API error: %s", string(body))
	}

	var swapResp JupiterSwapResponse
	if err := json.NewDecoder(resp.Body).Decode(&swapResp); err != nil {
		return nil, fmt.Errorf("failed to decode swap response: %w", err)
	}

	// In a real implementation, you would:
	// 1. Decode the transaction from swapResp.SwapTransaction
	// 2. Sign the transaction with the user's wallet
	// 3. Send the transaction to the Solana network
	// 4. Wait for confirmation

	// For now, simulate a successful transaction
	signature := solana.Signature{} // Would be the actual transaction signature

	result := &SwapResult{
		Signature:    signature,
		InputAmount:  quote.InputAmount,
		OutputAmount: quote.OutputAmount,
		PriceImpact:  quote.PriceImpact,
		Route:        quote.Route,
		Success:      true,
	}

	j.service.logger.Info(ctx, "Jupiter swap executed", map[string]interface{}{
		"signature":     signature.String(),
		"input_amount":  result.InputAmount.String(),
		"output_amount": result.OutputAmount.String(),
	})

	return result, nil
}

// GetTVL gets Jupiter's total value locked
func (j *JupiterClient) GetTVL(ctx context.Context) (decimal.Decimal, error) {
	// Jupiter is an aggregator, so TVL would be the sum of all integrated protocols
	// This is a simplified implementation
	return decimal.NewFromInt(2000000000), nil // $2B simulated TVL
}

// Helper methods

func (j *JupiterClient) getRawQuote(ctx context.Context, req SwapRequest) (*JupiterQuoteResponse, error) {
	amount := req.Amount.Mul(decimal.NewFromInt(1e9)).String()

	quoteReq := JupiterQuoteRequest{
		InputMint:   req.InputMint.String(),
		OutputMint:  req.OutputMint.String(),
		Amount:      amount,
		SlippageBps: int(req.SlippageBps),
	}

	url := fmt.Sprintf("%s/quote", j.baseURL)
	reqBody, err := json.Marshal(quoteReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := j.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quoteResp JupiterQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return nil, err
	}

	return &quoteResp, nil
}

func (j *JupiterClient) getComputeUnitPrice(slippageBps uint16) int {
	// Higher slippage tolerance = higher priority
	if slippageBps >= 1000 { // 10%+
		return 10000 // High priority
	} else if slippageBps >= 500 { // 5%+
		return 5000 // Medium priority
	} else {
		return 1000 // Low priority
	}
}
