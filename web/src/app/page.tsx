'use client'

import { motion } from 'framer-motion'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
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

            <motion.div
              className="flex flex-col sm:flex-row gap-4 justify-center flex-wrap"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.8 }}
            >
              <Button asChild size="lg" className="text-lg px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 shadow-lg hover:shadow-xl transition-all duration-300 animate-pulse-glow">
                <Link href="/trading">
                  HFT Trading <Zap className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-8 py-4 glass-card hover:bg-white/20 transition-all duration-300">
                <Link href="/performance">
                  Performance <BarChart3 className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-8 py-4 glass-card hover:bg-white/20 transition-all duration-300">
                <Link href="/compliance">
                  Compliance <Shield className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-8 py-4 glass-card hover:bg-white/20 transition-all duration-300">
                <Link href="/web3">
                  Web3 Features <Wallet className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-8 py-4 glass-card hover:bg-white/20 transition-all duration-300">
                <Link href="/dashboard">
                  Dashboard <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
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
              <Button asChild size="lg" className="text-lg px-10 py-4 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 shadow-xl hover:shadow-2xl transition-all duration-300 animate-pulse-glow">
                <Link href="/auth/register">
                  Start Free Trial <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-10 py-4 glass-card hover:bg-white/20 transition-all duration-300">
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
