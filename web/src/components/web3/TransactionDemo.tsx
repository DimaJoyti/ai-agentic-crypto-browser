'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { 
  Play, 
  Zap, 
  ArrowUpRight, 
  Repeat, 
  Shield,
  Coins,
  TestTube
} from 'lucide-react'
import { useTransactionMonitor } from '@/hooks/useTransactionMonitor'
import { TransactionType } from '@/lib/transaction-monitor'
import { SUPPORTED_CHAINS } from '@/lib/chains'
import { toast } from 'sonner'

export function TransactionDemo() {
  const [isGenerating, setIsGenerating] = useState(false)
  const [selectedChain, setSelectedChain] = useState('1')
  const [selectedType, setSelectedType] = useState<TransactionType>(TransactionType.SEND)
  const [amount, setAmount] = useState('0.1')
  const [recipient, setRecipient] = useState('0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1')

  const { trackTransaction } = useTransactionMonitor()

  // Generate a mock transaction hash
  const generateMockHash = (): `0x${string}` => {
    const chars = '0123456789abcdef'
    let hash = '0x'
    for (let i = 0; i < 64; i++) {
      hash += chars[Math.floor(Math.random() * chars.length)]
    }
    return hash as `0x${string}`
  }

  const generateDemoTransaction = async () => {
    setIsGenerating(true)

    try {
      const chainId = parseInt(selectedChain)
      const chain = SUPPORTED_CHAINS[chainId]
      const hash = generateMockHash()

      // Create metadata based on transaction type
      const metadata = {
        description: getTransactionDescription(selectedType, amount, chain?.gasToken || 'ETH'),
        tokenSymbol: chain?.gasToken,
        tokenAmount: selectedType === TransactionType.SEND ? amount : undefined,
        contractAddress: selectedType !== TransactionType.SEND ? recipient : undefined,
        methodName: getMethodName(selectedType)
      }

      // Track the demo transaction
      await trackTransaction(hash, chainId, selectedType, metadata)

      toast.success('Demo transaction created!', {
        description: `Tracking ${selectedType} transaction on ${chain?.shortName || 'Unknown Chain'}`
      })

      // Simulate transaction progression for demo purposes
      simulateTransactionProgress(hash, chainId)

    } catch (error) {
      toast.error('Failed to create demo transaction')
    } finally {
      setIsGenerating(false)
    }
  }

  const simulateTransactionProgress = (hash: `0x${string}`, chainId: number) => {
    // This is for demo purposes only - simulates transaction confirmation
    // In a real app, the transaction monitor would handle this automatically
    
    setTimeout(() => {
      // Simulate random success/failure for demo
      const willSucceed = Math.random() > 0.2 // 80% success rate
      
      if (willSucceed) {
        console.log(`Demo: Transaction ${hash} would be confirmed`)
      } else {
        console.log(`Demo: Transaction ${hash} would fail`)
      }
    }, 5000 + Math.random() * 10000) // Random delay between 5-15 seconds
  }

  const getTransactionDescription = (type: TransactionType, amount: string, token: string): string => {
    switch (type) {
      case TransactionType.SEND:
        return `Send ${amount} ${token}`
      case TransactionType.SWAP:
        return `Swap ${amount} ${token} for USDC`
      case TransactionType.APPROVE:
        return `Approve ${token} spending`
      case TransactionType.STAKE:
        return `Stake ${amount} ${token}`
      case TransactionType.UNSTAKE:
        return `Unstake ${amount} ${token}`
      case TransactionType.MINT:
        return `Mint NFT`
      default:
        return `${type} transaction`
    }
  }

  const getMethodName = (type: TransactionType): string => {
    switch (type) {
      case TransactionType.SEND:
        return 'transfer'
      case TransactionType.SWAP:
        return 'swapExactTokensForTokens'
      case TransactionType.APPROVE:
        return 'approve'
      case TransactionType.STAKE:
        return 'stake'
      case TransactionType.UNSTAKE:
        return 'unstake'
      case TransactionType.MINT:
        return 'mint'
      default:
        return 'unknown'
    }
  }

  const getTypeIcon = (type: TransactionType) => {
    switch (type) {
      case TransactionType.SEND:
        return <ArrowUpRight className="w-4 h-4" />
      case TransactionType.SWAP:
        return <Repeat className="w-4 h-4" />
      case TransactionType.APPROVE:
        return <Shield className="w-4 h-4" />
      case TransactionType.STAKE:
      case TransactionType.UNSTAKE:
        return <Zap className="w-4 h-4" />
      case TransactionType.MINT:
        return <Coins className="w-4 h-4" />
      default:
        return <Play className="w-4 h-4" />
    }
  }

  const transactionTypes = [
    { value: TransactionType.SEND, label: 'Send Tokens' },
    { value: TransactionType.SWAP, label: 'Token Swap' },
    { value: TransactionType.APPROVE, label: 'Token Approval' },
    { value: TransactionType.STAKE, label: 'Stake Tokens' },
    { value: TransactionType.UNSTAKE, label: 'Unstake Tokens' },
    { value: TransactionType.MINT, label: 'Mint NFT' }
  ]

  const supportedChains = Object.values(SUPPORTED_CHAINS).filter(chain => 
    !chain.isTestnet || chain.id === 11155111 // Include Sepolia for testing
  )

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <TestTube className="w-5 h-5" />
          Transaction Demo
        </CardTitle>
        <CardDescription>
          Generate demo transactions to test the real-time tracking system
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label>Blockchain Network</Label>
            <Select value={selectedChain} onValueChange={setSelectedChain}>
              <SelectTrigger>
                <SelectValue placeholder="Select chain" />
              </SelectTrigger>
              <SelectContent>
                {supportedChains.map(chain => (
                  <SelectItem key={chain.id} value={chain.id.toString()}>
                    <div className="flex items-center gap-2">
                      <span>{chain.icon}</span>
                      <span>{chain.name}</span>
                      <Badge variant="outline" className="text-xs">
                        {chain.shortName}
                      </Badge>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Transaction Type</Label>
            <Select value={selectedType} onValueChange={(value) => setSelectedType(value as TransactionType)}>
              <SelectTrigger>
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                {transactionTypes.map(type => (
                  <SelectItem key={type.value} value={type.value}>
                    <div className="flex items-center gap-2">
                      {getTypeIcon(type.value)}
                      <span>{type.label}</span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {selectedType === TransactionType.SEND && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Amount</Label>
              <Input
                type="number"
                step="0.001"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="0.1"
              />
            </div>

            <div className="space-y-2">
              <Label>Recipient Address</Label>
              <Input
                value={recipient}
                onChange={(e) => setRecipient(e.target.value)}
                placeholder="0x..."
              />
            </div>
          </div>
        )}

        <div className="bg-muted p-4 rounded-lg">
          <h4 className="font-medium mb-2">Demo Transaction Preview:</h4>
          <div className="space-y-1 text-sm text-muted-foreground">
            <p>• Network: {SUPPORTED_CHAINS[parseInt(selectedChain)]?.name}</p>
            <p>• Type: {transactionTypes.find(t => t.value === selectedType)?.label}</p>
            <p>• Description: {getTransactionDescription(selectedType, amount, SUPPORTED_CHAINS[parseInt(selectedChain)]?.gasToken || 'ETH')}</p>
            <p>• Status: Will simulate pending → confirmed/failed</p>
          </div>
        </div>

        <Button 
          onClick={generateDemoTransaction}
          disabled={isGenerating}
          className="w-full"
          size="lg"
        >
          {isGenerating ? (
            <>
              <Play className="w-4 h-4 mr-2 animate-pulse" />
              Generating Transaction...
            </>
          ) : (
            <>
              <Play className="w-4 h-4 mr-2" />
              Generate Demo Transaction
            </>
          )}
        </Button>

        <div className="text-xs text-muted-foreground bg-blue-50 dark:bg-blue-950 p-3 rounded-lg">
          <p className="font-medium mb-1">Note:</p>
          <p>
            This creates mock transactions for demonstration purposes. 
            Real transactions would be tracked automatically when you interact with your wallet.
            The demo will simulate transaction confirmation after 5-15 seconds.
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
