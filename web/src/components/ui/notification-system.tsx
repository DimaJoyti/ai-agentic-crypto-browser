'use client'

import { useState, useEffect, createContext, useContext, ReactNode } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  CheckCircle, 
  AlertCircle, 
  XCircle, 
  Info, 
  X, 
  Bell,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Zap,
  Shield,
  Activity
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

export type NotificationType = 'success' | 'error' | 'warning' | 'info' | 'trading' | 'security'

export interface Notification {
  id: string
  type: NotificationType
  title: string
  message: string
  duration?: number
  persistent?: boolean
  action?: {
    label: string
    onClick: () => void
  }
  metadata?: {
    symbol?: string
    price?: number
    change?: number
    volume?: number
  }
  timestamp: Date
}

interface NotificationContextType {
  notifications: Notification[]
  addNotification: (notification: Omit<Notification, 'id' | 'timestamp'>) => void
  removeNotification: (id: string) => void
  clearAll: () => void
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined)

export function useNotifications() {
  const context = useContext(NotificationContext)
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider')
  }
  return context
}

interface NotificationProviderProps {
  children: ReactNode
  maxNotifications?: number
}

export function NotificationProvider({ children, maxNotifications = 5 }: NotificationProviderProps) {
  const [notifications, setNotifications] = useState<Notification[]>([])

  const addNotification = (notification: Omit<Notification, 'id' | 'timestamp'>) => {
    const id = Math.random().toString(36).substr(2, 9)
    const newNotification: Notification = {
      ...notification,
      id,
      timestamp: new Date(),
      duration: notification.duration ?? 5000
    }

    setNotifications(prev => {
      const updated = [newNotification, ...prev]
      return updated.slice(0, maxNotifications)
    })

    // Auto-remove non-persistent notifications
    if (!notification.persistent && notification.duration !== 0) {
      setTimeout(() => {
        removeNotification(id)
      }, newNotification.duration)
    }
  }

  const removeNotification = (id: string) => {
    setNotifications(prev => prev.filter(n => n.id !== id))
  }

  const clearAll = () => {
    setNotifications([])
  }

  return (
    <NotificationContext.Provider value={{
      notifications,
      addNotification,
      removeNotification,
      clearAll
    }}>
      {children}
      <NotificationContainer />
    </NotificationContext.Provider>
  )
}

function NotificationContainer() {
  const { notifications, removeNotification } = useNotifications()

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2 max-w-sm w-full">
      <AnimatePresence>
        {notifications.map((notification) => (
          <NotificationItem
            key={notification.id}
            notification={notification}
            onRemove={() => removeNotification(notification.id)}
          />
        ))}
      </AnimatePresence>
    </div>
  )
}

interface NotificationItemProps {
  notification: Notification
  onRemove: () => void
}

function NotificationItem({ notification, onRemove }: NotificationItemProps) {
  const getIcon = () => {
    switch (notification.type) {
      case 'success':
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case 'error':
        return <XCircle className="h-5 w-5 text-red-500" />
      case 'warning':
        return <AlertCircle className="h-5 w-5 text-yellow-500" />
      case 'info':
        return <Info className="h-5 w-5 text-blue-500" />
      case 'trading':
        return <TrendingUp className="h-5 w-5 text-green-500" />
      case 'security':
        return <Shield className="h-5 w-5 text-purple-500" />
      default:
        return <Bell className="h-5 w-5 text-gray-500" />
    }
  }

  const getBorderColor = () => {
    switch (notification.type) {
      case 'success':
        return 'border-l-green-500'
      case 'error':
        return 'border-l-red-500'
      case 'warning':
        return 'border-l-yellow-500'
      case 'info':
        return 'border-l-blue-500'
      case 'trading':
        return 'border-l-green-500'
      case 'security':
        return 'border-l-purple-500'
      default:
        return 'border-l-gray-500'
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, x: 300, scale: 0.3 }}
      animate={{ opacity: 1, x: 0, scale: 1 }}
      exit={{ opacity: 0, x: 300, scale: 0.5, transition: { duration: 0.2 } }}
      className={cn(
        "bg-background border border-border rounded-lg shadow-lg p-4 border-l-4",
        getBorderColor()
      )}
    >
      <div className="flex items-start space-x-3">
        <div className="flex-shrink-0 mt-0.5">
          {getIcon()}
        </div>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <h4 className="font-medium text-sm">{notification.title}</h4>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-muted-foreground hover:text-foreground"
              onClick={onRemove}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
          
          <p className="text-sm text-muted-foreground mt-1">
            {notification.message}
          </p>

          {/* Trading-specific metadata */}
          {notification.type === 'trading' && notification.metadata && (
            <div className="flex items-center space-x-4 mt-2">
              {notification.metadata.symbol && (
                <Badge variant="outline" className="text-xs">
                  {notification.metadata.symbol}
                </Badge>
              )}
              {notification.metadata.price && (
                <span className="text-xs font-mono">
                  ${notification.metadata.price.toLocaleString()}
                </span>
              )}
              {notification.metadata.change && (
                <span className={cn(
                  "text-xs font-mono",
                  notification.metadata.change >= 0 ? "text-green-600" : "text-red-600"
                )}>
                  {notification.metadata.change >= 0 ? '+' : ''}
                  {notification.metadata.change.toFixed(2)}%
                </span>
              )}
            </div>
          )}

          {/* Action button */}
          {notification.action && (
            <Button
              variant="outline"
              size="sm"
              className="mt-2"
              onClick={notification.action.onClick}
            >
              {notification.action.label}
            </Button>
          )}

          {/* Timestamp */}
          <p className="text-xs text-muted-foreground mt-2">
            {notification.timestamp.toLocaleTimeString()}
          </p>
        </div>
      </div>
    </motion.div>
  )
}

// Notification Center Component
interface NotificationCenterProps {
  isOpen: boolean
  onClose: () => void
}

export function NotificationCenter({ isOpen, onClose }: NotificationCenterProps) {
  const { notifications, removeNotification, clearAll } = useNotifications()
  const [filter, setFilter] = useState<NotificationType | 'all'>('all')

  const filteredNotifications = notifications.filter(n => 
    filter === 'all' || n.type === filter
  )

  const notificationTypes: { type: NotificationType | 'all', label: string, icon: any }[] = [
    { type: 'all', label: 'All', icon: Bell },
    { type: 'trading', label: 'Trading', icon: TrendingUp },
    { type: 'security', label: 'Security', icon: Shield },
    { type: 'success', label: 'Success', icon: CheckCircle },
    { type: 'error', label: 'Errors', icon: XCircle },
    { type: 'warning', label: 'Warnings', icon: AlertCircle },
    { type: 'info', label: 'Info', icon: Info }
  ]

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Overlay */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-40"
            onClick={onClose}
          />

          {/* Notification Center Panel */}
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 20, stiffness: 300 }}
            className="fixed right-0 top-0 h-full w-96 bg-background border-l border-border z-50 overflow-hidden flex flex-col"
          >
            {/* Header */}
            <div className="p-6 border-b border-border">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <Bell className="h-5 w-5" />
                  <h2 className="font-semibold text-lg">Notifications</h2>
                  <Badge variant="secondary">{notifications.length}</Badge>
                </div>
                <div className="flex items-center space-x-2">
                  <Button variant="ghost" size="sm" onClick={clearAll}>
                    Clear All
                  </Button>
                  <Button variant="ghost" size="icon" onClick={onClose}>
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              {/* Filter Tabs */}
              <div className="flex space-x-1 mt-4 overflow-x-auto">
                {notificationTypes.map(({ type, label, icon: Icon }) => (
                  <Button
                    key={type}
                    variant={filter === type ? 'default' : 'ghost'}
                    size="sm"
                    onClick={() => setFilter(type)}
                    className="flex-shrink-0"
                  >
                    <Icon className="h-4 w-4 mr-1" />
                    {label}
                  </Button>
                ))}
              </div>
            </div>

            {/* Notifications List */}
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {filteredNotifications.length === 0 ? (
                <div className="text-center py-8">
                  <Bell className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                  <p className="text-muted-foreground">No notifications</p>
                </div>
              ) : (
                filteredNotifications.map((notification) => (
                  <motion.div
                    key={notification.id}
                    layout
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className={cn(
                      "p-3 rounded-lg border border-border bg-card hover:bg-muted/50 transition-colors cursor-pointer",
                      getBorderColor(notification.type)
                    )}
                    onClick={() => removeNotification(notification.id)}
                  >
                    <div className="flex items-start space-x-3">
                      <div className="flex-shrink-0 mt-0.5">
                        {getNotificationIcon(notification.type)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium text-sm">{notification.title}</h4>
                        <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                          {notification.message}
                        </p>
                        <p className="text-xs text-muted-foreground mt-2">
                          {notification.timestamp.toLocaleString()}
                        </p>
                      </div>
                    </div>
                  </motion.div>
                ))
              )}
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}

function getBorderColor(type: NotificationType) {
  switch (type) {
    case 'success': return 'border-l-4 border-l-green-500'
    case 'error': return 'border-l-4 border-l-red-500'
    case 'warning': return 'border-l-4 border-l-yellow-500'
    case 'info': return 'border-l-4 border-l-blue-500'
    case 'trading': return 'border-l-4 border-l-green-500'
    case 'security': return 'border-l-4 border-l-purple-500'
    default: return 'border-l-4 border-l-gray-500'
  }
}

function getNotificationIcon(type: NotificationType) {
  switch (type) {
    case 'success': return <CheckCircle className="h-4 w-4 text-green-500" />
    case 'error': return <XCircle className="h-4 w-4 text-red-500" />
    case 'warning': return <AlertCircle className="h-4 w-4 text-yellow-500" />
    case 'info': return <Info className="h-4 w-4 text-blue-500" />
    case 'trading': return <TrendingUp className="h-4 w-4 text-green-500" />
    case 'security': return <Shield className="h-4 w-4 text-purple-500" />
    default: return <Bell className="h-4 w-4 text-gray-500" />
  }
}
