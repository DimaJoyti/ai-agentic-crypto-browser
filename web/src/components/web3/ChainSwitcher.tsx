'use client'

import { useState } from 'react'
import { useChainId, useSwitchChain, useChains } from 'wagmi'
import { motion, AnimatePresence } from 'framer-motion'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  ChevronDown,
  Check,
  Loader2,
  AlertTriangle,
  Zap,
  Shield,
  DollarSign,
  Search,
  Star,
  StarOff,
  TrendingUp,
  Activity,
  Clock
} from 'lucide-react'
import { SUPPORTED_CHAINS, CHAIN_CATEGORIES, getChainsByCategory, getChainStatus } from '@/lib/chains'
import { toast } from 'sonner'

// Enhanced chain switcher with favorites and search
interface ChainSwitcherState {
  searchQuery: string
  activeTab: string
  favoriteChains: number[]
}

export function ChainSwitcher() {
  const [isOpen, setIsOpen] = useState(false)
  const [switchingTo, setSwitchingTo] = useState<number | null>(null)
  const [state, setState] = useState<ChainSwitcherState>({
    searchQuery: '',
    activeTab: 'popular',
    favoriteChains: [1, 137, 42161, 10] // Default favorites
  })

  const currentChainId = useChainId()
  const { switchChain, isPending, error } = useSwitchChain()
  const chains = useChains()

  const currentChain = SUPPORTED_CHAINS[currentChainId]
  const supportedChains = chains.filter(chain => SUPPORTED_CHAINS[chain.id])

  const handleChainSwitch = async (chainId: number) => {
    if (chainId === currentChainId) {
      setIsOpen(false)
      return
    }

    setSwitchingTo(chainId)

    try {
      await switchChain({ chainId })
      toast.success(`Switched to ${SUPPORTED_CHAINS[chainId]?.shortName}`)
      setIsOpen(false)
    } catch (err) {
      console.error('Chain switch error:', err)
      toast.error(`Failed to switch to ${SUPPORTED_CHAINS[chainId]?.shortName}`)
    } finally {
      setSwitchingTo(null)
    }
  }

  const toggleFavorite = (chainId: number) => {
    setState(prev => ({
      ...prev,
      favoriteChains: prev.favoriteChains.includes(chainId)
        ? prev.favoriteChains.filter(id => id !== chainId)
        : [...prev.favoriteChains, chainId]
    }))
  }

  const getFilteredChains = () => {
    let chainsToShow = supportedChains

    // Filter by search query
    if (state.searchQuery) {
      chainsToShow = chainsToShow.filter(chain => {
        const chainInfo = SUPPORTED_CHAINS[chain.id]
        return chainInfo?.name.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
               chainInfo?.shortName.toLowerCase().includes(state.searchQuery.toLowerCase())
      })
    }

    // Filter by tab
    switch (state.activeTab) {
      case 'favorites':
        chainsToShow = chainsToShow.filter(chain => state.favoriteChains.includes(chain.id))
        break
      case 'mainnet':
        chainsToShow = chainsToShow.filter(chain => !SUPPORTED_CHAINS[chain.id]?.isTestnet)
        break
      case 'layer2':
        chainsToShow = chainsToShow.filter(chain => SUPPORTED_CHAINS[chain.id]?.category === 'layer2')
        break
      case 'testnet':
        chainsToShow = chainsToShow.filter(chain => SUPPORTED_CHAINS[chain.id]?.isTestnet)
        break
      case 'popular':
      default:
        // Show popular chains first
        chainsToShow = chainsToShow.sort((a, b) => {
          const aInfo = SUPPORTED_CHAINS[a.id]
          const bInfo = SUPPORTED_CHAINS[b.id]
          const popularOrder = [1, 137, 42161, 10, 8453, 43114, 56]
          const aIndex = popularOrder.indexOf(a.id)
          const bIndex = popularOrder.indexOf(b.id)
          if (aIndex === -1 && bIndex === -1) return 0
          if (aIndex === -1) return 1
          if (bIndex === -1) return -1
          return aIndex - bIndex
        })
        break
    }

    return chainsToShow
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy': return 'text-green-500'
      case 'congested': return 'text-yellow-500'
      case 'degraded': return 'text-orange-500'
      case 'maintenance': return 'text-red-500'
      default: return 'text-gray-500'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'healthy': return <Shield className="w-3 h-3" />
      case 'congested': return <Activity className="w-3 h-3" />
      case 'degraded': return <AlertTriangle className="w-3 h-3" />
      case 'maintenance': return <AlertTriangle className="w-3 h-3" />
      default: return <div className="w-3 h-3" />
    }
  }

  const getChainIcon = (chainInfo: any) => {
    return (
      <div className={`w-4 h-4 rounded-full flex items-center justify-center text-xs ${chainInfo.color}`}>
        {chainInfo.icon}
      </div>
    )
  }

  if (!currentChain) {
    return (
      <Button variant="outline" size="sm" disabled>
        <AlertTriangle className="w-4 h-4 mr-2" />
        Unsupported Network
      </Button>
    )
  }

  const filteredChains = getFilteredChains()

  return (
    <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          {getChainIcon(currentChain)}
          <span className="hidden sm:inline">{currentChain.shortName}</span>
          {currentChain.isTestnet && (
            <Badge variant="secondary" className="text-xs">
              Testnet
            </Badge>
          )}
          <ChevronDown className="w-3 h-3" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-96 max-h-[600px] overflow-hidden">
        <DropdownMenuLabel className="flex items-center gap-2">
          <Zap className="w-4 h-4" />
          Switch Network
        </DropdownMenuLabel>
        <DropdownMenuSeparator />

        {/* Search */}
        <div className="p-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search networks..."
              value={state.searchQuery}
              onChange={(e) => setState(prev => ({ ...prev, searchQuery: e.target.value }))}
              className="pl-10"
            />
          </div>
        </div>

        {/* Tabs */}
        <Tabs value={state.activeTab} onValueChange={(value) => setState(prev => ({ ...prev, activeTab: value }))}>
          <div className="px-3">
            <TabsList className="grid w-full grid-cols-5">
              <TabsTrigger value="popular" className="text-xs">Popular</TabsTrigger>
              <TabsTrigger value="favorites" className="text-xs">
                <Star className="w-3 h-3" />
              </TabsTrigger>
              <TabsTrigger value="mainnet" className="text-xs">Mainnet</TabsTrigger>
              <TabsTrigger value="layer2" className="text-xs">L2</TabsTrigger>
              <TabsTrigger value="testnet" className="text-xs">Testnet</TabsTrigger>
            </TabsList>
          </div>

          <div className="max-h-80 overflow-y-auto">
            <TabsContent value={state.activeTab} className="mt-0">
              <div className="space-y-1 p-2">
                {filteredChains.length === 0 ? (
                  <div className="p-4 text-center text-sm text-muted-foreground">
                    No networks found
                  </div>
                ) : (
                  filteredChains.map((chain) => {
                    const info = SUPPORTED_CHAINS[chain.id]
                    const status = getChainStatus(chain.id)
                    const isCurrentChain = chain.id === currentChainId
                    const isSwitching = switchingTo === chain.id
                    const isFavorite = state.favoriteChains.includes(chain.id)

                    return (
                      <div
                        key={chain.id}
                        className="flex items-center justify-between p-3 rounded-lg hover:bg-accent cursor-pointer group"
                        onClick={() => handleChainSwitch(chain.id)}
                      >
                        <div className="flex items-center gap-3 flex-1">
                          <div className="relative">
                            {getChainIcon(info)}
                            {isCurrentChain && (
                              <div className="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full flex items-center justify-center">
                                <Check className="w-2 h-2 text-white" />
                              </div>
                            )}
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2">
                              <span className="font-medium truncate">{info.shortName}</span>
                              {info.isTestnet && (
                                <Badge variant="outline" className="text-xs">
                                  Testnet
                                </Badge>
                              )}
                              <Badge variant="secondary" className="text-xs">
                                {CHAIN_CATEGORIES[info.category]}
                              </Badge>
                            </div>
                            <div className="flex items-center gap-3 text-xs text-muted-foreground">
                              <span className="flex items-center gap-1">
                                <DollarSign className="w-3 h-3" />
                                {info.avgGasPrice}
                              </span>
                              <span className="flex items-center gap-1">
                                <Clock className="w-3 h-3" />
                                {info.blockTime}
                              </span>
                              <span className={`flex items-center gap-1 ${getStatusColor(status)}`}>
                                {getStatusIcon(status)}
                                {status}
                              </span>
                              {info.tvl && (
                                <span className="flex items-center gap-1">
                                  <TrendingUp className="w-3 h-3" />
                                  {info.tvl}
                                </span>
                              )}
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center gap-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-1 opacity-0 group-hover:opacity-100"
                            onClick={(e) => {
                              e.stopPropagation()
                              toggleFavorite(chain.id)
                            }}
                          >
                            {isFavorite ? (
                              <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                            ) : (
                              <StarOff className="w-3 h-3" />
                            )}
                          </Button>
                          {isSwitching && (
                            <Loader2 className="w-4 h-4 animate-spin" />
                          )}
                          {isCurrentChain && !isSwitching && (
                            <Check className="w-4 h-4 text-green-500" />
                          )}
                        </div>
                      </div>
                    )
                  })
                )}
              </div>
            </TabsContent>
          </div>
        </Tabs>

        <DropdownMenuSeparator />

        <div className="p-3 text-xs text-muted-foreground">
          <div className="flex items-center justify-between mb-1">
            <span>Current Network:</span>
            <span className="font-medium">{currentChain.name}</span>
          </div>
          <div className="flex items-center justify-between">
            <span>Gas Token:</span>
            <span className="font-medium">{currentChain.gasToken}</span>
          </div>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
