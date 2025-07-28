import { useState, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  nftMintingService, 
  type NFTMetadata,
  type CollectionMetadata,
  type NFTCollection,
  type MintResult,
  type ContractDeployment,
  type MintingProgress,
  type IPFSUploadResult,
  NFTStandard,
  MintingStatus 
} from '@/lib/nft-minting'
import { toast } from 'sonner'

export interface UseNFTMintingOptions {
  enableNotifications?: boolean
  onProgress?: (progress: MintingProgress) => void
  onComplete?: (result: MintResult | ContractDeployment) => void
  onError?: (error: string) => void
}

export interface NFTMintingState {
  collections: NFTCollection[]
  isLoading: boolean
  progress: MintingProgress | null
  error: string | null
  lastMintResult: MintResult | null
  lastDeployment: ContractDeployment | null
}

export function useNFTMinting(options: UseNFTMintingOptions = {}) {
  const {
    enableNotifications = true,
    onProgress,
    onComplete,
    onError
  } = options

  const [state, setState] = useState<NFTMintingState>({
    collections: [],
    isLoading: false,
    progress: null,
    error: null,
    lastMintResult: null,
    lastDeployment: null
  })

  // Update progress
  const updateProgress = useCallback((progress: MintingProgress) => {
    setState(prev => ({ ...prev, progress }))
    onProgress?.(progress)
  }, [onProgress])

  // Load collections
  const loadCollections = useCallback(async (ownerAddress?: Address) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const collections = await nftMintingService.getCollections(ownerAddress)
      setState(prev => ({
        ...prev,
        collections,
        isLoading: false
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load collections'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))

      if (enableNotifications) {
        toast.error('Load Error', {
          description: errorMessage
        })
      }
      onError?.(errorMessage)
    }
  }, [enableNotifications, onError])

  // Upload file to IPFS
  const uploadToIPFS = useCallback(async (file: File): Promise<IPFSUploadResult | null> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      updateProgress({
        status: MintingStatus.UPLOADING_METADATA,
        currentStep: 1,
        totalSteps: 1,
        message: 'Uploading file to IPFS...'
      })

      const result = await nftMintingService.uploadToIPFS(file)

      setState(prev => ({ ...prev, isLoading: false, progress: null }))

      if (enableNotifications) {
        toast.success('Upload Complete', {
          description: 'File uploaded to IPFS successfully'
        })
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to upload to IPFS'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.error('Upload Error', {
          description: errorMessage
        })
      }
      onError?.(errorMessage)
      return null
    }
  }, [enableNotifications, onError, updateProgress])

  // Upload metadata to IPFS
  const uploadMetadata = useCallback(async (metadata: NFTMetadata): Promise<IPFSUploadResult | null> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      // Validate metadata
      const validation = nftMintingService.validateMetadata(metadata)
      if (!validation.isValid) {
        throw new Error(`Invalid metadata: ${validation.errors.join(', ')}`)
      }

      updateProgress({
        status: MintingStatus.UPLOADING_METADATA,
        currentStep: 1,
        totalSteps: 1,
        message: 'Uploading metadata to IPFS...'
      })

      const result = await nftMintingService.uploadMetadataToIPFS(metadata)

      setState(prev => ({ ...prev, isLoading: false, progress: null }))

      if (enableNotifications) {
        toast.success('Metadata Uploaded', {
          description: 'NFT metadata uploaded to IPFS successfully'
        })
      }

      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to upload metadata'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.error('Metadata Error', {
          description: errorMessage
        })
      }
      onError?.(errorMessage)
      return null
    }
  }, [enableNotifications, onError, updateProgress])

  // Create collection
  const createCollection = useCallback(async (
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
  ): Promise<{ collection: NFTCollection; deployment: ContractDeployment } | null> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      updateProgress({
        status: MintingStatus.UPLOADING_METADATA,
        currentStep: 1,
        totalSteps: 3,
        message: 'Uploading collection image...'
      })

      updateProgress({
        status: MintingStatus.UPLOADING_METADATA,
        currentStep: 2,
        totalSteps: 3,
        message: 'Uploading collection metadata...'
      })

      updateProgress({
        status: MintingStatus.DEPLOYING_CONTRACT,
        currentStep: 3,
        totalSteps: 3,
        message: 'Deploying smart contract...'
      })

      const result = await nftMintingService.createCollection(
        name,
        symbol,
        description,
        image,
        maxSupply,
        royaltyPercentage,
        royaltyRecipient,
        chainId,
        ownerAddress,
        standard
      )

      setState(prev => ({
        ...prev,
        collections: [result.collection, ...prev.collections],
        lastDeployment: result.deployment,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.success('Collection Created', {
          description: `${name} collection deployed successfully!`
        })
      }

      onComplete?.(result.deployment)
      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create collection'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.error('Collection Error', {
          description: errorMessage
        })
      }
      onError?.(errorMessage)
      return null
    }
  }, [enableNotifications, onError, onComplete, updateProgress])

  // Mint single NFT
  const mintNFT = useCallback(async (
    contractAddress: Address,
    to: Address,
    metadata: NFTMetadata,
    chainId: number
  ): Promise<MintResult | null> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      // Validate metadata
      const validation = nftMintingService.validateMetadata(metadata)
      if (!validation.isValid) {
        throw new Error(`Invalid metadata: ${validation.errors.join(', ')}`)
      }

      updateProgress({
        status: MintingStatus.UPLOADING_METADATA,
        currentStep: 1,
        totalSteps: 2,
        message: 'Uploading NFT metadata...'
      })

      // Upload metadata to IPFS
      const metadataUpload = await nftMintingService.uploadMetadataToIPFS(metadata)

      updateProgress({
        status: MintingStatus.MINTING,
        currentStep: 2,
        totalSteps: 2,
        message: 'Minting NFT...',
        ipfsHash: metadataUpload.hash
      })

      // Mint NFT
      const result = await nftMintingService.mintNFT(
        contractAddress,
        to,
        metadataUpload.url,
        chainId
      )

      setState(prev => ({
        ...prev,
        lastMintResult: result,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.success('NFT Minted', {
          description: `NFT #${result.tokenId} minted successfully!`
        })
      }

      onComplete?.(result)
      return result
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to mint NFT'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.error('Minting Error', {
          description: errorMessage
        })
      }
      onError?.(errorMessage)
      return null
    }
  }, [enableNotifications, onError, onComplete, updateProgress])

  // Batch mint NFTs
  const batchMintNFTs = useCallback(async (
    contractAddress: Address,
    recipients: Address[],
    metadataList: NFTMetadata[],
    chainId: number
  ): Promise<MintResult[] | null> => {
    if (recipients.length !== metadataList.length) {
      const errorMessage = 'Recipients and metadata arrays must have the same length'
      setState(prev => ({ ...prev, error: errorMessage }))
      if (enableNotifications) {
        toast.error('Batch Mint Error', { description: errorMessage })
      }
      onError?.(errorMessage)
      return null
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const totalSteps = metadataList.length + 1
      const metadataUris: string[] = []

      // Upload all metadata
      for (let i = 0; i < metadataList.length; i++) {
        const validation = nftMintingService.validateMetadata(metadataList[i])
        if (!validation.isValid) {
          throw new Error(`Invalid metadata for NFT ${i + 1}: ${validation.errors.join(', ')}`)
        }

        updateProgress({
          status: MintingStatus.UPLOADING_METADATA,
          currentStep: i + 1,
          totalSteps,
          message: `Uploading metadata ${i + 1}/${metadataList.length}...`
        })

        const metadataUpload = await nftMintingService.uploadMetadataToIPFS(metadataList[i])
        metadataUris.push(metadataUpload.url)
      }

      updateProgress({
        status: MintingStatus.MINTING,
        currentStep: totalSteps,
        totalSteps,
        message: `Batch minting ${metadataList.length} NFTs...`
      })

      // Batch mint NFTs
      const results = await nftMintingService.batchMintNFTs(
        contractAddress,
        recipients,
        metadataUris,
        chainId
      )

      setState(prev => ({
        ...prev,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.success('Batch Mint Complete', {
          description: `Successfully minted ${results.length} NFTs!`
        })
      }

      return results
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to batch mint NFTs'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
        progress: null
      }))

      if (enableNotifications) {
        toast.error('Batch Mint Error', {
          description: errorMessage
        })
      }
      onError?.(errorMessage)
      return null
    }
  }, [enableNotifications, onError, updateProgress])

  // Get collection
  const getCollection = useCallback(async (id: string): Promise<NFTCollection | null> => {
    try {
      return await nftMintingService.getCollection(id)
    } catch (error) {
      console.error('Failed to get collection:', error)
      return null
    }
  }, [])

  // Validate metadata
  const validateMetadata = useCallback((metadata: NFTMetadata) => {
    return nftMintingService.validateMetadata(metadata)
  }, [])

  // Estimate gas cost
  const estimateGasCost = useCallback((operation: 'deploy' | 'mint' | 'batch_mint', count: number = 1) => {
    return nftMintingService.estimateGasCost(operation, count)
  }, [])

  // Generate metadata template
  const generateMetadataTemplate = useCallback(() => {
    return nftMintingService.generateMetadataTemplate()
  }, [])

  // Get contract code
  const getContractCode = useCallback((
    name: string,
    symbol: string,
    baseTokenURI: string,
    maxSupply: number,
    royaltyPercentage: number,
    royaltyRecipient: Address
  ) => {
    return nftMintingService.getERC721ContractCode(
      name,
      symbol,
      baseTokenURI,
      maxSupply,
      royaltyPercentage,
      royaltyRecipient
    )
  }, [])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Clear progress
  const clearProgress = useCallback(() => {
    setState(prev => ({ ...prev, progress: null }))
  }, [])

  return {
    // State
    ...state,

    // Actions
    loadCollections,
    uploadToIPFS,
    uploadMetadata,
    createCollection,
    mintNFT,
    batchMintNFTs,
    getCollection,

    // Utilities
    validateMetadata,
    estimateGasCost,
    generateMetadataTemplate,
    getContractCode,
    clearError,
    clearProgress,

    // Quick access helpers
    hasCollections: state.collections.length > 0,
    isUploading: state.progress?.status === MintingStatus.UPLOADING_METADATA,
    isDeploying: state.progress?.status === MintingStatus.DEPLOYING_CONTRACT,
    isMinting: state.progress?.status === MintingStatus.MINTING,
    isCompleted: state.progress?.status === MintingStatus.COMPLETED,
    isFailed: state.progress?.status === MintingStatus.FAILED,

    // Progress helpers
    progressPercentage: state.progress 
      ? Math.round((state.progress.currentStep / state.progress.totalSteps) * 100)
      : 0,
    
    // Collection helpers
    getCollectionsByStandard: (standard: NFTStandard) => 
      state.collections.filter(collection => collection.standard === standard),
    
    getCollectionsByChain: (chainId: number) =>
      state.collections.filter(collection => collection.chainId === chainId),

    // Recent activity
    recentCollections: state.collections.slice(0, 5),
    totalCollections: state.collections.length,
    totalNFTsMinted: state.collections.reduce((sum, collection) => sum + collection.totalSupply, 0)
  }
}
