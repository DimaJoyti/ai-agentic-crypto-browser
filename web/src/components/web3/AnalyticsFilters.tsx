'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { 
  Filter, 
  X, 
  Calendar as CalendarIcon,
  RotateCcw
} from 'lucide-react'
import { format } from 'date-fns'
import { 
  TransactionAnalytics,
  type AnalyticsFilters as Filters
} from '@/lib/transaction-analytics'
import { TransactionStatus, TransactionType } from '@/lib/transaction-monitor'
import { SUPPORTED_CHAINS } from '@/lib/chains'

interface AnalyticsFiltersProps {
  filters: Filters
  onFiltersChange: (filters: Filters) => void
}

export function AnalyticsFilters({ filters, onFiltersChange }: AnalyticsFiltersProps) {
  const [dateRange, setDateRange] = useState<{ from?: Date; to?: Date } | undefined>({
    from: filters.dateRange?.start,
    to: filters.dateRange?.end
  })

  const timeframes = TransactionAnalytics.getTimeframes()
  const supportedChains = Object.values(SUPPORTED_CHAINS).filter(chain => 
    !chain.isTestnet || chain.id === 11155111
  )

  const transactionTypes = [
    TransactionType.SEND,
    TransactionType.RECEIVE,
    TransactionType.SWAP,
    TransactionType.APPROVE,
    TransactionType.STAKE,
    TransactionType.UNSTAKE,
    TransactionType.MINT,
    TransactionType.BURN,
    TransactionType.CONTRACT_INTERACTION
  ]

  const transactionStatuses = [
    TransactionStatus.PENDING,
    TransactionStatus.CONFIRMED,
    TransactionStatus.FAILED,
    TransactionStatus.DROPPED
  ]

  const handleTimeframeChange = (timeframeLabel: string) => {
    const timeframe = timeframes.find(t => t.label === timeframeLabel)
    if (timeframe) {
      onFiltersChange({
        ...filters,
        timeframe,
        dateRange: {
          start: new Date(Date.now() - timeframe.days * 24 * 60 * 60 * 1000),
          end: new Date()
        }
      })
    }
  }

  const handleChainToggle = (chainId: number, checked: boolean) => {
    const newChains = checked
      ? [...filters.chains, chainId]
      : filters.chains.filter(id => id !== chainId)
    
    onFiltersChange({
      ...filters,
      chains: newChains
    })
  }

  const handleTypeToggle = (type: TransactionType, checked: boolean) => {
    const newTypes = checked
      ? [...filters.types, type]
      : filters.types.filter(t => t !== type)
    
    onFiltersChange({
      ...filters,
      types: newTypes
    })
  }

  const handleStatusToggle = (status: TransactionStatus, checked: boolean) => {
    const newStatuses = checked
      ? [...filters.status, status]
      : filters.status.filter(s => s !== status)
    
    onFiltersChange({
      ...filters,
      status: newStatuses
    })
  }

  const handleAmountChange = (field: 'minAmount' | 'maxAmount', value: string) => {
    const numValue = value === '' ? undefined : parseFloat(value)
    onFiltersChange({
      ...filters,
      [field]: numValue
    })
  }

  const handleDateRangeChange = (range: { from?: Date; to?: Date } | undefined) => {
    setDateRange(range)
    if (range?.from && range?.to) {
      onFiltersChange({
        ...filters,
        dateRange: {
          start: range.from,
          end: range.to
        }
      })
    }
  }

  const resetFilters = () => {
    const defaultFilters = TransactionAnalytics.getDefaultFilters()
    onFiltersChange(defaultFilters)
    setDateRange({
      from: defaultFilters.dateRange?.start,
      to: defaultFilters.dateRange?.end
    })
  }

  const getActiveFiltersCount = () => {
    let count = 0
    if (filters.chains.length > 0) count++
    if (filters.types.length > 0) count++
    if (filters.status.length > 0) count++
    if (filters.minAmount !== undefined) count++
    if (filters.maxAmount !== undefined) count++
    return count
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Filter className="w-5 h-5" />
              Analytics Filters
            </CardTitle>
            <CardDescription>
              Customize your transaction analytics view
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            {getActiveFiltersCount() > 0 && (
              <Badge variant="secondary">
                {getActiveFiltersCount()} active
              </Badge>
            )}
            <Button variant="outline" size="sm" onClick={resetFilters}>
              <RotateCcw className="w-4 h-4 mr-2" />
              Reset
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Timeframe */}
        <div className="space-y-2">
          <Label>Timeframe</Label>
          <Select value={filters.timeframe.label} onValueChange={handleTimeframeChange}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {timeframes.map(timeframe => (
                <SelectItem key={timeframe.label} value={timeframe.label}>
                  {timeframe.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Custom Date Range */}
        <div className="space-y-2">
          <Label>Custom Date Range</Label>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" className="w-full justify-start text-left font-normal">
                <CalendarIcon className="mr-2 h-4 w-4" />
                {dateRange?.from ? (
                  dateRange?.to ? (
                    <>
                      {format(dateRange.from, "LLL dd, y")} -{" "}
                      {format(dateRange.to, "LLL dd, y")}
                    </>
                  ) : (
                    format(dateRange.from, "LLL dd, y")
                  )
                ) : (
                  <span>Pick a date range</span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="start">
              <Calendar
                initialFocus
                mode="range"
                defaultMonth={dateRange?.from}
                selected={dateRange?.from && dateRange?.to ? { from: dateRange.from, to: dateRange.to } : undefined}
                onSelect={handleDateRangeChange}
                numberOfMonths={2}
              />
            </PopoverContent>
          </Popover>
        </div>

        {/* Amount Range */}
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Min Amount (ETH)</Label>
            <Input
              type="number"
              step="0.001"
              placeholder="0.0"
              value={filters.minAmount || ''}
              onChange={(e) => handleAmountChange('minAmount', e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label>Max Amount (ETH)</Label>
            <Input
              type="number"
              step="0.001"
              placeholder="No limit"
              value={filters.maxAmount || ''}
              onChange={(e) => handleAmountChange('maxAmount', e.target.value)}
            />
          </div>
        </div>

        {/* Chains */}
        <div className="space-y-3">
          <Label>Blockchain Networks</Label>
          <div className="grid grid-cols-2 gap-2">
            {supportedChains.map(chain => (
              <div key={chain.id} className="flex items-center space-x-2">
                <Checkbox
                  id={`chain-${chain.id}`}
                  checked={filters.chains.includes(chain.id)}
                  onCheckedChange={(checked) => handleChainToggle(chain.id, checked as boolean)}
                />
                <Label htmlFor={`chain-${chain.id}`} className="text-sm">
                  {chain.shortName}
                </Label>
              </div>
            ))}
          </div>
          {filters.chains.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {filters.chains.map(chainId => {
                const chain = SUPPORTED_CHAINS[chainId]
                return (
                  <Badge key={chainId} variant="secondary" className="text-xs">
                    {chain?.shortName || chainId}
                    <X 
                      className="w-3 h-3 ml-1 cursor-pointer" 
                      onClick={() => handleChainToggle(chainId, false)}
                    />
                  </Badge>
                )
              })}
            </div>
          )}
        </div>

        {/* Transaction Types */}
        <div className="space-y-3">
          <Label>Transaction Types</Label>
          <div className="grid grid-cols-2 gap-2">
            {transactionTypes.map(type => (
              <div key={type} className="flex items-center space-x-2">
                <Checkbox
                  id={`type-${type}`}
                  checked={filters.types.includes(type)}
                  onCheckedChange={(checked) => handleTypeToggle(type, checked as boolean)}
                />
                <Label htmlFor={`type-${type}`} className="text-sm capitalize">
                  {type.replace('_', ' ')}
                </Label>
              </div>
            ))}
          </div>
          {filters.types.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {filters.types.map(type => (
                <Badge key={type} variant="secondary" className="text-xs">
                  {type.replace('_', ' ')}
                  <X 
                    className="w-3 h-3 ml-1 cursor-pointer" 
                    onClick={() => handleTypeToggle(type, false)}
                  />
                </Badge>
              ))}
            </div>
          )}
        </div>

        {/* Transaction Status */}
        <div className="space-y-3">
          <Label>Transaction Status</Label>
          <div className="grid grid-cols-2 gap-2">
            {transactionStatuses.map(status => (
              <div key={status} className="flex items-center space-x-2">
                <Checkbox
                  id={`status-${status}`}
                  checked={filters.status.includes(status)}
                  onCheckedChange={(checked) => handleStatusToggle(status, checked as boolean)}
                />
                <Label htmlFor={`status-${status}`} className="text-sm capitalize">
                  {status}
                </Label>
              </div>
            ))}
          </div>
          {filters.status.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {filters.status.map(status => (
                <Badge key={status} variant="secondary" className="text-xs">
                  {status}
                  <X 
                    className="w-3 h-3 ml-1 cursor-pointer" 
                    onClick={() => handleStatusToggle(status, false)}
                  />
                </Badge>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
