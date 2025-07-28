import { createPublicClient, createWalletClient, http, type Address, type Hash } from 'viem'
import { privateKeyToAccount } from 'viem/accounts'
import { SUPPORTED_CHAINS } from './chains'

export enum NFTStandard {
  ERC721 = 'ERC721',
  ERC1155 = 'ERC1155'
}

export enum MintingStatus {
  IDLE = 'idle',
  UPLOADING_METADATA = 'uploading_metadata',
  DEPLOYING_CONTRACT = 'deploying_contract',
  MINTING = 'minting',
  COMPLETED = 'completed',
  FAILED = 'failed'
}

export interface NFTMetadata {
  name: string
  description: string
  image: string
  external_url?: string
  animation_url?: string
  attributes: {
    trait_type: string
    value: string | number
    display_type?: 'boost_number' | 'boost_percentage' | 'number' | 'date'
  }[]
  background_color?: string
  youtube_url?: string
}

export interface CollectionMetadata {
  name: string
  description: string
  image: string
  external_link?: string
  seller_fee_basis_points: number
  fee_recipient: Address
}

export interface IPFSUploadResult {
  hash: string
  url: string
  gateway_url: string
}

export interface ContractDeployment {
  contractAddress: Address
  transactionHash: Hash
  blockNumber: bigint
  gasUsed: bigint
  deploymentCost: string
}

export interface MintResult {
  tokenId: string
  transactionHash: Hash
  contractAddress: Address
  metadataUri: string
  gasUsed: bigint
  mintingCost: string
}

export interface NFTCollection {
  id: string
  name: string
  symbol: string
  description: string
  contractAddress: Address
  standard: NFTStandard
  chainId: number
  owner: Address
  totalSupply: number
  maxSupply?: number
  baseTokenURI: string
  royaltyPercentage: number
  royaltyRecipient: Address
  isRevealed: boolean
  createdAt: number
  deploymentTx: Hash
}

export interface MintingProgress {
  status: MintingStatus
  currentStep: number
  totalSteps: number
  message: string
  transactionHash?: Hash
  ipfsHash?: string
  contractAddress?: Address
  tokenId?: string
  error?: string
}

export class NFTMintingService {
  private static instance: NFTMintingService
  private clients: Map<number, any> = new Map()
  private collections: Map<string, NFTCollection> = new Map()
  private ipfsGateway = 'https://ipfs.io/ipfs/'
  private pinataApiKey = process.env.NEXT_PUBLIC_PINATA_API_KEY || ''
  private pinataSecretKey = process.env.NEXT_PUBLIC_PINATA_SECRET_KEY || ''

  private constructor() {
    this.initializeClients()
    this.initializeMockData()
  }

  static getInstance(): NFTMintingService {
    if (!NFTMintingService.instance) {
      NFTMintingService.instance = new NFTMintingService()
    }
    return NFTMintingService.instance
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
          console.warn(`Failed to initialize NFT minting client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeMockData() {
    // Mock NFT Collection
    const mockCollection: NFTCollection = {
      id: 'mock-collection-1',
      name: 'AI Art Collection',
      symbol: 'AIART',
      description: 'A collection of AI-generated artwork',
      contractAddress: '0x1234567890123456789012345678901234567890' as Address,
      standard: NFTStandard.ERC721,
      chainId: 1,
      owner: '0x0987654321098765432109876543210987654321' as Address,
      totalSupply: 100,
      maxSupply: 1000,
      baseTokenURI: 'https://ipfs.io/ipfs/QmYourCollectionHash/',
      royaltyPercentage: 5,
      royaltyRecipient: '0x0987654321098765432109876543210987654321' as Address,
      isRevealed: true,
      createdAt: Date.now() - 86400000 * 30,
      deploymentTx: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890' as Hash
    }

    this.collections.set(mockCollection.id, mockCollection)
  }

  // IPFS Methods
  async uploadToIPFS(file: File): Promise<IPFSUploadResult> {
    try {
      // Mock IPFS upload for demo
      const mockHash = `Qm${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}`
      
      // Simulate upload delay
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      return {
        hash: mockHash,
        url: `ipfs://${mockHash}`,
        gateway_url: `${this.ipfsGateway}${mockHash}`
      }
    } catch (error) {
      throw new Error(`Failed to upload to IPFS: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  async uploadMetadataToIPFS(metadata: NFTMetadata): Promise<IPFSUploadResult> {
    try {
      // Mock metadata upload
      const mockHash = `Qm${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}`
      
      // Simulate upload delay
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      return {
        hash: mockHash,
        url: `ipfs://${mockHash}`,
        gateway_url: `${this.ipfsGateway}${mockHash}`
      }
    } catch (error) {
      throw new Error(`Failed to upload metadata to IPFS: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  async uploadCollectionMetadata(metadata: CollectionMetadata): Promise<IPFSUploadResult> {
    try {
      // Mock collection metadata upload
      const mockHash = `Qm${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}`
      
      // Simulate upload delay
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      return {
        hash: mockHash,
        url: `ipfs://${mockHash}`,
        gateway_url: `${this.ipfsGateway}${mockHash}`
      }
    } catch (error) {
      throw new Error(`Failed to upload collection metadata to IPFS: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  // Contract Deployment Methods
  async deployNFTContract(
    name: string,
    symbol: string,
    baseTokenURI: string,
    maxSupply: number,
    royaltyPercentage: number,
    royaltyRecipient: Address,
    chainId: number,
    ownerAddress: Address,
    standard: NFTStandard = NFTStandard.ERC721
  ): Promise<ContractDeployment> {
    try {
      // Mock contract deployment
      const mockAddress = `0x${Math.random().toString(16).substring(2, 42).padStart(40, '0')}` as Address
      const mockTxHash = `0x${Math.random().toString(16).substring(2, 66).padStart(64, '0')}` as Hash
      
      // Simulate deployment delay
      await new Promise(resolve => setTimeout(resolve, 5000))
      
      return {
        contractAddress: mockAddress,
        transactionHash: mockTxHash,
        blockNumber: BigInt(Math.floor(Math.random() * 1000000) + 18000000),
        gasUsed: BigInt(Math.floor(Math.random() * 500000) + 1000000),
        deploymentCost: (Math.random() * 0.1 + 0.05).toFixed(4)
      }
    } catch (error) {
      throw new Error(`Failed to deploy contract: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  // Minting Methods
  async mintNFT(
    contractAddress: Address,
    to: Address,
    metadataUri: string,
    chainId: number
  ): Promise<MintResult> {
    try {
      // Mock NFT minting
      const mockTxHash = `0x${Math.random().toString(16).substring(2, 66).padStart(64, '0')}` as Hash
      const mockTokenId = Math.floor(Math.random() * 10000).toString()
      
      // Simulate minting delay
      await new Promise(resolve => setTimeout(resolve, 3000))
      
      return {
        tokenId: mockTokenId,
        transactionHash: mockTxHash,
        contractAddress,
        metadataUri,
        gasUsed: BigInt(Math.floor(Math.random() * 100000) + 50000),
        mintingCost: (Math.random() * 0.01 + 0.005).toFixed(4)
      }
    } catch (error) {
      throw new Error(`Failed to mint NFT: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  async batchMintNFTs(
    contractAddress: Address,
    recipients: Address[],
    metadataUris: string[],
    chainId: number
  ): Promise<MintResult[]> {
    try {
      if (recipients.length !== metadataUris.length) {
        throw new Error('Recipients and metadata URIs arrays must have the same length')
      }

      const results: MintResult[] = []
      
      // Mock batch minting
      for (let i = 0; i < recipients.length; i++) {
        const mockTxHash = `0x${Math.random().toString(16).substring(2, 66).padStart(64, '0')}` as Hash
        const mockTokenId = (Math.floor(Math.random() * 10000) + i).toString()
        
        results.push({
          tokenId: mockTokenId,
          transactionHash: mockTxHash,
          contractAddress,
          metadataUri: metadataUris[i],
          gasUsed: BigInt(Math.floor(Math.random() * 80000) + 40000),
          mintingCost: (Math.random() * 0.008 + 0.004).toFixed(4)
        })
      }
      
      // Simulate batch minting delay
      await new Promise(resolve => setTimeout(resolve, recipients.length * 1000))
      
      return results
    } catch (error) {
      throw new Error(`Failed to batch mint NFTs: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  // Collection Management
  async createCollection(
    name: string,
    symbol: string,
    description: string,
    image: File,
    maxSupply: number,
    royaltyPercentage: number,
    royaltyRecipient: Address,
    chainId: number,
    ownerAddress: Address,
    standard: NFTStandard = NFTStandard.ERC721
  ): Promise<{ collection: NFTCollection; deployment: ContractDeployment }> {
    try {
      // Upload collection image to IPFS
      const imageUpload = await this.uploadToIPFS(image)
      
      // Create collection metadata
      const collectionMetadata: CollectionMetadata = {
        name,
        description,
        image: imageUpload.url,
        seller_fee_basis_points: royaltyPercentage * 100,
        fee_recipient: royaltyRecipient
      }
      
      // Upload collection metadata to IPFS
      const metadataUpload = await this.uploadCollectionMetadata(collectionMetadata)
      
      // Deploy contract
      const deployment = await this.deployNFTContract(
        name,
        symbol,
        metadataUpload.url,
        maxSupply,
        royaltyPercentage,
        royaltyRecipient,
        chainId,
        ownerAddress,
        standard
      )
      
      // Create collection object
      const collection: NFTCollection = {
        id: `collection-${Date.now()}`,
        name,
        symbol,
        description,
        contractAddress: deployment.contractAddress,
        standard,
        chainId,
        owner: ownerAddress,
        totalSupply: 0,
        maxSupply,
        baseTokenURI: metadataUpload.url,
        royaltyPercentage,
        royaltyRecipient,
        isRevealed: false,
        createdAt: Date.now(),
        deploymentTx: deployment.transactionHash
      }
      
      // Store collection
      this.collections.set(collection.id, collection)
      
      return { collection, deployment }
    } catch (error) {
      throw new Error(`Failed to create collection: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  async getCollections(ownerAddress?: Address): Promise<NFTCollection[]> {
    let collections = Array.from(this.collections.values())
    
    if (ownerAddress) {
      collections = collections.filter(collection => 
        collection.owner.toLowerCase() === ownerAddress.toLowerCase()
      )
    }
    
    return collections.sort((a, b) => b.createdAt - a.createdAt)
  }

  async getCollection(id: string): Promise<NFTCollection | null> {
    return this.collections.get(id) || null
  }

  // Utility Methods
  validateMetadata(metadata: NFTMetadata): { isValid: boolean; errors: string[] } {
    const errors: string[] = []
    
    if (!metadata.name || metadata.name.trim().length === 0) {
      errors.push('Name is required')
    }
    
    if (!metadata.description || metadata.description.trim().length === 0) {
      errors.push('Description is required')
    }
    
    if (!metadata.image || metadata.image.trim().length === 0) {
      errors.push('Image is required')
    }
    
    if (metadata.attributes) {
      metadata.attributes.forEach((attr, index) => {
        if (!attr.trait_type || attr.trait_type.trim().length === 0) {
          errors.push(`Attribute ${index + 1}: trait_type is required`)
        }
        if (attr.value === undefined || attr.value === null || attr.value === '') {
          errors.push(`Attribute ${index + 1}: value is required`)
        }
      })
    }
    
    return {
      isValid: errors.length === 0,
      errors
    }
  }

  estimateGasCost(operation: 'deploy' | 'mint' | 'batch_mint', count: number = 1): {
    gasEstimate: number
    costInETH: string
    costInUSD: string
  } {
    let baseGas = 0
    
    switch (operation) {
      case 'deploy':
        baseGas = 2500000
        break
      case 'mint':
        baseGas = 100000
        break
      case 'batch_mint':
        baseGas = 80000 * count + 50000
        break
    }
    
    const gasPrice = 20 // 20 gwei
    const ethPrice = 2500 // $2500 per ETH
    
    const costInETH = (baseGas * gasPrice * 1e-9).toFixed(6)
    const costInUSD = (parseFloat(costInETH) * ethPrice).toFixed(2)
    
    return {
      gasEstimate: baseGas,
      costInETH,
      costInUSD
    }
  }

  generateMetadataTemplate(): NFTMetadata {
    return {
      name: '',
      description: '',
      image: '',
      attributes: [
        {
          trait_type: 'Background',
          value: ''
        },
        {
          trait_type: 'Eyes',
          value: ''
        },
        {
          trait_type: 'Mouth',
          value: ''
        }
      ]
    }
  }

  // Contract Templates
  getERC721ContractCode(
    name: string,
    symbol: string,
    baseTokenURI: string,
    maxSupply: number,
    royaltyPercentage: number,
    royaltyRecipient: Address
  ): string {
    return `
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/interfaces/IERC2981.sol";

contract ${name.replace(/\s+/g, '')} is ERC721, ERC721Enumerable, ERC721URIStorage, Ownable, IERC2981 {
    uint256 public constant MAX_SUPPLY = ${maxSupply};
    uint256 private _currentTokenId = 0;
    string private _baseTokenURI = "${baseTokenURI}";
    
    uint256 public royaltyPercentage = ${royaltyPercentage * 100}; // basis points
    address public royaltyRecipient = ${royaltyRecipient};
    
    constructor() ERC721("${name}", "${symbol}") {}
    
    function mint(address to) public onlyOwner {
        require(_currentTokenId < MAX_SUPPLY, "Max supply reached");
        _currentTokenId++;
        _safeMint(to, _currentTokenId);
    }
    
    function batchMint(address[] memory recipients) public onlyOwner {
        require(_currentTokenId + recipients.length <= MAX_SUPPLY, "Would exceed max supply");
        
        for (uint256 i = 0; i < recipients.length; i++) {
            _currentTokenId++;
            _safeMint(recipients[i], _currentTokenId);
        }
    }
    
    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
    }
    
    function setBaseURI(string memory newBaseURI) public onlyOwner {
        _baseTokenURI = newBaseURI;
    }
    
    function royaltyInfo(uint256, uint256 salePrice) external view override returns (address, uint256) {
        uint256 royaltyAmount = (salePrice * royaltyPercentage) / 10000;
        return (royaltyRecipient, royaltyAmount);
    }
    
    function supportsInterface(bytes4 interfaceId) public view override(ERC721, ERC721Enumerable, IERC165) returns (bool) {
        return interfaceId == type(IERC2981).interfaceId || super.supportsInterface(interfaceId);
    }
    
    function _beforeTokenTransfer(address from, address to, uint256 tokenId, uint256 batchSize) internal override(ERC721, ERC721Enumerable) {
        super._beforeTokenTransfer(from, to, tokenId, batchSize);
    }
    
    function _burn(uint256 tokenId) internal override(ERC721, ERC721URIStorage) {
        super._burn(tokenId);
    }
    
    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }
}
    `.trim()
  }
}

// Export singleton instance
export const nftMintingService = NFTMintingService.getInstance()
