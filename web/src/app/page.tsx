'use client'

import { motion } from 'framer-motion'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { FeatureStatus, FeatureBadge, QuickStats } from '@/components/ui/feature-status'
import {
  Bot,
  Globe,
  Wallet,
  Zap,
  Shield,
  Cpu,
  ArrowRight,
  MessageSquare,
  BarChart3
} from 'lucide-react'
import Link from 'next/link'

const features = [
  {
    icon: Bot,
    title: 'AI-Powered Automation',
    description: 'Intelligent agents that understand natural language and automate complex web interactions.',
    color: 'text-blue-500'
  },
  {
    icon: Globe,
    title: 'Smart Web Navigation',
    description: 'Autonomous browsing with context-aware navigation and content extraction.',
    color: 'text-green-500'
  },
  {
    icon: Wallet,
    title: 'Web3 Integration',
    description: 'Seamless cryptocurrency wallet connection and DeFi protocol interactions.',
    color: 'text-purple-500'
  },
  {
    icon: Zap,
    title: 'Real-time Processing',
    description: 'Lightning-fast response times with real-time WebSocket communication.',
    color: 'text-yellow-500'
  },
  {
    icon: Shield,
    title: 'Enterprise Security',
    description: 'Bank-grade security with JWT authentication and encrypted communications.',
    color: 'text-red-500'
  },
  {
    icon: Cpu,
    title: 'Scalable Architecture',
    description: 'Microservices architecture with comprehensive monitoring and observability.',
    color: 'text-indigo-500'
  }
]

const stats = [
  { label: 'AI Models Integrated', value: '3+' },
  { label: 'Blockchain Networks', value: '5+' },
  { label: 'DeFi Protocols', value: '10+' },
  { label: 'Response Time', value: '<2s' }
]

export default function HomePage() {
  const router = useRouter()

  const handleNavigation = (path: string) => {
    router.push(path)
  }

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden min-h-screen flex items-center">
        {/* Animated background */}
        <div className="absolute inset-0 bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 dark:from-gray-900 dark:via-blue-900/20 dark:to-purple-900/20" />
        <div className="absolute inset-0 bg-grid-pattern opacity-30" />
        <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-transparent to-background/80" />

        {/* Floating elements */}
        <div className="absolute top-20 left-10 w-72 h-72 bg-gradient-to-r from-blue-400/20 to-purple-400/20 rounded-full blur-3xl animate-float" />
        <div className="absolute bottom-20 right-10 w-96 h-96 bg-gradient-to-r from-purple-400/20 to-pink-400/20 rounded-full blur-3xl animate-float" style={{ animationDelay: '2s' }} />
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-64 h-64 bg-gradient-to-r from-cyan-400/20 to-blue-400/20 rounded-full blur-3xl animate-float" style={{ animationDelay: '4s' }} />

        <div className="relative container mx-auto px-4 py-20 lg:py-32 z-10">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 1, ease: "easeOut" }}
            className="text-center max-w-5xl mx-auto"
          >
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.8, delay: 0.2 }}
            >
              <Badge variant="outline" className="mb-8 px-4 py-2 text-sm font-medium bg-white/10 backdrop-blur-sm border-white/20 hover:bg-white/20 transition-all duration-300">
                🚀 Now in Beta - Experience the Future
              </Badge>
            </motion.div>

            <motion.h1
              className="text-5xl lg:text-7xl xl:text-8xl font-bold mb-8 gradient-text-crypto leading-tight"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.4 }}
            >
              AI-Powered Agentic
              <br />
              <span className="gradient-text">Crypto Browser</span>
            </motion.h1>

            <motion.p
              className="text-xl lg:text-2xl xl:text-3xl text-muted-foreground mb-12 leading-relaxed max-w-4xl mx-auto"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.6 }}
            >
              The future of web interaction is here. Combine the power of AI agents with
              Web3 technology for autonomous browsing and cryptocurrency management.
            </motion.p>

            {/* Enhanced Feature Cards */}
            <motion.div
              className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.8 }}
            >
              {/* HFT Trading - Primary Feature */}
              <motion.div
                whileHover={{ scale: 1.02, y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
                className="lg:col-span-3"
              >
                <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-blue-600 via-purple-600 to-indigo-700 p-8 text-white shadow-2xl hover:shadow-3xl transition-all duration-300 group cursor-pointer"
                     onClick={() => handleNavigation('/trading')}>
                    <div className="absolute inset-0 bg-grid-pattern opacity-10" />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent" />
                    <div className="relative flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="p-4 bg-white/20 rounded-xl backdrop-blur-sm group-hover:bg-white/30 transition-colors">
                          <Zap className="h-8 w-8" />
                        </div>
                        <div>
                          <h3 className="text-2xl font-bold mb-2">High-Frequency Trading</h3>
                          <p className="text-blue-100 text-lg">
                            Advanced algorithmic trading with microsecond execution
                          </p>
                        </div>
                      </div>
                      <div className="hidden md:flex items-center gap-2 text-blue-200">
                        <span className="text-sm">Start Trading</span>
                        <ArrowRight className="h-5 w-5 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </div>
              </motion.div>

              {/* Performance Analytics */}
              <motion.div
                whileHover={{ scale: 1.02, y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/50 dark:to-teal-950/50 p-6 shadow-xl hover:shadow-2xl transition-all duration-300 group border border-emerald-200/50 dark:border-emerald-800/50 h-full cursor-pointer"
                     onClick={() => handleNavigation('/performance')}>
                    <div className="absolute top-0 right-0 w-20 h-20 bg-emerald-500/10 rounded-full -translate-y-10 translate-x-10" />
                    <div className="relative">
                      <div className="flex items-center gap-3 mb-4">
                        <div className="p-3 bg-emerald-100 dark:bg-emerald-900/50 rounded-xl group-hover:bg-emerald-200 dark:group-hover:bg-emerald-800/70 transition-colors">
                          <BarChart3 className="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
                        </div>
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="text-xl font-bold text-emerald-800 dark:text-emerald-200">Performance</h3>
                            <FeatureBadge type="trending" className="text-xs" />
                          </div>
                          <FeatureStatus status="online" label="Live Monitoring" className="mt-1" />
                        </div>
                      </div>
                      <p className="text-emerald-700 dark:text-emerald-300 mb-4">
                        Real-time analytics and performance optimization insights
                      </p>
                      <QuickStats
                        stats={[
                          { label: 'Uptime', value: '99.9%', color: 'green', trend: 'up' },
                          { label: 'Response', value: '12ms', color: 'blue', trend: 'up' }
                        ]}
                        className="mb-4"
                      />
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <div className="text-xs bg-emerald-100 dark:bg-emerald-900/50 text-emerald-700 dark:text-emerald-300 px-2 py-1 rounded-full">
                            Analytics
                          </div>
                          <div className="text-xs bg-emerald-100 dark:bg-emerald-900/50 text-emerald-700 dark:text-emerald-300 px-2 py-1 rounded-full">
                            Metrics
                          </div>
                        </div>
                        <ArrowRight className="h-4 w-4 text-emerald-600 dark:text-emerald-400 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </div>
              </motion.div>

              {/* Compliance & Risk */}
              <motion.div
                whileHover={{ scale: 1.02, y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-orange-50 to-red-50 dark:from-orange-950/50 dark:to-red-950/50 p-6 shadow-xl hover:shadow-2xl transition-all duration-300 group border border-orange-200/50 dark:border-orange-800/50 h-full cursor-pointer"
                     onClick={() => handleNavigation('/compliance')}>
                    <div className="absolute top-0 right-0 w-20 h-20 bg-orange-500/10 rounded-full -translate-y-10 translate-x-10" />
                    <div className="relative">
                      <div className="flex items-center gap-3 mb-4">
                        <div className="p-3 bg-orange-100 dark:bg-orange-900/50 rounded-xl group-hover:bg-orange-200 dark:group-hover:bg-orange-800/70 transition-colors">
                          <Shield className="h-6 w-6 text-orange-600 dark:text-orange-400" />
                        </div>
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="text-xl font-bold text-orange-800 dark:text-orange-200">Compliance</h3>
                            <FeatureBadge type="premium" className="text-xs" />
                          </div>
                          <FeatureStatus status="online" label="98% Compliant" className="mt-1" />
                        </div>
                      </div>
                      <p className="text-orange-700 dark:text-orange-300 mb-4">
                        Regulatory compliance monitoring and risk management
                      </p>
                      <QuickStats
                        stats={[
                          { label: 'Risk Score', value: 'Low', color: 'green' },
                          { label: 'Violations', value: '0', color: 'green' }
                        ]}
                        className="mb-4"
                      />
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <div className="text-xs bg-orange-100 dark:bg-orange-900/50 text-orange-700 dark:text-orange-300 px-2 py-1 rounded-full">
                            Risk
                          </div>
                          <div className="text-xs bg-orange-100 dark:bg-orange-900/50 text-orange-700 dark:text-orange-300 px-2 py-1 rounded-full">
                            Audit
                          </div>
                        </div>
                        <ArrowRight className="h-4 w-4 text-orange-600 dark:text-orange-400 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </div>
              </motion.div>

              {/* Web3 Features */}
              <motion.div
                whileHover={{ scale: 1.02, y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-purple-50 to-pink-50 dark:from-purple-950/50 dark:to-pink-950/50 p-6 shadow-xl hover:shadow-2xl transition-all duration-300 group border border-purple-200/50 dark:border-purple-800/50 h-full cursor-pointer"
                     onClick={() => handleNavigation('/web3')}>
                    <div className="absolute top-0 right-0 w-20 h-20 bg-purple-500/10 rounded-full -translate-y-10 translate-x-10" />
                    <div className="relative">
                      <div className="flex items-center gap-3 mb-4">
                        <div className="p-3 bg-purple-100 dark:bg-purple-900/50 rounded-xl group-hover:bg-purple-200 dark:group-hover:bg-purple-800/70 transition-colors">
                          <Wallet className="h-6 w-6 text-purple-600 dark:text-purple-400" />
                        </div>
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="text-xl font-bold text-purple-800 dark:text-purple-200">Web3 Features</h3>
                            <FeatureBadge type="hot" className="text-xs" />
                          </div>
                          <FeatureStatus status="beta" label="Blockchain Ready" className="mt-1" />
                        </div>
                      </div>
                      <p className="text-purple-700 dark:text-purple-300 mb-4">
                        Decentralized finance and blockchain integration tools
                      </p>
                      <QuickStats
                        stats={[
                          { label: 'Networks', value: '5+', color: 'purple' },
                          { label: 'Protocols', value: '12', color: 'blue' }
                        ]}
                        className="mb-4"
                      />
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <div className="text-xs bg-purple-100 dark:bg-purple-900/50 text-purple-700 dark:text-purple-300 px-2 py-1 rounded-full">
                            DeFi
                          </div>
                          <div className="text-xs bg-purple-100 dark:bg-purple-900/50 text-purple-700 dark:text-purple-300 px-2 py-1 rounded-full">
                            NFTs
                          </div>
                        </div>
                        <ArrowRight className="h-4 w-4 text-purple-600 dark:text-purple-400 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </div>
              </motion.div>

              {/* Dashboard */}
              <motion.div
                whileHover={{ scale: 1.02, y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-950/50 dark:to-indigo-950/50 p-6 shadow-xl hover:shadow-2xl transition-all duration-300 group border border-blue-200/50 dark:border-blue-800/50 h-full cursor-pointer"
                     onClick={() => handleNavigation('/dashboard')}>
                    <div className="absolute top-0 right-0 w-20 h-20 bg-blue-500/10 rounded-full -translate-y-10 translate-x-10" />
                    <div className="relative">
                      <div className="flex items-center gap-3 mb-4">
                        <div className="p-3 bg-blue-100 dark:bg-blue-900/50 rounded-xl group-hover:bg-blue-200 dark:group-hover:bg-blue-800/70 transition-colors">
                          <BarChart3 className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                        </div>
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="text-xl font-bold text-blue-800 dark:text-blue-200">Dashboard</h3>
                            <FeatureBadge type="popular" className="text-xs" />
                          </div>
                          <FeatureStatus status="online" label="All Systems" className="mt-1" />
                        </div>
                      </div>
                      <p className="text-blue-700 dark:text-blue-300 mb-4">
                        Comprehensive overview of all trading activities
                      </p>
                      <QuickStats
                        stats={[
                          { label: 'Active', value: '24/7', color: 'blue' },
                          { label: 'Widgets', value: '15+', color: 'green' }
                        ]}
                        className="mb-4"
                      />
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <div className="text-xs bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 px-2 py-1 rounded-full">
                            Overview
                          </div>
                          <div className="text-xs bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 px-2 py-1 rounded-full">
                            Reports
                          </div>
                        </div>
                        <ArrowRight className="h-4 w-4 text-blue-600 dark:text-blue-400 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </div>
              </motion.div>

              {/* Solana Trading - Special Feature */}
              <motion.div
                whileHover={{ scale: 1.02, y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
                className="lg:col-span-2"
              >
                <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-600 p-6 text-white shadow-2xl hover:shadow-3xl transition-all duration-300 group cursor-pointer"
                     onClick={() => handleNavigation('/solana')}>
                    <div className="absolute inset-0 bg-grid-pattern opacity-10" />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent" />
                    <div className="relative flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="p-3 bg-white/20 rounded-xl backdrop-blur-sm group-hover:bg-white/30 transition-colors">
                          <Zap className="h-7 w-7" />
                        </div>
                        <div>
                          <h3 className="text-xl font-bold mb-2">Solana Trading</h3>
                          <p className="text-emerald-100">
                            High-speed trading on the Solana blockchain
                          </p>
                          <div className="flex items-center gap-3 mt-2">
                            <div className="text-xs bg-white/20 text-white px-2 py-1 rounded-full">
                              Ultra Fast
                            </div>
                            <div className="text-xs bg-white/20 text-white px-2 py-1 rounded-full">
                              Low Fees
                            </div>
                          </div>
                        </div>
                      </div>
                      <div className="hidden md:flex items-center gap-2 text-emerald-200">
                        <span className="text-sm">Trade Now</span>
                        <ArrowRight className="h-5 w-5 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </div>
              </motion.div>
            </motion.div>
          </motion.div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-20 relative">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-50/50 via-purple-50/50 to-pink-50/50 dark:from-blue-900/10 dark:via-purple-900/10 dark:to-pink-900/10" />
        <div className="container mx-auto px-4 relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl lg:text-4xl font-bold mb-4 gradient-text">
              Trusted by Thousands
            </h2>
            <p className="text-xl text-muted-foreground">
              Real metrics from our growing community
            </p>
          </motion.div>

          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, index) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 30, scale: 0.8 }}
                whileInView={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="text-center group"
              >
                <div className="glass-card p-6 rounded-2xl hover:scale-105 transition-all duration-300 hover:shadow-2xl">
                  <div className="text-4xl lg:text-5xl font-bold gradient-text mb-3 group-hover:animate-pulse">
                    {stat.value}
                  </div>
                  <div className="text-muted-foreground font-medium">{stat.label}</div>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-24 relative">
        <div className="absolute inset-0 bg-dot-pattern opacity-20" />
        <div className="container mx-auto px-4 relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="text-center mb-20"
          >
            <h2 className="text-4xl lg:text-5xl font-bold mb-6 gradient-text">
              Powerful Features
            </h2>
            <p className="text-xl lg:text-2xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
              Everything you need for intelligent web automation and cryptocurrency management
            </p>
          </motion.div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={feature.title}
                initial={{ opacity: 0, y: 30, scale: 0.9 }}
                whileInView={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                viewport={{ once: true }}
                whileHover={{ y: -5, scale: 1.02 }}
                className="group"
              >
                <Card className="h-full glass-card border-0 hover:shadow-2xl transition-all duration-500 overflow-hidden relative">
                  {/* Gradient overlay on hover */}
                  <div className="absolute inset-0 bg-gradient-to-br from-transparent via-transparent to-primary/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                  <CardHeader className="relative z-10">
                    <div className="flex items-center gap-4 mb-4">
                      <div className={`p-3 rounded-xl bg-gradient-to-br from-white/20 to-white/10 backdrop-blur-sm ${feature.color} group-hover:scale-110 transition-transform duration-300`}>
                        <feature.icon className="h-7 w-7" />
                      </div>
                      <CardTitle className="text-xl lg:text-2xl font-semibold group-hover:text-primary transition-colors duration-300">
                        {feature.title}
                      </CardTitle>
                    </div>
                  </CardHeader>
                  <CardContent className="relative z-10">
                    <CardDescription className="text-base lg:text-lg leading-relaxed text-muted-foreground group-hover:text-foreground transition-colors duration-300">
                      {feature.description}
                    </CardDescription>
                  </CardContent>

                  {/* Animated border */}
                  <div className="absolute inset-0 rounded-xl bg-gradient-to-r from-primary/20 via-purple-500/20 to-pink-500/20 opacity-0 group-hover:opacity-100 transition-opacity duration-500 -z-10"
                       style={{ padding: '1px' }}>
                    <div className="w-full h-full rounded-xl bg-background/80 backdrop-blur-sm" />
                  </div>
                </Card>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-blue-600/10 via-purple-600/10 to-pink-600/10" />
        <div className="absolute inset-0 bg-grid-pattern opacity-20" />

        {/* Floating elements */}
        <div className="absolute top-10 left-10 w-32 h-32 bg-gradient-to-r from-blue-400/30 to-purple-400/30 rounded-full blur-2xl animate-float" />
        <div className="absolute bottom-10 right-10 w-40 h-40 bg-gradient-to-r from-purple-400/30 to-pink-400/30 rounded-full blur-2xl animate-float" style={{ animationDelay: '2s' }} />

        <div className="container mx-auto px-4 text-center relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="max-w-4xl mx-auto"
          >
            <h2 className="text-4xl lg:text-5xl xl:text-6xl font-bold mb-8 gradient-text">
              Ready to Experience the Future?
            </h2>
            <p className="text-xl lg:text-2xl text-muted-foreground mb-12 leading-relaxed">
              Join thousands of users who are already using AI agents to automate their web interactions
              and manage their cryptocurrency portfolios.
            </p>
            <div className="flex flex-col sm:flex-row gap-6 justify-center">
              <Button asChild size="xl" variant="default" className="text-lg px-10 py-4 animate-pulse-glow">
                <Link href="/auth/register">
                  Start Free Trial <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button asChild variant="glass" size="xl" className="text-lg px-10 py-4">
                <Link href="/docs">
                  Read Documentation
                </Link>
              </Button>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-16 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-t from-secondary/50 to-background" />
        <div className="absolute inset-0 bg-dot-pattern opacity-10" />

        <div className="container mx-auto px-4 relative z-10">
          <div className="grid md:grid-cols-4 gap-12">
            <div className="md:col-span-2">
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6 }}
                viewport={{ once: true }}
              >
                <h3 className="text-2xl font-bold mb-6 gradient-text">AI Agentic Browser</h3>
                <p className="text-muted-foreground mb-6 text-lg leading-relaxed">
                  The next generation of web browsing powered by artificial intelligence
                  and blockchain technology.
                </p>
                <div className="flex gap-4">
                  <Button variant="ghost" size="sm" className="glass-card hover:bg-white/20 transition-all duration-300 hover:scale-110">
                    <Globe className="h-5 w-5" />
                  </Button>
                  <Button variant="ghost" size="sm" className="glass-card hover:bg-white/20 transition-all duration-300 hover:scale-110">
                    <Bot className="h-5 w-5" />
                  </Button>
                  <Button variant="ghost" size="sm" className="glass-card hover:bg-white/20 transition-all duration-300 hover:scale-110">
                    <MessageSquare className="h-5 w-5" />
                  </Button>
                </div>
              </motion.div>
            </div>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
              viewport={{ once: true }}
            >
              <h4 className="font-semibold mb-6 text-lg">Product</h4>
              <ul className="space-y-3 text-muted-foreground">
                <li><Link href="/features" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Features</Link></li>
                <li><Link href="/pricing" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Pricing</Link></li>
                <li><Link href="/demo" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Demo</Link></li>
                <li><Link href="/docs" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Documentation</Link></li>
              </ul>
            </motion.div>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              viewport={{ once: true }}
            >
              <h4 className="font-semibold mb-6 text-lg">Company</h4>
              <ul className="space-y-3 text-muted-foreground">
                <li><Link href="/about" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">About</Link></li>
                <li><Link href="/blog" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Blog</Link></li>
                <li><Link href="/careers" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Careers</Link></li>
                <li><Link href="/contact" className="hover:text-primary transition-colors duration-300 hover:translate-x-1 inline-block">Contact</Link></li>
              </ul>
            </motion.div>
          </div>
          <motion.div
            className="border-t border-border/50 mt-12 pt-8 text-center"
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            viewport={{ once: true }}
          >
            <p className="text-muted-foreground">
              &copy; 2024 AI Agentic Browser. All rights reserved. Built with ❤️ for the future of web interaction.
            </p>
          </motion.div>
        </div>
      </footer>
    </div>
  )
}
