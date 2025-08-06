'use client'

import { motion, AnimatePresence } from 'framer-motion'
import { Loader2, Zap, TrendingUp, BarChart3, Activity, DollarSign } from 'lucide-react'
import { cn } from '@/lib/utils'

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  variant?: 'default' | 'dots' | 'pulse' | 'bounce' | 'trading' | 'crypto'
  className?: string
  text?: string
}

export function LoadingSpinner({ 
  size = 'md', 
  variant = 'default', 
  className,
  text 
}: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12'
  }

  if (variant === 'dots') {
    return (
      <div className={cn("flex items-center space-x-1", className)}>
        {[0, 1, 2].map((i) => (
          <motion.div
            key={i}
            className={cn(
              "bg-primary rounded-full",
              size === 'sm' ? 'w-1 h-1' : size === 'md' ? 'w-2 h-2' : size === 'lg' ? 'w-3 h-3' : 'w-4 h-4'
            )}
            animate={{
              scale: [1, 1.5, 1],
              opacity: [0.5, 1, 0.5]
            }}
            transition={{
              duration: 1,
              repeat: Infinity,
              delay: i * 0.2
            }}
          />
        ))}
        {text && <span className="ml-2 text-sm text-muted-foreground">{text}</span>}
      </div>
    )
  }

  if (variant === 'pulse') {
    return (
      <div className={cn("flex items-center space-x-2", className)}>
        <motion.div
          className={cn(
            "bg-primary rounded-full",
            sizeClasses[size]
          )}
          animate={{
            scale: [1, 1.2, 1],
            opacity: [0.7, 1, 0.7]
          }}
          transition={{
            duration: 1.5,
            repeat: Infinity,
            ease: "easeInOut"
          }}
        />
        {text && <span className="text-sm text-muted-foreground">{text}</span>}
      </div>
    )
  }

  if (variant === 'bounce') {
    return (
      <div className={cn("flex items-center space-x-1", className)}>
        {[0, 1, 2].map((i) => (
          <motion.div
            key={i}
            className={cn(
              "bg-primary rounded-full",
              size === 'sm' ? 'w-2 h-2' : size === 'md' ? 'w-3 h-3' : size === 'lg' ? 'w-4 h-4' : 'w-5 h-5'
            )}
            animate={{
              y: [0, -10, 0]
            }}
            transition={{
              duration: 0.6,
              repeat: Infinity,
              delay: i * 0.1
            }}
          />
        ))}
        {text && <span className="ml-2 text-sm text-muted-foreground">{text}</span>}
      </div>
    )
  }

  if (variant === 'trading') {
    return (
      <div className={cn("flex items-center space-x-2", className)}>
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
        >
          <TrendingUp className={cn("text-green-500", sizeClasses[size])} />
        </motion.div>
        {text && <span className="text-sm text-muted-foreground">{text}</span>}
      </div>
    )
  }

  if (variant === 'crypto') {
    return (
      <div className={cn("flex items-center space-x-2", className)}>
        <motion.div
          animate={{ 
            rotate: [0, 180, 360],
            scale: [1, 1.1, 1]
          }}
          transition={{ 
            duration: 2, 
            repeat: Infinity, 
            ease: "easeInOut" 
          }}
        >
          <DollarSign className={cn("text-yellow-500", sizeClasses[size])} />
        </motion.div>
        {text && <span className="text-sm text-muted-foreground">{text}</span>}
      </div>
    )
  }

  return (
    <div className={cn("flex items-center space-x-2", className)}>
      <motion.div
        animate={{ rotate: 360 }}
        transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
      >
        <Loader2 className={cn("text-primary", sizeClasses[size])} />
      </motion.div>
      {text && <span className="text-sm text-muted-foreground">{text}</span>}
    </div>
  )
}

interface SkeletonProps {
  className?: string
  variant?: 'default' | 'card' | 'chart' | 'table' | 'trading'
}

export function Skeleton({ className, variant = 'default' }: SkeletonProps) {
  if (variant === 'card') {
    return (
      <div className={cn("space-y-3", className)}>
        <motion.div
          className="h-4 bg-muted rounded animate-pulse"
          initial={{ opacity: 0.5 }}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ duration: 1.5, repeat: Infinity }}
        />
        <motion.div
          className="h-3 bg-muted rounded animate-pulse w-3/4"
          initial={{ opacity: 0.5 }}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ duration: 1.5, repeat: Infinity, delay: 0.2 }}
        />
        <motion.div
          className="h-3 bg-muted rounded animate-pulse w-1/2"
          initial={{ opacity: 0.5 }}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ duration: 1.5, repeat: Infinity, delay: 0.4 }}
        />
      </div>
    )
  }

  if (variant === 'chart') {
    return (
      <div className={cn("space-y-4", className)}>
        <div className="flex space-x-2">
          {[...Array(7)].map((_, i) => (
            <motion.div
              key={i}
              className="bg-muted rounded w-8"
              style={{ height: Math.random() * 100 + 50 }}
              initial={{ opacity: 0.5, scaleY: 0 }}
              animate={{ 
                opacity: [0.5, 1, 0.5],
                scaleY: [0, 1, 0.8, 1]
              }}
              transition={{ 
                duration: 1.5, 
                repeat: Infinity, 
                delay: i * 0.1 
              }}
            />
          ))}
        </div>
      </div>
    )
  }

  if (variant === 'table') {
    return (
      <div className={cn("space-y-2", className)}>
        {[...Array(5)].map((_, i) => (
          <div key={i} className="flex space-x-4">
            <motion.div
              className="h-4 bg-muted rounded w-1/4"
              initial={{ opacity: 0.5 }}
              animate={{ opacity: [0.5, 1, 0.5] }}
              transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.1 }}
            />
            <motion.div
              className="h-4 bg-muted rounded w-1/3"
              initial={{ opacity: 0.5 }}
              animate={{ opacity: [0.5, 1, 0.5] }}
              transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.1 + 0.2 }}
            />
            <motion.div
              className="h-4 bg-muted rounded w-1/4"
              initial={{ opacity: 0.5 }}
              animate={{ opacity: [0.5, 1, 0.5] }}
              transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.1 + 0.4 }}
            />
          </div>
        ))}
      </div>
    )
  }

  if (variant === 'trading') {
    return (
      <div className={cn("space-y-4", className)}>
        {/* Price ticker skeleton */}
        <div className="flex justify-between items-center">
          <motion.div
            className="h-6 bg-muted rounded w-24"
            animate={{ opacity: [0.5, 1, 0.5] }}
            transition={{ duration: 1.5, repeat: Infinity }}
          />
          <motion.div
            className="h-6 bg-green-200 rounded w-16"
            animate={{ opacity: [0.5, 1, 0.5] }}
            transition={{ duration: 1.5, repeat: Infinity, delay: 0.2 }}
          />
        </div>
        
        {/* Chart skeleton */}
        <div className="h-32 bg-muted rounded relative overflow-hidden">
          <motion.div
            className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent"
            animate={{ x: [-100, 300] }}
            transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
          />
        </div>
        
        {/* Order book skeleton */}
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            {[...Array(3)].map((_, i) => (
              <motion.div
                key={i}
                className="h-3 bg-red-100 rounded"
                animate={{ opacity: [0.5, 1, 0.5] }}
                transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.1 }}
              />
            ))}
          </div>
          <div className="space-y-2">
            {[...Array(3)].map((_, i) => (
              <motion.div
                key={i}
                className="h-3 bg-green-100 rounded"
                animate={{ opacity: [0.5, 1, 0.5] }}
                transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.1 + 0.3 }}
              />
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <motion.div
      className={cn("h-4 bg-muted rounded", className)}
      initial={{ opacity: 0.5 }}
      animate={{ opacity: [0.5, 1, 0.5] }}
      transition={{ duration: 1.5, repeat: Infinity }}
    />
  )
}

interface ProgressIndicatorProps {
  progress: number
  variant?: 'default' | 'circular' | 'stepped'
  size?: 'sm' | 'md' | 'lg'
  className?: string
  showPercentage?: boolean
}

export function ProgressIndicator({ 
  progress, 
  variant = 'default', 
  size = 'md',
  className,
  showPercentage = true
}: ProgressIndicatorProps) {
  if (variant === 'circular') {
    const radius = size === 'sm' ? 20 : size === 'md' ? 30 : 40
    const circumference = 2 * Math.PI * radius
    const strokeDashoffset = circumference - (progress / 100) * circumference

    return (
      <div className={cn("relative inline-flex items-center justify-center", className)}>
        <svg
          className={cn(
            "transform -rotate-90",
            size === 'sm' ? 'w-12 h-12' : size === 'md' ? 'w-16 h-16' : 'w-20 h-20'
          )}
          viewBox={`0 0 ${(radius + 10) * 2} ${(radius + 10) * 2}`}
        >
          <circle
            cx={radius + 10}
            cy={radius + 10}
            r={radius}
            stroke="currentColor"
            strokeWidth="2"
            fill="transparent"
            className="text-muted"
          />
          <motion.circle
            cx={radius + 10}
            cy={radius + 10}
            r={radius}
            stroke="currentColor"
            strokeWidth="2"
            fill="transparent"
            strokeDasharray={circumference}
            strokeDashoffset={strokeDashoffset}
            strokeLinecap="round"
            className="text-primary"
            initial={{ strokeDashoffset: circumference }}
            animate={{ strokeDashoffset }}
            transition={{ duration: 0.5, ease: "easeInOut" }}
          />
        </svg>
        {showPercentage && (
          <div className="absolute inset-0 flex items-center justify-center">
            <span className={cn(
              "font-medium text-primary",
              size === 'sm' ? 'text-xs' : size === 'md' ? 'text-sm' : 'text-base'
            )}>
              {Math.round(progress)}%
            </span>
          </div>
        )}
      </div>
    )
  }

  if (variant === 'stepped') {
    const steps = 5
    const activeSteps = Math.ceil((progress / 100) * steps)

    return (
      <div className={cn("flex space-x-2", className)}>
        {[...Array(steps)].map((_, i) => (
          <motion.div
            key={i}
            className={cn(
              "rounded-full",
              size === 'sm' ? 'w-2 h-2' : size === 'md' ? 'w-3 h-3' : 'w-4 h-4',
              i < activeSteps ? 'bg-primary' : 'bg-muted'
            )}
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: i * 0.1 }}
          />
        ))}
        {showPercentage && (
          <span className="ml-2 text-sm text-muted-foreground">
            {Math.round(progress)}%
          </span>
        )}
      </div>
    )
  }

  return (
    <div className={cn("w-full", className)}>
      <div className={cn(
        "bg-muted rounded-full overflow-hidden",
        size === 'sm' ? 'h-1' : size === 'md' ? 'h-2' : 'h-3'
      )}>
        <motion.div
          className="h-full bg-primary rounded-full"
          initial={{ width: 0 }}
          animate={{ width: `${progress}%` }}
          transition={{ duration: 0.5, ease: "easeInOut" }}
        />
      </div>
      {showPercentage && (
        <div className="flex justify-between text-xs text-muted-foreground mt-1">
          <span>0%</span>
          <span>{Math.round(progress)}%</span>
          <span>100%</span>
        </div>
      )}
    </div>
  )
}
