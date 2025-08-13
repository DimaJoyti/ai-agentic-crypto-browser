'use client'

import React, { useEffect, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import { 
  Zap, 
  Settings, 
  Monitor, 
  Cpu, 
  MemoryStick,
  Timer,
  X,
  RefreshCw,
  TrendingUp,
  AlertTriangle
} from 'lucide-react'

interface PerformanceMetrics {
  renderTime: number
  memoryUsage: number
  componentCount: number
  reRenders: number
}

interface PerformanceSettings {
  disableAnimations: boolean
  reduceMotion: boolean
  lazyLoading: boolean
  virtualScrolling: boolean
  memoization: boolean
  debounceInputs: boolean
}

export function PerformanceOptimizer() {
  const [isOpen, setIsOpen] = useState(false)
  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    renderTime: 0,
    memoryUsage: 0,
    componentCount: 0,
    reRenders: 0
  })
  
  const [settings, setSettings] = useState<PerformanceSettings>({
    disableAnimations: false,
    reduceMotion: false,
    lazyLoading: true,
    virtualScrolling: false,
    memoization: true,
    debounceInputs: true
  })

  // Only show in development
  const isDev = process.env.NODE_ENV === 'development'

  useEffect(() => {
    if (!isDev) return

    // Monitor performance metrics
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()
      entries.forEach((entry) => {
        if (entry.entryType === 'measure') {
          setMetrics(prev => ({
            ...prev,
            renderTime: entry.duration
          }))
        }
      })
    })

    observer.observe({ entryTypes: ['measure'] })

    // Memory usage monitoring
    const memoryInterval = setInterval(() => {
      if ('memory' in performance) {
        const memory = (performance as any).memory
        setMetrics(prev => ({
          ...prev,
          memoryUsage: memory.usedJSHeapSize / 1024 / 1024 // MB
        }))
      }
    }, 1000)

    return () => {
      observer.disconnect()
      clearInterval(memoryInterval)
    }
  }, [isDev])

  // Apply performance settings
  useEffect(() => {
    if (settings.disableAnimations) {
      document.documentElement.style.setProperty('--animation-duration', '0s')
      document.documentElement.style.setProperty('--transition-duration', '0s')
    } else {
      document.documentElement.style.removeProperty('--animation-duration')
      document.documentElement.style.removeProperty('--transition-duration')
    }

    if (settings.reduceMotion) {
      document.documentElement.style.setProperty('--motion-reduce', '1')
    } else {
      document.documentElement.style.removeProperty('--motion-reduce')
    }
  }, [settings])

  const handleSettingChange = (key: keyof PerformanceSettings, value: boolean) => {
    setSettings(prev => ({ ...prev, [key]: value }))
    
    // Store in localStorage for persistence
    localStorage.setItem('dev-performance-settings', JSON.stringify({
      ...settings,
      [key]: value
    }))
  }

  const resetSettings = () => {
    const defaultSettings: PerformanceSettings = {
      disableAnimations: false,
      reduceMotion: false,
      lazyLoading: true,
      virtualScrolling: false,
      memoization: true,
      debounceInputs: true
    }
    setSettings(defaultSettings)
    localStorage.removeItem('dev-performance-settings')
  }

  // Load settings from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('dev-performance-settings')
    if (saved) {
      try {
        setSettings(JSON.parse(saved))
      } catch (e) {
        console.warn('Failed to load performance settings:', e)
      }
    }
  }, [])

  if (!isDev) return null

  return (
    <>
      {/* Floating Performance Button */}
      <motion.div
        className="fixed bottom-4 right-4 z-50"
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ delay: 1 }}
      >
        <Button
          onClick={() => setIsOpen(true)}
          size="icon"
          className="rounded-full shadow-lg bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700"
          title="Performance Optimizer (Dev Only)"
        >
          <Zap className="h-4 w-4" />
        </Button>
      </motion.div>

      {/* Performance Optimizer Panel */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
            onClick={() => setIsOpen(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="w-full max-w-2xl"
            >
              <Card className="glass-card border-0 shadow-2xl">
                <CardHeader className="flex flex-row items-center justify-between">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      <Zap className="h-5 w-5 text-orange-500" />
                      Performance Optimizer
                      <Badge variant="secondary">Development</Badge>
                    </CardTitle>
                    <CardDescription>
                      Optimize UI performance for faster development
                    </CardDescription>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setIsOpen(false)}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </CardHeader>

                <CardContent className="space-y-6">
                  {/* Performance Metrics */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="text-center p-3 bg-blue-50 dark:bg-blue-950/30 rounded-lg">
                      <Timer className="h-5 w-5 mx-auto mb-1 text-blue-600" />
                      <div className="text-sm font-medium">Render Time</div>
                      <div className="text-lg font-bold text-blue-600">
                        {metrics.renderTime.toFixed(1)}ms
                      </div>
                    </div>
                    
                    <div className="text-center p-3 bg-green-50 dark:bg-green-950/30 rounded-lg">
                      <MemoryStick className="h-5 w-5 mx-auto mb-1 text-green-600" />
                      <div className="text-sm font-medium">Memory</div>
                      <div className="text-lg font-bold text-green-600">
                        {metrics.memoryUsage.toFixed(1)}MB
                      </div>
                    </div>
                    
                    <div className="text-center p-3 bg-purple-50 dark:bg-purple-950/30 rounded-lg">
                      <Monitor className="h-5 w-5 mx-auto mb-1 text-purple-600" />
                      <div className="text-sm font-medium">Components</div>
                      <div className="text-lg font-bold text-purple-600">
                        {metrics.componentCount}
                      </div>
                    </div>
                    
                    <div className="text-center p-3 bg-orange-50 dark:bg-orange-950/30 rounded-lg">
                      <RefreshCw className="h-5 w-5 mx-auto mb-1 text-orange-600" />
                      <div className="text-sm font-medium">Re-renders</div>
                      <div className="text-lg font-bold text-orange-600">
                        {metrics.reRenders}
                      </div>
                    </div>
                  </div>

                  {/* Performance Settings */}
                  <div className="space-y-4">
                    <h3 className="text-lg font-semibold flex items-center gap-2">
                      <Settings className="h-5 w-5" />
                      Performance Settings
                    </h3>
                    
                    <div className="grid gap-4">
                      {Object.entries(settings).map(([key, value]) => (
                        <div key={key} className="flex items-center justify-between p-3 bg-muted/50 rounded-lg">
                          <div>
                            <div className="font-medium capitalize">
                              {key.replace(/([A-Z])/g, ' $1').trim()}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {getSettingDescription(key as keyof PerformanceSettings)}
                            </div>
                          </div>
                          <Switch
                            checked={value}
                            onCheckedChange={(checked) => 
                              handleSettingChange(key as keyof PerformanceSettings, checked)
                            }
                          />
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2">
                    <Button onClick={resetSettings} variant="outline" className="flex-1">
                      Reset to Defaults
                    </Button>
                    <Button onClick={() => window.location.reload()} className="flex-1">
                      <RefreshCw className="h-4 w-4 mr-2" />
                      Apply & Reload
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  )
}

function getSettingDescription(key: keyof PerformanceSettings): string {
  const descriptions = {
    disableAnimations: 'Disable all CSS animations for faster rendering',
    reduceMotion: 'Reduce motion effects for better performance',
    lazyLoading: 'Load components only when needed',
    virtualScrolling: 'Use virtual scrolling for large lists',
    memoization: 'Enable React.memo for component optimization',
    debounceInputs: 'Debounce input changes to reduce re-renders'
  }
  return descriptions[key]
}
