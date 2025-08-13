'use client'

import { useEffect, useState, useCallback } from 'react'

interface DevOptimizations {
  disableAnimations: boolean
  reduceMotion: boolean
  lazyLoading: boolean
  virtualScrolling: boolean
  memoization: boolean
  debounceInputs: boolean
}

interface PerformanceMetrics {
  renderTime: number
  memoryUsage: number
  componentCount: number
  reRenders: number
  fps: number
}

export function useDevOptimizations() {
  const [optimizations, setOptimizations] = useState<DevOptimizations>({
    disableAnimations: false,
    reduceMotion: false,
    lazyLoading: true,
    virtualScrolling: false,
    memoization: true,
    debounceInputs: true,
  })

  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    renderTime: 0,
    memoryUsage: 0,
    componentCount: 0,
    reRenders: 0,
    fps: 60,
  })

  const isDev = process.env.NODE_ENV === 'development'

  // Load optimizations from localStorage
  useEffect(() => {
    if (!isDev) return

    const saved = localStorage.getItem('dev-optimizations')
    if (saved) {
      try {
        setOptimizations(JSON.parse(saved))
      } catch (e) {
        console.warn('Failed to load dev optimizations:', e)
      }
    }
  }, [isDev])

  // Save optimizations to localStorage
  const updateOptimizations = useCallback((newOptimizations: Partial<DevOptimizations>) => {
    if (!isDev) return

    const updated = { ...optimizations, ...newOptimizations }
    setOptimizations(updated)
    localStorage.setItem('dev-optimizations', JSON.stringify(updated))
  }, [optimizations, isDev])

  // Performance monitoring
  useEffect(() => {
    if (!isDev) return

    let frameCount = 0
    let lastTime = performance.now()
    let animationId: number

    const measureFPS = () => {
      frameCount++
      const currentTime = performance.now()
      
      if (currentTime >= lastTime + 1000) {
        const fps = Math.round((frameCount * 1000) / (currentTime - lastTime))
        setMetrics(prev => ({ ...prev, fps }))
        frameCount = 0
        lastTime = currentTime
      }
      
      animationId = requestAnimationFrame(measureFPS)
    }

    measureFPS()

    // Memory monitoring
    const memoryInterval = setInterval(() => {
      if ('memory' in performance) {
        const memory = (performance as any).memory
        setMetrics(prev => ({
          ...prev,
          memoryUsage: Math.round(memory.usedJSHeapSize / 1024 / 1024)
        }))
      }
    }, 1000)

    // Performance observer for render times
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()
      entries.forEach((entry) => {
        if (entry.entryType === 'measure' && entry.name.includes('React')) {
          setMetrics(prev => ({
            ...prev,
            renderTime: Math.round(entry.duration)
          }))
        }
      })
    })

    try {
      observer.observe({ entryTypes: ['measure'] })
    } catch (e) {
      // Performance observer not supported
    }

    return () => {
      cancelAnimationFrame(animationId)
      clearInterval(memoryInterval)
      observer.disconnect()
    }
  }, [isDev])

  // Apply CSS optimizations
  useEffect(() => {
    if (!isDev) return

    const root = document.documentElement

    if (optimizations.disableAnimations) {
      root.style.setProperty('--animation-duration', '0s')
      root.style.setProperty('--transition-duration', '0s')
      root.classList.add('performance-mode')
    } else {
      root.style.removeProperty('--animation-duration')
      root.style.removeProperty('--transition-duration')
      root.classList.remove('performance-mode')
    }

    if (optimizations.reduceMotion) {
      root.classList.add('reduce-motion')
    } else {
      root.classList.remove('reduce-motion')
    }
  }, [optimizations, isDev])

  // Debounced input helper
  const createDebouncedCallback = useCallback((callback: Function, delay: number = 300) => {
    if (!optimizations.debounceInputs) return callback

    let timeoutId: NodeJS.Timeout
    return (...args: any[]) => {
      clearTimeout(timeoutId)
      timeoutId = setTimeout(() => callback(...args), delay)
    }
  }, [optimizations.debounceInputs])

  // Memoization helper
  const memoize = useCallback(<T extends (...args: any[]) => any>(fn: T): T => {
    if (!optimizations.memoization) return fn

    const cache = new Map()
    return ((...args: any[]) => {
      const key = JSON.stringify(args)
      if (cache.has(key)) {
        return cache.get(key)
      }
      const result = fn(...args)
      cache.set(key, result)
      return result
    }) as T
  }, [optimizations.memoization])

  // Lazy loading helper
  const useLazyLoading = useCallback((threshold: number = 0.1) => {
    const [isVisible, setIsVisible] = useState(!optimizations.lazyLoading)

    useEffect(() => {
      if (!optimizations.lazyLoading) {
        setIsVisible(true)
        return
      }

      const observer = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting) {
            setIsVisible(true)
            observer.disconnect()
          }
        },
        { threshold }
      )

      return () => observer.disconnect()
    }, [threshold])

    return { isVisible, setRef: useCallback((node: Element | null) => {
      if (node && optimizations.lazyLoading) {
        const observer = new IntersectionObserver(
          ([entry]) => {
            if (entry.isIntersecting) {
              setIsVisible(true)
              observer.disconnect()
            }
          },
          { threshold }
        )
        observer.observe(node)
      }
    }, [threshold]) }
  }, [optimizations.lazyLoading])

  // Performance-aware component wrapper class generator
  const getPerformanceWrapperClass = useCallback((className = '') => {
    return [
      className,
      optimizations.disableAnimations && 'performance-mode',
      optimizations.reduceMotion && 'reduce-motion',
      'performance-optimized'
    ].filter(Boolean).join(' ')
  }, [optimizations])

  // Reset all optimizations
  const resetOptimizations = useCallback(() => {
    if (!isDev) return

    const defaultOptimizations: DevOptimizations = {
      disableAnimations: false,
      reduceMotion: false,
      lazyLoading: true,
      virtualScrolling: false,
      memoization: true,
      debounceInputs: true,
    }

    setOptimizations(defaultOptimizations)
    localStorage.removeItem('dev-optimizations')
  }, [isDev])

  // Get performance status
  const getPerformanceStatus = useCallback(() => {
    const { fps, memoryUsage, renderTime } = metrics
    
    if (fps < 30 || memoryUsage > 100 || renderTime > 16) {
      return 'poor'
    } else if (fps < 50 || memoryUsage > 50 || renderTime > 8) {
      return 'fair'
    } else {
      return 'good'
    }
  }, [metrics])

  return {
    // State
    optimizations,
    metrics,
    isDev,
    
    // Actions
    updateOptimizations,
    resetOptimizations,
    
    // Helpers
    createDebouncedCallback,
    memoize,
    useLazyLoading,
    getPerformanceWrapperClass,
    getPerformanceStatus,
  }
}
