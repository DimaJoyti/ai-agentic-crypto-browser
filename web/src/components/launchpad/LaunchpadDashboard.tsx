'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Rocket, 
  TrendingUp,
  Clock,
  Users,
  DollarSign,
  Target,
  Star,
  Calendar,
  CheckCircle,
  AlertTriangle,
  Zap,
  Shield,
  Globe,
  Award,
  Eye,
  Heart,
  Share2,
  Filter,
  Search
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface LaunchpadProject {
  id: string
  name: string
  symbol: string
  description: string
  logo: string
  website: string
  twitter: string
  telegram: string
  category: 'defi' | 'gaming' | 'nft' | 'infrastructure' | 'dao' | 'metaverse'
  status: 'upcoming' | 'live' | 'ended' | 'completed'
  launchType: 'ido' | 'ino' | 'fair_launch' | 'dutch_auction' | 'lottery'
  
  // Financial Details
  totalRaise: number
  targetRaise: number
  tokenPrice: number
  totalSupply: number
  tokensForSale: number
  minAllocation: number
  maxAllocation: number
  
  // Timeline
  registrationStart: number
  registrationEnd: number
  saleStart: number
  saleEnd: number
  vestingStart: number
  
  // Progress
  currentRaised: number
  participantCount: number
  
  // Requirements
  kycRequired: boolean
  whitelistRequired: boolean
  minTier: number
  
  // Vesting
  vestingSchedule: Array<{
    percentage: number
    unlockDate: number
    description: string
  }>
  
  // Social Metrics
  likes: number
  views: number
  isLiked: boolean
  isWatchlisted: boolean
  
  // Risk Assessment
  riskScore: number
  auditStatus: 'pending' | 'in_progress' | 'completed' | 'failed'
  teamVerified: boolean
}

interface UserTier {
  level: number
  name: string
  requirements: {
    stakingAmount: number
    holdingPeriod: number
  }
  benefits: {
    allocationMultiplier: number
    earlyAccess: boolean
    feeDiscount: number
    guaranteedAllocation: boolean
  }
}

export function LaunchpadDashboard() {
  const [projects, setProjects] = useState<LaunchpadProject[]>([])
  const [userTier, setUserTier] = useState<UserTier | null>(null)
  const [activeTab, setActiveTab] = useState('live')
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const [searchQuery, setSearchQuery] = useState('')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock launchpad data
    const mockProjects: LaunchpadProject[] = [
      {
        id: 'project1',
        name: 'DeFiMax Protocol',
        symbol: 'DMAX',
        description: 'Next-generation DeFi protocol with automated yield optimization and cross-chain liquidity aggregation',
        logo: '/projects/defimax.png',
        website: 'https://defimax.io',
        twitter: '@defimax',
        telegram: 't.me/defimax',
        category: 'defi',
        status: 'live',
        launchType: 'ido',
        totalRaise: 2000000,
        targetRaise: 1500000,
        tokenPrice: 0.15,
        totalSupply: 100000000,
        tokensForSale: 10000000,
        minAllocation: 100,
        maxAllocation: 5000,
        registrationStart: Date.now() - 86400000 * 7,
        registrationEnd: Date.now() - 86400000 * 2,
        saleStart: Date.now() - 86400000,
        saleEnd: Date.now() + 86400000 * 2,
        vestingStart: Date.now() + 86400000 * 7,
        currentRaised: 1250000,
        participantCount: 2847,
        kycRequired: true,
        whitelistRequired: true,
        minTier: 2,
        vestingSchedule: [
          { percentage: 20, unlockDate: Date.now() + 86400000 * 7, description: 'TGE Unlock' },
          { percentage: 20, unlockDate: Date.now() + 86400000 * 37, description: '1 Month' },
          { percentage: 20, unlockDate: Date.now() + 86400000 * 67, description: '2 Months' },
          { percentage: 20, unlockDate: Date.now() + 86400000 * 97, description: '3 Months' },
          { percentage: 20, unlockDate: Date.now() + 86400000 * 127, description: '4 Months' }
        ],
        likes: 1247,
        views: 15632,
        isLiked: false,
        isWatchlisted: true,
        riskScore: 25,
        auditStatus: 'completed',
        teamVerified: true
      },
      {
        id: 'project2',
        name: 'MetaVerse Worlds',
        symbol: 'MVW',
        description: 'Immersive metaverse platform with play-to-earn gaming and virtual real estate',
        logo: '/projects/metaverse.png',
        website: 'https://metaverseworlds.io',
        twitter: '@metaverseworlds',
        telegram: 't.me/metaverseworlds',
        category: 'metaverse',
        status: 'upcoming',
        launchType: 'fair_launch',
        totalRaise: 5000000,
        targetRaise: 3000000,
        tokenPrice: 0.08,
        totalSupply: 500000000,
        tokensForSale: 37500000,
        minAllocation: 50,
        maxAllocation: 2500,
        registrationStart: Date.now() + 86400000 * 3,
        registrationEnd: Date.now() + 86400000 * 10,
        saleStart: Date.now() + 86400000 * 12,
        saleEnd: Date.now() + 86400000 * 15,
        vestingStart: Date.now() + 86400000 * 20,
        currentRaised: 0,
        participantCount: 0,
        kycRequired: true,
        whitelistRequired: false,
        minTier: 1,
        vestingSchedule: [
          { percentage: 25, unlockDate: Date.now() + 86400000 * 20, description: 'TGE Unlock' },
          { percentage: 25, unlockDate: Date.now() + 86400000 * 50, description: '1 Month' },
          { percentage: 25, unlockDate: Date.now() + 86400000 * 80, description: '2 Months' },
          { percentage: 25, unlockDate: Date.now() + 86400000 * 110, description: '3 Months' }
        ],
        likes: 892,
        views: 8934,
        isLiked: true,
        isWatchlisted: false,
        riskScore: 35,
        auditStatus: 'in_progress',
        teamVerified: true
      },
      {
        id: 'project3',
        name: 'GameFi Arena',
        symbol: 'GFA',
        description: 'Competitive gaming platform with NFT rewards and tournament infrastructure',
        logo: '/projects/gamefi.png',
        website: 'https://gamefiarena.io',
        twitter: '@gamefiarena',
        telegram: 't.me/gamefiarena',
        category: 'gaming',
        status: 'ended',
        launchType: 'lottery',
        totalRaise: 1000000,
        targetRaise: 800000,
        tokenPrice: 0.25,
        totalSupply: 200000000,
        tokensForSale: 3200000,
        minAllocation: 200,
        maxAllocation: 1000,
        registrationStart: Date.now() - 86400000 * 20,
        registrationEnd: Date.now() - 86400000 * 15,
        saleStart: Date.now() - 86400000 * 12,
        saleEnd: Date.now() - 86400000 * 10,
        vestingStart: Date.now() - 86400000 * 5,
        currentRaised: 950000,
        participantCount: 4521,
        kycRequired: false,
        whitelistRequired: true,
        minTier: 3,
        vestingSchedule: [
          { percentage: 30, unlockDate: Date.now() - 86400000 * 5, description: 'TGE Unlock' },
          { percentage: 35, unlockDate: Date.now() + 86400000 * 25, description: '1 Month' },
          { percentage: 35, unlockDate: Date.now() + 86400000 * 55, description: '2 Months' }
        ],
        likes: 2156,
        views: 23847,
        isLiked: true,
        isWatchlisted: true,
        riskScore: 15,
        auditStatus: 'completed',
        teamVerified: true
      }
    ]

    const mockUserTier: UserTier = {
      level: 3,
      name: 'Gold Tier',
      requirements: {
        stakingAmount: 10000,
        holdingPeriod: 30
      },
      benefits: {
        allocationMultiplier: 2.5,
        earlyAccess: true,
        feeDiscount: 25,
        guaranteedAllocation: true
      }
    }

    setProjects(mockProjects)
    setUserTier(mockUserTier)
  }, [isConnected])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'live': return 'text-green-500'
      case 'upcoming': return 'text-blue-500'
      case 'ended': return 'text-yellow-500'
      case 'completed': return 'text-gray-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'live': return 'default'
      case 'upcoming': return 'secondary'
      case 'ended': case 'completed': return 'outline'
      default: return 'outline'
    }
  }

  const getRiskScoreColor = (score: number) => {
    if (score <= 20) return 'text-green-500'
    if (score <= 40) return 'text-yellow-500'
    return 'text-red-500'
  }

  const getAuditStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'text-green-500'
      case 'in_progress': return 'text-yellow-500'
      case 'pending': return 'text-gray-500'
      case 'failed': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const filteredProjects = projects.filter(project => {
    const matchesTab = activeTab === 'all' || project.status === activeTab
    const matchesCategory = selectedCategory === 'all' || project.category === selectedCategory
    const matchesSearch = project.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         project.symbol.toLowerCase().includes(searchQuery.toLowerCase())
    
    return matchesTab && matchesCategory && matchesSearch
  })

  const toggleLike = (projectId: string) => {
    setProjects(prev => prev.map(project => 
      project.id === projectId 
        ? { 
            ...project, 
            isLiked: !project.isLiked,
            likes: project.isLiked ? project.likes - 1 : project.likes + 1
          }
        : project
    ))
  }

  const toggleWatchlist = (projectId: string) => {
    setProjects(prev => prev.map(project => 
      project.id === projectId 
        ? { ...project, isWatchlisted: !project.isWatchlisted }
        : project
    ))
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Rocket className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access the launchpad and participate in token launches
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
          <h2 className="text-2xl font-bold">Launchpad</h2>
          <p className="text-muted-foreground">
            Discover and participate in the next generation of crypto projects
          </p>
        </div>
        <div className="flex items-center gap-2">
          {userTier && (
            <Badge variant="outline">
              <Award className="w-3 h-3 mr-1" />
              {userTier.name}
            </Badge>
          )}
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            KYC Verified
          </Badge>
        </div>
      </div>

      {/* User Tier Info */}
      {userTier && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="w-5 h-5" />
              Your Tier: {userTier.name}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-primary">{userTier.benefits.allocationMultiplier}x</div>
                <div className="text-sm text-muted-foreground">Allocation Multiplier</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-green-500">{userTier.benefits.feeDiscount}%</div>
                <div className="text-sm text-muted-foreground">Fee Discount</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold">
                  {userTier.benefits.earlyAccess ? (
                    <CheckCircle className="w-8 h-8 text-green-500 mx-auto" />
                  ) : (
                    <AlertTriangle className="w-8 h-8 text-yellow-500 mx-auto" />
                  )}
                </div>
                <div className="text-sm text-muted-foreground">Early Access</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold">
                  {userTier.benefits.guaranteedAllocation ? (
                    <CheckCircle className="w-8 h-8 text-green-500 mx-auto" />
                  ) : (
                    <AlertTriangle className="w-8 h-8 text-yellow-500 mx-auto" />
                  )}
                </div>
                <div className="text-sm text-muted-foreground">Guaranteed Allocation</div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Filters and Search */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
            <input
              type="text"
              placeholder="Search projects..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>
        </div>
        <select
          value={selectedCategory}
          onChange={(e) => setSelectedCategory(e.target.value)}
          className="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
        >
          <option value="all">All Categories</option>
          <option value="defi">DeFi</option>
          <option value="gaming">Gaming</option>
          <option value="nft">NFT</option>
          <option value="metaverse">Metaverse</option>
          <option value="infrastructure">Infrastructure</option>
          <option value="dao">DAO</option>
        </select>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="live">Live</TabsTrigger>
          <TabsTrigger value="upcoming">Upcoming</TabsTrigger>
          <TabsTrigger value="ended">Ended</TabsTrigger>
          <TabsTrigger value="all">All Projects</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredProjects.map((project) => (
              <Card key={project.id} className="overflow-hidden">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                        <Rocket className="w-6 h-6" />
                      </div>
                      <div>
                        <h3 className="font-bold">{project.name}</h3>
                        <p className="text-sm text-muted-foreground">{project.symbol}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-1">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => toggleLike(project.id)}
                        className="p-1"
                      >
                        <Heart className={cn(
                          "w-4 h-4",
                          project.isLiked ? "fill-red-500 text-red-500" : "text-muted-foreground"
                        )} />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => toggleWatchlist(project.id)}
                        className="p-1"
                      >
                        <Star className={cn(
                          "w-4 h-4",
                          project.isWatchlisted ? "fill-yellow-500 text-yellow-500" : "text-muted-foreground"
                        )} />
                      </Button>
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-2 mt-2">
                    <Badge variant={getStatusBadgeVariant(project.status)}>
                      {project.status}
                    </Badge>
                    <Badge variant="outline" className="capitalize">
                      {project.launchType.replace('_', ' ')}
                    </Badge>
                    <Badge variant="outline" className="capitalize">
                      {project.category}
                    </Badge>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {project.description}
                  </p>

                  {/* Progress */}
                  {project.status === 'live' && (
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>Progress</span>
                        <span>{((project.currentRaised / project.targetRaise) * 100).toFixed(1)}%</span>
                      </div>
                      <Progress value={(project.currentRaised / project.targetRaise) * 100} className="h-2" />
                      <div className="flex justify-between text-xs text-muted-foreground">
                        <span>{formatCurrency(project.currentRaised)} raised</span>
                        <span>Goal: {formatCurrency(project.targetRaise)}</span>
                      </div>
                    </div>
                  )}

                  {/* Key Metrics */}
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Token Price</div>
                      <div className="font-medium">${project.tokenPrice}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Participants</div>
                      <div className="font-medium">{project.participantCount.toLocaleString()}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Min/Max</div>
                      <div className="font-medium">
                        ${project.minAllocation} - ${project.maxAllocation}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Risk Score</div>
                      <div className={cn("font-medium", getRiskScoreColor(project.riskScore))}>
                        {project.riskScore}/100
                      </div>
                    </div>
                  </div>

                  {/* Timeline */}
                  <div className="space-y-2">
                    <div className="text-sm font-medium">Timeline</div>
                    <div className="text-xs space-y-1">
                      {project.status === 'upcoming' && (
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Registration starts:</span>
                          <span>{formatTime(project.registrationStart)}</span>
                        </div>
                      )}
                      {project.status === 'live' && (
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Sale ends:</span>
                          <span>{formatTime(project.saleEnd)}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Vesting starts:</span>
                        <span>{formatTime(project.vestingStart)}</span>
                      </div>
                    </div>
                  </div>

                  {/* Requirements */}
                  <div className="flex flex-wrap gap-2">
                    {project.kycRequired && (
                      <Badge variant="outline" className="text-xs">
                        <Shield className="w-3 h-3 mr-1" />
                        KYC Required
                      </Badge>
                    )}
                    {project.whitelistRequired && (
                      <Badge variant="outline" className="text-xs">
                        <Users className="w-3 h-3 mr-1" />
                        Whitelist
                      </Badge>
                    )}
                    <Badge variant="outline" className="text-xs">
                      <Award className="w-3 h-3 mr-1" />
                      Tier {project.minTier}+
                    </Badge>
                  </div>

                  {/* Action Button */}
                  <Button className="w-full" size="sm">
                    {project.status === 'upcoming' ? 'Register Interest' :
                     project.status === 'live' ? 'Participate Now' :
                     project.status === 'ended' ? 'View Results' :
                     'View Details'}
                  </Button>

                  {/* Social Metrics */}
                  <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <div className="flex items-center gap-3">
                      <span className="flex items-center gap-1">
                        <Eye className="w-3 h-3" />
                        {project.views.toLocaleString()}
                      </span>
                      <span className="flex items-center gap-1">
                        <Heart className="w-3 h-3" />
                        {project.likes.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      {project.teamVerified && (
                        <CheckCircle className="w-3 h-3 text-green-500" />
                      )}
                      <span className={getAuditStatusColor(project.auditStatus)}>
                        {project.auditStatus === 'completed' ? 'Audited' : 
                         project.auditStatus === 'in_progress' ? 'Auditing' : 
                         'Pending Audit'}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {filteredProjects.length === 0 && (
            <div className="text-center py-12">
              <Rocket className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
              <h3 className="text-lg font-medium mb-2">No Projects Found</h3>
              <p className="text-muted-foreground">
                Try adjusting your filters or search criteria
              </p>
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
