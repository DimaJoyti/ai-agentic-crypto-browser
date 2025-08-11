package web3

// Common ERC-20 tokens by chain for balance reads
// Symbols and names are static; decimals are discovered on-chain and cached.
// Chain IDs covered: Ethereum (1), Polygon (137), Arbitrum (42161), Optimism (10)

type ERC20Token struct {
	Address     string
	Symbol      string
	Name        string
	CoinGeckoID string
}

var CommonERC20Tokens = map[int][]ERC20Token{
	1: { // Ethereum
		{Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", Symbol: "USDC", Name: "USD Coin", CoinGeckoID: "usd-coin"},
		{Address: "0xdAC17F958D2ee523a2206206994597C13D831ec7", Symbol: "USDT", Name: "Tether USD", CoinGeckoID: "tether"},
		{Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", Symbol: "WETH", Name: "Wrapped Ether", CoinGeckoID: "weth"},
	},
	137: { // Polygon
		{Address: "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", Symbol: "USDC", Name: "USD Coin (Bridged)", CoinGeckoID: "usd-coin"},
		{Address: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F", Symbol: "USDT", Name: "Tether USD", CoinGeckoID: "tether"},
		{Address: "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619", Symbol: "WETH", Name: "Wrapped Ether", CoinGeckoID: "weth"},
	},
	42161: { // Arbitrum
		{Address: "0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8", Symbol: "USDC", Name: "USD Coin (Bridged)", CoinGeckoID: "usd-coin"},
		{Address: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9", Symbol: "USDT", Name: "Tether USD", CoinGeckoID: "tether"},
		{Address: "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1", Symbol: "WETH", Name: "Wrapped Ether", CoinGeckoID: "weth"},
	},
	10: { // Optimism
		{Address: "0x7F5c764cBc14f9669B88837ca1490cCa17c31607", Symbol: "USDC", Name: "USD Coin", CoinGeckoID: "usd-coin"},
		{Address: "0x94b008aA00579c1307B0EF2c499aD98a8ce58e58", Symbol: "USDT", Name: "Tether USD", CoinGeckoID: "tether"},
		{Address: "0x4200000000000000000000000000000000000006", Symbol: "WETH", Name: "Wrapped Ether", CoinGeckoID: "weth"},
	},
}

// CoinGecko IDs for native asset pricing by chain.
// Arbitrum and Optimism native gas is ETH -> "ethereum"; Polygon native is MATIC -> "polygon".
var NativeCoinGeckoIDByChain = map[int]string{
	1:     "ethereum",
	137:   "polygon",
	42161: "ethereum",
	10:    "ethereum",
}
