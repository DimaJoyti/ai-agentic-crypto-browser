'use client'

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Clock,
  TrendingUp,
  Activity,
  CheckCircle,
  AlertTriangle
} from 'lucide-react'

interface FeatureTooltipProps {
  title: string
  description: string
  features: string[]
  stats?: Array<{
    label: string
    value: string
    icon?: React.ElementType
  }>
  status?: 'online' | 'offline' | 'maintenance' | 'beta'
  children: React.ReactNode
  className?: string
}

export function FeatureTooltip({
  title,
  description,
  features,
  stats,
  status = 'online',
  children,
  className = ''
}: FeatureTooltipProps) {
  const [isVisible, setIsVisible] = useState(false)

  const getStatusConfig = () => {
    switch (status) {
      case 'online':
        return {
          color: 'text-green-600',
          bgColor: 'bg-green-100 dark:bg-green-900/50',
          icon: CheckCircle,
          text: 'Online'
        }
      case 'offline':
        return {
          color: 'text-red-600',
          bgColor: 'bg-red-100 dark:bg-red-900/50',
          icon: AlertTriangle,
          text: 'Offline'
        }
      case 'maintenance':
        return {
          color: 'text-orange-600',
          bgColor: 'bg-orange-100 dark:bg-orange-900/50',
          icon: Clock,
          text: 'Maintenance'
        }
      default:
        return {
          color: 'text-gray-600',
          bgColor: 'bg-gray-100 dark:bg-gray-900/50',
          icon: Activity,
          text: 'Unknown'
        }
    }
  }

  const statusConfig = getStatusConfig()
  const StatusIcon = statusConfig.icon

  return (
    <div 
      className={`relative ${className}`}
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      {children}
      
      <AnimatePresence>
        {isVisible && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 10 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 10 }}
            transition={{ duration: 0.2 }}
            className="absolute top-full left-1/2 transform -translate-x-1/2 mt-2 z-50 w-80"
          >
            <Card className="glass-card border-0 shadow-2xl bg-white/95 dark:bg-black/95 backdrop-blur-md">
              <CardContent className="p-4">
                {/* Header */}
                <div className="flex items-center justify-between mb-3">
                  <h4 className="font-semibold text-lg">{title}</h4>
                  <div className={`flex items-center gap-1 px-2 py-1 rounded-full ${statusConfig.bgColor}`}>
                    <StatusIcon className={`h-3 w-3 ${statusConfig.color}`} />
                    <span className={`text-xs font-medium ${statusConfig.color}`}>
                      {statusConfig.text}
                    </span>
                  </div>
                </div>

                {/* Description */}
                <p className="text-sm text-muted-foreground mb-3">
                  {description}
                </p>

                {/* Stats */}
                {stats && stats.length > 0 && (
                  <div className="grid grid-cols-2 gap-3 mb-3">
                    {stats.map((stat, index) => {
                      const StatIcon = stat.icon || Activity
                      return (
                        <div key={index} className="flex items-center gap-2 p-2 bg-muted/50 rounded-lg">
                          <StatIcon className="h-4 w-4 text-primary" />
                          <div>
                            <div className="text-sm font-semibold">{stat.value}</div>
                            <div className="text-xs text-muted-foreground">{stat.label}</div>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                )}

                {/* Features */}
                <div className="space-y-2">
                  <h5 className="text-sm font-medium">Key Features:</h5>
                  <div className="flex flex-wrap gap-1">
                    {features.map((feature, index) => (
                      <Badge key={index} variant="secondary" className="text-xs">
                        {feature}
                      </Badge>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

interface QuickPreviewProps {
  title: string
  metrics: Array<{
    label: string
    value: string
    change?: string
    trend?: 'up' | 'down' | 'neutral'
  }>
  children: React.ReactNode
  className?: string
}

export function QuickPreview({
  title,
  metrics,
  children,
  className = ''
}: QuickPreviewProps) {
  const [isVisible, setIsVisible] = useState(false)

  const getTrendIcon = (trend?: 'up' | 'down' | 'neutral') => {
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
    <div 
      className={`relative ${className}`}
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      {children}
      
      <AnimatePresence>
        {isVisible && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95, x: -10 }}
            animate={{ opacity: 1, scale: 1, x: 0 }}
            exit={{ opacity: 0, scale: 0.95, x: -10 }}
            transition={{ duration: 0.2 }}
            className="absolute top-0 left-full ml-2 z-50 w-64"
          >
            <Card className="glass-card border-0 shadow-xl bg-white/95 dark:bg-black/95 backdrop-blur-md">
              <CardContent className="p-3">
                <h4 className="font-semibold text-sm mb-3">{title} - Live Metrics</h4>
                <div className="space-y-2">
                  {metrics.map((metric, index) => (
                    <div key={index} className="flex items-center justify-between p-2 bg-muted/30 rounded">
                      <div>
                        <div className="text-xs text-muted-foreground">{metric.label}</div>
                        <div className="text-sm font-semibold">{metric.value}</div>
                      </div>
                      <div className="flex items-center gap-1">
                        {metric.change && (
                          <span className={`text-xs ${
                            metric.trend === 'up' ? 'text-green-600' : 
                            metric.trend === 'down' ? 'text-red-600' : 
                            'text-gray-600'
                          }`}>
                            {metric.change}
                          </span>
                        )}
                        {getTrendIcon(metric.trend)}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}

interface FeatureCardProps {
  title: string
  description: string
  icon: React.ElementType
  href: string
  status: 'online' | 'offline' | 'maintenance' | 'beta'
  badge?: 'hot' | 'new' | 'popular' | 'premium'
  stats?: Array<{
    label: string
    value: string
    color?: string
  }>
  features: string[]
  gradient: string
  children?: React.ReactNode
}

export function EnhancedFeatureCard({
  title,
  description,
  icon: Icon,
  href,
  status,
  badge,
  stats,
  features,
  gradient,
  children
}: FeatureCardProps) {
  return (
    <FeatureTooltip
      title={title}
      description={description}
      features={features}
      stats={stats?.map(stat => ({
        label: stat.label,
        value: stat.value,
        icon: Activity
      }))}
      status={status}
    >
      <QuickPreview
        title={title}
        metrics={stats?.map(stat => ({
          label: stat.label,
          value: stat.value,
          trend: 'up' as const
        })) || []}
      >
        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
          className="h-full"
        >
          <a href={href} className="block h-full">
            <div className={`relative overflow-hidden rounded-2xl ${gradient} p-6 shadow-xl hover:shadow-2xl transition-all duration-300 group h-full`}>
              <div className="absolute top-0 right-0 w-20 h-20 bg-white/10 rounded-full -translate-y-10 translate-x-10" />
              <div className="relative h-full flex flex-col">
                <div className="flex items-center gap-3 mb-4">
                  <div className="p-3 bg-white/20 rounded-xl group-hover:bg-white/30 transition-colors">
                    <Icon className="h-6 w-6" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="text-xl font-bold">{title}</h3>
                      {badge && (
                        <Badge variant="secondary" className="text-xs">
                          {badge}
                        </Badge>
                      )}
                    </div>
                    <div className="flex items-center gap-2">
                      <div className={`w-2 h-2 rounded-full ${
                        status === 'online' ? 'bg-green-500 animate-pulse' :
                        status === 'offline' ? 'bg-red-500' :
                        status === 'maintenance' ? 'bg-orange-500' :
                        'bg-blue-500 animate-pulse'
                      }`} />
                      <span className="text-sm opacity-90">
                        {status === 'online' ? 'Online' :
                         status === 'offline' ? 'Offline' :
                         status === 'maintenance' ? 'Maintenance' :
                         'Beta'}
                      </span>
                    </div>
                  </div>
                </div>
                
                <p className="opacity-90 mb-4 flex-1">
                  {description}
                </p>
                
                {children}
              </div>
            </div>
          </a>
        </motion.div>
      </QuickPreview>
    </FeatureTooltip>
  )
}
