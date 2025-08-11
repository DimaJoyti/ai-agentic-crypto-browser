package web3

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ai-agentic-browser/pkg/database"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Minimal ERC-20 ABI with only required functions
const erc20ABIJSON = `[
  {"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},
  {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}
]`

var parsedERC20ABI abi.ABI

func init() {
	abiParsed, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		panic(fmt.Errorf("failed to parse ERC20 ABI: %w", err))
	}
	parsedERC20ABI = abiParsed
}

// getEthClient returns (and lazily initializes) an ethclient for the given chain.
func (s *Service) getEthClient(ctx context.Context, chainID int) (*ethclient.Client, error) {
	provider, ok := s.providers[chainID]
	if !ok {
		return nil, fmt.Errorf("no provider configured for chain ID: %d", chainID)
	}
	if provider.Client != nil {
		if c, ok := provider.Client.(*ethclient.Client); ok {
			return c, nil
		}
	}
	client, err := ethclient.DialContext(ctx, provider.RpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC for chain %d: %w", chainID, err)
	}
	provider.Client = client
	return client, nil
}

// getERC20Decimals reads token decimals with Redis L3 caching.
func (s *Service) getERC20Decimals(ctx context.Context, chainID int, tokenAddr string) (int, error) {
	key := fmt.Sprintf("token:decimals:%d:%s", chainID, strings.ToLower(tokenAddr))
	if s.redis != nil {
		if data, found, err := s.redis.GetLayered(ctx, key); err == nil && found {
			switch v := data.(type) {
			case float64:
				return int(v), nil
			case int:
				return v, nil
			case string:
				if v == "" {
					return 0, fmt.Errorf("empty cached decimals")
				}
				// attempt parse single-byte string fallbacks
				return int(v[0]), nil
			}
		}
	}

	client, err := s.getEthClient(ctx, chainID)
	if err != nil {
		return 0, err
	}
	to := common.HexToAddress(tokenAddr)
	callData, err := parsedERC20ABI.Pack("decimals")
	if err != nil {
		return 0, fmt.Errorf("abi pack decimals: %w", err)
	}
	res, err := client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: callData}, nil)
	if err != nil {
		return 0, fmt.Errorf("call decimals failed: %w", err)
	}
	var out []interface{}
	if err := parsedERC20ABI.UnpackIntoInterface(&out, "decimals", res); err != nil {
		return 0, fmt.Errorf("unpack decimals: %w", err)
	}
	if len(out) != 1 {
		return 0, fmt.Errorf("unexpected decimals output")
	}
	var dec int
	switch v := out[0].(type) {
	case uint8:
		dec = int(v)
	case uint16:
		dec = int(v)
	case *big.Int:
		dec = int(v.Int64())
	default:
		return 0, fmt.Errorf("unknown decimals type %T", v)
	}
	if s.redis != nil {
		_ = s.redis.SetLayered(ctx, key, dec, database.L3Cache)
	}
	return dec, nil
}

// getERC20Balance reads token balanceOf with Redis L1 caching (as string).
func (s *Service) getERC20Balance(ctx context.Context, chainID int, tokenAddr, walletAddr string) (*big.Int, error) {
	key := fmt.Sprintf("balance:%d:%s:%s", chainID, strings.ToLower(walletAddr), strings.ToLower(tokenAddr))
	if s.redis != nil {
		if data, found, err := s.redis.GetLayered(ctx, key); err == nil && found {
			switch v := data.(type) {
			case string:
				bi, ok := new(big.Int).SetString(v, 10)
				if !ok {
					return nil, fmt.Errorf("invalid cached big.Int")
				}
				return bi, nil
			}
		}
	}

	client, err := s.getEthClient(ctx, chainID)
	if err != nil {
		return nil, err
	}
	to := common.HexToAddress(tokenAddr)
	owner := common.HexToAddress(walletAddr)
	callData, err := parsedERC20ABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, fmt.Errorf("abi pack balanceOf: %w", err)
	}
	res, err := client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: callData}, nil)
	if err != nil {
		return nil, fmt.Errorf("call balanceOf failed: %w", err)
	}
	var out []interface{}
	if err := parsedERC20ABI.UnpackIntoInterface(&out, "balanceOf", res); err != nil {
		return nil, fmt.Errorf("unpack balanceOf: %w", err)
	}
	if len(out) != 1 {
		return nil, fmt.Errorf("unexpected balanceOf output")
	}
	bal, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected balance type %T", out[0])
	}
	if s.redis != nil {
		_ = s.redis.SetLayered(ctx, key, bal.String(), database.L1Cache)
	}
	return bal, nil
}

// getNativeBalance reads and caches native coin balance in L1
func (s *Service) getNativeBalance(ctx context.Context, chainID int, walletAddr string) (*big.Int, error) {
	key := fmt.Sprintf("balance:native:%d:%s", chainID, strings.ToLower(walletAddr))
	if s.redis != nil {
		if data, found, err := s.redis.GetLayered(ctx, key); err == nil && found {
			if v, ok := data.(string); ok {
				if bi, ok2 := new(big.Int).SetString(v, 10); ok2 {
					return bi, nil
				}
			}
		}
	}
	client, err := s.getEthClient(ctx, chainID)
	if err != nil {
		return nil, err
	}
	bal, err := client.BalanceAt(ctx, common.HexToAddress(walletAddr), nil)
	if err != nil {
		return nil, fmt.Errorf("native balance fetch failed: %w", err)
	}
	if s.redis != nil {
		_ = s.redis.SetLayered(ctx, key, bal.String(), database.L1Cache)
	}
	return bal, nil
}
