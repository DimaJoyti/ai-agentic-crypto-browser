'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Package, 
  Plus, 
  Trash2, 
  Play, 
  Clock, 
  DollarSign,
  Zap,
  CheckCircle,
  AlertTriangle,
  Target,
  ArrowRight,
  Layers
} from 'lucide-react'
import { toast } from 'sonner'

interface PendingTransaction {
  id: string
  type: 'send' | 'swap' | 'approve' | 'stake'
  description: string
  gasEstimate: string
  value?: string
  to?: string
  priority: 'low' | 'medium' | 'high'
  canBatch: boolean
}

interface TransactionBatch {
  id: string
  transactions: PendingTransaction[]
  estimatedGasSavings: string
  totalGasCost: string
  status: 'draft' | 'ready' | 'executing' | 'completed' | 'failed'
  createdAt: number
}

export function TransactionBatcher() {
  const [pendingTransactions] = useState<PendingTransaction[]>([
    {
      id: '1',
      type: 'approve',
      description: 'Approve USDC for Uniswap',
      gasEstimate: '0.002',
      priority: 'medium',
      canBatch: true
    },
    {
      id: '2',
      type: 'swap',
      description: 'Swap 100 USDC for ETH',
      gasEstimate: '0.008',
      value: '100',
      priority: 'high',
      canBatch: true
    },
    {
      id: '3',
      type: 'send',
      description: 'Send 0.5 ETH to wallet',
      gasEstimate: '0.001',
      value: '0.5',
      to: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1',
      priority: 'low',
      canBatch: false
    },
    {
      id: '4',
      type: 'stake',
      description: 'Stake 10 ETH in validator',
      gasEstimate: '0.003',
      value: '10',
      priority: 'medium',
      canBatch: true
    }
  ])

  const [selectedTransactions, setSelectedTransactions] = useState<string[]>([])
  const [batches, setBatches] = useState<TransactionBatch[]>([])

  const batchableTransactions = pendingTransactions.filter(tx => tx.canBatch)
  const selectedBatchable = selectedTransactions.filter(id => 
    batchableTransactions.some(tx => tx.id === id)
  )

  const calculateBatchSavings = (transactionIds: string[]) => {
    const transactions = pendingTransactions.filter(tx => transactionIds.includes(tx.id))
    const individualCost = transactions.reduce((sum, tx) => sum + parseFloat(tx.gasEstimate), 0)
    const batchCost = individualCost * 0.7 // Assume 30% savings
    const savings = individualCost - batchCost
    
    return {
      individualCost: individualCost.toFixed(4),
      batchCost: batchCost.toFixed(4),
      savings: savings.toFixed(4),
      savingsPercentage: Math.round((savings / individualCost) * 100)
    }
  }

  const handleTransactionSelect = (transactionId: string, checked: boolean) => {
    if (checked) {
      setSelectedTransactions(prev => [...prev, transactionId])
    } else {
      setSelectedTransactions(prev => prev.filter(id => id !== transactionId))
    }
  }

  const createBatch = () => {
    if (selectedBatchable.length < 2) {
      toast.error('Select at least 2 batchable transactions')
      return
    }

    const transactions = pendingTransactions.filter(tx => selectedBatchable.includes(tx.id))
    const savings = calculateBatchSavings(selectedBatchable)

    const newBatch: TransactionBatch = {
      id: Date.now().toString(),
      transactions,
      estimatedGasSavings: savings.savings,
      totalGasCost: savings.batchCost,
      status: 'draft',
      createdAt: Date.now()
    }

    setBatches(prev => [...prev, newBatch])
    setSelectedTransactions([])
    toast.success(`Batch created with ${savings.savingsPercentage}% gas savings`)
  }

  const executeBatch = (batchId: string) => {
    setBatches(prev => prev.map(batch => 
      batch.id === batchId 
        ? { ...batch, status: 'executing' }
        : batch
    ))

    // Simulate batch execution
    setTimeout(() => {
      setBatches(prev => prev.map(batch => 
        batch.id === batchId 
          ? { ...batch, status: 'completed' }
          : batch
      ))
      toast.success('Batch executed successfully!')
    }, 3000)
  }

  const deleteBatch = (batchId: string) => {
    setBatches(prev => prev.filter(batch => batch.id !== batchId))
    toast.success('Batch deleted')
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'send':
        return <ArrowRight className="w-4 h-4" />
      case 'swap':
        return <Zap className="w-4 h-4" />
      case 'approve':
        return <CheckCircle className="w-4 h-4" />
      case 'stake':
        return <Target className="w-4 h-4" />
      default:
        return <Package className="w-4 h-4" />
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'bg-red-100 text-red-800'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800'
      case 'low':
        return 'bg-green-100 text-green-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'draft':
        return <Clock className="w-4 h-4 text-gray-500" />
      case 'ready':
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'executing':
        return <Play className="w-4 h-4 text-blue-500 animate-pulse" />
      case 'completed':
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'failed':
        return <AlertTriangle className="w-4 h-4 text-red-500" />
      default:
        return <Clock className="w-4 h-4 text-gray-500" />
    }
  }

  const selectedSavings = selectedBatchable.length >= 2 
    ? calculateBatchSavings(selectedBatchable)
    : null

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold flex items-center gap-2">
          <Layers className="w-6 h-6" />
          Transaction Batcher
        </h2>
        <p className="text-muted-foreground">
          Combine multiple transactions to save on gas costs
        </p>
      </div>

      {/* Pending Transactions */}
      <Card>
        <CardHeader>
          <CardTitle>Pending Transactions</CardTitle>
          <CardDescription>
            Select transactions to batch together for gas savings
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {pendingTransactions.map((transaction) => (
              <motion.div
                key={transaction.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className={`border rounded-lg p-4 ${
                  selectedTransactions.includes(transaction.id)
                    ? 'ring-2 ring-blue-500 bg-blue-50 dark:bg-blue-950'
                    : ''
                } ${!transaction.canBatch ? 'opacity-60' : ''}`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {transaction.canBatch && (
                      <Checkbox
                        checked={selectedTransactions.includes(transaction.id)}
                        onCheckedChange={(checked) => 
                          handleTransactionSelect(transaction.id, checked as boolean)
                        }
                      />
                    )}
                    {getTypeIcon(transaction.type)}
                    <div>
                      <h4 className="font-medium">{transaction.description}</h4>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge 
                          variant="secondary" 
                          className={getPriorityColor(transaction.priority)}
                        >
                          {transaction.priority}
                        </Badge>
                        {!transaction.canBatch && (
                          <Badge variant="outline" className="text-xs">
                            Not batchable
                          </Badge>
                        )}
                      </div>
                    </div>
                  </div>
                  
                  <div className="text-right">
                    <p className="font-medium">{transaction.gasEstimate} ETH</p>
                    {transaction.value && (
                      <p className="text-sm text-muted-foreground">
                        {transaction.value} {transaction.type === 'send' ? 'ETH' : 'tokens'}
                      </p>
                    )}
                  </div>
                </div>
              </motion.div>
            ))}
          </div>

          {/* Batch Creation */}
          {selectedBatchable.length > 0 && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              className="mt-6 p-4 bg-muted rounded-lg"
            >
              <h4 className="font-medium mb-3">Batch Preview</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Selected transactions:</span>
                  <span>{selectedBatchable.length}</span>
                </div>
                {selectedSavings && (
                  <>
                    <div className="flex justify-between">
                      <span>Individual cost:</span>
                      <span>{selectedSavings.individualCost} ETH</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Batch cost:</span>
                      <span>{selectedSavings.batchCost} ETH</span>
                    </div>
                    <div className="flex justify-between font-medium text-green-600">
                      <span>Estimated savings:</span>
                      <span>{selectedSavings.savings} ETH ({selectedSavings.savingsPercentage}%)</span>
                    </div>
                  </>
                )}
              </div>
              
              <div className="flex gap-2 mt-4">
                <Button 
                  onClick={createBatch}
                  disabled={selectedBatchable.length < 2}
                  className="flex-1"
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Create Batch
                </Button>
                <Button 
                  variant="outline" 
                  onClick={() => setSelectedTransactions([])}
                >
                  Clear
                </Button>
              </div>
            </motion.div>
          )}
        </CardContent>
      </Card>

      {/* Active Batches */}
      {batches.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Transaction Batches</CardTitle>
            <CardDescription>
              Manage and execute your transaction batches
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <AnimatePresence>
                {batches.map((batch) => (
                  <motion.div
                    key={batch.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className="border rounded-lg p-4"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(batch.status)}
                        <h4 className="font-medium">
                          Batch #{batch.id.slice(-4)}
                        </h4>
                        <Badge variant="outline" className="capitalize">
                          {batch.status}
                        </Badge>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        {batch.status === 'draft' && (
                          <Button
                            size="sm"
                            onClick={() => executeBatch(batch.id)}
                          >
                            <Play className="w-3 h-3 mr-2" />
                            Execute
                          </Button>
                        )}
                        {batch.status !== 'executing' && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => deleteBatch(batch.id)}
                          >
                            <Trash2 className="w-3 h-3" />
                          </Button>
                        )}
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-3">
                      <div>
                        <p className="text-sm text-muted-foreground">Transactions</p>
                        <p className="font-medium">{batch.transactions.length}</p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">Total Cost</p>
                        <p className="font-medium">{batch.totalGasCost} ETH</p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">Savings</p>
                        <p className="font-medium text-green-600">
                          {batch.estimatedGasSavings} ETH
                        </p>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <p className="text-sm font-medium">Included Transactions:</p>
                      {batch.transactions.map((tx, index) => (
                        <div key={tx.id} className="flex items-center gap-2 text-sm">
                          <span className="text-muted-foreground">{index + 1}.</span>
                          {getTypeIcon(tx.type)}
                          <span>{tx.description}</span>
                        </div>
                      ))}
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Information */}
      <Alert>
        <Package className="h-4 w-4" />
        <AlertDescription>
          <strong>Transaction Batching Benefits:</strong> Combine multiple transactions into a single batch to save 20-50% on gas costs. 
          Not all transaction types can be batched together. Batching works best with DeFi operations like approvals and swaps.
        </AlertDescription>
      </Alert>
    </div>
  )
}
