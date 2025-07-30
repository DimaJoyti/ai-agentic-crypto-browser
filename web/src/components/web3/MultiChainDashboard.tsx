'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { useAccount } from 'wagmi'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import {
  TrendingUp,
  TrendingDown,
  RefreshCw,
  Wallet,
  ArrowUpDown,
  Globe,
  Zap,
  DollarSign
} from 'lucide-react'
import { useChainSwitching } from '@/hooks/useChainSwitching'
import { useBalanceAggregator, type PortfolioSummary } from '@/lib/balance-aggregator'
import { type CrossChainAsset } from '@/lib/multi-chain-manager'
import { toast } from 'sonner'

export function MultiChainDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [portfolioSummary, setPortfolioSummary] = useState<PortfolioSummary | null>(null)
  const [crossChainAssets, setCrossChainAssets] = useState<CrossChainAsset[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [lastUpdated, setLastUpdated] = useState<number | null>(null)

  const { address } = useAccount()
  const {
    currentChain,
    switchState,
    supportedChains,
    switchToChain,
    getRecommendedChains,
    isAutoSwitchEnabled,
    enableAutoSwitch,
    disableAutoSwitch
  } = useChainSwitching()

  const {
    getAggregatedBalances,
    getPortfolioSummary,
    clearCache,
    getCacheStats
  } = useBalanceAggregator()

  // Load portfolio data
  useEffect(() => {
    if (address) {
      loadPortfolioData()
    }
  }, [address])

  const loadPortfolioData = async (forceRefresh = false) => {
    if (!address) return

    setIsLoading(true)
    try {
      // Get aggregated balances
      const balanceData = await getAggregatedBalances(address, { forceRefresh })
      setCrossChainAssets(balanceData.crossChainAssets)

      // Get portfolio summary
      const summary = await getPortfolioSummary(address)
      setPortfolioSummary(summary)

      setLastUpdated(Date.now())
    } catch (error) {
      console.error('Failed to load portfolio data:', error)
      toast.error('Failed to load portfolio data')
    } finally {
      setIsLoading(false)
    }
  }

  const handleRefresh = () => {
    loadPortfolioData(true)
    toast.success('Portfolio data refreshed')
  }

  const handleChainSwitch = async (chainId: number) => {
    const success = await switchToChain(chainId)
    if (success) {
      // Refresh data after chain switch
      setTimeout(() => loadPortfolioData(), 1000)
    }
  }

  const formatCurrency = (value: string | number) => {
    const num = typeof value === 'string' ? parseFloat(value) : value
    if (num >= 1000000) {
      return `$${(num / 1000000).toFixed(2)}M`
    } else if (num >= 1000) {
      return `$${(num / 1000).toFixed(2)}K`
    }
    return `$${num.toFixed(2)}`
  }

  const formatPercentage = (value: number) => {
    const sign = value >= 0 ? '+' : ''
    return `${sign}${value.toFixed(2)}%`
  }

  const getChangeColor = (value: number) => {
    return value >= 0 ? 'text-green-600' : 'text-red-600'
  }

  const getChangeIcon = (value: number) => {
    return value >= 0 ? TrendingUp : TrendingDown
  }

  const recommendedChains = getRecommendedChains()

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Multi-Chain Portfolio</h2>
          <p className="text-muted-foreground">
            Manage your assets across multiple blockchain networks
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          {lastUpdated && (
            <span className="text-sm text-muted-foreground">
              Updated {new Date(lastUpdated).toLocaleTimeString()}
            </span>
          )}
        </div>
      </div>

      {/* Portfolio Overview */}
      {portfolioSummary && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioSummary.totalValueUSD)}</p>
                </div>
                <DollarSign className="w-8 h-8 text-green-500" />
              </div>
              <div className="mt-2 flex items-center text-sm">
                {(() => {
                  const ChangeIcon = getChangeIcon(portfolioSummary.totalChange24hPercent)
                  return (
                    <>
                      <ChangeIcon className="w-4 h-4 mr-1" />
                      <span className={getChangeColor(portfolioSummary.totalChange24hPercent)}>
                        {formatPercentage(portfolioSummary.totalChange24hPercent)}
                      </span>
                      <span className="text-muted-foreground ml-1">24h</span>
                    </>
                  )
                })()}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Active Chains</p>
                  <p className="text-2xl font-bold">{portfolioSummary.chainCount}</p>
                </div>
                <Globe className="w-8 h-8 text-blue-500" />
              </div>
              <div className="mt-2 text-sm text-muted-foreground">
                Across {supportedChains.length} supported networks
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Assets</p>
                  <p className="text-2xl font-bold">{portfolioSummary.tokenCount}</p>
                </div>
                <Wallet className="w-8 h-8 text-purple-500" />
              </div>
              <div className="mt-2 text-sm text-muted-foreground">
                Unique tokens and coins
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Current Chain</p>
                  <p className="text-lg font-bold">{currentChain?.name || 'Not Connected'}</p>
                </div>
                <div className="w-8 h-8 rounded-full flex items-center justify-center"
                     style={{ backgroundColor: currentChain?.color || '#666' }}>
                  <Zap className="w-4 h-4 text-white" />
                </div>
              </div>
              <div className="mt-2">
                {switchState.isSwitching ? (
                  <Badge variant="secondary">Switching...</Badge>
                ) : (
                  <Badge variant="outline">{currentChain?.category || 'Unknown'}</Badge>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="chains">Chains</TabsTrigger>
          <TabsTrigger value="assets">Assets</TabsTrigger>
          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {portfolioSummary && (
            <div className="grid gap-6 md:grid-cols-2">
              {/* Chain Distribution */}
              <Card>
                <CardHeader>
                  <CardTitle>Chain Distribution</CardTitle>
                  <CardDescription>Value distribution across networks</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {portfolioSummary.chainDistribution.map((chain) => (
                      <div key={chain.chainId} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <div 
                              className="w-3 h-3 rounded-full"
                              style={{ backgroundColor: chain.color }}
                            />
                            <span className="font-medium">{chain.chainName}</span>
                          </div>
                          <div className="text-right">
                            <div className="font-medium">{formatCurrency(chain.valueUSD)}</div>
                            <div className="text-sm text-muted-foreground">
                              {chain.percentage.toFixed(1)}%
                            </div>
                          </div>
                        </div>
                        <Progress value={chain.percentage} className="h-2" />
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              {/* Top Assets */}
              <Card>
                <CardHeader>
                  <CardTitle>Top Assets</CardTitle>
                  <CardDescription>Your largest holdings by value</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {portfolioSummary.topAssets.slice(0, 5).map((asset, index) => (
                      <div key={asset.symbol} className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-secondary rounded-full flex items-center justify-center">
                            <span className="text-sm font-medium">#{index + 1}</span>
                          </div>
                          <div>
                            <p className="font-medium">{asset.symbol}</p>
                            <p className="text-sm text-muted-foreground">
                              {asset.chains.length} chain{asset.chains.length > 1 ? 's' : ''}
                            </p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatCurrency(asset.totalValueUSD)}</p>
                          <p className="text-sm text-muted-foreground">
                            {parseFloat(asset.totalBalance).toFixed(4)} {asset.symbol}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        <TabsContent value="chains" className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {recommendedChains.map((chain) => (
              <motion.div
                key={chain.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 }}
              >
                <Card className={`cursor-pointer transition-all hover:shadow-md ${
                  currentChain?.id === chain.id ? 'ring-2 ring-primary' : ''
                }`}>
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between mb-4">
                      <div className="flex items-center gap-3">
                        <div 
                          className="w-10 h-10 rounded-lg flex items-center justify-center"
                          style={{ backgroundColor: chain.color }}
                        >
                          <Zap className="w-5 h-5 text-white" />
                        </div>
                        <div>
                          <h3 className="font-medium">{chain.name}</h3>
                          <p className="text-sm text-muted-foreground capitalize">
                            {chain.category}
                          </p>
                        </div>
                      </div>
                      {currentChain?.id === chain.id && (
                        <Badge variant="default">Active</Badge>
                      )}
                    </div>

                    <div className="space-y-2 mb-4">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">TVL</span>
                        <span>{chain.tvl || 'N/A'}</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Daily Volume</span>
                        <span>{chain.dailyVolume || 'N/A'}</span>
                      </div>
                    </div>

                    <div className="flex flex-wrap gap-1 mb-4">
                      {chain.features.slice(0, 3).map((feature) => (
                        <Badge key={feature} variant="secondary" className="text-xs">
                          {feature}
                        </Badge>
                      ))}
                    </div>

                    <Button
                      variant={currentChain?.id === chain.id ? "secondary" : "outline"}
                      size="sm"
                      className="w-full"
                      onClick={() => handleChainSwitch(chain.id)}
                      disabled={switchState.isSwitching || currentChain?.id === chain.id}
                    >
                      {switchState.isSwitching && switchState.targetChainId === chain.id ? (
                        <>
                          <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                          Switching...
                        </>
                      ) : currentChain?.id === chain.id ? (
                        'Connected'
                      ) : (
                        <>
                          <ArrowUpDown className="w-4 h-4 mr-2" />
                          Switch to {chain.name}
                        </>
                      )}
                    </Button>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="assets" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Cross-Chain Assets</CardTitle>
              <CardDescription>Your assets across all connected chains</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {crossChainAssets.map((asset) => (
                  <div key={asset.symbol} className="border rounded-lg p-4">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-secondary rounded-lg flex items-center justify-center">
                          <Wallet className="w-5 h-5" />
                        </div>
                        <div>
                          <h4 className="font-medium">{asset.symbol}</h4>
                          <p className="text-sm text-muted-foreground">{asset.name}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">{formatCurrency(asset.totalValueUSD)}</p>
                        <p className="text-sm text-muted-foreground">
                          {parseFloat(asset.totalBalance).toFixed(4)} {asset.symbol}
                        </p>
                      </div>
                    </div>

                    <div className="grid gap-2 md:grid-cols-2 lg:grid-cols-3">
                      {asset.chains.map((chain) => (
                        <div key={`${asset.symbol}-${chain.chainId}`} 
                             className="flex items-center justify-between p-2 bg-secondary/50 rounded">
                          <span className="text-sm font-medium">{chain.chainName}</span>
                          <div className="text-right">
                            <p className="text-sm font-medium">{formatCurrency(chain.valueUSD)}</p>
                            <p className="text-xs text-muted-foreground">
                              {parseFloat(chain.balance).toFixed(4)}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}

                {crossChainAssets.length === 0 && (
                  <div className="text-center py-8 text-muted-foreground">
                    No cross-chain assets found
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="settings" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Multi-Chain Settings</CardTitle>
              <CardDescription>Configure your multi-chain experience</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium">Auto Chain Switching</h4>
                  <p className="text-sm text-muted-foreground">
                    Automatically switch to optimal chains for dApps
                  </p>
                </div>
                <Button
                  variant={isAutoSwitchEnabled ? "default" : "outline"}
                  onClick={isAutoSwitchEnabled ? disableAutoSwitch : enableAutoSwitch}
                >
                  {isAutoSwitchEnabled ? 'Enabled' : 'Disabled'}
                </Button>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium">Clear Cache</h4>
                  <p className="text-sm text-muted-foreground">
                    Clear cached balance and price data
                  </p>
                </div>
                <Button variant="outline" onClick={() => {
                  clearCache()
                  toast.success('Cache cleared')
                }}>
                  Clear Cache
                </Button>
              </div>

              <div>
                <h4 className="font-medium mb-2">Cache Statistics</h4>
                <div className="grid gap-2 md:grid-cols-2">
                  {(() => {
                    const stats = getCacheStats()
                    return (
                      <>
                        <div className="text-sm">
                          <span className="text-muted-foreground">Balance Cache:</span>
                          <span className="ml-2">{stats.balanceCacheSize} entries</span>
                        </div>
                        <div className="text-sm">
                          <span className="text-muted-foreground">Price Cache:</span>
                          <span className="ml-2">{stats.priceCacheSize} entries</span>
                        </div>
                        <div className="text-sm">
                          <span className="text-muted-foreground">History Cache:</span>
                          <span className="ml-2">{stats.historyCacheSize} entries</span>
                        </div>
                        <div className="text-sm">
                          <span className="text-muted-foreground">Memory Usage:</span>
                          <span className="ml-2">{stats.totalMemoryUsage}</span>
                        </div>
                      </>
                    )
                  })()}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
