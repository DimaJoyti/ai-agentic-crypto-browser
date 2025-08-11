'use client'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  BarChart3,
  TrendingUp,
  Users,
  Shield,
  Zap,
  Target,
  Coins,
  Activity,
  PieChart,
  ArrowRight,
  Building2,
  CreditCard
} from 'lucide-react'

// Import our new components
import { TradingInterface } from '@/components/trading/TradingInterface'
import { EnhancedDeFiDashboard } from '@/components/defi/EnhancedDeFiDashboard'
import { AdvancedPortfolioAnalytics } from '@/components/analytics/AdvancedPortfolioAnalytics'
import { StakingYieldPlatform } from '@/components/staking/StakingYieldPlatform'
import { SocialTradingPlatform } from '@/components/social/SocialTradingPlatform'
import { AdvancedSecurityDashboard } from '@/components/security/AdvancedSecurityDashboard'
import { AdvancedOrderTypes } from '@/components/trading/AdvancedOrderTypes'
import { InstitutionalTradingHub } from '@/components/institutional/InstitutionalTradingHub'
import { FiatGatewayHub } from '@/components/fiat/FiatGatewayHub'
import { LaunchpadDashboard } from '@/components/launchpad/LaunchpadDashboard'
import { ProjectDetails } from '@/components/launchpad/ProjectDetails'
import { LaunchpadKYC } from '@/components/launchpad/LaunchpadKYC'
import { VestingDashboard } from '@/components/launchpad/VestingDashboard'
import { ProjectSubmission } from '@/components/launchpad/ProjectSubmission'

export default function KabePage() {
  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white">
        <div className="container mx-auto px-4 py-12">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold mb-4">
              Web3 Platform
            </h1>
            <p className="text-xl md:text-2xl mb-8 opacity-90">
              Professional cryptocurrency trading with advanced DeFi features
            </p>
            <div className="flex flex-wrap justify-center gap-4 mb-8">
              <Badge variant="secondary" className="text-lg px-4 py-2">
                <BarChart3 className="w-4 h-4 mr-2" />
                Advanced Trading
              </Badge>
              <Badge variant="secondary" className="text-lg px-4 py-2">
                <Coins className="w-4 h-4 mr-2" />
                DeFi Integration
              </Badge>
              <Badge variant="secondary" className="text-lg px-4 py-2">
                <Users className="w-4 h-4 mr-2" />
                Social Trading
              </Badge>
              <Badge variant="secondary" className="text-lg px-4 py-2">
                <Shield className="w-4 h-4 mr-2" />
                Portfolio Analytics
              </Badge>
            </div>
            <Button size="lg" variant="secondary" className="text-lg px-8 py-3">
              Get Started
              <ArrowRight className="w-5 h-5 ml-2" />
            </Button>
          </div>
        </div>
      </div>

      {/* Features Overview */}
      <div className="container mx-auto px-4 py-12">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold mb-4">Platform Features</h2>
          <p className="text-muted-foreground text-lg">
            Everything you need for professional cryptocurrency trading and DeFi participation
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
          <Card className="text-center hover:shadow-lg transition-shadow">
            <CardContent className="p-6">
              <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <BarChart3 className="w-6 h-6 text-blue-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">Advanced Trading</h3>
              <p className="text-muted-foreground text-sm">
                Professional trading interface with real-time charts, order book, and advanced order types
              </p>
            </CardContent>
          </Card>

          <Card className="text-center hover:shadow-lg transition-shadow">
            <CardContent className="p-6">
              <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Coins className="w-6 h-6 text-green-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">DeFi Hub</h3>
              <p className="text-muted-foreground text-sm">
                Comprehensive DeFi dashboard with yield farming, staking, and liquidity mining
              </p>
            </CardContent>
          </Card>

          <Card className="text-center hover:shadow-lg transition-shadow">
            <CardContent className="p-6">
              <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Users className="w-6 h-6 text-purple-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">Social Trading</h3>
              <p className="text-muted-foreground text-sm">
                Follow top traders, copy strategies, and share insights with the community
              </p>
            </CardContent>
          </Card>

          <Card className="text-center hover:shadow-lg transition-shadow">
            <CardContent className="p-6">
              <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <PieChart className="w-6 h-6 text-orange-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">Portfolio Analytics</h3>
              <p className="text-muted-foreground text-sm">
                Advanced analytics with performance metrics, risk assessment, and insights
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Main Platform Interface */}
        <Tabs defaultValue="trading" className="space-y-6">
          <TabsList className="grid w-full grid-cols-10">
            <TabsTrigger value="trading" className="flex items-center gap-2">
              <BarChart3 className="w-4 h-4" />
              Trading
            </TabsTrigger>
            <TabsTrigger value="advanced" className="flex items-center gap-2">
              <Target className="w-4 h-4" />
              Advanced
            </TabsTrigger>
            <TabsTrigger value="fiat" className="flex items-center gap-2">
              <CreditCard className="w-4 h-4" />
              Fiat
            </TabsTrigger>
            <TabsTrigger value="defi" className="flex items-center gap-2">
              <Coins className="w-4 h-4" />
              DeFi
            </TabsTrigger>
            <TabsTrigger value="staking" className="flex items-center gap-2">
              <Zap className="w-4 h-4" />
              Staking
            </TabsTrigger>
            <TabsTrigger value="launchpad" className="flex items-center gap-2">
              <TrendingUp className="w-4 h-4" />
              Launchpad
            </TabsTrigger>
            <TabsTrigger value="social" className="flex items-center gap-2">
              <Users className="w-4 h-4" />
              Social
            </TabsTrigger>
            <TabsTrigger value="analytics" className="flex items-center gap-2">
              <PieChart className="w-4 h-4" />
              Analytics
            </TabsTrigger>
            <TabsTrigger value="security" className="flex items-center gap-2">
              <Shield className="w-4 h-4" />
              Security
            </TabsTrigger>
            <TabsTrigger value="institutional" className="flex items-center gap-2">
              <Building2 className="w-4 h-4" />
              Enterprise
            </TabsTrigger>
          </TabsList>

          <TabsContent value="trading">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Advanced Trading Interface
                </CardTitle>
                <p className="text-muted-foreground">
                  Professional-grade trading platform with real-time data and advanced features
                </p>
              </CardHeader>
              <CardContent className="p-4">
                <TradingInterface embedded={true} />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="advanced">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Target className="w-5 h-5" />
                  Advanced Order Types
                </CardTitle>
                <p className="text-muted-foreground">
                  Professional order types including OCO, trailing stops, iceberg orders, and algorithmic strategies
                </p>
              </CardHeader>
              <CardContent>
                <AdvancedOrderTypes />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="fiat">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CreditCard className="w-5 h-5" />
                  Fiat Gateway
                </CardTitle>
                <p className="text-muted-foreground">
                  Buy and sell cryptocurrency with fiat currency using credit cards, bank transfers, and global payment methods
                </p>
              </CardHeader>
              <CardContent>
                <FiatGatewayHub />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="defi">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Coins className="w-5 h-5" />
                  DeFi Dashboard
                </CardTitle>
                <p className="text-muted-foreground">
                  Comprehensive DeFi platform with protocol integration and yield optimization
                </p>
              </CardHeader>
              <CardContent>
                <EnhancedDeFiDashboard />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="staking">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Staking & Yield Platform
                </CardTitle>
                <p className="text-muted-foreground">
                  Earn rewards through staking, yield farming, and liquidity provision
                </p>
              </CardHeader>
              <CardContent>
                <StakingYieldPlatform />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="social">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Users className="w-5 h-5" />
                  Social Trading Platform
                </CardTitle>
                <p className="text-muted-foreground">
                  Connect with traders, copy strategies, and share market insights
                </p>
              </CardHeader>
              <CardContent>
                <SocialTradingPlatform />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="analytics">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Portfolio Analytics
                </CardTitle>
                <p className="text-muted-foreground">
                  Advanced portfolio analysis with performance metrics and risk assessment
                </p>
              </CardHeader>
              <CardContent>
                <AdvancedPortfolioAnalytics />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="security">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Security Dashboard
                </CardTitle>
                <p className="text-muted-foreground">
                  Advanced security features including 2FA, withdrawal whitelist, and activity monitoring
                </p>
              </CardHeader>
              <CardContent>
                <AdvancedSecurityDashboard />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="launchpad">
            <LaunchpadDashboard />
          </TabsContent>

          <TabsContent value="institutional">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Building2 className="w-5 h-5" />
                  Institutional Trading
                </CardTitle>
                <p className="text-muted-foreground">
                  Enterprise-grade features including API management, OTC desk, and custody solutions
                </p>
              </CardHeader>
              <CardContent>
                <InstitutionalTradingHub />
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Key Statistics */}
        <div className="mt-16 grid grid-cols-1 md:grid-cols-4 gap-6">
          <Card className="text-center">
            <CardContent className="p-6">
              <div className="text-3xl font-bold text-blue-600 mb-2">50+</div>
              <div className="text-muted-foreground">Trading Pairs</div>
            </CardContent>
          </Card>
          
          <Card className="text-center">
            <CardContent className="p-6">
              <div className="text-3xl font-bold text-green-600 mb-2">15+</div>
              <div className="text-muted-foreground">DeFi Protocols</div>
            </CardContent>
          </Card>
          
          <Card className="text-center">
            <CardContent className="p-6">
              <div className="text-3xl font-bold text-purple-600 mb-2">24/7</div>
              <div className="text-muted-foreground">Real-time Data</div>
            </CardContent>
          </Card>
          
          <Card className="text-center">
            <CardContent className="p-6">
              <div className="text-3xl font-bold text-orange-600 mb-2">99.9%</div>
              <div className="text-muted-foreground">Uptime</div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Footer */}
      <div className="bg-muted/50 mt-16">
        <div className="container mx-auto px-4 py-8">
          <div className="text-center">
            <h3 className="text-xl font-bold mb-2">Ready to Start Trading?</h3>
            <p className="text-muted-foreground mb-4">
              Join thousands of traders using our advanced Web3 platform
            </p>
            <Button size="lg">
              Connect Wallet
              <Zap className="w-4 h-4 ml-2" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
