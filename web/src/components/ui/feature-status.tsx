'use client'

import React from 'react'
import { motion } from 'framer-motion'
import { Badge } from '@/components/ui/badge'
import { 
  CheckCircle, 
  AlertTriangle, 
  Clock, 
  Zap, 
  TrendingUp, 
  Shield,
  Activity,
  Wifi,
  WifiOff
} from 'lucide-react'

interface FeatureStatusProps {
  status: 'online' | 'offline' | 'warning' | 'maintenance' | 'beta' | 'new'
  label?: string
  showIcon?: boolean
  animated?: boolean
  className?: string
}

export function FeatureStatus({ 
  status, 
  label, 
  showIcon = true, 
  animated = true,
  className = '' 
}: FeatureStatusProps) {
  const getStatusConfig = () => {
    switch (status) {
      case 'online':
        return {
          color: 'bg-green-500',
          textColor: 'text-green-600 dark:text-green-400',
          bgColor: 'bg-green-100 dark:bg-green-900/50',
          icon: CheckCircle,
          text: label || 'Online',
          pulse: true
        }
      case 'offline':
        return {
          color: 'bg-red-500',
          textColor: 'text-red-600 dark:text-red-400',
          bgColor: 'bg-red-100 dark:bg-red-900/50',
          icon: WifiOff,
          text: label || 'Offline',
          pulse: false
        }
      case 'warning':
        return {
          color: 'bg-yellow-500',
          textColor: 'text-yellow-600 dark:text-yellow-400',
          bgColor: 'bg-yellow-100 dark:bg-yellow-900/50',
          icon: AlertTriangle,
          text: label || 'Warning',
          pulse: true
        }
      case 'maintenance':
        return {
          color: 'bg-orange-500',
          textColor: 'text-orange-600 dark:text-orange-400',
          bgColor: 'bg-orange-100 dark:bg-orange-900/50',
          icon: Clock,
          text: label || 'Maintenance',
          pulse: false
        }
      case 'beta':
        return {
          color: 'bg-blue-500',
          textColor: 'text-blue-600 dark:text-blue-400',
          bgColor: 'bg-blue-100 dark:bg-blue-900/50',
          icon: Zap,
          text: label || 'Beta',
          pulse: true
        }
      case 'new':
        return {
          color: 'bg-purple-500',
          textColor: 'text-purple-600 dark:text-purple-400',
          bgColor: 'bg-purple-100 dark:bg-purple-900/50',
          icon: TrendingUp,
          text: label || 'New',
          pulse: true
        }
      default:
        return {
          color: 'bg-gray-500',
          textColor: 'text-gray-600 dark:text-gray-400',
          bgColor: 'bg-gray-100 dark:bg-gray-900/50',
          icon: Activity,
          text: label || 'Unknown',
          pulse: false
        }
    }
  }

  const config = getStatusConfig()
  const Icon = config.icon

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <div className="relative">
        <div 
          className={`w-2 h-2 rounded-full ${config.color} ${
            config.pulse && animated ? 'animate-pulse' : ''
          }`} 
        />
        {config.pulse && animated && (
          <motion.div
            className={`absolute inset-0 w-2 h-2 rounded-full ${config.color} opacity-75`}
            animate={{
              scale: [1, 1.5, 1],
              opacity: [0.75, 0, 0.75],
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: "easeInOut",
            }}
          />
        )}
      </div>
      <span className={`text-sm font-medium ${config.textColor}`}>
        {config.text}
      </span>
      {showIcon && (
        <Icon className={`h-3 w-3 ${config.textColor}`} />
      )}
    </div>
  )
}

interface FeatureBadgeProps {
  type: 'hot' | 'trending' | 'popular' | 'premium' | 'free' | 'pro'
  className?: string
}

export function FeatureBadge({ type, className = '' }: FeatureBadgeProps) {
  const getBadgeConfig = () => {
    switch (type) {
      case 'hot':
        return {
          text: '🔥 Hot',
          variant: 'destructive' as const,
          className: 'bg-gradient-to-r from-red-500 to-orange-500 text-white animate-pulse'
        }
      case 'trending':
        return {
          text: '📈 Trending',
          variant: 'default' as const,
          className: 'bg-gradient-to-r from-blue-500 to-purple-500 text-white'
        }
      case 'popular':
        return {
          text: '⭐ Popular',
          variant: 'secondary' as const,
          className: 'bg-gradient-to-r from-yellow-400 to-orange-400 text-black'
        }
      case 'premium':
        return {
          text: '👑 Premium',
          variant: 'outline' as const,
          className: 'bg-gradient-to-r from-purple-600 to-pink-600 text-white border-0'
        }
      case 'free':
        return {
          text: '🆓 Free',
          variant: 'secondary' as const,
          className: 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
        }
      case 'pro':
        return {
          text: '⚡ Pro',
          variant: 'default' as const,
          className: 'bg-gradient-to-r from-indigo-500 to-blue-500 text-white'
        }
      default:
        return {
          text: type,
          variant: 'secondary' as const,
          className: ''
        }
    }
  }

  const config = getBadgeConfig()

  return (
    <Badge 
      variant={config.variant}
      className={`${config.className} ${className} text-xs font-semibold shadow-sm`}
    >
      {config.text}
    </Badge>
  )
}

interface MetricDisplayProps {
  label: string
  value: string | number
  trend?: 'up' | 'down' | 'neutral'
  color?: 'green' | 'red' | 'blue' | 'purple' | 'orange'
  className?: string
}

export function MetricDisplay({ 
  label, 
  value, 
  trend, 
  color = 'blue',
  className = '' 
}: MetricDisplayProps) {
  const getColorClasses = () => {
    switch (color) {
      case 'green':
        return 'text-green-600 dark:text-green-400'
      case 'red':
        return 'text-red-600 dark:text-red-400'
      case 'blue':
        return 'text-blue-600 dark:text-blue-400'
      case 'purple':
        return 'text-purple-600 dark:text-purple-400'
      case 'orange':
        return 'text-orange-600 dark:text-orange-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const getTrendIcon = () => {
    switch (trend) {
      case 'up':
        return <TrendingUp className="h-3 w-3 text-green-500" />
      case 'down':
        return <TrendingUp className="h-3 w-3 text-red-500 rotate-180" />
      default:
        return null
    }
  }

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <div className="text-center">
        <div className={`text-lg font-bold ${getColorClasses()}`}>
          {value}
        </div>
        <div className="text-xs text-muted-foreground">
          {label}
        </div>
      </div>
      {trend && getTrendIcon()}
    </div>
  )
}

interface QuickStatsProps {
  stats: Array<{
    label: string
    value: string | number
    color?: 'green' | 'red' | 'blue' | 'purple' | 'orange'
    trend?: 'up' | 'down' | 'neutral'
  }>
  className?: string
}

export function QuickStats({ stats, className = '' }: QuickStatsProps) {
  return (
    <div className={`flex items-center gap-4 ${className}`}>
      {stats.map((stat, index) => (
        <MetricDisplay
          key={index}
          label={stat.label}
          value={stat.value}
          color={stat.color}
          trend={stat.trend}
        />
      ))}
    </div>
  )
}
