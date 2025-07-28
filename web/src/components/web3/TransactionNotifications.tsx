'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { 
  Bell, 
  X, 
  ExternalLink, 
  Clock, 
  CheckCircle, 
  XCircle, 
  AlertTriangle,
  Settings,
  Volume2,
  VolumeX
} from 'lucide-react'
import { useTransactionMonitor } from '@/hooks/useTransactionMonitor'
import { TransactionStatus, type TransactionData } from '@/lib/transaction-monitor'
import { SUPPORTED_CHAINS } from '@/lib/chains'

interface NotificationSettings {
  enabled: boolean
  sound: boolean
  desktop: boolean
  showPending: boolean
  showConfirmed: boolean
  showFailed: boolean
  autoHide: boolean
  autoHideDelay: number
}

export function TransactionNotifications() {
  const [isOpen, setIsOpen] = useState(false)
  const [settings, setSettings] = useState<NotificationSettings>({
    enabled: true,
    sound: true,
    desktop: false,
    showPending: true,
    showConfirmed: true,
    showFailed: true,
    autoHide: true,
    autoHideDelay: 5000
  })
  const [dismissedNotifications, setDismissedNotifications] = useState<Set<string>>(new Set())

  const { pendingTransactions, confirmedTransactions, failedTransactions } = useTransactionMonitor()

  // Request desktop notification permission
  useEffect(() => {
    if (settings.desktop && 'Notification' in window && Notification.permission === 'default') {
      Notification.requestPermission()
    }
  }, [settings.desktop])

  // Show desktop notifications for new transactions
  useEffect(() => {
    if (!settings.enabled || !settings.desktop || !('Notification' in window)) return

    const allTransactions = [...pendingTransactions, ...confirmedTransactions, ...failedTransactions]
    
    allTransactions.forEach(tx => {
      const notificationId = `${tx.hash}-${tx.status}`
      
      if (dismissedNotifications.has(notificationId)) return

      const shouldShow = 
        (tx.status === TransactionStatus.PENDING && settings.showPending) ||
        (tx.status === TransactionStatus.CONFIRMED && settings.showConfirmed) ||
        (tx.status === TransactionStatus.FAILED && settings.showFailed)

      if (shouldShow && Notification.permission === 'granted') {
        const chain = SUPPORTED_CHAINS[tx.chainId]
        const title = getNotificationTitle(tx.status)
        const body = `Transaction ${formatHash(tx.hash)} on ${chain?.shortName || 'Unknown Chain'}`
        
        const notification = new Notification(title, {
          body,
          icon: '/favicon.ico',
          tag: notificationId
        })

        notification.onclick = () => {
          window.focus()
          openBlockExplorer(tx)
          notification.close()
        }

        if (settings.autoHide) {
          setTimeout(() => notification.close(), settings.autoHideDelay)
        }
      }
    })
  }, [pendingTransactions, confirmedTransactions, failedTransactions, settings, dismissedNotifications])

  const getNotificationTitle = (status: TransactionStatus): string => {
    switch (status) {
      case TransactionStatus.PENDING:
        return 'Transaction Pending'
      case TransactionStatus.CONFIRMED:
        return 'Transaction Confirmed'
      case TransactionStatus.FAILED:
        return 'Transaction Failed'
      case TransactionStatus.DROPPED:
        return 'Transaction Dropped'
      default:
        return 'Transaction Update'
    }
  }

  const getStatusIcon = (status: TransactionStatus) => {
    switch (status) {
      case TransactionStatus.PENDING:
        return <Clock className="w-4 h-4 text-yellow-500" />
      case TransactionStatus.CONFIRMED:
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case TransactionStatus.FAILED:
        return <XCircle className="w-4 h-4 text-red-500" />
      case TransactionStatus.DROPPED:
        return <AlertTriangle className="w-4 h-4 text-orange-500" />
      default:
        return <Bell className="w-4 h-4" />
    }
  }

  const formatHash = (hash: string) => {
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`
  }

  const openBlockExplorer = (tx: TransactionData) => {
    const chain = SUPPORTED_CHAINS[tx.chainId]
    if (chain?.blockExplorers?.default?.url) {
      window.open(`${chain.blockExplorers.default.url}/tx/${tx.hash}`, '_blank')
    }
  }

  const dismissNotification = (tx: TransactionData) => {
    const notificationId = `${tx.hash}-${tx.status}`
    setDismissedNotifications(prev => new Set(Array.from(prev).concat(notificationId)))
  }

  const getVisibleTransactions = () => {
    const allTransactions = [...pendingTransactions, ...confirmedTransactions, ...failedTransactions]
    
    return allTransactions.filter(tx => {
      const notificationId = `${tx.hash}-${tx.status}`
      if (dismissedNotifications.has(notificationId)) return false

      return (
        (tx.status === TransactionStatus.PENDING && settings.showPending) ||
        (tx.status === TransactionStatus.CONFIRMED && settings.showConfirmed) ||
        (tx.status === TransactionStatus.FAILED && settings.showFailed)
      )
    }).slice(0, 5) // Show max 5 notifications
  }

  const visibleTransactions = getVisibleTransactions()
  const hasNotifications = visibleTransactions.length > 0

  return (
    <div className="fixed top-4 right-4 z-50">
      {/* Notification Bell */}
      <div className="relative">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setIsOpen(!isOpen)}
          className={`relative ${hasNotifications ? 'bg-blue-50 border-blue-200' : ''}`}
        >
          <Bell className="w-4 h-4" />
          {hasNotifications && (
            <div className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full flex items-center justify-center">
              <span className="text-xs text-white font-bold">{visibleTransactions.length}</span>
            </div>
          )}
        </Button>

        {/* Notification Panel */}
        <AnimatePresence>
          {isOpen && (
            <motion.div
              initial={{ opacity: 0, y: -10, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              className="absolute top-12 right-0 w-80 bg-background border rounded-lg shadow-lg"
            >
              {/* Header */}
              <div className="flex items-center justify-between p-4 border-b">
                <h3 className="font-semibold">Transaction Notifications</h3>
                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSettings(prev => ({ ...prev, sound: !prev.sound }))}
                  >
                    {settings.sound ? <Volume2 className="w-4 h-4" /> : <VolumeX className="w-4 h-4" />}
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setIsOpen(false)}
                  >
                    <X className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              {/* Notifications */}
              <div className="max-h-96 overflow-y-auto">
                {visibleTransactions.length === 0 ? (
                  <div className="p-8 text-center">
                    <Bell className="w-8 h-8 text-muted-foreground mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">No new notifications</p>
                  </div>
                ) : (
                  <div className="space-y-2 p-2">
                    {visibleTransactions.map((tx, index) => {
                      const chain = SUPPORTED_CHAINS[tx.chainId]
                      const progress = tx.status === TransactionStatus.PENDING 
                        ? (tx.confirmations / tx.maxConfirmations) * 100 
                        : 100

                      return (
                        <motion.div
                          key={`${tx.hash}-${tx.status}`}
                          initial={{ opacity: 0, x: 20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.1 }}
                        >
                          <Card className="hover:bg-accent/50 transition-colors">
                            <CardContent className="p-3">
                              <div className="flex items-start justify-between">
                                <div className="flex items-start gap-3 flex-1">
                                  {getStatusIcon(tx.status)}
                                  <div className="flex-1 min-w-0">
                                    <div className="flex items-center gap-2 mb-1">
                                      <span className="text-sm font-medium">
                                        {getNotificationTitle(tx.status)}
                                      </span>
                                      <Badge variant="outline" className="text-xs">
                                        {chain?.shortName}
                                      </Badge>
                                    </div>
                                    <p className="text-xs text-muted-foreground mb-2">
                                      {formatHash(tx.hash)}
                                    </p>
                                    
                                    {tx.status === TransactionStatus.PENDING && (
                                      <div className="space-y-1">
                                        <div className="flex justify-between text-xs">
                                          <span>Confirmations</span>
                                          <span>{tx.confirmations}/{tx.maxConfirmations}</span>
                                        </div>
                                        <Progress value={progress} className="h-1" />
                                      </div>
                                    )}

                                    {tx.metadata?.description && (
                                      <p className="text-xs text-muted-foreground mt-1">
                                        {tx.metadata.description}
                                      </p>
                                    )}
                                  </div>
                                </div>

                                <div className="flex items-center gap-1 ml-2">
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => openBlockExplorer(tx)}
                                    className="h-6 w-6 p-0"
                                  >
                                    <ExternalLink className="w-3 h-3" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => dismissNotification(tx)}
                                    className="h-6 w-6 p-0"
                                  >
                                    <X className="w-3 h-3" />
                                  </Button>
                                </div>
                              </div>
                            </CardContent>
                          </Card>
                        </motion.div>
                      )
                    })}
                  </div>
                )}
              </div>

              {/* Settings */}
              <div className="border-t p-3">
                <div className="flex items-center justify-between text-sm">
                  <span>Desktop notifications</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSettings(prev => ({ ...prev, desktop: !prev.desktop }))}
                    className="h-6 px-2"
                  >
                    {settings.desktop ? 'On' : 'Off'}
                  </Button>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Floating Notifications */}
      <div className="fixed top-16 right-4 space-y-2 pointer-events-none">
        <AnimatePresence>
          {!isOpen && visibleTransactions.slice(0, 3).map((tx, index) => {
            const chain = SUPPORTED_CHAINS[tx.chainId]
            
            return (
              <motion.div
                key={`${tx.hash}-${tx.status}-float`}
                initial={{ opacity: 0, x: 100 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: 100 }}
                transition={{ delay: index * 0.2 }}
                className="pointer-events-auto"
              >
                <Card className="w-72 bg-background/95 backdrop-blur border shadow-lg">
                  <CardContent className="p-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(tx.status)}
                        <div>
                          <p className="text-sm font-medium">
                            {getNotificationTitle(tx.status)}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {formatHash(tx.hash)} on {chain?.shortName}
                          </p>
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => dismissNotification(tx)}
                        className="h-6 w-6 p-0"
                      >
                        <X className="w-3 h-3" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            )
          })}
        </AnimatePresence>
      </div>
    </div>
  )
}
