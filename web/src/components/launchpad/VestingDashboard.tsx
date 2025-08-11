'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Lock,
  Unlock,
  Clock,
  Calendar,
  DollarSign,
  TrendingUp,
  Award,
  CheckCircle,
  AlertTriangle,
  Download,
  ExternalLink,
  BarChart3,
  Zap,
  Target,
  Gift,
  History,
  Wallet
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface VestingSchedule {
  id: string
  projectName: string
  projectSymbol: string
  projectLogo: string
  totalTokens: number
  claimedTokens: number
  nextUnlockAmount: number
  nextUnlockDate: number
  vestingPeriods: Array<{
    id: string
    percentage: number
    amount: number
    unlockDate: number
    description: string
    isUnlocked: boolean
    isClaimed: boolean
    claimedAt?: number
    transactionHash?: string
  }>
  tokenPrice: number
  currentPrice: number
  totalValue: number
  unrealizedGains: number
  status: 'active' | 'completed' | 'cancelled'
  vestingType: 'linear' | 'cliff' | 'milestone'
  cliffPeriod?: number
  linearDuration?: number
}

interface ClaimHistory {
  id: string
  projectName: string
  projectSymbol: string
  amount: number
  value: number
  claimedAt: number
  transactionHash: string
  gasUsed: number
  gasFee: number
}

interface VestingStats {
  totalProjects: number
  totalTokensVesting: number
  totalValueLocked: number
  totalClaimable: number
  totalClaimed: number
  nextUnlockValue: number
  nextUnlockDate: number
  averageAPY: number
}

export function VestingDashboard() {
  const [vestingSchedules, setVestingSchedules] = useState<VestingSchedule[]>([])
  const [claimHistory, setClaimHistory] = useState<ClaimHistory[]>([])
  const [vestingStats, setVestingStats] = useState<VestingStats | null>(null)
  const [selectedProject, setSelectedProject] = useState<string | null>(null)

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock vesting data
    const mockVestingSchedules: VestingSchedule[] = [
      {
        id: 'vesting1',
        projectName: 'DeFiMax Protocol',
        projectSymbol: 'DMAX',
        projectLogo: '/projects/defimax.png',
        totalTokens: 6666.67,
        claimedTokens: 1333.33,
        nextUnlockAmount: 1333.33,
        nextUnlockDate: Date.now() + 86400000 * 7,
        tokenPrice: 0.15,
        currentPrice: 0.28,
        totalValue: 1866.67,
        unrealizedGains: 866.67,
        status: 'active',
        vestingType: 'linear',
        linearDuration: 120,
        vestingPeriods: [
          {
            id: 'period1',
            percentage: 20,
            amount: 1333.33,
            unlockDate: Date.now() - 86400000 * 7,
            description: 'TGE Unlock',
            isUnlocked: true,
            isClaimed: true,
            claimedAt: Date.now() - 86400000 * 6,
            transactionHash: '0x1234...5678'
          },
          {
            id: 'period2',
            percentage: 20,
            amount: 1333.33,
            unlockDate: Date.now() + 86400000 * 7,
            description: '1 Month',
            isUnlocked: false,
            isClaimed: false
          },
          {
            id: 'period3',
            percentage: 20,
            amount: 1333.33,
            unlockDate: Date.now() + 86400000 * 37,
            description: '2 Months',
            isUnlocked: false,
            isClaimed: false
          },
          {
            id: 'period4',
            percentage: 20,
            amount: 1333.33,
            unlockDate: Date.now() + 86400000 * 67,
            description: '3 Months',
            isUnlocked: false,
            isClaimed: false
          },
          {
            id: 'period5',
            percentage: 20,
            amount: 1333.33,
            unlockDate: Date.now() + 86400000 * 97,
            description: '4 Months',
            isUnlocked: false,
            isClaimed: false
          }
        ]
      },
      {
        id: 'vesting2',
        projectName: 'GameFi Arena',
        projectSymbol: 'GFA',
        projectLogo: '/projects/gamefi.png',
        totalTokens: 4000,
        claimedTokens: 1200,
        nextUnlockAmount: 1400,
        nextUnlockDate: Date.now() + 86400000 * 25,
        tokenPrice: 0.25,
        currentPrice: 0.42,
        totalValue: 1680,
        unrealizedGains: 680,
        status: 'active',
        vestingType: 'cliff',
        cliffPeriod: 30,
        vestingPeriods: [
          {
            id: 'period1',
            percentage: 30,
            amount: 1200,
            unlockDate: Date.now() - 86400000 * 5,
            description: 'TGE Unlock',
            isUnlocked: true,
            isClaimed: true,
            claimedAt: Date.now() - 86400000 * 4,
            transactionHash: '0xabcd...efgh'
          },
          {
            id: 'period2',
            percentage: 35,
            amount: 1400,
            unlockDate: Date.now() + 86400000 * 25,
            description: '1 Month Cliff',
            isUnlocked: false,
            isClaimed: false
          },
          {
            id: 'period3',
            percentage: 35,
            amount: 1400,
            unlockDate: Date.now() + 86400000 * 55,
            description: '2 Months',
            isUnlocked: false,
            isClaimed: false
          }
        ]
      }
    ]

    const mockClaimHistory: ClaimHistory[] = [
      {
        id: 'claim1',
        projectName: 'DeFiMax Protocol',
        projectSymbol: 'DMAX',
        amount: 1333.33,
        value: 373.33,
        claimedAt: Date.now() - 86400000 * 6,
        transactionHash: '0x1234567890abcdef',
        gasUsed: 45000,
        gasFee: 12.5
      },
      {
        id: 'claim2',
        projectName: 'GameFi Arena',
        projectSymbol: 'GFA',
        amount: 1200,
        value: 504,
        claimedAt: Date.now() - 86400000 * 4,
        transactionHash: '0xabcdef1234567890',
        gasUsed: 38000,
        gasFee: 8.7
      }
    ]

    const mockVestingStats: VestingStats = {
      totalProjects: 2,
      totalTokensVesting: 10666.67,
      totalValueLocked: 3546.67,
      totalClaimable: 0,
      totalClaimed: 2533.33,
      nextUnlockValue: 373.33,
      nextUnlockDate: Date.now() + 86400000 * 7,
      averageAPY: 145.2
    }

    setVestingSchedules(mockVestingSchedules)
    setClaimHistory(mockClaimHistory)
    setVestingStats(mockVestingStats)
  }, [isConnected])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString()
  }

  const formatTimeDetailed = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getTimeUntilUnlock = (timestamp: number) => {
    const now = Date.now()
    const diff = timestamp - now
    
    if (diff <= 0) return 'Available now'
    
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    
    if (days > 0) return `${days}d ${hours}h`
    return `${hours}h`
  }

  const handleClaim = (vestingId: string, periodId: string) => {
    setVestingSchedules(prev => prev.map(vesting => 
      vesting.id === vestingId 
        ? {
            ...vesting,
            vestingPeriods: vesting.vestingPeriods.map(period =>
              period.id === periodId
                ? { ...period, isClaimed: true, claimedAt: Date.now(), transactionHash: '0x' + Math.random().toString(16).substr(2, 8) }
                : period
            ),
            claimedTokens: vesting.claimedTokens + vesting.vestingPeriods.find(p => p.id === periodId)!.amount
          }
        : vesting
    ))
  }

  const getClaimableAmount = (vesting: VestingSchedule) => {
    return vesting.vestingPeriods
      .filter(period => period.isUnlocked && !period.isClaimed)
      .reduce((sum, period) => sum + period.amount, 0)
  }

  const getClaimableValue = (vesting: VestingSchedule) => {
    return getClaimableAmount(vesting) * vesting.currentPrice
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Lock className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to view your vesting schedules and claim tokens
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
          <h2 className="text-2xl font-bold">Vesting Dashboard</h2>
          <p className="text-muted-foreground">
            Track and claim your vested tokens from launchpad participations
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Gift className="w-3 h-3 mr-1" />
            {vestingStats?.totalProjects} Projects
          </Badge>
          <Button variant="outline" size="sm">
            <Download className="w-4 h-4 mr-2" />
            Export Report
          </Button>
        </div>
      </div>

      {/* Vesting Stats */}
      {vestingStats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <Wallet className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Total Value Locked</span>
              </div>
              <div className="text-2xl font-bold">{formatCurrency(vestingStats.totalValueLocked)}</div>
              <div className="text-xs text-muted-foreground">
                {vestingStats.totalTokensVesting.toFixed(0)} tokens
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <Unlock className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Claimable Now</span>
              </div>
              <div className="text-2xl font-bold text-green-500">
                {formatCurrency(vestingStats.totalClaimable)}
              </div>
              <div className="text-xs text-muted-foreground">
                Ready to claim
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <Clock className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Next Unlock</span>
              </div>
              <div className="text-2xl font-bold">{formatCurrency(vestingStats.nextUnlockValue)}</div>
              <div className="text-xs text-muted-foreground">
                {getTimeUntilUnlock(vestingStats.nextUnlockDate)}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <TrendingUp className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">Average APY</span>
              </div>
              <div className="text-2xl font-bold text-green-500">
                {vestingStats.averageAPY.toFixed(1)}%
              </div>
              <div className="text-xs text-muted-foreground">
                Unrealized gains
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Vesting Schedules */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Active Vesting Schedules</h3>
        
        {vestingSchedules.map((vesting) => (
          <Card key={vesting.id}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                    <Award className="w-6 h-6" />
                  </div>
                  <div>
                    <h4 className="font-bold">{vesting.projectName}</h4>
                    <p className="text-sm text-muted-foreground">{vesting.projectSymbol}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className="capitalize">
                    {vesting.vestingType}
                  </Badge>
                  <Badge variant={vesting.status === 'active' ? 'default' : 'secondary'}>
                    {vesting.status}
                  </Badge>
                </div>
              </div>
            </CardHeader>

            <CardContent className="space-y-4">
              {/* Overview Stats */}
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                  <div className="text-muted-foreground">Total Tokens</div>
                  <div className="font-bold">{vesting.totalTokens.toFixed(2)} {vesting.projectSymbol}</div>
                </div>
                <div>
                  <div className="text-muted-foreground">Claimed</div>
                  <div className="font-bold">{vesting.claimedTokens.toFixed(2)} {vesting.projectSymbol}</div>
                </div>
                <div>
                  <div className="text-muted-foreground">Current Value</div>
                  <div className="font-bold">{formatCurrency(vesting.totalValue)}</div>
                </div>
                <div>
                  <div className="text-muted-foreground">Unrealized P&L</div>
                  <div className={cn(
                    "font-bold",
                    vesting.unrealizedGains >= 0 ? "text-green-500" : "text-red-500"
                  )}>
                    {vesting.unrealizedGains >= 0 ? '+' : ''}{formatCurrency(vesting.unrealizedGains)}
                  </div>
                </div>
              </div>

              {/* Progress Bar */}
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Vesting Progress</span>
                  <span>{((vesting.claimedTokens / vesting.totalTokens) * 100).toFixed(1)}%</span>
                </div>
                <Progress value={(vesting.claimedTokens / vesting.totalTokens) * 100} className="h-2" />
              </div>

              {/* Claimable Amount */}
              {getClaimableAmount(vesting) > 0 && (
                <Alert>
                  <Unlock className="h-4 w-4" />
                  <AlertDescription>
                    You have {getClaimableAmount(vesting).toFixed(2)} {vesting.projectSymbol} 
                    ({formatCurrency(getClaimableValue(vesting))}) ready to claim!
                  </AlertDescription>
                </Alert>
              )}

              {/* Vesting Periods */}
              <div className="space-y-3">
                <h5 className="font-medium">Vesting Schedule</h5>
                {vesting.vestingPeriods.map((period) => (
                  <div key={period.id} className="flex items-center justify-between p-3 border rounded">
                    <div className="flex items-center gap-3">
                      {period.isClaimed ? (
                        <CheckCircle className="w-5 h-5 text-green-500" />
                      ) : period.isUnlocked ? (
                        <Unlock className="w-5 h-5 text-blue-500" />
                      ) : (
                        <Lock className="w-5 h-5 text-muted-foreground" />
                      )}
                      <div>
                        <div className="font-medium">{period.description}</div>
                        <div className="text-sm text-muted-foreground">
                          {period.amount.toFixed(2)} {vesting.projectSymbol} ({period.percentage}%)
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {period.isClaimed ? `Claimed: ${formatTimeDetailed(period.claimedAt!)}` :
                           period.isUnlocked ? 'Available now' :
                           `Unlocks: ${formatTime(period.unlockDate)}`}
                        </div>
                      </div>
                    </div>

                    <div className="flex items-center gap-2">
                      <div className="text-right">
                        <div className="font-medium">
                          {formatCurrency(period.amount * vesting.currentPrice)}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          ${vesting.currentPrice.toFixed(3)} per token
                        </div>
                      </div>

                      {period.isUnlocked && !period.isClaimed && (
                        <Button 
                          size="sm"
                          onClick={() => handleClaim(vesting.id, period.id)}
                        >
                          <Unlock className="w-3 h-3 mr-1" />
                          Claim
                        </Button>
                      )}

                      {period.isClaimed && period.transactionHash && (
                        <Button variant="outline" size="sm">
                          <ExternalLink className="w-3 h-3" />
                        </Button>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              {/* Quick Actions */}
              <div className="flex gap-2">
                <Button 
                  variant="outline" 
                  size="sm"
                  disabled={getClaimableAmount(vesting) === 0}
                >
                  <Unlock className="w-3 h-3 mr-1" />
                  Claim All Available
                </Button>
                <Button variant="outline" size="sm">
                  <BarChart3 className="w-3 h-3 mr-1" />
                  View Analytics
                </Button>
                <Button variant="outline" size="sm">
                  <ExternalLink className="w-3 h-3 mr-1" />
                  View Project
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Claim History */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <History className="w-5 h-5" />
            Claim History
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {claimHistory.map((claim) => (
              <div key={claim.id} className="flex items-center justify-between p-3 border rounded">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 bg-green-100 dark:bg-green-900/20 rounded-full flex items-center justify-center">
                    <CheckCircle className="w-4 h-4 text-green-500" />
                  </div>
                  <div>
                    <div className="font-medium">
                      {claim.amount.toFixed(2)} {claim.projectSymbol}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {claim.projectName}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {formatTimeDetailed(claim.claimedAt)}
                    </div>
                  </div>
                </div>

                <div className="text-right">
                  <div className="font-medium">{formatCurrency(claim.value)}</div>
                  <div className="text-xs text-muted-foreground">
                    Gas: {formatCurrency(claim.gasFee)}
                  </div>
                  <Button variant="ghost" size="sm" className="p-0 h-auto">
                    <ExternalLink className="w-3 h-3" />
                  </Button>
                </div>
              </div>
            ))}
          </div>

          {claimHistory.length === 0 && (
            <div className="text-center py-8">
              <History className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
              <h3 className="text-lg font-medium mb-2">No Claims Yet</h3>
              <p className="text-muted-foreground">
                Your token claims will appear here once you start claiming vested tokens
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
