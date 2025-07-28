import { createPublicClient, http, type Address, type Hash } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum TransactionType {
  SEND = 'send',
  RECEIVE = 'receive',
  SWAP = 'swap',
  APPROVE = 'approve',
  MINT = 'mint',
  BURN = 'burn',
  STAKE = 'stake',
  UNSTAKE = 'unstake',
  CLAIM = 'claim',
  DEPOSIT = 'deposit',
  WITHDRAW = 'withdraw',
  NFT_TRANSFER = 'nft_transfer',
  NFT_MINT = 'nft_mint',
  NFT_SALE = 'nft_sale',
  CONTRACT_INTERACTION = 'contract_interaction',
  BRIDGE = 'bridge',
  UNKNOWN = 'unknown'
}

export enum TransactionStatus {
  PENDING = 'pending',
  CONFIRMED = 'confirmed',
  FAILED = 'failed',
  DROPPED = 'dropped',
  REPLACED = 'replaced'
}

export enum TransactionCategory {
  DEFI = 'defi',
  NFT = 'nft',
  TOKEN = 'token',
  ETH = 'eth',
  CONTRACT = 'contract',
  BRIDGE = 'bridge'
}

export interface TokenTransfer {
  tokenAddress: Address
  tokenSymbol: string
  tokenName: string
  tokenDecimals: number
  from: Address
  to: Address
  value: string
  valueFormatted: string
  logoURI?: string
}

export interface NFTTransfer {
  contractAddress: Address
  tokenId: string
  from: Address
  to: Address
  tokenName?: string
  tokenSymbol?: string
  metadata?: {
    name: string
    description: string
    image: string
  }
}

export interface TransactionDetails {
  hash: Hash
  blockNumber: bigint
  blockHash: string
  transactionIndex: number
  from: Address
  to: Address | null
  value: string
  gasPrice: string
  gasLimit: string
  gasUsed: string
  gasEfficiency: number // gasUsed / gasLimit * 100
  nonce: number
  input: string
  status: TransactionStatus
  timestamp: number
  confirmations: number
  type: TransactionType
  category: TransactionCategory
  chainId: number
  
  // Enhanced data
  valueUSD?: string
  gasCostETH: string
  gasCostUSD?: string
  tokenTransfers: TokenTransfer[]
  nftTransfers: NFTTransfer[]
  
  // Contract interaction details
  contractAddress?: Address
  contractName?: string
  methodName?: string
  methodSignature?: string
  decodedInput?: any
  
  // DeFi specific
  protocol?: string
  protocolAction?: string
  
  // Error details
  errorReason?: string
  revertReason?: string
  
  // Metadata
  description?: string
  tags: string[]
  isInternal: boolean
  replacedBy?: Hash
  replacementFor?: Hash
}

export interface TransactionSearchFilters {
  address?: Address
  chainId?: number
  type?: TransactionType[]
  category?: TransactionCategory[]
  status?: TransactionStatus[]
  fromDate?: Date
  toDate?: Date
  minValue?: string
  maxValue?: string
  tokenAddress?: Address
  contractAddress?: Address
  protocol?: string
  tags?: string[]
  searchQuery?: string
}

export interface TransactionSearchResult {
  transactions: TransactionDetails[]
  totalCount: number
  hasMore: boolean
  nextCursor?: string
}

export interface TransactionSummary {
  totalTransactions: number
  totalValueETH: string
  totalValueUSD: string
  totalGasCostETH: string
  totalGasCostUSD: string
  successRate: number
  averageGasPrice: string
  averageConfirmationTime: number
  typeDistribution: Record<TransactionType, number>
  categoryDistribution: Record<TransactionCategory, number>
  monthlyActivity: {
    month: string
    count: number
    volume: string
  }[]
}

export class TransactionHistoryService {
  private static instance: TransactionHistoryService
  private clients: Map<number, any> = new Map()
  private transactionCache: Map<string, TransactionDetails> = new Map()
  private searchIndex: Map<string, Set<string>> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeMockData()
  }

  static getInstance(): TransactionHistoryService {
    if (!TransactionHistoryService.instance) {
      TransactionHistoryService.instance = new TransactionHistoryService()
    }
    return TransactionHistoryService.instance
  }

  private initializeClients() {
    Object.values(SUPPORTED_CHAINS).forEach(chain => {
      if (!chain.isTestnet || chain.id === 11155111) {
        try {
          const client = createPublicClient({
            chain: {
              id: chain.id,
              name: chain.name,
              network: chain.shortName.toLowerCase(),
              nativeCurrency: chain.nativeCurrency,
              rpcUrls: chain.rpcUrls
            } as any,
            transport: http()
          })
          this.clients.set(chain.id, client)
        } catch (error) {
          console.warn(`Failed to initialize transaction history client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeMockData() {
    // Mock transaction data
    const mockTransactions: TransactionDetails[] = [
      {
        hash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890' as Hash,
        blockNumber: BigInt(18500000),
        blockHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
        transactionIndex: 45,
        from: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
        to: '0xA0b86a33E6441b8dB4C6b4b8b4C6b4b8b4C6b4b8' as Address,
        value: '1500000000000000000', // 1.5 ETH
        gasPrice: '20000000000', // 20 gwei
        gasLimit: '21000',
        gasUsed: '21000',
        gasEfficiency: 100,
        nonce: 42,
        input: '0x',
        status: TransactionStatus.CONFIRMED,
        timestamp: Date.now() - 86400000 * 2, // 2 days ago
        confirmations: 1250,
        type: TransactionType.SEND,
        category: TransactionCategory.ETH,
        chainId: 1,
        valueUSD: '3750.00',
        gasCostETH: '0.00042',
        gasCostUSD: '1.05',
        tokenTransfers: [],
        nftTransfers: [],
        description: 'ETH transfer to wallet',
        tags: ['transfer', 'eth'],
        isInternal: false
      },
      {
        hash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef' as Hash,
        blockNumber: BigInt(18499500),
        blockHash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890',
        transactionIndex: 123,
        from: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
        to: '0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D' as Address, // Uniswap V2 Router
        value: '0',
        gasPrice: '25000000000', // 25 gwei
        gasLimit: '150000',
        gasUsed: '142350',
        gasEfficiency: 94.9,
        nonce: 41,
        input: '0x38ed1739000000000000000000000000000000000000000000000000016345785d8a0000',
        status: TransactionStatus.CONFIRMED,
        timestamp: Date.now() - 86400000 * 3, // 3 days ago
        confirmations: 1750,
        type: TransactionType.SWAP,
        category: TransactionCategory.DEFI,
        chainId: 1,
        gasCostETH: '0.0035588',
        gasCostUSD: '8.90',
        tokenTransfers: [
          {
            tokenAddress: '0xA0b86a33E6441b8dB4C6b4b8b4C6b4b8b4C6b4b8' as Address,
            tokenSymbol: 'USDC',
            tokenName: 'USD Coin',
            tokenDecimals: 6,
            from: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
            to: '0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D' as Address,
            value: '1000000000', // 1000 USDC
            valueFormatted: '1000.00'
          }
        ],
        nftTransfers: [],
        contractAddress: '0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D' as Address,
        contractName: 'Uniswap V2 Router',
        methodName: 'swapExactTokensForETH',
        protocol: 'Uniswap V2',
        protocolAction: 'Swap USDC for ETH',
        description: 'Swapped 1000 USDC for ETH on Uniswap',
        tags: ['swap', 'uniswap', 'defi', 'usdc', 'eth'],
        isInternal: false
      },
      {
        hash: '0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba' as Hash,
        blockNumber: BigInt(18498000),
        blockHash: '0xfedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210',
        transactionIndex: 67,
        from: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
        to: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address, // BAYC
        value: '0',
        gasPrice: '30000000000', // 30 gwei
        gasLimit: '100000',
        gasUsed: '85420',
        gasEfficiency: 85.4,
        nonce: 40,
        input: '0xa22cb46500000000000000000000000074250d5630b4cf539739df2c5dacb4c659f2488d0000000000000000000000000000000000000000000000000000000000000001',
        status: TransactionStatus.CONFIRMED,
        timestamp: Date.now() - 86400000 * 5, // 5 days ago
        confirmations: 2500,
        type: TransactionType.NFT_TRANSFER,
        category: TransactionCategory.NFT,
        chainId: 1,
        gasCostETH: '0.002563',
        gasCostUSD: '6.41',
        tokenTransfers: [],
        nftTransfers: [
          {
            contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
            tokenId: '1234',
            from: '0x0000000000000000000000000000000000000000' as Address,
            to: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
            tokenName: 'Bored Ape Yacht Club',
            tokenSymbol: 'BAYC',
            metadata: {
              name: 'Bored Ape #1234',
              description: 'A unique Bored Ape NFT',
              image: 'https://ipfs.io/ipfs/QmYourImageHash'
            }
          }
        ],
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
        contractName: 'Bored Ape Yacht Club',
        methodName: 'setApprovalForAll',
        description: 'Set approval for all BAYC tokens',
        tags: ['nft', 'bayc', 'approval'],
        isInternal: false
      },
      {
        hash: '0x5555666677778888999900001111222233334444555566667777888899990000' as Hash,
        blockNumber: BigInt(18497000),
        blockHash: '0x0000111122223333444455556666777788889999aaaabbbbccccddddeeeeffff',
        transactionIndex: 89,
        from: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
        to: '0x7d2768dE32b0b80b7a3454c06BdAc94A69DDc7A9' as Address, // Aave V2 LendingPool
        value: '0',
        gasPrice: '18000000000', // 18 gwei
        gasLimit: '200000',
        gasUsed: '185420',
        gasEfficiency: 92.7,
        nonce: 39,
        input: '0xe8eda9df000000000000000000000000a0b86a33e6441b8db4c6b4b8b4c6b4b8b4c6b4b8',
        status: TransactionStatus.CONFIRMED,
        timestamp: Date.now() - 86400000 * 7, // 7 days ago
        confirmations: 3500,
        type: TransactionType.DEPOSIT,
        category: TransactionCategory.DEFI,
        chainId: 1,
        gasCostETH: '0.0033376',
        gasCostUSD: '8.34',
        tokenTransfers: [
          {
            tokenAddress: '0xA0b86a33E6441b8dB4C6b4b8b4C6b4b8b4C6b4b8' as Address,
            tokenSymbol: 'USDC',
            tokenName: 'USD Coin',
            tokenDecimals: 6,
            from: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
            to: '0x7d2768dE32b0b80b7a3454c06BdAc94A69DDc7A9' as Address,
            value: '5000000000', // 5000 USDC
            valueFormatted: '5000.00'
          }
        ],
        nftTransfers: [],
        contractAddress: '0x7d2768dE32b0b80b7a3454c06BdAc94A69DDc7A9' as Address,
        contractName: 'Aave V2 LendingPool',
        methodName: 'deposit',
        protocol: 'Aave V2',
        protocolAction: 'Deposit USDC',
        description: 'Deposited 5000 USDC to Aave V2',
        tags: ['deposit', 'aave', 'defi', 'usdc', 'lending'],
        isInternal: false
      }
    ]

    // Cache transactions and build search index
    mockTransactions.forEach(tx => {
      this.transactionCache.set(tx.hash, tx)
      this.indexTransaction(tx)
    })
  }

  private indexTransaction(tx: TransactionDetails) {
    // Index by various fields for search
    const indexFields = [
      tx.hash.toLowerCase(),
      tx.from.toLowerCase(),
      tx.to?.toLowerCase() || '',
      tx.contractAddress?.toLowerCase() || '',
      tx.contractName?.toLowerCase() || '',
      tx.methodName?.toLowerCase() || '',
      tx.protocol?.toLowerCase() || '',
      tx.description?.toLowerCase() || '',
      ...tx.tags.map(tag => tag.toLowerCase()),
      ...tx.tokenTransfers.map(transfer => transfer.tokenSymbol.toLowerCase()),
      ...tx.nftTransfers.map(transfer => transfer.tokenSymbol?.toLowerCase() || '')
    ].filter(Boolean)

    indexFields.forEach(field => {
      if (!this.searchIndex.has(field)) {
        this.searchIndex.set(field, new Set())
      }
      this.searchIndex.get(field)!.add(tx.hash)
    })
  }

  // Public methods
  async getTransactionHistory(
    address: Address,
    chainId: number = 1,
    limit: number = 50,
    cursor?: string
  ): Promise<TransactionSearchResult> {
    try {
      // Mock implementation - in real app, this would query blockchain APIs
      const allTransactions = Array.from(this.transactionCache.values())
        .filter(tx => 
          tx.chainId === chainId && 
          (tx.from.toLowerCase() === address.toLowerCase() || 
           tx.to?.toLowerCase() === address.toLowerCase())
        )
        .sort((a, b) => b.timestamp - a.timestamp)

      const startIndex = cursor ? parseInt(cursor) : 0
      const endIndex = Math.min(startIndex + limit, allTransactions.length)
      const transactions = allTransactions.slice(startIndex, endIndex)

      return {
        transactions,
        totalCount: allTransactions.length,
        hasMore: endIndex < allTransactions.length,
        nextCursor: endIndex < allTransactions.length ? endIndex.toString() : undefined
      }
    } catch (error) {
      throw new Error(`Failed to get transaction history: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  async searchTransactions(
    filters: TransactionSearchFilters,
    limit: number = 50,
    cursor?: string
  ): Promise<TransactionSearchResult> {
    try {
      let transactions = Array.from(this.transactionCache.values())

      // Apply filters
      if (filters.address) {
        const address = filters.address.toLowerCase()
        transactions = transactions.filter(tx => 
          tx.from.toLowerCase() === address || 
          tx.to?.toLowerCase() === address
        )
      }

      if (filters.chainId) {
        transactions = transactions.filter(tx => tx.chainId === filters.chainId)
      }

      if (filters.type && filters.type.length > 0) {
        transactions = transactions.filter(tx => filters.type!.includes(tx.type))
      }

      if (filters.category && filters.category.length > 0) {
        transactions = transactions.filter(tx => filters.category!.includes(tx.category))
      }

      if (filters.status && filters.status.length > 0) {
        transactions = transactions.filter(tx => filters.status!.includes(tx.status))
      }

      if (filters.fromDate) {
        transactions = transactions.filter(tx => tx.timestamp >= filters.fromDate!.getTime())
      }

      if (filters.toDate) {
        transactions = transactions.filter(tx => tx.timestamp <= filters.toDate!.getTime())
      }

      if (filters.minValue) {
        const minValue = parseFloat(filters.minValue)
        transactions = transactions.filter(tx => parseFloat(tx.value) >= minValue)
      }

      if (filters.maxValue) {
        const maxValue = parseFloat(filters.maxValue)
        transactions = transactions.filter(tx => parseFloat(tx.value) <= maxValue)
      }

      if (filters.tokenAddress) {
        const tokenAddress = filters.tokenAddress.toLowerCase()
        transactions = transactions.filter(tx => 
          tx.tokenTransfers.some(transfer => 
            transfer.tokenAddress.toLowerCase() === tokenAddress
          )
        )
      }

      if (filters.contractAddress) {
        const contractAddress = filters.contractAddress.toLowerCase()
        transactions = transactions.filter(tx => 
          tx.contractAddress?.toLowerCase() === contractAddress
        )
      }

      if (filters.protocol) {
        const protocol = filters.protocol.toLowerCase()
        transactions = transactions.filter(tx => 
          tx.protocol?.toLowerCase().includes(protocol)
        )
      }

      if (filters.tags && filters.tags.length > 0) {
        transactions = transactions.filter(tx => 
          filters.tags!.some(tag => 
            tx.tags.some(txTag => 
              txTag.toLowerCase().includes(tag.toLowerCase())
            )
          )
        )
      }

      if (filters.searchQuery) {
        const query = filters.searchQuery.toLowerCase()
        transactions = transactions.filter(tx => 
          tx.hash.toLowerCase().includes(query) ||
          tx.from.toLowerCase().includes(query) ||
          tx.to?.toLowerCase().includes(query) ||
          tx.contractName?.toLowerCase().includes(query) ||
          tx.methodName?.toLowerCase().includes(query) ||
          tx.protocol?.toLowerCase().includes(query) ||
          tx.description?.toLowerCase().includes(query) ||
          tx.tags.some(tag => tag.toLowerCase().includes(query)) ||
          tx.tokenTransfers.some(transfer => 
            transfer.tokenSymbol.toLowerCase().includes(query) ||
            transfer.tokenName.toLowerCase().includes(query)
          )
        )
      }

      // Sort by timestamp (newest first)
      transactions.sort((a, b) => b.timestamp - a.timestamp)

      const startIndex = cursor ? parseInt(cursor) : 0
      const endIndex = Math.min(startIndex + limit, transactions.length)
      const paginatedTransactions = transactions.slice(startIndex, endIndex)

      return {
        transactions: paginatedTransactions,
        totalCount: transactions.length,
        hasMore: endIndex < transactions.length,
        nextCursor: endIndex < transactions.length ? endIndex.toString() : undefined
      }
    } catch (error) {
      throw new Error(`Failed to search transactions: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  async getTransaction(hash: Hash): Promise<TransactionDetails | null> {
    return this.transactionCache.get(hash) || null
  }

  async getTransactionSummary(
    address: Address,
    chainId: number = 1,
    fromDate?: Date,
    toDate?: Date
  ): Promise<TransactionSummary> {
    try {
      let transactions = Array.from(this.transactionCache.values())
        .filter(tx => 
          tx.chainId === chainId && 
          (tx.from.toLowerCase() === address.toLowerCase() || 
           tx.to?.toLowerCase() === address.toLowerCase())
        )

      if (fromDate) {
        transactions = transactions.filter(tx => tx.timestamp >= fromDate.getTime())
      }

      if (toDate) {
        transactions = transactions.filter(tx => tx.timestamp <= toDate.getTime())
      }

      const totalTransactions = transactions.length
      const confirmedTransactions = transactions.filter(tx => tx.status === TransactionStatus.CONFIRMED)
      
      const totalValueETH = transactions
        .reduce((sum, tx) => sum + parseFloat(tx.value), 0)
        .toString()

      const totalValueUSD = transactions
        .reduce((sum, tx) => sum + parseFloat(tx.valueUSD || '0'), 0)
        .toString()

      const totalGasCostETH = transactions
        .reduce((sum, tx) => sum + parseFloat(tx.gasCostETH), 0)
        .toString()

      const totalGasCostUSD = transactions
        .reduce((sum, tx) => sum + parseFloat(tx.gasCostUSD || '0'), 0)
        .toString()

      const successRate = totalTransactions > 0 
        ? (confirmedTransactions.length / totalTransactions) * 100 
        : 0

      const averageGasPrice = transactions.length > 0
        ? (transactions.reduce((sum, tx) => sum + parseFloat(tx.gasPrice), 0) / transactions.length).toString()
        : '0'

      const averageConfirmationTime = 15 // Mock average confirmation time in seconds

      // Type distribution
      const typeDistribution: Record<TransactionType, number> = {} as any
      Object.values(TransactionType).forEach(type => {
        typeDistribution[type] = transactions.filter(tx => tx.type === type).length
      })

      // Category distribution
      const categoryDistribution: Record<TransactionCategory, number> = {} as any
      Object.values(TransactionCategory).forEach(category => {
        categoryDistribution[category] = transactions.filter(tx => tx.category === category).length
      })

      // Monthly activity (mock data)
      const monthlyActivity = [
        { month: '2024-01', count: 45, volume: '12.5' },
        { month: '2024-02', count: 38, volume: '8.2' },
        { month: '2024-03', count: 52, volume: '15.8' },
        { month: '2024-04', count: 41, volume: '11.3' }
      ]

      return {
        totalTransactions,
        totalValueETH,
        totalValueUSD,
        totalGasCostETH,
        totalGasCostUSD,
        successRate,
        averageGasPrice,
        averageConfirmationTime,
        typeDistribution,
        categoryDistribution,
        monthlyActivity
      }
    } catch (error) {
      throw new Error(`Failed to get transaction summary: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  // Utility methods
  formatTransactionValue(value: string, decimals: number = 18): string {
    const valueNum = parseFloat(value) / Math.pow(10, decimals)
    if (valueNum >= 1000000) return `${(valueNum / 1000000).toFixed(2)}M`
    if (valueNum >= 1000) return `${(valueNum / 1000).toFixed(2)}K`
    return valueNum.toFixed(4)
  }

  formatGasPrice(gasPrice: string): string {
    const gasPriceGwei = parseFloat(gasPrice) / 1e9
    return `${gasPriceGwei.toFixed(2)} gwei`
  }

  getTransactionTypeIcon(type: TransactionType): string {
    switch (type) {
      case TransactionType.SEND:
        return '↗️'
      case TransactionType.RECEIVE:
        return '↙️'
      case TransactionType.SWAP:
        return '🔄'
      case TransactionType.APPROVE:
        return '✅'
      case TransactionType.MINT:
        return '🪙'
      case TransactionType.BURN:
        return '🔥'
      case TransactionType.STAKE:
        return '🔒'
      case TransactionType.UNSTAKE:
        return '🔓'
      case TransactionType.CLAIM:
        return '🎁'
      case TransactionType.DEPOSIT:
        return '📥'
      case TransactionType.WITHDRAW:
        return '📤'
      case TransactionType.NFT_TRANSFER:
        return '🖼️'
      case TransactionType.NFT_MINT:
        return '🎨'
      case TransactionType.NFT_SALE:
        return '💰'
      case TransactionType.CONTRACT_INTERACTION:
        return '⚙️'
      case TransactionType.BRIDGE:
        return '🌉'
      default:
        return '❓'
    }
  }

  getTransactionStatusColor(status: TransactionStatus): string {
    switch (status) {
      case TransactionStatus.CONFIRMED:
        return 'text-green-600'
      case TransactionStatus.PENDING:
        return 'text-yellow-600'
      case TransactionStatus.FAILED:
        return 'text-red-600'
      case TransactionStatus.DROPPED:
        return 'text-gray-600'
      case TransactionStatus.REPLACED:
        return 'text-blue-600'
      default:
        return 'text-gray-600'
    }
  }

  getCategoryColor(category: TransactionCategory): string {
    switch (category) {
      case TransactionCategory.DEFI:
        return 'bg-purple-100 text-purple-800'
      case TransactionCategory.NFT:
        return 'bg-pink-100 text-pink-800'
      case TransactionCategory.TOKEN:
        return 'bg-blue-100 text-blue-800'
      case TransactionCategory.ETH:
        return 'bg-gray-100 text-gray-800'
      case TransactionCategory.CONTRACT:
        return 'bg-orange-100 text-orange-800'
      case TransactionCategory.BRIDGE:
        return 'bg-green-100 text-green-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }
}

// Export singleton instance
export const transactionHistoryService = TransactionHistoryService.getInstance()
