'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Activity, 
  Clock, 
  CheckCircle, 
  XCircle, 
  AlertTriangle,
  ExternalLink,
  Copy,
  Trash2,
  RefreshCw,
  Zap,
  ArrowUpRight,
  ArrowDownLeft,
  Repeat,
  Shield
} from 'lucide-react'
import { useTransactionMonitor } from '@/hooks/useTransactionMonitor'
import { TransactionStatus, TransactionType, type TransactionData } from '@/lib/transaction-monitor'
import { SUPPORTED_CHAINS } from '@/lib/chains'
import { TransactionDemo } from './TransactionDemo'
import { toast } from 'sonner'

export function TransactionTracker() {
  const [activeTab, setActiveTab] = useState('all')
  const [autoRefresh, setAutoRefresh] = useState(true)
  
  const {
    transactions,
    isLoading,
    error,
    pendingTransactions,
    confirmedTransactions,
    failedTransactions,
    stopTracking,
    clearAllTransactions
  } = useTransactionMonitor({
    showNotifications: true,
    autoRemoveCompleted: false // Keep completed transactions for viewing
  })

  // Auto-refresh every 10 seconds
  useEffect(() => {
    if (!autoRefresh) return

    const interval = setInterval(() => {
      // The hook automatically updates, this is just for UI feedback
    }, 10000)

    return () => clearInterval(interval)
  }, [autoRefresh])

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
        return 'bg-yellow-500'
      case TransactionStatus.CONFIRMED:
        return 'bg-green-500'
      case TransactionStatus.FAILED:
        return 'bg-red-500'
      case TransactionStatus.DROPPED:
        return 'bg-orange-500'
      default:
        return 'bg-gray-500'
    }
  }

  const getTypeIcon = (type: TransactionType) => {
    switch (type) {
      case TransactionType.SEND:
        return <ArrowUpRight className="w-4 h-4" />
      case TransactionType.RECEIVE:
        return <ArrowDownLeft className="w-4 h-4" />
      case TransactionType.SWAP:
        return <Repeat className="w-4 h-4" />
      case TransactionType.APPROVE:
        return <Shield className="w-4 h-4" />
      case TransactionType.STAKE:
      case TransactionType.UNSTAKE:
        return <Zap className="w-4 h-4" />
      default:
        return <Activity className="w-4 h-4" />
    }
  }

  const formatHash = (hash: string) => {
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleTimeString()
  }

  const copyHash = (hash: string) => {
    navigator.clipboard.writeText(hash)
    toast.success('Transaction hash copied to clipboard')
  }

  const openBlockExplorer = (tx: TransactionData) => {
    const chain = SUPPORTED_CHAINS[tx.chainId]
    if (chain?.blockExplorers?.default?.url) {
      window.open(`${chain.blockExplorers.default.url}/tx/${tx.hash}`, '_blank')
    }
  }

  const getConfirmationProgress = (tx: TransactionData) => {
    if (tx.status !== TransactionStatus.PENDING && tx.status !== TransactionStatus.CONFIRMED) {
      return 100
    }
    return Math.min((tx.confirmations / tx.maxConfirmations) * 100, 100)
  }

  const getFilteredTransactions = () => {
    switch (activeTab) {
      case 'pending':
        return pendingTransactions
      case 'confirmed':
        return confirmedTransactions
      case 'failed':
        return failedTransactions
      default:
        return transactions
    }
  }

  const filteredTransactions = getFilteredTransactions()

  return (
    <div className="space-y-6">
      {/* Demo Component */}
      <TransactionDemo />
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Activity className="w-6 h-6" />
            Transaction Tracker
          </h2>
          <p className="text-muted-foreground">
            Real-time monitoring of your blockchain transactions
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setAutoRefresh(!autoRefresh)}
            className={autoRefresh ? 'bg-green-50 border-green-200' : ''}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${autoRefresh ? 'animate-spin' : ''}`} />
            Auto Refresh
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={clearAllTransactions}
            disabled={transactions.length === 0}
          >
            <Trash2 className="w-4 h-4 mr-2" />
            Clear All
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total</p>
                <p className="text-2xl font-bold">{transactions.length}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Pending</p>
                <p className="text-2xl font-bold">{pendingTransactions.length}</p>
              </div>
              <Clock className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Confirmed</p>
                <p className="text-2xl font-bold">{confirmedTransactions.length}</p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Failed</p>
                <p className="text-2xl font-bold">{failedTransactions.length}</p>
              </div>
              <XCircle className="w-8 h-8 text-red-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Transaction List */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Transactions</CardTitle>
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList>
                <TabsTrigger value="all">All</TabsTrigger>
                <TabsTrigger value="pending">Pending</TabsTrigger>
                <TabsTrigger value="confirmed">Confirmed</TabsTrigger>
                <TabsTrigger value="failed">Failed</TabsTrigger>
              </TabsList>
            </Tabs>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading && transactions.length === 0 ? (
            <div className="text-center py-8">
              <RefreshCw className="w-8 h-8 animate-spin mx-auto mb-4 text-muted-foreground" />
              <p className="text-muted-foreground">Loading transactions...</p>
            </div>
          ) : filteredTransactions.length === 0 ? (
            <div className="text-center py-8">
              <Activity className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-semibold mb-2">No transactions found</h3>
              <p className="text-muted-foreground">
                {activeTab === 'all' 
                  ? 'Your transactions will appear here when you start using the wallet'
                  : `No ${activeTab} transactions at the moment`
                }
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              <AnimatePresence>
                {filteredTransactions.map((tx, index) => {
                  const chain = SUPPORTED_CHAINS[tx.chainId]
                  const progress = getConfirmationProgress(tx)
                  
                  return (
                    <motion.div
                      key={tx.hash}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -20 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                    >
                      <div className="flex items-center gap-4 flex-1">
                        <div className="flex items-center gap-2">
                          {getStatusIcon(tx.status)}
                          {getTypeIcon(tx.type)}
                        </div>
                        
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1">
                            <span className="font-medium">{formatHash(tx.hash)}</span>
                            <Badge variant="outline" className="text-xs">
                              {chain?.shortName || `Chain ${tx.chainId}`}
                            </Badge>
                            <Badge variant="secondary" className="text-xs">
                              {tx.type.replace('_', ' ')}
                            </Badge>
                          </div>
                          
                          <div className="flex items-center gap-4 text-sm text-muted-foreground">
                            <span>{formatTime(tx.timestamp)}</span>
                            {tx.value && parseFloat(tx.value) > 0 && (
                              <span>{parseFloat(tx.value).toFixed(6)} {chain?.gasToken}</span>
                            )}
                            {tx.status === TransactionStatus.PENDING && (
                              <span>{tx.confirmations}/{tx.maxConfirmations} confirmations</span>
                            )}
                            {tx.gasUsed && (
                              <span>Gas: {parseInt(tx.gasUsed).toLocaleString()}</span>
                            )}
                          </div>

                          {tx.status === TransactionStatus.PENDING && (
                            <div className="mt-2">
                              <Progress value={progress} className="h-1" />
                            </div>
                          )}
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <div className={`w-2 h-2 rounded-full ${getStatusColor(tx.status)}`} />
                        
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => copyHash(tx.hash)}
                        >
                          <Copy className="w-4 h-4" />
                        </Button>
                        
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => openBlockExplorer(tx)}
                        >
                          <ExternalLink className="w-4 h-4" />
                        </Button>
                        
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => stopTracking(tx.hash)}
                          className="text-red-500 hover:text-red-600"
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </motion.div>
                  )
                })}
              </AnimatePresence>
            </div>
          )}
        </CardContent>
      </Card>

      {error && (
        <Card className="border-red-200 bg-red-50 dark:bg-red-950">
          <CardContent className="p-4">
            <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
              <AlertTriangle className="w-4 h-4" />
              <span className="font-medium">Error:</span>
              <span>{error}</span>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
