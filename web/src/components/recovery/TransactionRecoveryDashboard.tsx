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
  RotateCcw,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Clock,
  Zap,
  DollarSign,
  Activity,
  Play,
  Pause,
  X,
  RefreshCw,
  Settings,
  TrendingUp,
  Shield,
  Target,
  Lightbulb
} from 'lucide-react'
import { useTransactionRecovery, useRecoveryStats } from '@/hooks/useTransactionRecovery'
import { 
  FailureReason, 
  RecoveryStatus, 
  RecoveryType,
  type FailedTransaction,
  type RecoveryStrategy 
} from '@/lib/transaction-recovery'
import { cn } from '@/lib/utils'

export function TransactionRecoveryDashboard() {
  const [activeTab, setActiveTab] = useState('failed')
  const [selectedTransaction, setSelectedTransaction] = useState<FailedTransaction | null>(null)

  const {
    state,
    startRecovery,
    cancelRecovery,
    retryRecovery,
    clearFailedTransaction,
    clearAllFailedTransactions,
    clearError
  } = useTransactionRecovery({
    autoAnalyze: true,
    enableNotifications: true,
    requireUserConfirmation: true
  })

  const stats = useRecoveryStats()

  const getFailureReasonIcon = (reason: FailureReason) => {
    switch (reason) {
      case FailureReason.OUT_OF_GAS:
        return <Zap className="w-4 h-4 text-orange-500" />
      case FailureReason.INSUFFICIENT_FUNDS:
        return <DollarSign className="w-4 h-4 text-red-500" />
      case FailureReason.GAS_PRICE_TOO_LOW:
        return <TrendingUp className="w-4 h-4 text-yellow-500" />
      case FailureReason.NETWORK_ERROR:
        return <Activity className="w-4 h-4 text-blue-500" />
      case FailureReason.TIMEOUT:
        return <Clock className="w-4 h-4 text-gray-500" />
      default:
        return <AlertTriangle className="w-4 h-4 text-red-500" />
    }
  }

  const getRecoveryStatusIcon = (status: RecoveryStatus) => {
    switch (status) {
      case RecoveryStatus.RECOVERY_SUCCESS:
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case RecoveryStatus.RECOVERY_FAILED:
        return <XCircle className="w-4 h-4 text-red-500" />
      case RecoveryStatus.RECOVERY_IN_PROGRESS:
        return <Activity className="w-4 h-4 text-blue-500 animate-pulse" />
      case RecoveryStatus.RECOVERY_AVAILABLE:
        return <RotateCcw className="w-4 h-4 text-orange-500" />
      default:
        return <Clock className="w-4 h-4 text-gray-500" />
    }
  }

  const getRecoveryTypeIcon = (type: RecoveryType) => {
    switch (type) {
      case RecoveryType.INCREASE_GAS_PRICE:
        return <TrendingUp className="w-4 h-4" />
      case RecoveryType.INCREASE_GAS_LIMIT:
        return <Zap className="w-4 h-4" />
      case RecoveryType.SPEED_UP_TRANSACTION:
        return <Activity className="w-4 h-4" />
      case RecoveryType.REPLACE_TRANSACTION:
        return <RefreshCw className="w-4 h-4" />
      case RecoveryType.RETRY_TRANSACTION:
        return <RotateCcw className="w-4 h-4" />
      default:
        return <Target className="w-4 h-4" />
    }
  }

  const getFailureReasonName = (reason: FailureReason) => {
    return reason.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
  }

  const getRecoveryStatusName = (status: RecoveryStatus) => {
    return status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
  }

  const getRecoveryTypeName = (type: RecoveryType) => {
    return type.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
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

  const formatHash = (hash: string) => {
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`
  }

  const handleStartRecovery = async (transaction: FailedTransaction, strategy?: RecoveryStrategy) => {
    try {
      await startRecovery(transaction.hash, strategy)
    } catch (error) {
      console.error('Failed to start recovery:', error)
    }
  }

  const handleCancelRecovery = (transaction: FailedTransaction) => {
    cancelRecovery(transaction.hash)
  }

  const handleRetryRecovery = async (transaction: FailedTransaction) => {
    try {
      await retryRecovery(transaction.hash)
    } catch (error) {
      console.error('Failed to retry recovery:', error)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Transaction Recovery</h2>
          <p className="text-muted-foreground">
            Intelligent recovery system for failed transactions
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={clearAllFailedTransactions}
            disabled={state.failedTransactions.length === 0}
          >
            <X className="w-4 h-4 mr-2" />
            Clear All
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
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

      {/* Recovery Statistics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Failed</p>
                <p className="text-2xl font-bold">{stats.totalFailed}</p>
              </div>
              <XCircle className="w-8 h-8 text-red-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Failed transactions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Recoverable</p>
                <p className="text-2xl font-bold text-orange-600">{stats.recoverable}</p>
              </div>
              <RotateCcw className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Can be recovered
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">In Progress</p>
                <p className="text-2xl font-bold text-blue-600">{stats.inProgress}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Currently recovering
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold text-green-600">
                  {stats.successRate.toFixed(1)}%
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Recovery success rate
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="failed">Failed Transactions ({state.failedTransactions.length})</TabsTrigger>
          <TabsTrigger value="recoverable">Recoverable ({state.recoverableTransactions.length})</TabsTrigger>
          <TabsTrigger value="progress">In Progress ({state.recoveryQueue.length})</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="space-y-4">
          {state.failedTransactions.length > 0 ? (
            <div className="space-y-3">
              <AnimatePresence>
                {state.failedTransactions
                  .filter(tx => {
                    switch (activeTab) {
                      case 'failed':
                        return true
                      case 'recoverable':
                        return tx.canRecover
                      case 'progress':
                        return tx.recoveryStatus === RecoveryStatus.RECOVERY_IN_PROGRESS
                      default:
                        return true
                    }
                  })
                  .map((transaction, index) => (
                    <motion.div
                      key={transaction.hash}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -20 }}
                      transition={{ delay: index * 0.05 }}
                    >
                      <Card className="transition-all duration-200 hover:shadow-md">
                        <CardContent className="p-4">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                              <div className="flex items-center gap-2">
                                {getFailureReasonIcon(transaction.failureReason)}
                                {getRecoveryStatusIcon(transaction.recoveryStatus)}
                              </div>
                              
                              <div>
                                <div className="flex items-center gap-2 mb-1">
                                  <p className="font-medium">{formatHash(transaction.hash)}</p>
                                  <Badge variant="outline" className="text-xs">
                                    {transaction.metadata?.type || 'unknown'}
                                  </Badge>
                                  <Badge 
                                    variant={transaction.canRecover ? "default" : "secondary"}
                                    className="text-xs"
                                  >
                                    {transaction.canRecover ? 'Recoverable' : 'Not Recoverable'}
                                  </Badge>
                                </div>
                                
                                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                                  <span>Reason: {getFailureReasonName(transaction.failureReason)}</span>
                                  <span>Status: {getRecoveryStatusName(transaction.recoveryStatus)}</span>
                                  <span>{formatTime(transaction.timestamp)}</span>
                                </div>

                                <p className="text-sm text-muted-foreground mt-1">
                                  {transaction.failureDetails.errorMessage}
                                </p>

                                {transaction.recoveryAttempts.length > 0 && (
                                  <div className="flex items-center gap-1 mt-1">
                                    <RotateCcw className="w-3 h-3 text-blue-500" />
                                    <span className="text-xs text-blue-600">
                                      {transaction.recoveryAttempts.length} recovery attempts
                                    </span>
                                  </div>
                                )}
                              </div>
                            </div>

                            <div className="flex items-center gap-2">
                              {/* Recovery Strategy Info */}
                              {transaction.canRecover && (
                                <div className="text-right mr-4">
                                  <div className="flex items-center gap-1 mb-1">
                                    {getRecoveryTypeIcon(transaction.suggestedFix.type)}
                                    <span className="text-sm font-medium">
                                      {getRecoveryTypeName(transaction.suggestedFix.type)}
                                    </span>
                                  </div>
                                  <div className="text-xs text-muted-foreground">
                                    {transaction.suggestedFix.confidence}% confidence
                                  </div>
                                  <div className="text-xs text-muted-foreground">
                                    Cost: {transaction.suggestedFix.estimatedCost} ETH
                                  </div>
                                </div>
                              )}

                              {/* Action Buttons */}
                              {transaction.recoveryStatus === RecoveryStatus.RECOVERY_AVAILABLE && (
                                <Button
                                  size="sm"
                                  onClick={() => handleStartRecovery(transaction)}
                                  disabled={state.isRecovering}
                                >
                                  <Play className="w-4 h-4 mr-1" />
                                  Recover
                                </Button>
                              )}

                              {transaction.recoveryStatus === RecoveryStatus.RECOVERY_IN_PROGRESS && (
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() => handleCancelRecovery(transaction)}
                                >
                                  <Pause className="w-4 h-4 mr-1" />
                                  Cancel
                                </Button>
                              )}

                              {transaction.recoveryStatus === RecoveryStatus.RECOVERY_FAILED && (
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() => handleRetryRecovery(transaction)}
                                  disabled={state.isRecovering}
                                >
                                  <RotateCcw className="w-4 h-4 mr-1" />
                                  Retry
                                </Button>
                              )}

                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => clearFailedTransaction(transaction.hash)}
                              >
                                <X className="w-4 h-4" />
                              </Button>
                            </div>
                          </div>

                          {/* Recovery Progress */}
                          {transaction.recoveryStatus === RecoveryStatus.RECOVERY_IN_PROGRESS && (
                            <div className="mt-4">
                              <div className="flex items-center justify-between text-sm mb-2">
                                <span>Recovery Progress</span>
                                <span>Processing...</span>
                              </div>
                              <Progress value={50} className="h-2" />
                            </div>
                          )}

                          {/* Recovery Strategy Details */}
                          {transaction.canRecover && selectedTransaction?.hash === transaction.hash && (
                            <motion.div
                              initial={{ opacity: 0, height: 0 }}
                              animate={{ opacity: 1, height: 'auto' }}
                              exit={{ opacity: 0, height: 0 }}
                              className="mt-4 p-4 border rounded-lg bg-muted/50"
                            >
                              <h4 className="font-medium mb-2 flex items-center gap-2">
                                <Lightbulb className="w-4 h-4" />
                                Recovery Strategy: {transaction.suggestedFix.title}
                              </h4>
                              <p className="text-sm text-muted-foreground mb-3">
                                {transaction.suggestedFix.description}
                              </p>
                              
                              <div className="grid gap-2 md:grid-cols-2">
                                <div className="text-sm">
                                  <span className="font-medium">Confidence:</span> {transaction.suggestedFix.confidence}%
                                </div>
                                <div className="text-sm">
                                  <span className="font-medium">Estimated Cost:</span> {transaction.suggestedFix.estimatedCost} ETH
                                </div>
                                <div className="text-sm">
                                  <span className="font-medium">Estimated Time:</span> {transaction.suggestedFix.estimatedTime}s
                                </div>
                                <div className="text-sm">
                                  <span className="font-medium">Risk Level:</span> 
                                  <Badge variant={transaction.suggestedFix.riskLevel === 'low' ? 'default' : 
                                               transaction.suggestedFix.riskLevel === 'medium' ? 'secondary' : 'destructive'}
                                         className="ml-1 text-xs">
                                    {transaction.suggestedFix.riskLevel}
                                  </Badge>
                                </div>
                              </div>

                              {transaction.suggestedFix.warnings && transaction.suggestedFix.warnings.length > 0 && (
                                <div className="mt-3">
                                  <h5 className="text-sm font-medium text-orange-600 mb-1">Warnings:</h5>
                                  <ul className="text-sm text-muted-foreground space-y-1">
                                    {transaction.suggestedFix.warnings.map((warning, idx) => (
                                      <li key={idx} className="flex items-start gap-1">
                                        <AlertTriangle className="w-3 h-3 text-orange-500 mt-0.5 flex-shrink-0" />
                                        {warning}
                                      </li>
                                    ))}
                                  </ul>
                                </div>
                              )}

                              <div className="flex items-center gap-2 mt-4">
                                <Button
                                  size="sm"
                                  onClick={() => handleStartRecovery(transaction)}
                                  disabled={state.isRecovering}
                                >
                                  Execute Recovery
                                </Button>
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() => setSelectedTransaction(null)}
                                >
                                  Close
                                </Button>
                              </div>
                            </motion.div>
                          )}

                          {/* Toggle Details Button */}
                          {transaction.canRecover && (
                            <div className="mt-3 text-center">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => setSelectedTransaction(
                                  selectedTransaction?.hash === transaction.hash ? null : transaction
                                )}
                              >
                                {selectedTransaction?.hash === transaction.hash ? 'Hide Details' : 'Show Recovery Details'}
                              </Button>
                            </div>
                          )}
                        </CardContent>
                      </Card>
                    </motion.div>
                  ))}
              </AnimatePresence>
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <CheckCircle className="w-12 h-12 mx-auto text-green-500 mb-4" />
                <h3 className="text-lg font-medium mb-2">No Failed Transactions</h3>
                <p className="text-muted-foreground">
                  {activeTab === 'failed' && 'All your transactions are successful!'}
                  {activeTab === 'recoverable' && 'No recoverable transactions at the moment'}
                  {activeTab === 'progress' && 'No recovery operations in progress'}
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
