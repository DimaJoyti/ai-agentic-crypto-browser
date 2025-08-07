'use client'

import { useState } from 'react'
import { useConnect, useAccount } from 'wagmi'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Wallet, Monitor, Smartphone } from 'lucide-react'

interface SimpleWalletModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: (address: string, chainId: number) => void
}

export function SimpleWalletModal({ isOpen, onClose, onSuccess }: SimpleWalletModalProps) {
  const [isConnecting, setIsConnecting] = useState(false)
  const { connect, connectors, isPending } = useConnect()
  const { address, isConnected } = useAccount()

  const handleConnect = async (connector: any) => {
    try {
      setIsConnecting(true)
      await connect({ connector })
      
      if (isConnected && address) {
        onSuccess?.(address, 1) // Default to mainnet
        onClose()
      }
    } catch (error) {
      console.error('Connection failed:', error)
    } finally {
      setIsConnecting(false)
    }
  }

  const getConnectorIcon = (connectorName: string) => {
    switch (connectorName.toLowerCase()) {
      case 'metamask':
        return <Monitor className="w-6 h-6" />
      case 'walletconnect':
        return <Smartphone className="w-6 h-6" />
      default:
        return <Wallet className="w-6 h-6" />
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Wallet className="w-5 h-5" />
            Connect Wallet
          </DialogTitle>
          <DialogDescription>
            Choose a wallet to connect to the application
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-3">
          {connectors.map((connector) => (
            <Card
              key={connector.id}
              className="cursor-pointer transition-all hover:shadow-md border-2 border-transparent hover:border-primary"
              onClick={() => handleConnect(connector)}
            >
              <CardContent className="p-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-secondary rounded-lg flex items-center justify-center">
                    {getConnectorIcon(connector.name)}
                  </div>
                  <div className="flex-1">
                    <h4 className="font-medium">{connector.name}</h4>
                    <p className="text-sm text-muted-foreground">
                      Connect using {connector.name}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {(isPending || isConnecting) && (
          <div className="text-center py-4">
            <p className="text-sm text-muted-foreground">Connecting...</p>
          </div>
        )}

        <div className="text-center">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
