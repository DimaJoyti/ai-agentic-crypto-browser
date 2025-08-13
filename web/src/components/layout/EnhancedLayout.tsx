'use client'

import { useState, useEffect, ReactNode } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { usePathname } from 'next/navigation'
import { EnhancedNavigation } from './EnhancedNavigation'
import { ThemeCustomizer } from '@/components/ui/theme-customizer'
import { NotificationProvider, NotificationCenter, useNotifications } from '@/components/ui/notification-system'
import { LoadingSpinner } from '@/components/ui/enhanced-loading'
import { PerformanceOptimizer } from '@/components/dev/PerformanceOptimizer'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { 
  Settings, 
  Bell, 
  Search, 
  Menu,
  X,
  Palette,
  Zap,
  Activity,
  TrendingUp,
  BarChart3,
  Globe,
  User,
  HelpCircle,
  Maximize2,
  Minimize2,
  RefreshCw
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface EnhancedLayoutProps {
  children: ReactNode
  showNavigation?: boolean
  showHeader?: boolean
  className?: string
}

function LayoutContent({ children, showNavigation = true, showHeader = true, className }: EnhancedLayoutProps) {
  const [isThemeCustomizerOpen, setIsThemeCustomizerOpen] = useState(false)
  const [isNotificationCenterOpen, setIsNotificationCenterOpen] = useState(false)
  const [isSearchOpen, setIsSearchOpen] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const pathname = usePathname()
  const { notifications, addNotification } = useNotifications()

  // Page transition loading
  useEffect(() => {
    setIsLoading(true)
    const timer = setTimeout(() => setIsLoading(false), 500)
    return () => clearTimeout(timer)
  }, [pathname])

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.metaKey || e.ctrlKey) {
        switch (e.key) {
          case 'k':
            e.preventDefault()
            setIsSearchOpen(true)
            break
          case ',':
            e.preventDefault()
            setIsThemeCustomizerOpen(true)
            break
          case 'n':
            e.preventDefault()
            setIsNotificationCenterOpen(true)
            break
          case 'f':
            e.preventDefault()
            toggleFullscreen()
            break
        }
      }
      if (e.key === 'Escape') {
        setIsSearchOpen(false)
        setIsThemeCustomizerOpen(false)
        setIsNotificationCenterOpen(false)
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [])

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen()
      setIsFullscreen(true)
    } else {
      document.exitFullscreen()
      setIsFullscreen(false)
    }
  }

  const refreshPage = () => {
    setIsLoading(true)
    window.location.reload()
  }

  // Demo notification
  const showDemoNotification = () => {
    addNotification({
      type: 'trading',
      title: 'Order Executed',
      message: 'Your BTC buy order has been successfully executed',
      metadata: {
        symbol: 'BTC/USD',
        price: 45000,
        change: 2.5
      },
      action: {
        label: 'View Details',
        onClick: () => console.log('View order details')
      }
    })
  }

  return (
    <div className={cn("min-h-screen bg-background", className)}>
      {/* Page Loading Overlay */}
      <AnimatePresence>
        {isLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center"
          >
            <div className="text-center space-y-4">
              <LoadingSpinner size="lg" variant="trading" />
              <p className="text-sm text-muted-foreground">Loading...</p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="flex h-screen overflow-hidden">
        {/* Navigation Sidebar */}
        {showNavigation && (
          <motion.div
            initial={false}
            animate={{ width: 'auto' }}
            className="flex-shrink-0"
          >
            <EnhancedNavigation />
          </motion.div>
        )}

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Enhanced Header */}
          {showHeader && (
            <motion.header
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              className="bg-background/95 backdrop-blur-xl border-b border-border px-6 py-4 flex items-center justify-between"
            >
              {/* Left Section */}
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                  <span className="text-sm text-muted-foreground">Live</span>
                </div>
                
                {/* Breadcrumb */}
                <nav className="hidden md:flex items-center space-x-2 text-sm">
                  <span className="text-muted-foreground">AI Browser</span>
                  <span className="text-muted-foreground">/</span>
                  <span className="font-medium capitalize">
                    {pathname?.split('/').filter(Boolean).pop() || 'Home'}
                  </span>
                </nav>
              </div>

              {/* Center Section - Search */}
              <div className="flex-1 max-w-md mx-4">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <input
                    type="text"
                    placeholder="Search... (⌘K)"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onFocus={() => setIsSearchOpen(true)}
                    className="w-full pl-10 pr-4 py-2 bg-muted/50 border border-border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:bg-background transition-all"
                  />
                </div>
              </div>

              {/* Right Section */}
              <div className="flex items-center space-x-2">
                {/* Quick Actions */}
                <div className="hidden lg:flex items-center space-x-1">
                  <Button variant="ghost" size="icon" onClick={refreshPage} title="Refresh (⌘R)">
                    <RefreshCw className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="icon" onClick={toggleFullscreen} title="Fullscreen (⌘F)">
                    {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
                  </Button>
                </div>

                {/* Notifications */}
                <Button
                  variant="ghost"
                  size="icon"
                  className="relative"
                  onClick={() => setIsNotificationCenterOpen(true)}
                  title="Notifications (⌘N)"
                >
                  <Bell className="h-4 w-4" />
                  {notifications.length > 0 && (
                    <Badge 
                      variant="destructive" 
                      className="absolute -top-1 -right-1 h-5 w-5 p-0 flex items-center justify-center text-xs"
                    >
                      {notifications.length}
                    </Badge>
                  )}
                </Button>

                {/* Theme Customizer */}
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setIsThemeCustomizerOpen(true)}
                  title="Customize Theme (⌘,)"
                >
                  <Palette className="h-4 w-4" />
                </Button>

                {/* Demo Notification Button */}
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={showDemoNotification}
                  className="hidden lg:flex"
                >
                  <Zap className="h-4 w-4 mr-2" />
                  Demo Alert
                </Button>

                {/* User Menu */}
                <Button variant="ghost" size="icon" title="User Settings">
                  <User className="h-4 w-4" />
                </Button>
              </div>
            </motion.header>
          )}

          {/* Main Content */}
          <main className="flex-1 overflow-auto">
            <motion.div
              key={pathname}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
              className="h-full"
            >
              {children}
            </motion.div>
          </main>

          {/* Status Bar */}
          <motion.footer
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            className="bg-muted/30 border-t border-border px-6 py-2 flex items-center justify-between text-xs text-muted-foreground"
          >
            <div className="flex items-center space-x-4">
              <span>AI Agentic Browser v2.0.0</span>
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-green-500 rounded-full" />
                <span>Connected</span>
              </div>
              <span>Last updated: {new Date().toLocaleTimeString()}</span>
            </div>
            <div className="flex items-center space-x-4">
              <span>Memory: 45.2MB</span>
              <span>CPU: 12%</span>
              <Button variant="ghost" size="sm" className="h-6 px-2">
                <HelpCircle className="h-3 w-3 mr-1" />
                Help
              </Button>
            </div>
          </motion.footer>
        </div>
      </div>

      {/* Global Search Modal */}
      <AnimatePresence>
        {isSearchOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 bg-black/50 z-50"
              onClick={() => setIsSearchOpen(false)}
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: -20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: -20 }}
              className="fixed top-20 left-1/2 transform -translate-x-1/2 w-full max-w-2xl bg-background border border-border rounded-lg shadow-xl z-50 p-6"
            >
              <div className="space-y-4">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-muted-foreground" />
                  <input
                    type="text"
                    placeholder="Search everything..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full pl-12 pr-4 py-3 bg-muted/50 border border-border rounded-lg text-lg focus:outline-none focus:ring-2 focus:ring-primary/20"
                    autoFocus
                  />
                </div>
                <div className="text-sm text-muted-foreground">
                  Search across trading data, settings, documentation, and more...
                </div>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Theme Customizer */}
      <ThemeCustomizer
        isOpen={isThemeCustomizerOpen}
        onClose={() => setIsThemeCustomizerOpen(false)}
      />

      {/* Notification Center */}
      <NotificationCenter
        isOpen={isNotificationCenterOpen}
        onClose={() => setIsNotificationCenterOpen(false)}
      />

      {/* Performance Optimizer (Development Only) */}
      <PerformanceOptimizer />
    </div>
  )
}

export function EnhancedLayout(props: EnhancedLayoutProps) {
  return (
    <NotificationProvider>
      <LayoutContent {...props} />
    </NotificationProvider>
  )
}
