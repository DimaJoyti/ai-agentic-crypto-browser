'use client'

import { useState, useRef } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  Palette, 
  Upload, 
  Zap,
  Plus,
  Image as ImageIcon,
  FileText,
  Coins,
  Users,
  RefreshCw,
  CheckCircle,
  AlertCircle,
  Copy,
  ExternalLink,
  Trash2,
  Settings,
  Code,
  Sparkles
} from 'lucide-react'
import { useNFTMinting } from '@/hooks/useNFTMinting'
import { NFTStandard, type NFTMetadata } from '@/lib/nft-minting'
import { type Address } from 'viem'

interface NFTMintingDashboardProps {
  userAddress?: Address
  chainId?: number
}

export function NFTMintingDashboard({ userAddress, chainId = 1 }: NFTMintingDashboardProps) {
  const [activeTab, setActiveTab] = useState('create')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [metadata, setMetadata] = useState<NFTMetadata>({
    name: '',
    description: '',
    image: '',
    attributes: [
      { trait_type: 'Background', value: '' },
      { trait_type: 'Eyes', value: '' },
      { trait_type: 'Mouth', value: '' }
    ]
  })
  const [collectionForm, setCollectionForm] = useState({
    name: '',
    symbol: '',
    description: '',
    maxSupply: 1000,
    royaltyPercentage: 5,
    royaltyRecipient: userAddress || ''
  })

  const fileInputRef = useRef<HTMLInputElement>(null)

  const {
    collections,
    isLoading,
    progress,
    error,
    lastMintResult,
    lastDeployment,
    loadCollections,
    uploadToIPFS,
    uploadMetadata,
    createCollection,
    mintNFT,
    batchMintNFTs,
    validateMetadata,
    estimateGasCost,
    generateMetadataTemplate,
    getContractCode,
    clearError,
    hasCollections,
    isUploading,
    isDeploying,
    isMinting,
    progressPercentage,
    recentCollections,
    totalCollections,
    totalNFTsMinted
  } = useNFTMinting({
    enableNotifications: true,
    onProgress: (progress) => {
      console.log('Minting progress:', progress)
    },
    onComplete: (result) => {
      console.log('Minting complete:', result)
    },
    onError: (error) => {
      console.error('Minting error:', error)
    }
  })

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      setSelectedFile(file)
      // Create preview URL
      const previewUrl = URL.createObjectURL(file)
      setMetadata(prev => ({ ...prev, image: previewUrl }))
    }
  }

  const handleUploadImage = async () => {
    if (!selectedFile) return

    const result = await uploadToIPFS(selectedFile)
    if (result) {
      setMetadata(prev => ({ ...prev, image: result.url }))
    }
  }

  const handleCreateCollection = async () => {
    if (!userAddress || !selectedFile) return

    await createCollection(
      collectionForm.name,
      collectionForm.symbol,
      collectionForm.description,
      selectedFile,
      collectionForm.maxSupply,
      collectionForm.royaltyPercentage,
      collectionForm.royaltyRecipient as Address,
      chainId,
      userAddress,
      NFTStandard.ERC721
    )
  }

  const handleMintNFT = async () => {
    if (!userAddress || !collections[0]) return

    await mintNFT(
      collections[0].contractAddress,
      userAddress,
      metadata,
      chainId
    )
  }

  const addAttribute = () => {
    setMetadata(prev => ({
      ...prev,
      attributes: [...prev.attributes, { trait_type: '', value: '' }]
    }))
  }

  const removeAttribute = (index: number) => {
    setMetadata(prev => ({
      ...prev,
      attributes: prev.attributes.filter((_, i) => i !== index)
    }))
  }

  const updateAttribute = (index: number, field: 'trait_type' | 'value', value: string) => {
    setMetadata(prev => ({
      ...prev,
      attributes: prev.attributes.map((attr, i) => 
        i === index ? { ...attr, [field]: value } : attr
      )
    }))
  }

  const validation = validateMetadata(metadata)
  const deployGasCost = estimateGasCost('deploy')
  const mintGasCost = estimateGasCost('mint')

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Palette className="w-6 h-6" />
            NFT Minting & Creation
          </h2>
          <p className="text-muted-foreground">
            Create, deploy, and mint NFT collections with IPFS metadata storage
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={() => loadCollections(userAddress)}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Collections</p>
                <p className="text-2xl font-bold">{totalCollections}</p>
              </div>
              <Coins className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">NFTs Minted</p>
                <p className="text-2xl font-bold">{totalNFTsMinted}</p>
              </div>
              <Sparkles className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Deploy Cost</p>
                <p className="text-2xl font-bold">{deployGasCost.costInETH} ETH</p>
                <p className="text-xs text-muted-foreground">${deployGasCost.costInUSD}</p>
              </div>
              <Zap className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Mint Cost</p>
                <p className="text-2xl font-bold">{mintGasCost.costInETH} ETH</p>
                <p className="text-xs text-muted-foreground">${mintGasCost.costInUSD}</p>
              </div>
              <Users className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Progress Indicator */}
      {progress && (
        <Card>
          <CardContent className="p-6">
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="font-medium">{progress.message}</h3>
                <Badge variant="outline">
                  {progress.currentStep}/{progress.totalSteps}
                </Badge>
              </div>
              <Progress value={progressPercentage} className="h-2" />
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                {isUploading && <Upload className="w-4 h-4 animate-pulse" />}
                {isDeploying && <Settings className="w-4 h-4 animate-spin" />}
                {isMinting && <Zap className="w-4 h-4 animate-pulse" />}
                <span>{progressPercentage}% complete</span>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Error Display */}
      {error && (
        <Card className="border-red-200 bg-red-50">
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <AlertCircle className="w-5 h-5 text-red-600" />
              <p className="text-red-800">{error}</p>
              <Button variant="ghost" size="sm" onClick={clearError}>
                <Trash2 className="w-4 h-4" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="create">Create Collection</TabsTrigger>
          <TabsTrigger value="mint">Mint NFT</TabsTrigger>
          <TabsTrigger value="collections">My Collections</TabsTrigger>
          <TabsTrigger value="tools">Tools</TabsTrigger>
        </TabsList>

        <TabsContent value="create" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Create New Collection</CardTitle>
              <CardDescription>
                Deploy a new NFT smart contract with custom metadata
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="collection-name">Collection Name</Label>
                    <Input
                      id="collection-name"
                      placeholder="My Awesome Collection"
                      value={collectionForm.name}
                      onChange={(e) => setCollectionForm(prev => ({ ...prev, name: e.target.value }))}
                    />
                  </div>

                  <div>
                    <Label htmlFor="collection-symbol">Symbol</Label>
                    <Input
                      id="collection-symbol"
                      placeholder="MAC"
                      value={collectionForm.symbol}
                      onChange={(e) => setCollectionForm(prev => ({ ...prev, symbol: e.target.value }))}
                    />
                  </div>

                  <div>
                    <Label htmlFor="collection-description">Description</Label>
                    <Textarea
                      id="collection-description"
                      placeholder="Describe your collection..."
                      value={collectionForm.description}
                      onChange={(e) => setCollectionForm(prev => ({ ...prev, description: e.target.value }))}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="max-supply">Max Supply</Label>
                      <Input
                        id="max-supply"
                        type="number"
                        value={collectionForm.maxSupply}
                        onChange={(e) => setCollectionForm(prev => ({ ...prev, maxSupply: parseInt(e.target.value) }))}
                      />
                    </div>

                    <div>
                      <Label htmlFor="royalty">Royalty %</Label>
                      <Input
                        id="royalty"
                        type="number"
                        min="0"
                        max="10"
                        value={collectionForm.royaltyPercentage}
                        onChange={(e) => setCollectionForm(prev => ({ ...prev, royaltyPercentage: parseFloat(e.target.value) }))}
                      />
                    </div>
                  </div>

                  <div>
                    <Label htmlFor="royalty-recipient">Royalty Recipient</Label>
                    <Input
                      id="royalty-recipient"
                      placeholder="0x..."
                      value={collectionForm.royaltyRecipient}
                      onChange={(e) => setCollectionForm(prev => ({ ...prev, royaltyRecipient: e.target.value }))}
                    />
                  </div>
                </div>

                <div className="space-y-4">
                  <div>
                    <Label>Collection Image</Label>
                    <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6 text-center">
                      {selectedFile ? (
                        <div className="space-y-4">
                          <img
                            src={URL.createObjectURL(selectedFile)}
                            alt="Preview"
                            className="w-32 h-32 object-cover rounded-lg mx-auto"
                          />
                          <p className="text-sm text-muted-foreground">{selectedFile.name}</p>
                          <Button variant="outline" size="sm" onClick={() => fileInputRef.current?.click()}>
                            Change Image
                          </Button>
                        </div>
                      ) : (
                        <div className="space-y-4">
                          <ImageIcon className="w-12 h-12 text-muted-foreground mx-auto" />
                          <div>
                            <Button variant="outline" onClick={() => fileInputRef.current?.click()}>
                              <Upload className="w-4 h-4 mr-2" />
                              Upload Image
                            </Button>
                          </div>
                        </div>
                      )}
                      <input
                        ref={fileInputRef}
                        type="file"
                        accept="image/*"
                        onChange={handleFileSelect}
                        className="hidden"
                      />
                    </div>
                  </div>

                  <div className="p-4 bg-muted rounded-lg">
                    <h4 className="font-medium mb-2">Gas Estimate</h4>
                    <div className="space-y-1 text-sm">
                      <div className="flex justify-between">
                        <span>Gas Limit:</span>
                        <span>{deployGasCost.gasEstimate.toLocaleString()}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Cost (ETH):</span>
                        <span>{deployGasCost.costInETH}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Cost (USD):</span>
                        <span>${deployGasCost.costInUSD}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <Button 
                onClick={handleCreateCollection}
                disabled={!userAddress || !selectedFile || !collectionForm.name || isLoading}
                className="w-full"
              >
                {isLoading ? (
                  <>
                    <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                    Creating Collection...
                  </>
                ) : (
                  <>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Collection
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="mint" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Mint NFT</CardTitle>
              <CardDescription>
                Create and mint a new NFT with custom metadata
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="nft-name">NFT Name</Label>
                    <Input
                      id="nft-name"
                      placeholder="My Awesome NFT"
                      value={metadata.name}
                      onChange={(e) => setMetadata(prev => ({ ...prev, name: e.target.value }))}
                    />
                  </div>

                  <div>
                    <Label htmlFor="nft-description">Description</Label>
                    <Textarea
                      id="nft-description"
                      placeholder="Describe your NFT..."
                      value={metadata.description}
                      onChange={(e) => setMetadata(prev => ({ ...prev, description: e.target.value }))}
                    />
                  </div>

                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <Label>Attributes</Label>
                      <Button variant="outline" size="sm" onClick={addAttribute}>
                        <Plus className="w-4 h-4 mr-2" />
                        Add
                      </Button>
                    </div>
                    <div className="space-y-2">
                      {metadata.attributes.map((attr, index) => (
                        <div key={index} className="flex gap-2">
                          <Input
                            placeholder="Trait type"
                            value={attr.trait_type}
                            onChange={(e) => updateAttribute(index, 'trait_type', e.target.value)}
                          />
                          <Input
                            placeholder="Value"
                            value={attr.value}
                            onChange={(e) => updateAttribute(index, 'value', e.target.value)}
                          />
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => removeAttribute(index)}
                          >
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>

                <div className="space-y-4">
                  <div>
                    <Label>NFT Image</Label>
                    <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6 text-center">
                      {metadata.image ? (
                        <div className="space-y-4">
                          <img
                            src={metadata.image}
                            alt="NFT Preview"
                            className="w-32 h-32 object-cover rounded-lg mx-auto"
                          />
                          <Button variant="outline" size="sm" onClick={handleUploadImage}>
                            <Upload className="w-4 h-4 mr-2" />
                            Upload to IPFS
                          </Button>
                        </div>
                      ) : (
                        <div className="space-y-4">
                          <ImageIcon className="w-12 h-12 text-muted-foreground mx-auto" />
                          <Button variant="outline" onClick={() => fileInputRef.current?.click()}>
                            <Upload className="w-4 h-4 mr-2" />
                            Select Image
                          </Button>
                        </div>
                      )}
                    </div>
                  </div>

                  <div className="p-4 bg-muted rounded-lg">
                    <h4 className="font-medium mb-2">Metadata Validation</h4>
                    {validation.isValid ? (
                      <div className="flex items-center gap-2 text-green-600">
                        <CheckCircle className="w-4 h-4" />
                        <span className="text-sm">Metadata is valid</span>
                      </div>
                    ) : (
                      <div className="space-y-1">
                        {validation.errors.map((error, index) => (
                          <div key={index} className="flex items-center gap-2 text-red-600">
                            <AlertCircle className="w-4 h-4" />
                            <span className="text-sm">{error}</span>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>

              <Button 
                onClick={handleMintNFT}
                disabled={!userAddress || !validation.isValid || !hasCollections || isLoading}
                className="w-full"
              >
                {isLoading ? (
                  <>
                    <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                    Minting NFT...
                  </>
                ) : (
                  <>
                    <Zap className="w-4 h-4 mr-2" />
                    Mint NFT
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="collections" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>My Collections</CardTitle>
              <CardDescription>
                Manage your deployed NFT collections
              </CardDescription>
            </CardHeader>
            <CardContent>
              {hasCollections ? (
                <div className="space-y-4">
                  {recentCollections.map((collection, index) => (
                    <motion.div
                      key={collection.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="border rounded-lg p-4"
                    >
                      <div className="flex items-center justify-between mb-4">
                        <div>
                          <h4 className="font-medium">{collection.name}</h4>
                          <p className="text-sm text-muted-foreground">{collection.symbol}</p>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge variant="outline">{collection.standard}</Badge>
                          <Badge variant="secondary">
                            {collection.totalSupply}/{collection.maxSupply}
                          </Badge>
                        </div>
                      </div>

                      <p className="text-sm text-muted-foreground mb-4">{collection.description}</p>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                        <div>
                          <p className="text-xs text-muted-foreground">Contract</p>
                          <p className="text-sm font-mono">
                            {collection.contractAddress.slice(0, 6)}...{collection.contractAddress.slice(-4)}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Royalty</p>
                          <p className="text-sm">{collection.royaltyPercentage}%</p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Chain</p>
                          <p className="text-sm">{collection.chainId === 1 ? 'Ethereum' : `Chain ${collection.chainId}`}</p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">Created</p>
                          <p className="text-sm">{new Date(collection.createdAt).toLocaleDateString()}</p>
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <Button variant="outline" size="sm">
                          <Copy className="w-3 h-3 mr-2" />
                          Copy Address
                        </Button>
                        <Button variant="outline" size="sm">
                          <ExternalLink className="w-3 h-3 mr-2" />
                          View on Explorer
                        </Button>
                        <Button variant="outline" size="sm">
                          <Code className="w-3 h-3 mr-2" />
                          View Contract
                        </Button>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Coins className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-semibold mb-2">No Collections Yet</h3>
                  <p className="text-muted-foreground mb-4">
                    Create your first NFT collection to get started
                  </p>
                  <Button onClick={() => setActiveTab('create')}>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Collection
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="tools" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Metadata Generator</CardTitle>
                <CardDescription>
                  Generate NFT metadata templates
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Button 
                  onClick={() => setMetadata(generateMetadataTemplate())}
                  className="w-full"
                >
                  <FileText className="w-4 h-4 mr-2" />
                  Generate Template
                </Button>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Gas Calculator</CardTitle>
                <CardDescription>
                  Estimate deployment and minting costs
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span>Deploy Contract:</span>
                    <span>{deployGasCost.costInETH} ETH</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Mint NFT:</span>
                    <span>{mintGasCost.costInETH} ETH</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Batch Mint (10):</span>
                    <span>{estimateGasCost('batch_mint', 10).costInETH} ETH</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {lastMintResult && (
            <Card>
              <CardHeader>
                <CardTitle>Last Mint Result</CardTitle>
                <CardDescription>
                  Details of your most recent NFT mint
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span>Token ID:</span>
                    <span className="font-mono">#{lastMintResult.tokenId}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Transaction:</span>
                    <span className="font-mono">
                      {lastMintResult.transactionHash.slice(0, 10)}...
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span>Gas Used:</span>
                    <span>{lastMintResult.gasUsed.toString()}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Cost:</span>
                    <span>{lastMintResult.mintingCost} ETH</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
