'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Activity,
  Clock,
  CheckCircle,
  XCircle,
  AlertTriangle,
  ExternalLink,
  RefreshCw,
  Zap,
  X,
  RotateCcw,
  TrendingUp,
  Eye,
  Settings
} from 'lucide-react'
import { useTransactionMonitor } from '@/hooks/useTransactionMonitor'
import { TransactionStatus, TransactionType } from '@/lib/transaction-monitor'
import { type Hash } from 'viem'
import { cn } from '@/lib/utils'
import { formatEther } from 'viem'

interface TransactionStatusProps {
  hash: Hash
  compact?: boolean
}

export function TransactionStatusComponent({ hash, compact = false }: TransactionStatusProps) {
  const { getTransaction } = useTransactionMonitor()
  const transaction = getTransaction(hash)

  if (!transaction) {
    return (
      <div className="flex items-center gap-2 text-muted-foreground">
        <AlertTriangle className="w-4 h-4" />
        <span className="text-sm">Transaction not found</span>
      </div>
    )
  }

  const getStatusIcon = (status: TransactionStatus) => {
    switch (status) {
      case TransactionStatus.PENDING:
        return <Clock className="w-4 h-4 text-yellow-500 animate-pulse" />
      case TransactionStatus.CONFIRMED:
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case TransactionStatus.FAILED:
        return <XCircle className="w-4 h-4 text-red-500" />
      case TransactionStatus.DROPPED:
        return <AlertTriangle className="w-4 h-4 text-orange-500" />
      default:
        return <Activity className="w-4 h-4 text-gray-500" />
    }
  }

  const getStatusText = (status: TransactionStatus) => {
    switch (status) {
      case TransactionStatus.PENDING:
        return 'Pending'
      case TransactionStatus.CONFIRMED:
        return 'Confirmed'
      case TransactionStatus.FAILED:
        return 'Failed'
      case TransactionStatus.DROPPED:
        return 'Dropped'
      default:
        return 'Unknown'
    }
  }

  const formatHash = (hash: Hash) => {
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`
  }

  if (compact) {
    return (
      <div className="flex items-center gap-2">
        {getStatusIcon(transaction.status)}
        <span className="text-sm font-medium">{getStatusText(transaction.status)}</span>
        {transaction.status === TransactionStatus.CONFIRMED && (
          <span className="text-xs text-muted-foreground">
            ({transaction.confirmations}/{transaction.maxConfirmations})
          </span>
        )}
      </div>
    )
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Transaction Status</CardTitle>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => window.open(`https://etherscan.io/tx/${hash}`, '_blank')}
          >
            <ExternalLink className="w-4 h-4" />
          </Button>
        </div>
        <CardDescription>{formatHash(hash)}</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium">Status</span>
            <div className="flex items-center gap-2">
              {getStatusIcon(transaction.status)}
              <span className="text-sm font-medium">{getStatusText(transaction.status)}</span>
            </div>
          </div>

          {transaction.status === TransactionStatus.CONFIRMED && (
            <div>
              <div className="flex items-center justify-between text-sm mb-2">
                <span>Confirmations</span>
                <span>{transaction.confirmations}/{transaction.maxConfirmations}</span>
              </div>
              <Progress
                value={(transaction.confirmations / transaction.maxConfirmations) * 100}
                className="h-2"
              />
            </div>
          )}

          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-muted-foreground">Value</span>
              <p className="font-medium">{parseFloat(transaction.value).toFixed(4)} ETH</p>
            </div>
            <div>
              <span className="text-muted-foreground">Chain</span>
              <p className="font-medium">Chain {transaction.chainId}</p>
            </div>
            {transaction.gasUsed && (
              <div>
                <span className="text-muted-foreground">Gas Used</span>
                <p className="font-medium">{parseInt(transaction.gasUsed).toLocaleString()}</p>
              </div>
            )}
            {transaction.blockNumber && (
              <div>
                <span className="text-muted-foreground">Block</span>
                <p className="font-medium">{transaction.blockNumber.toLocaleString()}</p>
              </div>
            )}
          </div>

          {transaction.error && (
            <Alert variant="destructive">
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription className="text-sm">
                {transaction.error}
              </AlertDescription>
            </Alert>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export function TransactionTracker() {
  const [activeTab, setActiveTab] = useState('all')
  const [showCompleted, setShowCompleted] = useState(true)

  const {
    transactions,
    pendingTransactions,
    confirmedTransactions,
    failedTransactions,
    isLoading,
    error,
    stopTracking,
    retryTransaction,
    cancelTransaction,
    speedUpTransaction,
    startRealtimeMonitoring,
    stopRealtimeMonitoring,
    getStats
  } = useTransactionMonitor({
    showNotifications: true,
    autoRemoveCompleted: false
  })

  const stats = getStats()

  const getStatusIcon = (status: TransactionStatus) => {
    switch (status) {
      case TransactionStatus.PENDING:
        return <Clock className="w-4 h-4 text-yellow-500 animate-pulse" />
      case TransactionStatus.CONFIRMED:
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case TransactionStatus.FAILED:
        return <XCircle className="w-4 h-4 text-red-500" />
      case TransactionStatus.DROPPED:
        return <AlertTriangle className="w-4 h-4 text-orange-500" />
      default:
        return <Activity className="w-4 h-4 text-gray-500" />
    }
  }

  const getStatusColor = (status: TransactionStatus) => {
    switch (status) {
      case TransactionStatus.PENDING:
        return 'border-yellow-200 bg-yellow-50 dark:border-yellow-800 dark:bg-yellow-950'
      case TransactionStatus.CONFIRMED:
        return 'border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950'
      case TransactionStatus.FAILED:
        return 'border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-950'
      case TransactionStatus.DROPPED:
        return 'border-orange-200 bg-orange-50 dark:border-orange-800 dark:bg-orange-950'
      default:
        return 'border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-950'
    }
  }

  const getTypeIcon = (type: TransactionType) => {
    switch (type) {
      case TransactionType.SEND:
        return <TrendingUp className="w-4 h-4" />
      case TransactionType.RECEIVE:
        return <TrendingUp className="w-4 h-4 rotate-180" />
      case TransactionType.SWAP:
        return <RefreshCw className="w-4 h-4" />
      case TransactionType.APPROVE:
        return <CheckCircle className="w-4 h-4" />
      default:
        return <Activity className="w-4 h-4" />
    }
  }

  const formatHash = (hash: Hash) => {
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`
  }

  const formatTime = (timestamp: number) => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (days > 0) return `${days}d ago`
    if (hours > 0) return `${hours}h ago`
    if (minutes > 0) return `${minutes}m ago`
    return 'Just now'
  }

  const getConfirmationProgress = (confirmations: number, maxConfirmations: number) => {
    return Math.min((confirmations / maxConfirmations) * 100, 100)
  }

  const openBlockExplorer = (hash: Hash, chainId: number) => {
    const explorers: Record<number, string> = {
      1: 'https://etherscan.io/tx/',
      137: 'https://polygonscan.com/tx/',
      42161: 'https://arbiscan.io/tx/',
      10: 'https://optimistic.etherscan.io/tx/',
      8453: 'https://basescan.org/tx/',
      11155111: 'https://sepolia.etherscan.io/tx/'
    }
    
    const baseUrl = explorers[chainId] || 'https://etherscan.io/tx/'
    window.open(`${baseUrl}${hash}`, '_blank')
  }

  const filteredTransactions = transactions.filter(tx => {
    if (activeTab === 'pending') return tx.status === TransactionStatus.PENDING
    if (activeTab === 'confirmed') return tx.status === TransactionStatus.CONFIRMED
    if (activeTab === 'failed') return tx.status === TransactionStatus.FAILED
    if (!showCompleted && tx.status === TransactionStatus.CONFIRMED) return false
    return true
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Transaction Tracker</h2>
          <p className="text-muted-foreground">
            Monitor your blockchain transactions in real-time
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowCompleted(!showCompleted)}
          >
            <Eye className="w-4 h-4 mr-2" />
            {showCompleted ? 'Hide' : 'Show'} Completed
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => startRealtimeMonitoring(1)}
          >
            <Activity className="w-4 h-4 mr-2" />
            Real-time
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Statistics Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total</p>
                <p className="text-2xl font-bold">{stats.total}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Pending</p>
                <p className="text-2xl font-bold text-yellow-600">{stats.pending}</p>
              </div>
              <Clock className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Confirmed</p>
                <p className="text-2xl font-bold text-green-600">{stats.confirmed}</p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Failed</p>
                <p className="text-2xl font-bold text-red-600">{stats.failed}</p>
              </div>
              <XCircle className="w-8 h-8 text-red-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Transaction Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="all">All ({transactions.length})</TabsTrigger>
          <TabsTrigger value="pending">Pending ({pendingTransactions.length})</TabsTrigger>
          <TabsTrigger value="confirmed">Confirmed ({confirmedTransactions.length})</TabsTrigger>
          <TabsTrigger value="failed">Failed ({failedTransactions.length})</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {filteredTransactions.length > 0 ? (
            <div className="space-y-3">
              <AnimatePresence>
                {filteredTransactions.map((transaction, index) => (
                  <motion.div
                    key={transaction.hash}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ delay: index * 0.05 }}
                  >
                    <Card className={cn("transition-all duration-200", getStatusColor(transaction.status))}>
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="flex items-center gap-2">
                              {getStatusIcon(transaction.status)}
                              {getTypeIcon(transaction.type)}
                            </div>
                            
                            <div>
                              <div className="flex items-center gap-2">
                                <p className="font-medium">{formatHash(transaction.hash)}</p>
                                <Badge variant="outline" className="text-xs">
                                  {transaction.type}
                                </Badge>
                                <Badge variant="secondary" className="text-xs">
                                  Chain {transaction.chainId}
                                </Badge>
                              </div>
                              
                              <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                                <span>Value: {parseFloat(transaction.value).toFixed(4)} ETH</span>
                                <span>{formatTime(transaction.timestamp)}</span>
                                {transaction.gasUsed && (
                                  <span>Gas: {parseInt(transaction.gasUsed).toLocaleString()}</span>
                                )}
                              </div>

                              {transaction.metadata?.description && (
                                <p className="text-sm text-muted-foreground mt-1">
                                  {transaction.metadata.description}
                                </p>
                              )}

                              {/* Confirmation Progress */}
                              {transaction.status === TransactionStatus.CONFIRMED && (
                                <div className="mt-2">
                                  <div className="flex items-center justify-between text-xs text-muted-foreground mb-1">
                                    <span>Confirmations</span>
                                    <span>{transaction.confirmations}/{transaction.maxConfirmations}</span>
                                  </div>
                                  <Progress 
                                    value={getConfirmationProgress(transaction.confirmations, transaction.maxConfirmations)} 
                                    className="h-1"
                                  />
                                </div>
                              )}
                            </div>
                          </div>

                          <div className="flex items-center gap-2">
                            {/* Action Buttons */}
                            {transaction.status === TransactionStatus.PENDING && (
                              <>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => speedUpTransaction(transaction.hash, '20000000000')}
                                  disabled={isLoading}
                                >
                                  <Zap className="w-4 h-4" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => cancelTransaction(transaction.hash)}
                                  disabled={isLoading}
                                >
                                  <X className="w-4 h-4" />
                                </Button>
                              </>
                            )}

                            {transaction.status === TransactionStatus.FAILED && (
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => retryTransaction(transaction.hash)}
                                disabled={isLoading}
                              >
                                <RotateCcw className="w-4 h-4" />
                              </Button>
                            )}

                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => openBlockExplorer(transaction.hash, transaction.chainId)}
                            >
                              <ExternalLink className="w-4 h-4" />
                            </Button>

                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => stopTracking(transaction.hash)}
                            >
                              <X className="w-4 h-4" />
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Activity className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Transactions</h3>
                <p className="text-muted-foreground">
                  {activeTab === 'all' 
                    ? 'No transactions are being tracked'
                    : `No ${activeTab} transactions found`
                  }
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
