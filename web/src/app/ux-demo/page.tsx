'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { EnhancedLayout } from '@/components/layout/EnhancedLayout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { LoadingSpinner, Skeleton, ProgressIndicator } from '@/components/ui/enhanced-loading'
import { useNotifications } from '@/components/ui/notification-system'
import { useAccessibility } from '@/components/accessibility/AccessibilityProvider'
import { UserPreferences } from '@/components/settings/user-preferences'
import { 
  Palette, 
  Zap, 
  Bell, 
  Eye, 
  Smartphone,
  Monitor,
  Accessibility,
  Settings,
  TrendingUp,
  BarChart3,
  DollarSign,
  Activity,
  Sparkles,
  Rocket,
  Star,
  Heart,
  ThumbsUp
} from 'lucide-react'

export default function UXDemoPage() {
  const [loadingStates, setLoadingStates] = useState({
    spinner: false,
    skeleton: false,
    progress: 0
  })
  const { addNotification } = useNotifications()
  const { announceMessage, setIsAccessibilityPanelOpen } = useAccessibility()

  const showDemoNotifications = () => {
    const notifications = [
      {
        type: 'success' as const,
        title: 'Order Executed',
        message: 'Your BTC buy order has been successfully executed',
        metadata: { symbol: 'BTC/USD', price: 45000, change: 2.5 }
      },
      {
        type: 'trading' as const,
        title: 'Price Alert',
        message: 'ETH has reached your target price of $3,000',
        metadata: { symbol: 'ETH/USD', price: 3000, change: 5.2 }
      },
      {
        type: 'warning' as const,
        title: 'Risk Warning',
        message: 'Your portfolio exposure to volatile assets is high'
      },
      {
        type: 'info' as const,
        title: 'Market Update',
        message: 'Bitcoin dominance has increased to 45.2%'
      },
      {
        type: 'security' as const,
        title: 'Security Alert',
        message: 'New login detected from Chrome on Windows'
      }
    ]

    notifications.forEach((notification, index) => {
      setTimeout(() => {
        addNotification(notification)
      }, index * 1000)
    })
  }

  const simulateLoading = (type: 'spinner' | 'skeleton' | 'progress') => {
    if (type === 'progress') {
      setLoadingStates(prev => ({ ...prev, progress: 0 }))
      const interval = setInterval(() => {
        setLoadingStates(prev => {
          const newProgress = prev.progress + 10
          if (newProgress >= 100) {
            clearInterval(interval)
            return { ...prev, progress: 100 }
          }
          return { ...prev, progress: newProgress }
        })
      }, 200)
    } else {
      setLoadingStates(prev => ({ ...prev, [type]: true }))
      setTimeout(() => {
        setLoadingStates(prev => ({ ...prev, [type]: false }))
      }, 3000)
    }
  }

  return (
    <EnhancedLayout>
      <div className="container mx-auto p-6 space-y-8">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center space-y-4"
        >
          <div className="flex items-center justify-center space-x-2">
            <Sparkles className="h-8 w-8 text-primary" />
            <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-transparent">
              UX Enhancement Demo
            </h1>
            <Sparkles className="h-8 w-8 text-primary" />
          </div>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Experience the enhanced user interface with advanced accessibility, animations, and interactive features
          </p>
          <div className="flex flex-wrap justify-center gap-2">
            <Badge variant="secondary" className="animate-pulse">
              <Rocket className="h-3 w-3 mr-1" />
              Enhanced UX
            </Badge>
            <Badge variant="outline">
              <Accessibility className="h-3 w-3 mr-1" />
              Accessible
            </Badge>
            <Badge variant="outline">
              <Smartphone className="h-3 w-3 mr-1" />
              Responsive
            </Badge>
            <Badge variant="outline">
              <Zap className="h-3 w-3 mr-1" />
              Fast
            </Badge>
          </div>
        </motion.div>

        <Tabs defaultValue="components" className="w-full">
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="components">Components</TabsTrigger>
            <TabsTrigger value="notifications">Notifications</TabsTrigger>
            <TabsTrigger value="loading">Loading States</TabsTrigger>
            <TabsTrigger value="accessibility">Accessibility</TabsTrigger>
            <TabsTrigger value="preferences">Preferences</TabsTrigger>
          </TabsList>

          {/* Enhanced Components */}
          <TabsContent value="components" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {/* Animated Cards */}
              <motion.div
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Card className="glass-card">
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <TrendingUp className="h-5 w-5 text-green-500" />
                      <span>Portfolio Value</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-3xl font-bold text-green-500">$125,430.50</div>
                    <div className="text-sm text-muted-foreground">+5.2% today</div>
                  </CardContent>
                </Card>
              </motion.div>

              <motion.div
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Card className="border-primary/20 bg-primary/5">
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <BarChart3 className="h-5 w-5 text-primary" />
                      <span>Active Trades</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-3xl font-bold">12</div>
                    <div className="text-sm text-muted-foreground">3 pending orders</div>
                  </CardContent>
                </Card>
              </motion.div>

              <motion.div
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Card className="border-yellow-500/20 bg-yellow-500/5">
                  <CardHeader>
                    <CardTitle className="flex items-center space-x-2">
                      <DollarSign className="h-5 w-5 text-yellow-500" />
                      <span>P&L Today</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-3xl font-bold text-green-500">+$2,340.80</div>
                    <div className="text-sm text-muted-foreground">+1.9% return</div>
                  </CardContent>
                </Card>
              </motion.div>
            </div>

            {/* Interactive Buttons */}
            <Card>
              <CardHeader>
                <CardTitle>Interactive Elements</CardTitle>
                <CardDescription>Hover and click effects with smooth animations</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex flex-wrap gap-4">
                  <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                    <Button className="animate-pulse-glow">
                      <Star className="h-4 w-4 mr-2" />
                      Primary Action
                    </Button>
                  </motion.div>
                  
                  <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                    <Button variant="outline">
                      <Heart className="h-4 w-4 mr-2" />
                      Secondary
                    </Button>
                  </motion.div>
                  
                  <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                    <Button variant="ghost">
                      <ThumbsUp className="h-4 w-4 mr-2" />
                      Ghost Button
                    </Button>
                  </motion.div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Notifications Demo */}
          <TabsContent value="notifications" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Bell className="h-5 w-5" />
                  <span>Enhanced Notification System</span>
                </CardTitle>
                <CardDescription>
                  Advanced notifications with different types, actions, and metadata
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <Button onClick={showDemoNotifications} className="w-full">
                  <Bell className="h-4 w-4 mr-2" />
                  Show Demo Notifications
                </Button>
                <div className="text-sm text-muted-foreground">
                  This will show 5 different types of notifications with various features:
                  <ul className="list-disc list-inside mt-2 space-y-1">
                    <li>Trading notifications with price data</li>
                    <li>Success/error states with actions</li>
                    <li>Security alerts with timestamps</li>
                    <li>Auto-dismiss and persistent options</li>
                    <li>Rich metadata and interactive elements</li>
                  </ul>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Loading States Demo */}
          <TabsContent value="loading" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Loading Spinners</CardTitle>
                  <CardDescription>Various animated loading indicators</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span>Default Spinner:</span>
                      <LoadingSpinner />
                    </div>
                    <div className="flex items-center justify-between">
                      <span>Trading Spinner:</span>
                      <LoadingSpinner variant="trading" />
                    </div>
                    <div className="flex items-center justify-between">
                      <span>Crypto Spinner:</span>
                      <LoadingSpinner variant="crypto" />
                    </div>
                    <div className="flex items-center justify-between">
                      <span>Dots Animation:</span>
                      <LoadingSpinner variant="dots" />
                    </div>
                    <div className="flex items-center justify-between">
                      <span>Pulse Effect:</span>
                      <LoadingSpinner variant="pulse" />
                    </div>
                  </div>
                  <Button 
                    onClick={() => simulateLoading('spinner')}
                    disabled={loadingStates.spinner}
                    className="w-full"
                  >
                    {loadingStates.spinner ? (
                      <LoadingSpinner size="sm" className="mr-2" />
                    ) : (
                      <Activity className="h-4 w-4 mr-2" />
                    )}
                    {loadingStates.spinner ? 'Loading...' : 'Simulate Loading'}
                  </Button>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Skeleton Loading</CardTitle>
                  <CardDescription>Content placeholders while loading</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {loadingStates.skeleton ? (
                    <div className="space-y-4">
                      <Skeleton variant="card" />
                      <Skeleton variant="chart" />
                      <Skeleton variant="table" />
                    </div>
                  ) : (
                    <div className="space-y-4">
                      <div className="p-4 border rounded-lg">
                        <h4 className="font-medium">Sample Content</h4>
                        <p className="text-sm text-muted-foreground">This content will be replaced with skeleton loading</p>
                      </div>
                      <div className="h-24 bg-muted/30 rounded-lg flex items-center justify-center">
                        <span className="text-muted-foreground">Chart Area</span>
                      </div>
                    </div>
                  )}
                  <Button 
                    onClick={() => simulateLoading('skeleton')}
                    disabled={loadingStates.skeleton}
                    className="w-full"
                  >
                    Show Skeleton Loading
                  </Button>
                </CardContent>
              </Card>
            </div>

            <Card>
              <CardHeader>
                <CardTitle>Progress Indicators</CardTitle>
                <CardDescription>Different progress visualization styles</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium mb-2">Linear Progress</h4>
                    <ProgressIndicator progress={loadingStates.progress} />
                  </div>
                  <div>
                    <h4 className="font-medium mb-2">Circular Progress</h4>
                    <ProgressIndicator progress={loadingStates.progress} variant="circular" />
                  </div>
                  <div>
                    <h4 className="font-medium mb-2">Stepped Progress</h4>
                    <ProgressIndicator progress={loadingStates.progress} variant="stepped" />
                  </div>
                </div>
                <Button 
                  onClick={() => simulateLoading('progress')}
                  className="w-full"
                >
                  Simulate Progress
                </Button>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Accessibility Demo */}
          <TabsContent value="accessibility" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Accessibility className="h-5 w-5" />
                  <span>Accessibility Features</span>
                </CardTitle>
                <CardDescription>
                  Comprehensive accessibility support for all users
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <Button 
                    onClick={() => setIsAccessibilityPanelOpen(true)}
                    className="w-full"
                  >
                    <Settings className="h-4 w-4 mr-2" />
                    Open Accessibility Panel
                  </Button>
                  <Button 
                    onClick={() => announceMessage('This is a test announcement for screen readers')}
                    variant="outline"
                    className="w-full"
                  >
                    <Eye className="h-4 w-4 mr-2" />
                    Test Screen Reader
                  </Button>
                </div>
                
                <div className="text-sm text-muted-foreground space-y-2">
                  <h4 className="font-medium text-foreground">Keyboard Shortcuts:</h4>
                  <ul className="space-y-1">
                    <li><kbd className="px-2 py-1 bg-muted rounded text-xs">Alt + A</kbd> - Open accessibility panel</li>
                    <li><kbd className="px-2 py-1 bg-muted rounded text-xs">Alt + H</kbd> - Toggle high contrast</li>
                    <li><kbd className="px-2 py-1 bg-muted rounded text-xs">Alt + T</kbd> - Toggle large text</li>
                    <li><kbd className="px-2 py-1 bg-muted rounded text-xs">Alt + M</kbd> - Toggle reduced motion</li>
                    <li><kbd className="px-2 py-1 bg-muted rounded text-xs">⌘/Ctrl + K</kbd> - Open search</li>
                    <li><kbd className="px-2 py-1 bg-muted rounded text-xs">⌘/Ctrl + ,</kbd> - Open theme customizer</li>
                  </ul>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* User Preferences */}
          <TabsContent value="preferences">
            <UserPreferences />
          </TabsContent>
        </Tabs>
      </div>
    </EnhancedLayout>
  )
}
