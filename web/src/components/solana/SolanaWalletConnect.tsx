'use client'

import React, { useState, useCallback } from 'react'
import { useWallet, useConnection } from '@solana/wallet-adapter-react'
import { WalletMultiButton, WalletDisconnectButton } from '@solana/wallet-adapter-react-ui'
import { LAMPORTS_PER_SOL, PublicKey } from '@solana/web3.js'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Separator } from '@/components/ui/separator'
import { 
  Wallet, 
  Copy, 
  ExternalLink, 
  RefreshCw, 
  CheckCircle, 
  AlertCircle,
  Loader2,
  Eye,
  EyeOff
} from 'lucide-react'
import { toast } from 'sonner'
import { formatCurrency, formatNumber, cn } from '@/lib/utils'
import { useSolanaWallet } from '@/hooks/useSolanaWallet'

interface SolanaWalletConnectProps {
  className?: string
  showBalance?: boolean
  showTokens?: boolean
  compact?: boolean
}

export function SolanaWalletConnect({ 
  className, 
  showBalance = true, 
  showTokens = false,
  compact = false 
}: SolanaWalletConnectProps) {
  const { connection } = useConnection()
  const { 
    publicKey, 
    connected, 
    connecting, 
    disconnecting, 
    wallet,
    connect,
    disconnect
  } = useWallet()

  const {
    balance,
    tokens,
    isLoading: walletLoading,
    error: walletError,
    refresh: refreshWallet
  } = useSolanaWallet()

  const [showPrivateInfo, setShowPrivateInfo] = useState(false)
  const [isRefreshing, setIsRefreshing] = useState(false)

  const handleCopyAddress = useCallback(() => {
    if (publicKey) {
      navigator.clipboard.writeText(publicKey.toString())
      toast.success('Address copied to clipboard')
    }
  }, [publicKey])

  const handleViewOnExplorer = useCallback(() => {
    if (publicKey) {
      const url = `https://explorer.solana.com/address/${publicKey.toString()}`
      window.open(url, '_blank')
    }
  }, [publicKey])

  const handleRefresh = useCallback(async () => {
    setIsRefreshing(true)
    try {
      await refreshWallet()
      toast.success('Wallet data refreshed')
    } catch (error) {
      toast.error('Failed to refresh wallet data')
    } finally {
      setIsRefreshing(false)
    }
  }, [refreshWallet])

  const formatAddress = (address: string, length: number = 8) => {
    return `${address.slice(0, length)}...${address.slice(-length)}`
  }

  if (compact) {
    return (
      <div className={cn('flex items-center space-x-2', className)}>
        {connected ? (
          <>
            <Badge variant="outline" className="text-green-600 border-green-600">
              <CheckCircle className="h-3 w-3 mr-1" />
              Connected
            </Badge>
            <WalletDisconnectButton />
          </>
        ) : (
          <WalletMultiButton />
        )}
      </div>
    )
  }

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <div className="flex items-center">
            <Wallet className="h-5 w-5 mr-2" />
            Solana Wallet
          </div>
          {connected && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleRefresh}
              disabled={isRefreshing}
            >
              <RefreshCw className={cn("h-4 w-4", isRefreshing && "animate-spin")} />
            </Button>
          )}
        </CardTitle>
        <CardDescription>
          Connect your Solana wallet to start trading
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Connection Status */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Status:</span>
          <div className="flex items-center space-x-2">
            {connecting && (
              <Badge variant="outline" className="text-yellow-600 border-yellow-600">
                <Loader2 className="h-3 w-3 mr-1 animate-spin" />
                Connecting
              </Badge>
            )}
            {connected && (
              <Badge variant="outline" className="text-green-600 border-green-600">
                <CheckCircle className="h-3 w-3 mr-1" />
                Connected
              </Badge>
            )}
            {!connected && !connecting && (
              <Badge variant="outline" className="text-gray-600 border-gray-600">
                <AlertCircle className="h-3 w-3 mr-1" />
                Disconnected
              </Badge>
            )}
          </div>
        </div>

        {/* Wallet Info */}
        {connected && wallet && (
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Wallet:</span>
              <div className="flex items-center space-x-2">
                <img 
                  src={wallet.adapter.icon} 
                  alt={wallet.adapter.name}
                  className="h-4 w-4"
                />
                <span className="text-sm font-medium">{wallet.adapter.name}</span>
              </div>
            </div>

            {publicKey && (
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Address:</span>
                  <div className="flex items-center space-x-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setShowPrivateInfo(!showPrivateInfo)}
                    >
                      {showPrivateInfo ? <EyeOff className="h-3 w-3" /> : <Eye className="h-3 w-3" />}
                    </Button>
                  </div>
                </div>
                
                {showPrivateInfo && (
                  <div className="flex items-center space-x-2 p-2 bg-muted rounded-md">
                    <code className="text-xs flex-1 break-all">
                      {publicKey.toString()}
                    </code>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleCopyAddress}
                    >
                      <Copy className="h-3 w-3" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleViewOnExplorer}
                    >
                      <ExternalLink className="h-3 w-3" />
                    </Button>
                  </div>
                )}
              </div>
            )}

            {/* Balance */}
            {showBalance && (
              <>
                <Separator />
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">SOL Balance:</span>
                    <div className="text-right">
                      {walletLoading ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <div>
                          <div className="font-semibold">
                            {formatNumber(balance, 4)} SOL
                          </div>
                          <div className="text-xs text-muted-foreground">
                            ≈ {formatCurrency(balance * 100)} {/* Placeholder price */}
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </>
            )}

            {/* Token Holdings */}
            {showTokens && tokens.length > 0 && (
              <>
                <Separator />
                <div className="space-y-2">
                  <span className="text-sm text-muted-foreground">Token Holdings:</span>
                  <div className="space-y-1 max-h-32 overflow-y-auto">
                    {tokens.slice(0, 5).map((token, index) => (
                      <div key={index} className="flex items-center justify-between text-xs">
                        <span>{token.symbol}</span>
                        <span>{formatNumber(token.balance, 2)}</span>
                      </div>
                    ))}
                    {tokens.length > 5 && (
                      <div className="text-xs text-muted-foreground text-center">
                        +{tokens.length - 5} more tokens
                      </div>
                    )}
                  </div>
                </div>
              </>
            )}
          </div>
        )}

        {/* Error Display */}
        {walletError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{walletError}</AlertDescription>
          </Alert>
        )}

        {/* Connection Buttons */}
        <div className="flex space-x-2">
          {!connected ? (
            <WalletMultiButton className="flex-1" />
          ) : (
            <WalletDisconnectButton className="flex-1" />
          )}
        </div>

        {/* Network Info */}
        <div className="text-xs text-muted-foreground text-center">
          Connected to Solana Mainnet
        </div>
      </CardContent>
    </Card>
  )
}
