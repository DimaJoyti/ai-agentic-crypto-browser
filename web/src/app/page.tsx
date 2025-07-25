'use client'

import { useState } from 'react'
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
  Github,
  Twitter,
  MessageSquare
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
  const [isHovered, setIsHovered] = useState<number | null>(null)

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-br from-background via-background to-secondary/20">
        <div className="absolute inset-0 bg-grid-pattern opacity-5" />
        <div className="relative container mx-auto px-4 py-20 lg:py-32">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center max-w-4xl mx-auto"
          >
            <Badge variant="outline" className="mb-6">
              🚀 Now in Beta
            </Badge>
            <h1 className="text-4xl lg:text-6xl font-bold mb-6 gradient-text">
              AI-Powered Agentic Crypto Browser
            </h1>
            <p className="text-xl lg:text-2xl text-muted-foreground mb-8 leading-relaxed">
              The future of web interaction is here. Combine the power of AI agents with 
              Web3 technology for autonomous browsing and cryptocurrency management.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button asChild size="lg" className="text-lg px-8">
                <Link href="/auth/login">
                  Get Started <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-8">
                <Link href="/demo">
                  View Demo
                </Link>
              </Button>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-16 bg-secondary/30">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, index) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                className="text-center"
              >
                <div className="text-3xl lg:text-4xl font-bold text-primary mb-2">
                  {stat.value}
                </div>
                <div className="text-muted-foreground">{stat.label}</div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl lg:text-4xl font-bold mb-4">
              Powerful Features
            </h2>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Everything you need for intelligent web automation and cryptocurrency management
            </p>
          </motion.div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={feature.title}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                onHoverStart={() => setIsHovered(index)}
                onHoverEnd={() => setIsHovered(null)}
              >
                <Card className="h-full transition-all duration-300 hover:shadow-lg hover:scale-105">
                  <CardHeader>
                    <div className="flex items-center gap-3 mb-2">
                      <div className={`p-2 rounded-lg bg-secondary ${feature.color}`}>
                        <feature.icon className="h-6 w-6" />
                      </div>
                      <CardTitle className="text-xl">{feature.title}</CardTitle>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <CardDescription className="text-base leading-relaxed">
                      {feature.description}
                    </CardDescription>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-primary/10 via-secondary/10 to-accent/10">
        <div className="container mx-auto px-4 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="max-w-3xl mx-auto"
          >
            <h2 className="text-3xl lg:text-4xl font-bold mb-6">
              Ready to Experience the Future?
            </h2>
            <p className="text-xl text-muted-foreground mb-8">
              Join thousands of users who are already using AI agents to automate their web interactions
              and manage their cryptocurrency portfolios.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button asChild size="lg" className="text-lg px-8">
                <Link href="/auth/register">
                  Start Free Trial <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="text-lg px-8">
                <Link href="/docs">
                  Read Documentation
                </Link>
              </Button>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 border-t border-border">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8">
            <div className="md:col-span-2">
              <h3 className="text-lg font-semibold mb-4">AI Agentic Browser</h3>
              <p className="text-muted-foreground mb-4">
                The next generation of web browsing powered by artificial intelligence
                and blockchain technology.
              </p>
              <div className="flex gap-4">
                <Button variant="ghost" size="sm">
                  <Github className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="sm">
                  <Twitter className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="sm">
                  <MessageSquare className="h-4 w-4" />
                </Button>
              </div>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Product</h4>
              <ul className="space-y-2 text-muted-foreground">
                <li><Link href="/features" className="hover:text-foreground">Features</Link></li>
                <li><Link href="/pricing" className="hover:text-foreground">Pricing</Link></li>
                <li><Link href="/demo" className="hover:text-foreground">Demo</Link></li>
                <li><Link href="/docs" className="hover:text-foreground">Documentation</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Company</h4>
              <ul className="space-y-2 text-muted-foreground">
                <li><Link href="/about" className="hover:text-foreground">About</Link></li>
                <li><Link href="/blog" className="hover:text-foreground">Blog</Link></li>
                <li><Link href="/careers" className="hover:text-foreground">Careers</Link></li>
                <li><Link href="/contact" className="hover:text-foreground">Contact</Link></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-border mt-8 pt-8 text-center text-muted-foreground">
            <p>&copy; 2024 AI Agentic Browser. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
