'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Progress } from '@/components/ui/progress'
import { 
  Users, 
  TrendingUp, 
  TrendingDown, 
  Star, 
  Copy,
  MessageCircle,
  Heart,
  Share2,
  Trophy,
  Target,
  DollarSign,
  Percent,
  Calendar,
  Eye,
  UserPlus,
  Settings
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface Trader {
  id: string
  username: string
  avatar: string
  verified: boolean
  followers: number
  following: number
  totalReturn: number
  monthlyReturn: number
  winRate: number
  totalTrades: number
  copiers: number
  riskScore: number
  maxDrawdown: number
  sharpeRatio: number
  isFollowing: boolean
  isCopying: boolean
  copyAmount: number
}

interface Strategy {
  id: string
  name: string
  description: string
  trader: Trader
  performance: number
  followers: number
  riskLevel: 'low' | 'medium' | 'high'
  minInvestment: number
  maxDrawdown: number
  avgMonthlyReturn: number
  isActive: boolean
}

interface SocialPost {
  id: string
  trader: Trader
  content: string
  timestamp: number
  likes: number
  comments: number
  shares: number
  isLiked: boolean
  trade?: {
    pair: string
    side: 'buy' | 'sell'
    price: number
    amount: number
    pnl?: number
  }
}

export function SocialTradingPlatform() {
  const [topTraders, setTopTraders] = useState<Trader[]>([])
  const [strategies, setStrategies] = useState<Strategy[]>([])
  const [socialFeed, setSocialFeed] = useState<SocialPost[]>([])
  const [activeTab, setActiveTab] = useState('leaderboard')
  const [searchTerm, setSearchTerm] = useState('')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    // Generate mock data
    const mockTraders: Trader[] = [
      {
        id: 'trader1',
        username: 'CryptoKing',
        avatar: '/avatars/trader1.jpg',
        verified: true,
        followers: 15420,
        following: 234,
        totalReturn: 245.8,
        monthlyReturn: 18.5,
        winRate: 72.4,
        totalTrades: 1247,
        copiers: 892,
        riskScore: 65,
        maxDrawdown: -12.3,
        sharpeRatio: 2.1,
        isFollowing: false,
        isCopying: false,
        copyAmount: 0
      },
      {
        id: 'trader2',
        username: 'DeFiMaster',
        avatar: '/avatars/trader2.jpg',
        verified: true,
        followers: 8930,
        following: 156,
        totalReturn: 189.2,
        monthlyReturn: 15.2,
        winRate: 68.9,
        totalTrades: 892,
        copiers: 567,
        riskScore: 58,
        maxDrawdown: -8.7,
        sharpeRatio: 1.8,
        isFollowing: true,
        isCopying: true,
        copyAmount: 5000
      },
      {
        id: 'trader3',
        username: 'YieldHunter',
        avatar: '/avatars/trader3.jpg',
        verified: false,
        followers: 3420,
        following: 89,
        totalReturn: 156.7,
        monthlyReturn: 12.8,
        winRate: 75.2,
        totalTrades: 634,
        copiers: 234,
        riskScore: 42,
        maxDrawdown: -6.2,
        sharpeRatio: 2.3,
        isFollowing: false,
        isCopying: false,
        copyAmount: 0
      }
    ]

    const mockStrategies: Strategy[] = [
      {
        id: 'strategy1',
        name: 'DeFi Yield Maximizer',
        description: 'Automated strategy focusing on high-yield DeFi protocols with risk management',
        trader: mockTraders[1],
        performance: 156.8,
        followers: 1240,
        riskLevel: 'medium',
        minInvestment: 1000,
        maxDrawdown: -8.5,
        avgMonthlyReturn: 12.4,
        isActive: true
      },
      {
        id: 'strategy2',
        name: 'Conservative Staking',
        description: 'Low-risk staking strategy across multiple blue-chip cryptocurrencies',
        trader: mockTraders[2],
        performance: 89.2,
        followers: 892,
        riskLevel: 'low',
        minInvestment: 500,
        maxDrawdown: -4.2,
        avgMonthlyReturn: 6.8,
        isActive: true
      }
    ]

    const mockSocialFeed: SocialPost[] = [
      {
        id: 'post1',
        trader: mockTraders[0],
        content: 'Just opened a long position on ETH. The technical indicators are looking bullish and we might see a breakout above $2600. Risk management is key! 📈',
        timestamp: Date.now() - 3600000,
        likes: 234,
        comments: 45,
        shares: 12,
        isLiked: false,
        trade: {
          pair: 'ETH/USDT',
          side: 'buy',
          price: 2520,
          amount: 5.2
        }
      },
      {
        id: 'post2',
        trader: mockTraders[1],
        content: 'New DeFi yield farming opportunity spotted! 🌾 AAVE lending pool showing 15.2% APY with low impermanent loss risk. Perfect for conservative portfolios.',
        timestamp: Date.now() - 7200000,
        likes: 189,
        comments: 28,
        shares: 34,
        isLiked: true
      },
      {
        id: 'post3',
        trader: mockTraders[2],
        content: 'Closed my BTC position with +8.5% profit. Sometimes it\'s better to take profits early than to be greedy. Patience and discipline always win! 💪',
        timestamp: Date.now() - 10800000,
        likes: 156,
        comments: 19,
        shares: 8,
        isLiked: false,
        trade: {
          pair: 'BTC/USDT',
          side: 'sell',
          price: 45200,
          amount: 0.25,
          pnl: 850
        }
      }
    ]

    setTopTraders(mockTraders)
    setStrategies(mockStrategies)
    setSocialFeed(mockSocialFeed)
  }, [])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTimeAgo = (timestamp: number) => {
    const diff = Date.now() - timestamp
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours < 1) return 'Just now'
    if (hours < 24) return `${hours}h ago`
    return `${Math.floor(hours / 24)}d ago`
  }

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'low': return 'text-green-500'
      case 'medium': return 'text-yellow-500'
      case 'high': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskBadgeVariant = (level: string) => {
    switch (level) {
      case 'low': return 'default'
      case 'medium': return 'secondary'
      case 'high': return 'destructive'
      default: return 'outline'
    }
  }

  const handleFollow = (traderId: string) => {
    setTopTraders(prev => prev.map(trader => 
      trader.id === traderId 
        ? { ...trader, isFollowing: !trader.isFollowing, followers: trader.followers + (trader.isFollowing ? -1 : 1) }
        : trader
    ))
  }

  const handleCopyTrade = (traderId: string) => {
    setTopTraders(prev => prev.map(trader => 
      trader.id === traderId 
        ? { ...trader, isCopying: !trader.isCopying, copiers: trader.copiers + (trader.isCopying ? -1 : 1) }
        : trader
    ))
  }

  const handleLike = (postId: string) => {
    setSocialFeed(prev => prev.map(post => 
      post.id === postId 
        ? { ...post, isLiked: !post.isLiked, likes: post.likes + (post.isLiked ? -1 : 1) }
        : post
    ))
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Users className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Your Wallet</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access social trading features and follow top traders
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
          <h1 className="text-3xl font-bold">Social Trading</h1>
          <p className="text-muted-foreground">Follow, copy, and learn from top traders</p>
        </div>
        <div className="flex items-center gap-2">
          <Input
            placeholder="Search traders..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-64"
          />
          <Button variant="outline">
            <Settings className="w-4 h-4" />
          </Button>
        </div>
      </div>

      {/* Main Platform */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
          <TabsTrigger value="strategies">Strategies</TabsTrigger>
          <TabsTrigger value="feed">Social Feed</TabsTrigger>
          <TabsTrigger value="portfolio">My Copies</TabsTrigger>
        </TabsList>

        <TabsContent value="leaderboard" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
            {topTraders.map((trader, index) => (
              <Card key={trader.id} className="hover:shadow-md transition-shadow">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="relative">
                        <Avatar className="w-12 h-12">
                          <AvatarImage src={trader.avatar} />
                          <AvatarFallback>{trader.username.slice(0, 2).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        {index < 3 && (
                          <div className="absolute -top-1 -right-1 w-6 h-6 bg-yellow-500 rounded-full flex items-center justify-center">
                            <Trophy className="w-3 h-3 text-white" />
                          </div>
                        )}
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <h3 className="font-bold">{trader.username}</h3>
                          {trader.verified && (
                            <Badge variant="default" className="text-xs">
                              ✓ Verified
                            </Badge>
                          )}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {trader.followers.toLocaleString()} followers
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-lg font-bold text-green-500">
                        +{trader.totalReturn}%
                      </div>
                      <div className="text-xs text-muted-foreground">Total Return</div>
                    </div>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Monthly Return</div>
                      <div className="font-medium text-green-500">+{trader.monthlyReturn}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Win Rate</div>
                      <div className="font-medium">{trader.winRate}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Total Trades</div>
                      <div className="font-medium">{trader.totalTrades.toLocaleString()}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Copiers</div>
                      <div className="font-medium">{trader.copiers.toLocaleString()}</div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Risk Score</span>
                      <span className="font-medium">{trader.riskScore}/100</span>
                    </div>
                    <Progress value={trader.riskScore} className="h-2" />
                  </div>

                  <div className="flex gap-2">
                    <Button
                      variant={trader.isFollowing ? "default" : "outline"}
                      size="sm"
                      onClick={() => handleFollow(trader.id)}
                      className="flex-1"
                    >
                      <UserPlus className="w-3 h-3 mr-1" />
                      {trader.isFollowing ? 'Following' : 'Follow'}
                    </Button>
                    <Button
                      variant={trader.isCopying ? "default" : "outline"}
                      size="sm"
                      onClick={() => handleCopyTrade(trader.id)}
                      className="flex-1"
                    >
                      <Copy className="w-3 h-3 mr-1" />
                      {trader.isCopying ? 'Copying' : 'Copy'}
                    </Button>
                  </div>

                  {trader.isCopying && (
                    <div className="p-2 bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800 rounded text-sm">
                      <div className="flex justify-between">
                        <span className="text-blue-700 dark:text-blue-300">Copy Amount:</span>
                        <span className="font-medium">{formatCurrency(trader.copyAmount)}</span>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="strategies" className="space-y-4">
          <div className="space-y-4">
            {strategies.map((strategy) => (
              <Card key={strategy.id} className="hover:shadow-md transition-shadow">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-start gap-4">
                      <Avatar className="w-12 h-12">
                        <AvatarImage src={strategy.trader.avatar} />
                        <AvatarFallback>{strategy.trader.username.slice(0, 2).toUpperCase()}</AvatarFallback>
                      </Avatar>
                      <div>
                        <h3 className="font-bold text-lg">{strategy.name}</h3>
                        <p className="text-sm text-muted-foreground mb-2">{strategy.description}</p>
                        <div className="flex items-center gap-2">
                          <span className="text-sm text-muted-foreground">by</span>
                          <span className="font-medium">{strategy.trader.username}</span>
                          {strategy.trader.verified && (
                            <Badge variant="default" className="text-xs">✓</Badge>
                          )}
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-2xl font-bold text-green-500">
                        +{strategy.performance}%
                      </div>
                      <div className="text-sm text-muted-foreground">Performance</div>
                    </div>
                  </div>

                  <div className="grid grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Followers</div>
                      <div className="font-medium">{strategy.followers.toLocaleString()}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Risk Level</div>
                      <Badge variant={getRiskBadgeVariant(strategy.riskLevel)} className="text-xs">
                        {strategy.riskLevel}
                      </Badge>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Min Investment</div>
                      <div className="font-medium">{formatCurrency(strategy.minInvestment)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Monthly Avg</div>
                      <div className="font-medium text-green-500">+{strategy.avgMonthlyReturn}%</div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <span>Max Drawdown: {strategy.maxDrawdown}%</span>
                      <Badge variant={strategy.isActive ? "default" : "secondary"}>
                        {strategy.isActive ? "Active" : "Inactive"}
                      </Badge>
                    </div>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm">
                        <Eye className="w-3 h-3 mr-1" />
                        View Details
                      </Button>
                      <Button size="sm">
                        <Copy className="w-3 h-3 mr-1" />
                        Copy Strategy
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="feed" className="space-y-4">
          <div className="space-y-4">
            {socialFeed.map((post) => (
              <Card key={post.id}>
                <CardContent className="p-6">
                  <div className="flex items-start gap-4">
                    <Avatar className="w-10 h-10">
                      <AvatarImage src={post.trader.avatar} />
                      <AvatarFallback>{post.trader.username.slice(0, 2).toUpperCase()}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="font-medium">{post.trader.username}</span>
                        {post.trader.verified && (
                          <Badge variant="default" className="text-xs">✓</Badge>
                        )}
                        <span className="text-sm text-muted-foreground">
                          {formatTimeAgo(post.timestamp)}
                        </span>
                      </div>
                      
                      <p className="text-sm mb-3">{post.content}</p>

                      {post.trade && (
                        <div className="p-3 bg-muted/50 rounded-lg mb-3">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <Badge variant={post.trade.side === 'buy' ? 'default' : 'secondary'}>
                                {post.trade.side.toUpperCase()}
                              </Badge>
                              <span className="font-medium">{post.trade.pair}</span>
                            </div>
                            <div className="text-right">
                              <div className="font-medium">${post.trade.price.toLocaleString()}</div>
                              <div className="text-sm text-muted-foreground">
                                {post.trade.amount} {post.trade.pair.split('/')[0]}
                              </div>
                            </div>
                          </div>
                          {post.trade.pnl && (
                            <div className="mt-2 text-sm">
                              <span className="text-muted-foreground">P&L: </span>
                              <span className={cn(
                                "font-medium",
                                post.trade.pnl >= 0 ? "text-green-500" : "text-red-500"
                              )}>
                                {post.trade.pnl >= 0 ? '+' : ''}{formatCurrency(post.trade.pnl)}
                              </span>
                            </div>
                          )}
                        </div>
                      )}

                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        <button
                          onClick={() => handleLike(post.id)}
                          className={cn(
                            "flex items-center gap-1 hover:text-red-500 transition-colors",
                            post.isLiked && "text-red-500"
                          )}
                        >
                          <Heart className={cn("w-4 h-4", post.isLiked && "fill-current")} />
                          {post.likes}
                        </button>
                        <button className="flex items-center gap-1 hover:text-blue-500 transition-colors">
                          <MessageCircle className="w-4 h-4" />
                          {post.comments}
                        </button>
                        <button className="flex items-center gap-1 hover:text-green-500 transition-colors">
                          <Share2 className="w-4 h-4" />
                          {post.shares}
                        </button>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="portfolio" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Copy className="w-5 h-5" />
                  Active Copy Trades
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {topTraders.filter(t => t.isCopying).map((trader) => (
                    <div key={trader.id} className="flex items-center justify-between p-3 border rounded">
                      <div className="flex items-center gap-3">
                        <Avatar className="w-8 h-8">
                          <AvatarImage src={trader.avatar} />
                          <AvatarFallback>{trader.username.slice(0, 2).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="font-medium">{trader.username}</div>
                          <div className="text-sm text-muted-foreground">
                            {formatCurrency(trader.copyAmount)} allocated
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-medium text-green-500">+{trader.monthlyReturn}%</div>
                        <Button variant="outline" size="sm" className="mt-1">
                          Manage
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Target className="w-5 h-5" />
                  Copy Trading Stats
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Total Allocated</span>
                    <span className="font-medium">
                      {formatCurrency(topTraders.filter(t => t.isCopying).reduce((sum, t) => sum + t.copyAmount, 0))}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Active Copies</span>
                    <span className="font-medium">{topTraders.filter(t => t.isCopying).length}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Total Return</span>
                    <span className="font-medium text-green-500">+12.8%</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Best Performer</span>
                    <span className="font-medium">DeFiMaster</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
