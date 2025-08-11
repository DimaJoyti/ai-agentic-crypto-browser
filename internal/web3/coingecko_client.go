package web3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ai-agentic-browser/pkg/database"
)

// CoinGeckoClient fetches market prices with Redis caching and simple rate limiting.
type CoinGeckoClient struct {
	httpClient *http.Client
	redis      *database.RedisClient
	limiter    chan struct{}
}

func NewCoinGeckoClient(redis *database.RedisClient) *CoinGeckoClient {
	c := &CoinGeckoClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		redis:      redis,
		limiter:    make(chan struct{}, 5), // allow ~5 concurrent or rps
	}
	// refill tokens periodically for gentle rate limiting
	go func() {
		t := time.NewTicker(200 * time.Millisecond)
		defer t.Stop()
		for range t.C {
			select {
			case c.limiter <- struct{}{}:
			default:
			}
		}
	}()
	return c
}

var allowedCurrencies = map[string]struct{}{"usd": {}, "eur": {}, "btc": {}}

// GetPrices fetches prices for the given CoinGecko token IDs.
// tokenIds are CoinGecko IDs like "ethereum", "bitcoin", "polygon".
func (c *CoinGeckoClient) GetPrices(ctx context.Context, currency string, tokenIds []string) (map[string]TokenPrice, error) {
	if len(tokenIds) == 0 {
		return map[string]TokenPrice{}, nil
	}
	cur := strings.ToLower(strings.TrimSpace(currency))
	if _, ok := allowedCurrencies[cur]; !ok {
		cur = "usd"
	}
	ids := make([]string, 0, len(tokenIds))
	for _, id := range tokenIds {
		id = strings.ToLower(strings.TrimSpace(id))
		if id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return map[string]TokenPrice{}, nil
	}
	sort.Strings(ids)
	cacheKey := fmt.Sprintf("prices:%s:%s", cur, strings.Join(ids, ","))
	// Try Redis cached JSON
	if c.redis != nil {
		if s, err := c.redis.GetString(ctx, cacheKey); err == nil && s != "" {
			var cached map[string]TokenPrice
			if json.Unmarshal([]byte(s), &cached) == nil {
				return cached, nil
			}
		}
	}

	// Consume limiter token (or timeout)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case c.limiter <- struct{}{}:
		// proceed
	case <-time.After(2 * time.Second):
		return nil, errors.New("coingecko rate limiter timeout")
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&ids=%s", cur, strings.Join(ids, ","))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("coingecko error: status %d", resp.StatusCode)
	}
	var arr []struct {
		ID                         string    `json:"id"`
		Symbol                     string    `json:"symbol"`
		Name                       string    `json:"name"`
		CurrentPrice               float64   `json:"current_price"`
		PriceChange24h             float64   `json:"price_change_24h"`
		PriceChangePercentage24h   float64   `json:"price_change_percentage_24h"`
		MarketCap                  float64   `json:"market_cap"`
		TotalVolume                float64   `json:"total_volume"`
		LastUpdated                time.Time `json:"last_updated"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, err
	}
	out := make(map[string]TokenPrice, len(arr))
	for _, it := range arr {
		out[it.ID] = TokenPrice{
			Token:           it.ID,
			Symbol:          strings.ToUpper(it.Symbol),
			Name:            it.Name,
			Price:           it.CurrentPrice,
			PriceChange24h:  it.PriceChange24h,
			PriceChangePerc: it.PriceChangePercentage24h,
			MarketCap:       it.MarketCap,
			Volume24h:       it.TotalVolume,
			Currency:        strings.ToUpper(cur),
			LastUpdated:     it.LastUpdated,
		}
	}
	if c.redis != nil {
		b, _ := json.Marshal(out)
		_ = c.redis.SetWithExpiry(ctx, cacheKey, string(b), 60*time.Second)
	}
	return out, nil
}

