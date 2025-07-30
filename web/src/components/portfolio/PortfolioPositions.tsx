'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  Plus,
  TrendingUp,
  TrendingDown,
  Edit,
  Trash2,
  Eye,
  EyeOff,
  ArrowUpDown,
  Filter,
  Search
} from 'lucide-react'
import { usePortfolioAnalytics } from '@/hooks/usePortfolioAnalytics'
import { type PortfolioPosition } from '@/lib/portfolio-analytics'
import { cn } from '@/lib/utils'
import { toast } from 'sonner'

interface PortfolioPositionsProps {
  showValues?: boolean
}

export function PortfolioPositions({ showValues = true }: PortfolioPositionsProps) {
  const [hideValues, setHideValues] = useState(!showValues)
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [editingPosition, setEditingPosition] = useState<PortfolioPosition | null>(null)
  const [sortBy, setSortBy] = useState<'symbol' | 'value' | 'pnl' | 'allocation'>('value')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc')
  const [filterText, setFilterText] = useState('')
  const [newPosition, setNewPosition] = useState({
    symbol: '',
    amount: '',
    averageCost: ''
  })

  const {
    state,
    addPosition,
    updatePosition,
    removePosition,
    addTransaction
  } = usePortfolioAnalytics({
    autoSync: true,
    trackPriceUpdates: true
  })

  const filteredPositions = state.positions
    .filter(position => 
      position.symbol.toLowerCase().includes(filterText.toLowerCase()) ||
      position.amount > 0
    )
    .sort((a, b) => {
      let aValue: number, bValue: number
      
      switch (sortBy) {
        case 'symbol':
          return sortOrder === 'asc' 
            ? a.symbol.localeCompare(b.symbol)
            : b.symbol.localeCompare(a.symbol)
        case 'value':
          aValue = a.marketValue
          bValue = b.marketValue
          break
        case 'pnl':
          aValue = a.unrealizedPnLPercent
          bValue = b.unrealizedPnLPercent
          break
        case 'allocation':
          aValue = a.allocation
          bValue = b.allocation
          break
        default:
          aValue = a.marketValue
          bValue = b.marketValue
      }
      
      return sortOrder === 'asc' ? aValue - bValue : bValue - aValue
    })

  const formatCurrency = (value: number) => {
    if (hideValues) return '****'
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value)
  }

  const formatPercent = (value: number) => {
    if (hideValues) return '**%'
    const sign = value >= 0 ? '+' : ''
    return `${sign}${value.toFixed(2)}%`
  }

  const formatAmount = (amount: number, decimals = 4) => {
    if (hideValues) return '****'
    return amount.toFixed(decimals)
  }

  const getChangeColor = (value: number) => {
    if (value > 0) return 'text-green-600 dark:text-green-400'
    if (value < 0) return 'text-red-600 dark:text-red-400'
    return 'text-gray-600 dark:text-gray-400'
  }

  const handleAddPosition = () => {
    if (!newPosition.symbol || !newPosition.amount || !newPosition.averageCost) {
      toast.error('Please fill in all fields')
      return
    }

    const amount = parseFloat(newPosition.amount)
    const averageCost = parseFloat(newPosition.averageCost)

    if (isNaN(amount) || isNaN(averageCost) || amount <= 0 || averageCost <= 0) {
      toast.error('Please enter valid amounts')
      return
    }

    addPosition({
      symbol: newPosition.symbol.toUpperCase(),
      amount,
      averageCost,
      currentPrice: averageCost // Will be updated by price feed
    })

    // Add buy transaction
    addTransaction({
      type: 'buy',
      symbol: newPosition.symbol.toUpperCase(),
      amount,
      price: averageCost,
      value: amount * averageCost,
      fee: 0,
      timestamp: Date.now()
    })

    setNewPosition({ symbol: '', amount: '', averageCost: '' })
    setIsAddDialogOpen(false)
  }

  const handleEditPosition = () => {
    if (!editingPosition) return

    const amount = parseFloat((document.getElementById('edit-amount') as HTMLInputElement)?.value || '0')
    const averageCost = parseFloat((document.getElementById('edit-cost') as HTMLInputElement)?.value || '0')

    if (isNaN(amount) || isNaN(averageCost) || amount < 0 || averageCost <= 0) {
      toast.error('Please enter valid amounts')
      return
    }

    updatePosition(editingPosition.symbol, {
      amount,
      averageCost
    })

    setEditingPosition(null)
    setIsEditDialogOpen(false)
  }

  const handleRemovePosition = (symbol: string) => {
    if (window.confirm(`Are you sure you want to remove ${symbol} from your portfolio?`)) {
      removePosition(symbol)
    }
  }

  const handleSort = (field: typeof sortBy) => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')
    } else {
      setSortBy(field)
      setSortOrder('desc')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">Portfolio Positions</h3>
          <p className="text-sm text-muted-foreground">
            Manage your cryptocurrency holdings and track performance
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setHideValues(!hideValues)}
          >
            {hideValues ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
          </Button>
          
          <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
            <DialogTrigger asChild>
              <Button size="sm">
                <Plus className="w-4 h-4 mr-2" />
                Add Position
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Position</DialogTitle>
                <DialogDescription>
                  Add a new cryptocurrency position to your portfolio
                </DialogDescription>
              </DialogHeader>
              
              <div className="space-y-4">
                <div>
                  <Label htmlFor="symbol">Symbol</Label>
                  <Input
                    id="symbol"
                    placeholder="e.g., BTC, ETH"
                    value={newPosition.symbol}
                    onChange={(e) => setNewPosition(prev => ({ ...prev, symbol: e.target.value }))}
                  />
                </div>
                
                <div>
                  <Label htmlFor="amount">Amount</Label>
                  <Input
                    id="amount"
                    type="number"
                    step="0.00000001"
                    placeholder="0.00000000"
                    value={newPosition.amount}
                    onChange={(e) => setNewPosition(prev => ({ ...prev, amount: e.target.value }))}
                  />
                </div>
                
                <div>
                  <Label htmlFor="cost">Average Cost (USD)</Label>
                  <Input
                    id="cost"
                    type="number"
                    step="0.01"
                    placeholder="0.00"
                    value={newPosition.averageCost}
                    onChange={(e) => setNewPosition(prev => ({ ...prev, averageCost: e.target.value }))}
                  />
                </div>
              </div>

              <DialogFooter>
                <Button variant="outline" onClick={() => setIsAddDialogOpen(false)}>
                  Cancel
                </Button>
                <Button onClick={handleAddPosition}>
                  Add Position
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Filters and Search */}
      <div className="flex items-center gap-4">
        <div className="flex-1">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search positions..."
              value={filterText}
              onChange={(e) => setFilterText(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>
        
        <Select value={sortBy} onValueChange={(value: any) => setSortBy(value)}>
          <SelectTrigger className="w-40">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="value">Market Value</SelectItem>
            <SelectItem value="pnl">P&L %</SelectItem>
            <SelectItem value="allocation">Allocation</SelectItem>
            <SelectItem value="symbol">Symbol</SelectItem>
          </SelectContent>
        </Select>
        
        <Button
          variant="outline"
          size="sm"
          onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
        >
          <ArrowUpDown className="w-4 h-4" />
        </Button>
      </div>

      {/* Positions Table */}
      <Card>
        <CardContent className="p-0">
          {filteredPositions.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left p-4 font-medium cursor-pointer hover:bg-muted/50" 
                        onClick={() => handleSort('symbol')}>
                      Asset
                    </th>
                    <th className="text-right p-4 font-medium">Amount</th>
                    <th className="text-right p-4 font-medium">Avg Cost</th>
                    <th className="text-right p-4 font-medium">Current Price</th>
                    <th className="text-right p-4 font-medium cursor-pointer hover:bg-muted/50"
                        onClick={() => handleSort('value')}>
                      Market Value
                    </th>
                    <th className="text-right p-4 font-medium cursor-pointer hover:bg-muted/50"
                        onClick={() => handleSort('pnl')}>
                      P&L
                    </th>
                    <th className="text-right p-4 font-medium cursor-pointer hover:bg-muted/50"
                        onClick={() => handleSort('allocation')}>
                      Allocation
                    </th>
                    <th className="text-right p-4 font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <AnimatePresence>
                    {filteredPositions.map((position, index) => (
                      <motion.tr
                        key={position.symbol}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        transition={{ delay: index * 0.05 }}
                        className="border-b hover:bg-muted/50"
                      >
                        <td className="p-4">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                              <span className="text-sm font-bold">{position.symbol.slice(0, 2)}</span>
                            </div>
                            <div>
                              <p className="font-medium">{position.symbol}</p>
                              <p className="text-xs text-muted-foreground">
                                Since {new Date(position.firstPurchaseDate).toLocaleDateString()}
                              </p>
                            </div>
                          </div>
                        </td>
                        <td className="p-4 text-right font-mono">
                          {formatAmount(position.amount)}
                        </td>
                        <td className="p-4 text-right">
                          {formatCurrency(position.averageCost)}
                        </td>
                        <td className="p-4 text-right">
                          {formatCurrency(position.currentPrice)}
                        </td>
                        <td className="p-4 text-right font-medium">
                          {formatCurrency(position.marketValue)}
                        </td>
                        <td className="p-4 text-right">
                          <div className={cn("font-medium", getChangeColor(position.unrealizedPnL))}>
                            {formatCurrency(position.unrealizedPnL)}
                          </div>
                          <div className={cn("text-sm", getChangeColor(position.unrealizedPnLPercent))}>
                            {formatPercent(position.unrealizedPnLPercent)}
                          </div>
                        </td>
                        <td className="p-4 text-right">
                          <Badge variant="outline">
                            {position.allocation.toFixed(1)}%
                          </Badge>
                        </td>
                        <td className="p-4 text-right">
                          <div className="flex items-center justify-end gap-1">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => {
                                setEditingPosition(position)
                                setIsEditDialogOpen(true)
                              }}
                            >
                              <Edit className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleRemovePosition(position.symbol)}
                            >
                              <Trash2 className="w-4 h-4" />
                            </Button>
                          </div>
                        </td>
                      </motion.tr>
                    ))}
                  </AnimatePresence>
                </tbody>
              </table>
            </div>
          ) : (
            <div className="text-center py-12">
              <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mx-auto mb-4">
                <Plus className="w-8 h-8 text-muted-foreground" />
              </div>
              <h3 className="text-lg font-medium mb-2">No Positions Found</h3>
              <p className="text-muted-foreground mb-4">
                {filterText ? 'No positions match your search criteria' : 'Start by adding your first position'}
              </p>
              <Button onClick={() => setIsAddDialogOpen(true)}>
                <Plus className="w-4 h-4 mr-2" />
                Add Position
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Edit Position Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Position</DialogTitle>
            <DialogDescription>
              Update your {editingPosition?.symbol} position details
            </DialogDescription>
          </DialogHeader>
          
          {editingPosition && (
            <div className="space-y-4">
              <div>
                <Label htmlFor="edit-amount">Amount</Label>
                <Input
                  id="edit-amount"
                  type="number"
                  step="0.00000001"
                  defaultValue={editingPosition.amount.toString()}
                />
              </div>
              
              <div>
                <Label htmlFor="edit-cost">Average Cost (USD)</Label>
                <Input
                  id="edit-cost"
                  type="number"
                  step="0.01"
                  defaultValue={editingPosition.averageCost.toString()}
                />
              </div>
              
              <div className="p-4 bg-muted rounded-lg">
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">Current Price:</span>
                    <p className="font-medium">{formatCurrency(editingPosition.currentPrice)}</p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Market Value:</span>
                    <p className="font-medium">{formatCurrency(editingPosition.marketValue)}</p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Unrealized P&L:</span>
                    <p className={cn("font-medium", getChangeColor(editingPosition.unrealizedPnL))}>
                      {formatCurrency(editingPosition.unrealizedPnL)}
                    </p>
                  </div>
                  <div>
                    <span className="text-muted-foreground">P&L %:</span>
                    <p className={cn("font-medium", getChangeColor(editingPosition.unrealizedPnLPercent))}>
                      {formatPercent(editingPosition.unrealizedPnLPercent)}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleEditPosition}>
              Update Position
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
