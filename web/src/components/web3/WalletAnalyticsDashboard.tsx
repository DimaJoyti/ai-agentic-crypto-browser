'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  BarChart3,
  TrendingUp,
  Clock,
  Wallet,
  Activity,
  Download,
  Upload,
  Trash2,
  RefreshCw,
  PieChart,
  Calendar,
  Shield,
  Zap
} from 'lucide-react'
import { useWalletStore } from '@/stores/walletStore'
import { useWalletPersistence } from '@/lib/wallet-persistence'
import { useSessionRecovery } from '@/lib/session-recovery'
import { toast } from 'sonner'

export function WalletAnalyticsDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [storageStats, setStorageStats] = useState<any>(null)
  const [isLoading, setIsLoading] = useState(false)

  const walletStore = useWalletStore()
  const persistence = useWalletPersistence()
  const sessionRecovery = useSessionRecovery()

  // Load storage stats
  useEffect(() => {
    const loadStats = async () => {
      try {
        const stats = await persistence.getStats()
        setStorageStats(stats)
      } catch (error) {
        console.error('Failed to load storage stats:', error)
      }
    }
    loadStats()
  }, [persistence])

  const handleExportData = async () => {
    setIsLoading(true)
    try {
      const data = await persistence.exportData()
      const blob = new Blob([data], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `wallet-backup-${new Date().toISOString().split('T')[0]}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      toast.success('Wallet data exported successfully')
    } catch (error) {
      toast.error('Failed to export wallet data')
    } finally {
      setIsLoading(false)
    }
  }

  const handleImportData = () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.json'
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0]
      if (!file) return

      setIsLoading(true)
      try {
        const text = await file.text()
        const success = await persistence.importData(text)
        if (success) {
          toast.success('Wallet data imported successfully')
          window.location.reload() // Refresh to show imported data
        } else {
          toast.error('Failed to import wallet data')
        }
      } catch (error) {
        toast.error('Invalid backup file')
      } finally {
        setIsLoading(false)
      }
    }
    input.click()
  }

  const handleClearData = () => {
    if (window.confirm('Are you sure you want to clear all wallet data? This action cannot be undone.')) {
      persistence.clearData()
      toast.success('All wallet data cleared')
      window.location.reload()
    }
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const formatTimeAgo = (timestamp: number) => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`
    if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`
    if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`
    return 'Just now'
  }

  const getMostUsedChain = () => {
    const chains = Object.entries(walletStore.analytics.favoriteChains)
    if (chains.length === 0) return null
    return chains.sort(([,a], [,b]) => b - a)[0]
  }

  const getConnectionTrend = () => {
    const monthlyStats = Object.entries(walletStore.analytics.monthlyStats)
    if (monthlyStats.length < 2) return 0
    
    const sorted = monthlyStats.sort(([a], [b]) => a.localeCompare(b))
    const current = sorted[sorted.length - 1]?.[1]?.connections || 0
    const previous = sorted[sorted.length - 2]?.[1]?.connections || 0
    
    if (previous === 0) return 100
    return ((current - previous) / previous) * 100
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Wallet Analytics</h2>
          <p className="text-muted-foreground">
            Insights into your wallet usage and performance
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={handleExportData} disabled={isLoading}>
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button variant="outline" onClick={handleImportData} disabled={isLoading}>
            <Upload className="w-4 h-4 mr-2" />
            Import
          </Button>
          <Button variant="destructive" onClick={handleClearData} disabled={isLoading}>
            <Trash2 className="w-4 h-4 mr-2" />
            Clear
          </Button>
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="wallets">Wallets</TabsTrigger>
          <TabsTrigger value="sessions">Sessions</TabsTrigger>
          <TabsTrigger value="storage">Storage</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Total Connections</p>
                    <p className="text-2xl font-bold">{walletStore.analytics.totalConnections}</p>
                  </div>
                  <Wallet className="w-8 h-8 text-blue-500" />
                </div>
                <div className="mt-2 flex items-center text-sm">
                  <TrendingUp className="w-4 h-4 mr-1 text-green-500" />
                  <span className="text-green-600">
                    {getConnectionTrend() > 0 ? '+' : ''}{getConnectionTrend().toFixed(1)}%
                  </span>
                  <span className="text-muted-foreground ml-1">from last month</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Total Transactions</p>
                    <p className="text-2xl font-bold">{walletStore.analytics.totalTransactions}</p>
                  </div>
                  <Activity className="w-8 h-8 text-green-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Connected Wallets</p>
                    <p className="text-2xl font-bold">{walletStore.connectedWallets.length}</p>
                  </div>
                  <Shield className="w-8 h-8 text-purple-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Active Sessions</p>
                    <p className="text-2xl font-bold">{sessionRecovery.sessionStats.activeSessions}</p>
                  </div>
                  <Clock className="w-8 h-8 text-orange-500" />
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Most Used Wallets</CardTitle>
                <CardDescription>Your wallet usage statistics</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {Object.entries(walletStore.analytics.walletUsageStats)
                    .sort(([,a], [,b]) => b.connections - a.connections)
                    .slice(0, 5)
                    .map(([walletId, stats]) => (
                      <div key={walletId} className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-secondary rounded-lg flex items-center justify-center">
                            <Wallet className="w-4 h-4" />
                          </div>
                          <div>
                            <p className="font-medium capitalize">{walletId}</p>
                            <p className="text-sm text-muted-foreground">
                              Last used {formatTimeAgo(stats.lastUsed)}
                            </p>
                          </div>
                        </div>
                        <Badge variant="secondary">{stats.connections} connections</Badge>
                      </div>
                    ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Favorite Networks</CardTitle>
                <CardDescription>Your most used blockchain networks</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {Object.entries(walletStore.analytics.favoriteChains)
                    .sort(([,a], [,b]) => b - a)
                    .slice(0, 5)
                    .map(([chainId, count]) => {
                      const chainNames: Record<string, string> = {
                        '1': 'Ethereum',
                        '137': 'Polygon',
                        '42161': 'Arbitrum',
                        '10': 'Optimism',
                        '56': 'BSC'
                      }
                      return (
                        <div key={chainId} className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 bg-secondary rounded-lg flex items-center justify-center">
                              <Zap className="w-4 h-4" />
                            </div>
                            <span className="font-medium">{chainNames[chainId] || `Chain ${chainId}`}</span>
                          </div>
                          <Badge variant="secondary">{count} transactions</Badge>
                        </div>
                      )
                    })}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="wallets" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Connected Wallets</CardTitle>
              <CardDescription>Manage your connected wallet accounts</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {walletStore.connectedWallets.map((wallet) => (
                  <motion.div
                    key={wallet.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-secondary rounded-lg flex items-center justify-center">
                        <Wallet className="w-5 h-5" />
                      </div>
                      <div>
                        <p className="font-medium">{wallet.provider.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {wallet.address.slice(0, 6)}...{wallet.address.slice(-4)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          Connected {formatTimeAgo(wallet.connectedAt)}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {wallet.isActive && (
                        <Badge variant="default">Active</Badge>
                      )}
                      <Badge variant="outline">{wallet.connector}</Badge>
                    </div>
                  </motion.div>
                ))}
                
                {walletStore.connectedWallets.length === 0 && (
                  <div className="text-center py-8 text-muted-foreground">
                    No wallets connected yet
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="sessions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Session History</CardTitle>
              <CardDescription>Your wallet connection sessions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {walletStore.sessions.slice(0, 10).map((session) => (
                  <div key={session.id} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-secondary rounded-lg flex items-center justify-center">
                        <Clock className="w-4 h-4" />
                      </div>
                      <div>
                        <p className="font-medium">
                          {session.address.slice(0, 6)}...{session.address.slice(-4)}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {formatTimeAgo(session.startTime)}
                          {session.duration && ` • ${Math.round(session.duration / 60000)} min`}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium">{session.transactionCount} transactions</p>
                      <p className="text-xs text-muted-foreground">Chain {session.chainId}</p>
                    </div>
                  </div>
                ))}
                
                {walletStore.sessions.length === 0 && (
                  <div className="text-center py-8 text-muted-foreground">
                    No sessions recorded yet
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="storage" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Storage Usage</CardTitle>
              <CardDescription>Monitor your local storage usage</CardDescription>
            </CardHeader>
            <CardContent>
              {storageStats && (
                <div className="space-y-6">
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium">Storage Used</span>
                      <span className="text-sm text-muted-foreground">
                        {formatBytes(storageStats.quota.used)} / {formatBytes(storageStats.quota.total)}
                      </span>
                    </div>
                    <Progress value={storageStats.quota.percentage} className="h-2" />
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <h4 className="font-medium mb-2">Storage Breakdown</h4>
                      <div className="space-y-2">
                        {Object.entries(storageStats.itemSizes).map(([key, size]) => (
                          <div key={key} className="flex justify-between text-sm">
                            <span className="capitalize">{key.replace('wallet_', '').replace('_', ' ')}</span>
                            <span>{formatBytes(size as number)}</span>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div>
                      <h4 className="font-medium mb-2">Statistics</h4>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span>Total Items</span>
                          <span>{storageStats.itemCount}</span>
                        </div>
                        <div className="flex justify-between">
                          <span>Total Size</span>
                          <span>{formatBytes(storageStats.totalSize)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span>Available Space</span>
                          <span>{formatBytes(storageStats.quota.available)}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
