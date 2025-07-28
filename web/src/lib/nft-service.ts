import { createPublicClient, http, type Address } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum NFTStandard {
  ERC721 = 'ERC721',
  ERC1155 = 'ERC1155'
}

export enum NFTRarity {
  COMMON = 'common',
  UNCOMMON = 'uncommon',
  RARE = 'rare',
  EPIC = 'epic',
  LEGENDARY = 'legendary'
}

export interface NFTAttribute {
  trait_type: string
  value: string | number
  display_type?: string
  max_value?: number
  rarity_score?: number
}

export interface NFTMetadata {
  name: string
  description: string
  image: string
  external_url?: string
  animation_url?: string
  attributes: NFTAttribute[]
  background_color?: string
}

export interface NFT {
  id: string
  tokenId: string
  contractAddress: Address
  chainId: number
  standard: NFTStandard
  metadata: NFTMetadata
  owner: Address
  creator?: Address
  mintedAt: number
  lastTransferAt: number
  currentPrice?: string
  floorPrice?: string
  lastSalePrice?: string
  rarity?: NFTRarity
  rarityRank?: number
  totalSupply?: number
  isListed: boolean
  listingPrice?: string
  marketplace?: string
}

export interface NFTCollection {
  id: string
  contractAddress: Address
  chainId: number
  name: string
  symbol: string
  description: string
  image: string
  bannerImage?: string
  website?: string
  discord?: string
  twitter?: string
  totalSupply: number
  ownersCount: number
  floorPrice: string
  volumeTraded: string
  volume24h: string
  volume7d: string
  volume30d: string
  marketCap: string
  averagePrice: string
  createdAt: number
  verified: boolean
  featured: boolean
  category: string
  royaltyFee: number
  creator: Address
}

export interface UserNFTPortfolio {
  userAddress: Address
  totalNFTs: number
  totalCollections: number
  totalValue: string
  topCollection: string
  nfts: NFT[]
  collections: NFTCollection[]
  recentActivity: NFTActivity[]
}

export interface NFTActivity {
  id: string
  type: 'mint' | 'transfer' | 'sale' | 'listing' | 'offer'
  nftId: string
  from: Address
  to: Address
  price?: string
  marketplace?: string
  timestamp: number
  transactionHash: string
}

export interface NFTValuation {
  nftId: string
  estimatedValue: string
  confidence: number
  lastUpdated: number
  factors: {
    floorPrice: string
    rarityScore: number
    recentSales: string[]
    marketTrend: 'up' | 'down' | 'stable'
  }
}

export class NFTService {
  private static instance: NFTService
  private clients: Map<number, any> = new Map()
  private collections: Map<string, NFTCollection> = new Map()
  private nfts: Map<string, NFT> = new Map()
  private userPortfolios: Map<string, UserNFTPortfolio> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeMockData()
  }

  static getInstance(): NFTService {
    if (!NFTService.instance) {
      NFTService.instance = new NFTService()
    }
    return NFTService.instance
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
          console.warn(`Failed to initialize NFT client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeMockData() {
    // Mock NFT Collections
    const boredApes: NFTCollection = {
      id: 'bored-ape-yacht-club',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      chainId: 1,
      name: 'Bored Ape Yacht Club',
      symbol: 'BAYC',
      description: 'A collection of 10,000 unique Bored Ape NFTs— unique digital collectibles living on the Ethereum blockchain.',
      image: 'https://i.seadn.io/gae/Ju9CkWtV-1Okvf45wo8UctR-M9He2PjILP0oOvxE89AyiPPGtrR3gysu1Zgy0hjd2xKIgjJJtWIc0ybj4Vd7wv8t3pxDGHoJBzDB?auto=format&w=384',
      bannerImage: 'https://i.seadn.io/gae/i5dYZRkVCUK97bfprQ3WXyrT9BnLSZtVKGJlKQ919uaUB0sxbngVCioaiyu9r6snqfi2aaTyIvv6DHm4m2R3y7hMajbsv14pSZK8mhs?auto=format&w=1920',
      website: 'https://boredapeyachtclub.com',
      discord: 'https://discord.gg/3P5K3dzgdB',
      twitter: 'https://twitter.com/BoredApeYC',
      totalSupply: 10000,
      ownersCount: 5432,
      floorPrice: '15.5',
      volumeTraded: '1250000',
      volume24h: '125.5',
      volume7d: '850.2',
      volume30d: '3200.8',
      marketCap: '155000',
      averagePrice: '18.2',
      createdAt: 1619827200000, // April 30, 2021
      verified: true,
      featured: true,
      category: 'PFP',
      royaltyFee: 2.5,
      creator: '0x6B175474E89094C44Da98b954EedeAC495271d0F' as Address
    }

    const cryptoPunks: NFTCollection = {
      id: 'cryptopunks',
      contractAddress: '0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB' as Address,
      chainId: 1,
      name: 'CryptoPunks',
      symbol: 'PUNK',
      description: '10,000 unique collectible characters with proof of ownership stored on the Ethereum blockchain.',
      image: 'https://i.seadn.io/gae/BdxvLseXcfl57BiuQcQYdJ64v-aI8din7WPk0Pgo3qQFhAUH-B6i-dCqqc_mCkRIzULmwzwecnohLhrcH8A9mpWIZqA7ygc52Sr81hE?auto=format&w=384',
      totalSupply: 10000,
      ownersCount: 3281,
      floorPrice: '65.8',
      volumeTraded: '2800000',
      volume24h: '85.2',
      volume7d: '520.1',
      volume30d: '2100.5',
      marketCap: '658000',
      averagePrice: '75.5',
      createdAt: 1498867200000, // June 30, 2017
      verified: true,
      featured: true,
      category: 'PFP',
      royaltyFee: 0,
      creator: '0x6B175474E89094C44Da98b954EedeAC495271d0F' as Address
    }

    const azuki: NFTCollection = {
      id: 'azuki',
      contractAddress: '0xED5AF388653567Af2F388E6224dC7C4b3241C544' as Address,
      chainId: 1,
      name: 'Azuki',
      symbol: 'AZUKI',
      description: 'A collection of 10,000 avatars that give you membership access to The Garden.',
      image: 'https://i.seadn.io/gae/H8jOCJuQokNqGBpkBN5wk1oZwO7LM8bNnrHCaekV2nKjnCqw6UB5oaH8XyNeBDj6bA_n1mjejzhFQUP3O1NfjFLHr3FOaeHcTOOT?auto=format&w=384',
      totalSupply: 10000,
      ownersCount: 4521,
      floorPrice: '8.2',
      volumeTraded: '450000',
      volume24h: '45.8',
      volume7d: '280.5',
      volume30d: '1200.2',
      marketCap: '82000',
      averagePrice: '12.1',
      createdAt: 1642636800000, // January 20, 2022
      verified: true,
      featured: true,
      category: 'PFP',
      royaltyFee: 5.0,
      creator: '0x6B175474E89094C44Da98b954EedeAC495271d0F' as Address
    }

    this.collections.set(boredApes.id, boredApes)
    this.collections.set(cryptoPunks.id, cryptoPunks)
    this.collections.set(azuki.id, azuki)

    // Mock NFTs
    const mockNFT1: NFT = {
      id: 'bayc-1234',
      tokenId: '1234',
      contractAddress: boredApes.contractAddress,
      chainId: 1,
      standard: NFTStandard.ERC721,
      metadata: {
        name: 'Bored Ape #1234',
        description: 'A unique Bored Ape with rare traits',
        image: 'https://i.seadn.io/gae/example-ape.png',
        attributes: [
          { trait_type: 'Background', value: 'Blue', rarity_score: 15.2 },
          { trait_type: 'Fur', value: 'Golden Brown', rarity_score: 8.5 },
          { trait_type: 'Eyes', value: 'Laser Eyes', rarity_score: 2.1 },
          { trait_type: 'Mouth', value: 'Bored', rarity_score: 45.8 },
          { trait_type: 'Hat', value: 'Crown', rarity_score: 1.2 }
        ]
      },
      owner: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      creator: boredApes.creator,
      mintedAt: 1619827200000,
      lastTransferAt: Date.now() - 86400000 * 30,
      currentPrice: '18.5',
      floorPrice: '15.5',
      lastSalePrice: '16.8',
      rarity: NFTRarity.RARE,
      rarityRank: 245,
      totalSupply: 10000,
      isListed: false
    }

    const mockNFT2: NFT = {
      id: 'azuki-5678',
      tokenId: '5678',
      contractAddress: azuki.contractAddress,
      chainId: 1,
      standard: NFTStandard.ERC721,
      metadata: {
        name: 'Azuki #5678',
        description: 'A beautiful Azuki with unique traits',
        image: 'https://i.seadn.io/gae/example-azuki.png',
        attributes: [
          { trait_type: 'Type', value: 'Human', rarity_score: 65.2 },
          { trait_type: 'Hair', value: 'Pink Hairband', rarity_score: 12.5 },
          { trait_type: 'Eyes', value: 'Closed', rarity_score: 18.1 },
          { trait_type: 'Mouth', value: 'Frown', rarity_score: 25.8 },
          { trait_type: 'Clothing', value: 'Red Jacket', rarity_score: 8.2 }
        ]
      },
      owner: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      creator: azuki.creator,
      mintedAt: 1642636800000,
      lastTransferAt: Date.now() - 86400000 * 15,
      currentPrice: '10.2',
      floorPrice: '8.2',
      lastSalePrice: '9.5',
      rarity: NFTRarity.UNCOMMON,
      rarityRank: 1250,
      totalSupply: 10000,
      isListed: true,
      listingPrice: '10.2',
      marketplace: 'OpenSea'
    }

    this.nfts.set(mockNFT1.id, mockNFT1)
    this.nfts.set(mockNFT2.id, mockNFT2)

    // Mock user portfolio
    const mockPortfolio: UserNFTPortfolio = {
      userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      totalNFTs: 2,
      totalCollections: 2,
      totalValue: '28.7',
      topCollection: 'Bored Ape Yacht Club',
      nfts: [mockNFT1, mockNFT2],
      collections: [boredApes, azuki],
      recentActivity: [
        {
          id: 'activity-1',
          type: 'transfer',
          nftId: mockNFT2.id,
          from: '0x1234567890123456789012345678901234567890' as Address,
          to: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
          timestamp: Date.now() - 86400000 * 15,
          transactionHash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890'
        },
        {
          id: 'activity-2',
          type: 'sale',
          nftId: mockNFT1.id,
          from: '0x0987654321098765432109876543210987654321' as Address,
          to: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
          price: '16.8',
          marketplace: 'OpenSea',
          timestamp: Date.now() - 86400000 * 30,
          transactionHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef'
        }
      ]
    }

    this.userPortfolios.set(mockPortfolio.userAddress.toLowerCase(), mockPortfolio)
  }

  // Public methods
  getAllCollections(chainId?: number): NFTCollection[] {
    const collections = Array.from(this.collections.values())
    return chainId ? collections.filter(collection => collection.chainId === chainId) : collections
  }

  getCollection(collectionId: string): NFTCollection | undefined {
    return this.collections.get(collectionId)
  }

  getCollectionByAddress(contractAddress: Address, chainId: number): NFTCollection | undefined {
    return Array.from(this.collections.values()).find(
      collection =>
        collection.contractAddress.toLowerCase() === contractAddress.toLowerCase() &&
        collection.chainId === chainId
    )
  }

  getUserPortfolio(userAddress: Address): UserNFTPortfolio | null {
    return this.userPortfolios.get(userAddress.toLowerCase()) || null
  }

  getUserNFTs(userAddress: Address, chainId?: number): NFT[] {
    const portfolio = this.getUserPortfolio(userAddress)
    if (!portfolio) return []

    return chainId
      ? portfolio.nfts.filter(nft => nft.chainId === chainId)
      : portfolio.nfts
  }

  getNFT(nftId: string): NFT | undefined {
    return this.nfts.get(nftId)
  }

  async getNFTValuation(nftId: string): Promise<NFTValuation | null> {
    const nft = this.getNFT(nftId)
    if (!nft) return null

    // Mock valuation calculation
    const basePrice = parseFloat(nft.floorPrice || '0')
    const rarityMultiplier = nft.rarityRank ? (1 + (10000 - nft.rarityRank) / 10000) : 1
    const estimatedValue = (basePrice * rarityMultiplier).toFixed(2)

    return {
      nftId,
      estimatedValue,
      confidence: 85,
      lastUpdated: Date.now(),
      factors: {
        floorPrice: nft.floorPrice || '0',
        rarityScore: nft.rarityRank || 0,
        recentSales: ['16.8', '17.2', '15.9'],
        marketTrend: 'stable'
      }
    }
  }

  getTopCollections(limit: number = 10): NFTCollection[] {
    return Array.from(this.collections.values())
      .sort((a, b) => parseFloat(b.volumeTraded) - parseFloat(a.volumeTraded))
      .slice(0, limit)
  }

  getTrendingCollections(limit: number = 10): NFTCollection[] {
    return Array.from(this.collections.values())
      .sort((a, b) => parseFloat(b.volume24h) - parseFloat(a.volume24h))
      .slice(0, limit)
  }

  searchCollections(query: string): NFTCollection[] {
    const lowercaseQuery = query.toLowerCase()
    return Array.from(this.collections.values()).filter(collection =>
      collection.name.toLowerCase().includes(lowercaseQuery) ||
      collection.symbol.toLowerCase().includes(lowercaseQuery) ||
      collection.description.toLowerCase().includes(lowercaseQuery)
    )
  }

  getCollectionsByCategory(category: string): NFTCollection[] {
    return Array.from(this.collections.values()).filter(
      collection => collection.category.toLowerCase() === category.toLowerCase()
    )
  }

  calculatePortfolioValue(userAddress: Address): {
    totalValue: number
    topNFT: NFT | null
    averageValue: number
    totalGainLoss: number
  } {
    const portfolio = this.getUserPortfolio(userAddress)
    if (!portfolio) {
      return { totalValue: 0, topNFT: null, averageValue: 0, totalGainLoss: 0 }
    }

    const totalValue = portfolio.nfts.reduce((sum, nft) => {
      return sum + parseFloat(nft.currentPrice || nft.floorPrice || '0')
    }, 0)

    const topNFT = portfolio.nfts.reduce((top, nft) => {
      const nftValue = parseFloat(nft.currentPrice || nft.floorPrice || '0')
      const topValue = parseFloat(top?.currentPrice || top?.floorPrice || '0')
      return nftValue > topValue ? nft : top
    }, portfolio.nfts[0] || null)

    const averageValue = portfolio.nfts.length > 0 ? totalValue / portfolio.nfts.length : 0

    // Mock gain/loss calculation
    const totalGainLoss = totalValue * 0.15 // Assume 15% gain

    return {
      totalValue,
      topNFT,
      averageValue,
      totalGainLoss
    }
  }
}

// Export singleton instance
export const nftService = NFTService.getInstance()