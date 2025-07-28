'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Image as ImageIcon, 
  ExternalLink, 
  Heart,
  Share2,
  TrendingUp,
  Clock,
  User,
  Crown,
  Star,
  Zap,
  Shield,
  Copy,
  Eye,
  BarChart3,
  Activity,
  DollarSign
} from 'lucide-react'
import { useNFTCollection } from '@/hooks/useNFTCollection'
import { NFTRarity, type NFT, type NFTValuation } from '@/lib/nft-service'
import { toast } from 'sonner'

interface NFTDetailsProps {
  nftId: string
  onClose?: () => void
}

export function NFTDetails({ nftId, onClose }: NFTDetailsProps) {
  const [valuation, setValuation] = useState<NFTValuation | null>(null)
  const [isLoadingValuation, setIsLoadingValuation] = useState(false)

  const {
    getNFT,
    getCollection,
    getNFTValuation,
    formatETH,
    formatCurrency
  } = useNFTCollection()

  const nft = getNFT(nftId)
  const collection = nft ? getCollection(nft.contractAddress) : null

  useEffect(() => {
    if (nft) {
      loadValuation()
    }
  }, [nft])

  const loadValuation = async () => {
    if (!nft) return
    
    setIsLoadingValuation(true)
    try {
      const val = await getNFTValuation(nft.id)
      setValuation(val)
    } catch (error) {
      console.error('Failed to load valuation:', error)
    } finally {
      setIsLoadingValuation(false)
    }
  }

  const getRarityColor = (rarity: NFTRarity) => {
    switch (rarity) {
      case NFTRarity.COMMON:
        return 'bg-gray-100 text-gray-800'
      case NFTRarity.UNCOMMON:
        return 'bg-green-100 text-green-800'
      case NFTRarity.RARE:
        return 'bg-blue-100 text-blue-800'
      case NFTRarity.EPIC:
        return 'bg-purple-100 text-purple-800'
      case NFTRarity.LEGENDARY:
        return 'bg-yellow-100 text-yellow-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getRarityIcon = (rarity: NFTRarity) => {
    switch (rarity) {
      case NFTRarity.LEGENDARY:
        return <Crown className="w-4 h-4" />
      case NFTRarity.EPIC:
        return <Star className="w-4 h-4" />
      case NFTRarity.RARE:
        return <Zap className="w-4 h-4" />
      default:
        return null
    }
  }

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text)
    toast.success(`${label} copied to clipboard`)
  }

  const formatDate = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  if (!nft) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <ImageIcon className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">NFT Not Found</h3>
          <p className="text-muted-foreground">
            The requested NFT could not be found
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">{nft.metadata.name}</h2>
          <p className="text-muted-foreground">
            {collection?.name} #{nft.tokenId}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <Heart className="w-4 h-4 mr-2" />
            Favorite
          </Button>
          <Button variant="outline" size="sm">
            <Share2 className="w-4 h-4 mr-2" />
            Share
          </Button>
          {onClose && (
            <Button variant="outline" size="sm" onClick={onClose}>
              Close
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* NFT Image and Basic Info */}
        <div className="space-y-6">
          <Card>
            <CardContent className="p-6">
              <div className="aspect-square bg-muted rounded-lg flex items-center justify-center mb-4">
                <ImageIcon className="w-24 h-24 text-muted-foreground" />
              </div>
              
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Current Price</span>
                  <span className="font-bold text-lg">{formatETH(nft.currentPrice || nft.floorPrice || '0')}</span>
                </div>
                
                {nft.lastSalePrice && (
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Last Sale</span>
                    <span className="font-medium">{formatETH(nft.lastSalePrice)}</span>
                  </div>
                )}

                {nft.rarity && (
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Rarity</span>
                    <Badge className={getRarityColor(nft.rarity)}>
                      {getRarityIcon(nft.rarity)}
                      <span className="ml-1 capitalize">{nft.rarity}</span>
                    </Badge>
                  </div>
                )}

                {nft.rarityRank && (
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Rarity Rank</span>
                    <span className="font-medium">#{nft.rarityRank} / {nft.totalSupply}</span>
                  </div>
                )}

                {nft.isListed && (
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Listed Price</span>
                    <span className="font-medium text-green-600">{formatETH(nft.listingPrice || '0')}</span>
                  </div>
                )}
              </div>

              <div className="flex gap-2 mt-4">
                <Button className="flex-1">
                  <DollarSign className="w-4 h-4 mr-2" />
                  Make Offer
                </Button>
                <Button variant="outline" className="flex-1">
                  <ExternalLink className="w-4 h-4 mr-2" />
                  View on OpenSea
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Collection Info */}
          {collection && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <ImageIcon className="w-5 h-5" />
                  Collection
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                    <ImageIcon className="w-6 h-6" />
                  </div>
                  <div>
                    <h4 className="font-medium flex items-center gap-2">
                      {collection.name}
                      {collection.verified && <Shield className="w-4 h-4 text-blue-500" />}
                    </h4>
                    <p className="text-sm text-muted-foreground">{collection.symbol}</p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">Floor Price</span>
                    <p className="font-medium">{formatETH(collection.floorPrice)}</p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Total Supply</span>
                    <p className="font-medium">{collection.totalSupply.toLocaleString()}</p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Owners</span>
                    <p className="font-medium">{collection.ownersCount.toLocaleString()}</p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Royalty</span>
                    <p className="font-medium">{collection.royaltyFee}%</p>
                  </div>
                </div>

                <Button variant="outline" className="w-full mt-4">
                  <Eye className="w-4 h-4 mr-2" />
                  View Collection
                </Button>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Detailed Information */}
        <div className="space-y-6">
          <Tabs defaultValue="details">
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="details">Details</TabsTrigger>
              <TabsTrigger value="traits">Traits</TabsTrigger>
              <TabsTrigger value="history">History</TabsTrigger>
              <TabsTrigger value="valuation">Valuation</TabsTrigger>
            </TabsList>

            <TabsContent value="details" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle>Description</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    {nft.metadata.description}
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Details</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Contract Address</span>
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-mono">
                          {nft.contractAddress.slice(0, 6)}...{nft.contractAddress.slice(-4)}
                        </span>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => copyToClipboard(nft.contractAddress, 'Contract address')}
                        >
                          <Copy className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Token ID</span>
                      <span className="text-sm font-medium">{nft.tokenId}</span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Token Standard</span>
                      <span className="text-sm font-medium">{nft.standard}</span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Chain</span>
                      <span className="text-sm font-medium">Ethereum</span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Minted</span>
                      <span className="text-sm font-medium">{formatDate(nft.mintedAt)}</span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Last Transfer</span>
                      <span className="text-sm font-medium">{formatDate(nft.lastTransferAt)}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="traits" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle>Attributes</CardTitle>
                  <CardDescription>
                    Traits and properties of this NFT
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 gap-3">
                    {nft.metadata.attributes.map((attribute, index) => (
                      <motion.div
                        key={index}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.05 }}
                        className="border rounded-lg p-3 text-center"
                      >
                        <p className="text-xs text-muted-foreground uppercase tracking-wide">
                          {attribute.trait_type}
                        </p>
                        <p className="font-medium mt-1">{attribute.value}</p>
                        {attribute.rarity_score && (
                          <p className="text-xs text-blue-600 mt-1">
                            {attribute.rarity_score.toFixed(1)}% rarity
                          </p>
                        )}
                      </motion.div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="history" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle>Transaction History</CardTitle>
                  <CardDescription>
                    Recent activity for this NFT
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center gap-3 p-3 border rounded-lg">
                      <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                        <TrendingUp className="w-4 h-4 text-green-600" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium">Sale</p>
                        <p className="text-xs text-muted-foreground">
                          Sold for {formatETH(nft.lastSalePrice || '0')} • {formatDate(nft.lastTransferAt)}
                        </p>
                      </div>
                      <Button variant="ghost" size="sm">
                        <ExternalLink className="w-3 h-3" />
                      </Button>
                    </div>

                    <div className="flex items-center gap-3 p-3 border rounded-lg">
                      <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                        <User className="w-4 h-4 text-blue-600" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium">Transfer</p>
                        <p className="text-xs text-muted-foreground">
                          Transferred • {formatDate(nft.lastTransferAt)}
                        </p>
                      </div>
                      <Button variant="ghost" size="sm">
                        <ExternalLink className="w-3 h-3" />
                      </Button>
                    </div>

                    <div className="flex items-center gap-3 p-3 border rounded-lg">
                      <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center">
                        <Zap className="w-4 h-4 text-purple-600" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium">Mint</p>
                        <p className="text-xs text-muted-foreground">
                          Minted • {formatDate(nft.mintedAt)}
                        </p>
                      </div>
                      <Button variant="ghost" size="sm">
                        <ExternalLink className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="valuation" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <BarChart3 className="w-5 h-5" />
                    AI Valuation
                  </CardTitle>
                  <CardDescription>
                    Estimated value based on market data and rarity
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {isLoadingValuation ? (
                    <div className="text-center py-8">
                      <div className="animate-spin w-8 h-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4" />
                      <p className="text-sm text-muted-foreground">Calculating valuation...</p>
                    </div>
                  ) : valuation ? (
                    <div className="space-y-4">
                      <div className="text-center p-4 bg-muted rounded-lg">
                        <p className="text-2xl font-bold">{formatETH(valuation.estimatedValue)}</p>
                        <p className="text-sm text-muted-foreground">
                          Estimated Value • {valuation.confidence}% confidence
                        </p>
                      </div>

                      <div className="space-y-3">
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-muted-foreground">Floor Price</span>
                          <span className="font-medium">{formatETH(valuation.factors.floorPrice)}</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-muted-foreground">Rarity Score</span>
                          <span className="font-medium">#{valuation.factors.rarityScore}</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-muted-foreground">Market Trend</span>
                          <Badge variant="outline" className="capitalize">
                            {valuation.factors.marketTrend}
                          </Badge>
                        </div>
                      </div>

                      <div>
                        <h5 className="text-sm font-medium mb-2">Recent Sales</h5>
                        <div className="space-y-1">
                          {valuation.factors.recentSales.map((sale, index) => (
                            <div key={index} className="flex items-center justify-between text-sm">
                              <span className="text-muted-foreground">Sale {index + 1}</span>
                              <span>{formatETH(sale)}</span>
                            </div>
                          ))}
                        </div>
                      </div>

                      <p className="text-xs text-muted-foreground">
                        Last updated: {formatDate(valuation.lastUpdated)}
                      </p>
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <BarChart3 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                      <p className="text-sm text-muted-foreground">Valuation not available</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  )
}
