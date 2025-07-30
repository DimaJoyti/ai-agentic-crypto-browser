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
  List,
  Play,
  Pause,
  Trash2,
  X,
  ArrowUp,
  ArrowDown,
  Clock,
  CheckCircle,
  XCircle,
  AlertTriangle,
  RotateCcw,
  Settings,
  Activity,
  TrendingUp,
  Timer,
  Zap
} from 'lucide-react'
import { useTransactionPool } from '@/hooks/useTransactionPool'
import { TransactionPriority, QueueStatus, type QueuedTransaction } from '@/lib/transaction-pool'
import { cn } from '@/lib/utils'

export function TransactionQueueDashboard() {
  const [activeTab, setActiveTab] = useState('queue')
  const [showCompleted, setShowCompleted] = useState(false)

  const {
    state,
    addTransaction,
    removeTransaction,
    cancelTransaction,
    updatePriority,
    startProcessing,
    stopProcessing,
    clearCompleted,
    clearError
  } = useTransactionPool({
    autoStart: true,
    enableNotifications: true,
    filterByAddress: true
  })

  const getStatusIcon = (status: QueueStatus) => {
    switch (status) {
      case QueueStatus.QUEUED:
        return <List className="w-4 h-4 text-blue-500" />
      case QueueStatus.PENDING:
        return <Clock className="w-4 h-4 text-yellow-500 animate-pulse" />
      case QueueStatus.SUBMITTED:
        return <Activity className="w-4 h-4 text-orange-500 animate-pulse" />
      case QueueStatus.CONFIRMED:
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case QueueStatus.FAILED:
        return <XCircle className="w-4 h-4 text-red-500" />
      case QueueStatus.CANCELLED:
        return <X className="w-4 h-4 text-gray-500" />
      case QueueStatus.EXPIRED:
        return <AlertTriangle className="w-4 h-4 text-orange-500" />
      default:
        return <Activity className="w-4 h-4 text-gray-500" />
    }
  }

  const getStatusColor = (status: QueueStatus) => {
    switch (status) {
      case QueueStatus.QUEUED:
        return 'border-blue-200 bg-blue-50 dark:border-blue-800 dark:bg-blue-950'
      case QueueStatus.PENDING:
        return 'border-yellow-200 bg-yellow-50 dark:border-yellow-800 dark:bg-yellow-950'
      case QueueStatus.SUBMITTED:
        return 'border-orange-200 bg-orange-50 dark:border-orange-800 dark:bg-orange-950'
      case QueueStatus.CONFIRMED:
        return 'border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950'
      case QueueStatus.FAILED:
        return 'border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-950'
      case QueueStatus.CANCELLED:
        return 'border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-950'
      case QueueStatus.EXPIRED:
        return 'border-orange-200 bg-orange-50 dark:border-orange-800 dark:bg-orange-950'
      default:
        return 'border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-950'
    }
  }

  const getPriorityIcon = (priority: TransactionPriority) => {
    switch (priority) {
      case TransactionPriority.LOW:
        return <ArrowDown className="w-4 h-4 text-blue-500" />
      case TransactionPriority.NORMAL:
        return <Activity className="w-4 h-4 text-green-500" />
      case TransactionPriority.HIGH:
        return <ArrowUp className="w-4 h-4 text-orange-500" />
      case TransactionPriority.URGENT:
        return <Zap className="w-4 h-4 text-red-500" />
      default:
        return <Activity className="w-4 h-4 text-gray-500" />
    }
  }

  const getPriorityName = (priority: TransactionPriority) => {
    const names = ['Low', 'Normal', 'High', 'Urgent']
    return names[priority] || 'Unknown'
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

  const formatHash = (hash?: string) => {
    if (!hash) return 'N/A'
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`
  }

  const handlePriorityChange = (id: string, currentPriority: TransactionPriority, increase: boolean) => {
    const newPriority = increase 
      ? Math.min(currentPriority + 1, TransactionPriority.URGENT)
      : Math.max(currentPriority - 1, TransactionPriority.LOW)
    
    updatePriority(id, newPriority)
  }

  const filteredTransactions = showCompleted 
    ? state.transactions
    : state.transactions.filter(tx => 
        ![QueueStatus.CONFIRMED, QueueStatus.CANCELLED, QueueStatus.EXPIRED].includes(tx.status)
      )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Transaction Queue</h2>
          <p className="text-muted-foreground">
            Manage and monitor your transaction queue with intelligent processing
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowCompleted(!showCompleted)}
          >
            {showCompleted ? 'Hide' : 'Show'} Completed
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={clearCompleted}
          >
            <Trash2 className="w-4 h-4 mr-2" />
            Clear Completed
          </Button>
          <Button
            variant={state.isProcessing ? "destructive" : "default"}
            size="sm"
            onClick={state.isProcessing ? stopProcessing : startProcessing}
          >
            {state.isProcessing ? (
              <>
                <Pause className="w-4 h-4 mr-2" />
                Stop Processing
              </>
            ) : (
              <>
                <Play className="w-4 h-4 mr-2" />
                Start Processing
              </>
            )}
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Queue Statistics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total</p>
                <p className="text-2xl font-bold">{state.stats.totalTransactions}</p>
              </div>
              <List className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              All transactions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Queued</p>
                <p className="text-2xl font-bold text-blue-600">{state.stats.queuedTransactions}</p>
              </div>
              <Clock className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Waiting to process
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Processing</p>
                <p className="text-2xl font-bold text-yellow-600">
                  {state.stats.pendingTransactions + state.stats.submittedTransactions}
                </p>
              </div>
              <Activity className="w-8 h-8 text-yellow-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Currently processing
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold text-green-600">
                  {state.stats.successRate.toFixed(1)}%
                </p>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Confirmation rate
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Processing Status */}
      {state.isProcessing && (
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="font-medium">Queue Processing Active</h3>
                <p className="text-sm text-muted-foreground">
                  Processing up to {state.config.maxConcurrentTransactions} transactions concurrently
                </p>
              </div>
              <div className="flex items-center gap-2">
                <Activity className="w-5 h-5 text-green-500 animate-pulse" />
                <span className="text-sm font-medium text-green-600">Active</span>
              </div>
            </div>
            
            {state.stats.queuedTransactions > 0 && (
              <div>
                <div className="flex items-center justify-between text-sm mb-2">
                  <span>Queue Progress</span>
                  <span>
                    {state.stats.totalTransactions - state.stats.queuedTransactions} / {state.stats.totalTransactions}
                  </span>
                </div>
                <Progress 
                  value={((state.stats.totalTransactions - state.stats.queuedTransactions) / state.stats.totalTransactions) * 100} 
                  className="h-2"
                />
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Transaction Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="queue">Queue ({state.queuedTransactions.length})</TabsTrigger>
          <TabsTrigger value="processing">Processing ({state.pendingTransactions.length + state.submittedTransactions.length})</TabsTrigger>
          <TabsTrigger value="completed">Completed ({state.confirmedTransactions.length})</TabsTrigger>
          <TabsTrigger value="failed">Failed ({state.failedTransactions.length})</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {filteredTransactions.length > 0 ? (
            <div className="space-y-3">
              <AnimatePresence>
                {filteredTransactions
                  .filter(tx => {
                    switch (activeTab) {
                      case 'queue':
                        return tx.status === QueueStatus.QUEUED
                      case 'processing':
                        return [QueueStatus.PENDING, QueueStatus.SUBMITTED].includes(tx.status)
                      case 'completed':
                        return tx.status === QueueStatus.CONFIRMED
                      case 'failed':
                        return [QueueStatus.FAILED, QueueStatus.CANCELLED, QueueStatus.EXPIRED].includes(tx.status)
                      default:
                        return true
                    }
                  })
                  .map((transaction, index) => (
                    <motion.div
                      key={transaction.id}
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
                                {getPriorityIcon(transaction.priority)}
                              </div>
                              
                              <div>
                                <div className="flex items-center gap-2 mb-1">
                                  <p className="font-medium">{transaction.id.slice(0, 8)}...</p>
                                  <Badge variant="outline" className="text-xs">
                                    {transaction.metadata?.type || 'unknown'}
                                  </Badge>
                                  <Badge variant="secondary" className="text-xs">
                                    {getPriorityName(transaction.priority)}
                                  </Badge>
                                </div>
                                
                                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                                  <span>To: {transaction.to ? formatHash(transaction.to) : 'Contract'}</span>
                                  <span>Value: {parseFloat(transaction.value).toFixed(4)} ETH</span>
                                  <span>{formatTime(transaction.createdAt)}</span>
                                  {transaction.hash && (
                                    <span>Hash: {formatHash(transaction.hash)}</span>
                                  )}
                                </div>

                                {transaction.metadata?.description && (
                                  <p className="text-sm text-muted-foreground mt-1">
                                    {transaction.metadata.description}
                                  </p>
                                )}

                                {transaction.retryCount > 0 && (
                                  <div className="flex items-center gap-1 mt-1">
                                    <RotateCcw className="w-3 h-3 text-orange-500" />
                                    <span className="text-xs text-orange-600">
                                      Retry {transaction.retryCount}/{transaction.maxRetries}
                                    </span>
                                  </div>
                                )}

                                {transaction.error && (
                                  <p className="text-xs text-red-600 mt-1">
                                    Error: {transaction.error}
                                  </p>
                                )}
                              </div>
                            </div>

                            <div className="flex items-center gap-2">
                              {/* Priority Controls */}
                              {transaction.status === QueueStatus.QUEUED && (
                                <>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handlePriorityChange(transaction.id, transaction.priority, true)}
                                    disabled={transaction.priority === TransactionPriority.URGENT}
                                  >
                                    <ArrowUp className="w-4 h-4" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handlePriorityChange(transaction.id, transaction.priority, false)}
                                    disabled={transaction.priority === TransactionPriority.LOW}
                                  >
                                    <ArrowDown className="w-4 h-4" />
                                  </Button>
                                </>
                              )}

                              {/* Cancel Button */}
                              {[QueueStatus.QUEUED, QueueStatus.PENDING].includes(transaction.status) && (
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => cancelTransaction(transaction.id)}
                                >
                                  <X className="w-4 h-4" />
                                </Button>
                              )}

                              {/* Remove Button */}
                              {[QueueStatus.FAILED, QueueStatus.CANCELLED, QueueStatus.EXPIRED].includes(transaction.status) && (
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => removeTransaction(transaction.id)}
                                >
                                  <Trash2 className="w-4 h-4" />
                                </Button>
                              )}
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
                <List className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Transactions</h3>
                <p className="text-muted-foreground">
                  {activeTab === 'queue' && 'No transactions in queue'}
                  {activeTab === 'processing' && 'No transactions being processed'}
                  {activeTab === 'completed' && 'No completed transactions'}
                  {activeTab === 'failed' && 'No failed transactions'}
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
