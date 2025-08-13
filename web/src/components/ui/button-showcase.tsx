'use client'

import React, { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { EnhancedButton } from '@/components/ui/enhanced-button'
import { Badge } from '@/components/ui/badge'
import { 
  Sparkles, 
  Zap, 
  Heart, 
  Star, 
  Download, 
  Play, 
  Pause, 
  Settings,
  ShoppingCart,
  Crown,
  Rocket,
  Gift,
  Shield,
  AlertTriangle,
  CheckCircle,
  Info
} from 'lucide-react'

export function ButtonShowcase() {
  const [loading, setLoading] = useState<string | null>(null)

  const handleLoadingDemo = (buttonId: string) => {
    setLoading(buttonId)
    setTimeout(() => setLoading(null), 3000)
  }

  return (
    <div className="space-y-8 p-6">
      <div className="text-center space-y-4">
        <motion.h1 
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-4xl font-bold gradient-text"
        >
          Enhanced Button Showcase
        </motion.h1>
        <motion.p 
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="text-muted-foreground text-lg"
        >
          Discover our collection of beautiful, interactive buttons with animations and effects
        </motion.p>
      </div>

      {/* Main Button Variants */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
      >
        <Card className="glass-card">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Sparkles className="h-5 w-5" />
              Primary Button Variants
            </CardTitle>
            <CardDescription>
              Enhanced buttons with gradients, animations, and interactive effects
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              
              {/* Default Enhanced */}
              <div className="space-y-3">
                <Badge variant="secondary">Default Enhanced</Badge>
                <EnhancedButton 
                  size="lg" 
                  className="w-full"
                  loading={loading === 'default'}
                  onClick={() => handleLoadingDemo('default')}
                >
                  <Zap className="mr-2 h-4 w-4" />
                  Primary Action
                </EnhancedButton>
              </div>

              {/* Gradient */}
              <div className="space-y-3">
                <Badge variant="secondary">Gradient</Badge>
                <EnhancedButton 
                  variant="gradient" 
                  size="lg" 
                  className="w-full"
                  glow="emerald"
                >
                  <Heart className="mr-2 h-4 w-4" />
                  Gradient Style
                </EnhancedButton>
              </div>

              {/* Glass */}
              <div className="space-y-3">
                <Badge variant="secondary">Glass Morphism</Badge>
                <EnhancedButton 
                  variant="glass" 
                  size="lg" 
                  className="w-full"
                >
                  <Star className="mr-2 h-4 w-4" />
                  Glass Effect
                </EnhancedButton>
              </div>

              {/* Premium */}
              <div className="space-y-3">
                <Badge variant="secondary">Premium</Badge>
                <EnhancedButton 
                  variant="premium" 
                  size="lg" 
                  className="w-full"
                  glow="yellow"
                >
                  <Crown className="mr-2 h-4 w-4" />
                  Premium Feature
                </EnhancedButton>
              </div>

              {/* Neon */}
              <div className="space-y-3">
                <Badge variant="secondary">Neon Glow</Badge>
                <EnhancedButton 
                  variant="neon" 
                  size="lg" 
                  className="w-full"
                  pulse={true}
                  glow="cyan"
                >
                  <Rocket className="mr-2 h-4 w-4" />
                  Neon Effect
                </EnhancedButton>
              </div>

              {/* Success */}
              <div className="space-y-3">
                <Badge variant="secondary">Success</Badge>
                <EnhancedButton 
                  variant="success" 
                  size="lg" 
                  className="w-full"
                  loading={loading === 'success'}
                  loadingText="Processing..."
                  onClick={() => handleLoadingDemo('success')}
                >
                  <CheckCircle className="mr-2 h-4 w-4" />
                  Success Action
                </EnhancedButton>
              </div>

            </div>
          </CardContent>
        </Card>
      </motion.div>

      {/* Action Buttons */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
      >
        <Card className="glass-card">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Play className="h-5 w-5" />
              Action & Utility Buttons
            </CardTitle>
            <CardDescription>
              Specialized buttons for different use cases and contexts
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              
              <EnhancedButton variant="outline" size="sm">
                <Download className="mr-2 h-4 w-4" />
                Download
              </EnhancedButton>

              <EnhancedButton variant="secondary" size="sm">
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </EnhancedButton>

              <EnhancedButton variant="destructive" size="sm">
                <AlertTriangle className="mr-2 h-4 w-4" />
                Delete
              </EnhancedButton>

              <EnhancedButton variant="ghost" size="sm">
                <Info className="mr-2 h-4 w-4" />
                Info
              </EnhancedButton>

              <EnhancedButton variant="warning" size="sm" glow="red">
                <Shield className="mr-2 h-4 w-4" />
                Warning
              </EnhancedButton>

              <EnhancedButton variant="glass" size="sm">
                <Gift className="mr-2 h-4 w-4" />
                Gift
              </EnhancedButton>

            </div>
          </CardContent>
        </Card>
      </motion.div>

      {/* Size Variations */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5 }}
      >
        <Card className="glass-card">
          <CardHeader>
            <CardTitle>Size Variations</CardTitle>
            <CardDescription>
              Different button sizes for various use cases
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap items-center gap-4">
              
              <EnhancedButton size="sm" variant="outline">
                Small
              </EnhancedButton>

              <EnhancedButton size="default" variant="default">
                Default
              </EnhancedButton>

              <EnhancedButton size="lg" variant="gradient">
                Large
              </EnhancedButton>

              <EnhancedButton size="xl" variant="premium" glow="yellow">
                Extra Large
              </EnhancedButton>

            </div>
          </CardContent>
        </Card>
      </motion.div>

      {/* Icon Buttons */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.6 }}
      >
        <Card className="glass-card">
          <CardHeader>
            <CardTitle>Icon Buttons</CardTitle>
            <CardDescription>
              Compact icon-only buttons for toolbars and actions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap items-center gap-4">
              
              <EnhancedButton size="icon-sm" variant="ghost">
                <Play className="h-4 w-4" />
              </EnhancedButton>

              <EnhancedButton size="icon" variant="default">
                <Pause className="h-4 w-4" />
              </EnhancedButton>

              <EnhancedButton size="icon-lg" variant="gradient" glow="emerald">
                <ShoppingCart className="h-5 w-5" />
              </EnhancedButton>

              <EnhancedButton size="icon" variant="neon" pulse={true}>
                <Zap className="h-4 w-4" />
              </EnhancedButton>

            </div>
          </CardContent>
        </Card>
      </motion.div>

    </div>
  )
}
