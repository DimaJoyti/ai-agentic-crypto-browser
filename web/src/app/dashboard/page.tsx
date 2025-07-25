'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  Bot, 
  Globe, 
  Wallet, 
  Activity, 
  TrendingUp, 
  Zap,
  Eye,
  MessageSquare,
  Settings,
  BarChart3,
  DollarSign,
  Users,
  Clock
} from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { useWebSocket } from '@/hooks/useWebSocket'
import { formatCurrency, formatRelativeTime } from '@/lib/utils'

interface DashboardStats {
  totalSessions: number
  activeTasks: number
  completedTasks: number
  totalValue: number
  portfolioChange: number
  aiInteractions: number
}

interface RecentActivity {
  id: string
  type: 'ai_chat' | 'browser_action' | 'web3_transaction' | 'task_completed'
  title: string
  description: string
  timestamp: Date
  status: 'success' | 'pending' | 'failed'
  metadata?: any
}

export default function DashboardPage() {
  const { user } = useAuth()
  const { isConnected, lastMessage } = useWebSocket()
  const [stats, setStats] = useState<DashboardStats>({
    totalSessions: 0,
    activeTasks: 0,
    completedTasks: 0,
    totalValue: 0,
    portfolioChange: 0,
    aiInteractions: 0,
  })
  const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    // Simulate loading dashboard data
    const loadDashboardData = async () => {
      setIsLoading(true)
      
      // Simulate API calls
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setStats({
        totalSessions: 24,
        activeTasks: 3,
        completedTasks: 156,
        totalValue: 12450.75,
        portfolioChange: 8.2,
        aiInteractions: 89,
      })

      setRecentActivity([
        {
          id: '1',
          type: 'ai_chat',
          title: 'AI Chat Session',
          description: 'Automated web scraping task completed',
          timestamp: new Date(Date.now() - 5 * 60 * 1000),
          status: 'success',
        },
        {
          id: '2',
          type: 'web3_transaction',
          title: 'DeFi Transaction',
          description: 'Added liquidity to ETH-USDC pool on Uniswap',
          timestamp: new Date(Date.now() - 15 * 60 * 1000),
          status: 'success',
        },
        {
          id: '3',
          type: 'browser_action',
          title: 'Browser Automation',
          description: 'Form submission on contact page',
          timestamp: new Date(Date.now() - 30 * 60 * 1000),
          status: 'pending',
        },
        {
          id: '4',
          type: 'task_completed',
          title: 'Task Completed',
          description: 'Content extraction from 5 websites',
          timestamp: new Date(Date.now() - 45 * 60 * 1000),
          status: 'success',
        },
      ])

      setIsLoading(false)
    }

    loadDashboardData()
  }, [])

  const getActivityIcon = (type: RecentActivity['type']) => {
    switch (type) {
      case 'ai_chat':
        return <Bot className="h-4 w-4" />
      case 'browser_action':
        return <Globe className="h-4 w-4" />
      case 'web3_transaction':
        return <Wallet className="h-4 w-4" />
      case 'task_completed':
        return <Zap className="h-4 w-4" />
      default:
        return <Activity className="h-4 w-4" />
    }
  }

  const getStatusColor = (status: RecentActivity['status']) => {
    switch (status) {
      case 'success':
        return 'bg-green-500'
      case 'pending':
        return 'bg-yellow-500'
      case 'failed':
        return 'bg-red-500'
      default:
        return 'bg-gray-500'
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="spinner mb-4" />
          <p className="text-muted-foreground">Loading dashboard...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b border-border bg-card">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold">Dashboard</h1>
              <p className="text-muted-foreground mt-1">
                Welcome back, {user?.firstName || 'User'}
              </p>
            </div>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                <span className="text-sm text-muted-foreground">
                  {isConnected ? 'Connected' : 'Disconnected'}
                </span>
              </div>
              <Button variant="outline" size="sm">
                <Settings className="h-4 w-4 mr-2" />
                Settings
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Active Tasks</CardTitle>
                <Zap className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.activeTasks}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.completedTasks} completed this month
                </p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Portfolio Value</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatCurrency(stats.totalValue)}</div>
                <p className="text-xs text-green-600">
                  +{stats.portfolioChange}% from last month
                </p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">AI Interactions</CardTitle>
                <Bot className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.aiInteractions}</div>
                <p className="text-xs text-muted-foreground">
                  +12% from last week
                </p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.4 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Browser Sessions</CardTitle>
                <Globe className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.totalSessions}</div>
                <p className="text-xs text-muted-foreground">
                  3 active sessions
                </p>
              </CardContent>
            </Card>
          </motion.div>
        </div>

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Column - Main Content */}
          <div className="lg:col-span-2 space-y-6">
            <Tabs defaultValue="overview" className="w-full">
              <TabsList className="grid w-full grid-cols-4">
                <TabsTrigger value="overview">Overview</TabsTrigger>
                <TabsTrigger value="ai">AI Agent</TabsTrigger>
                <TabsTrigger value="browser">Browser</TabsTrigger>
                <TabsTrigger value="web3">Web3</TabsTrigger>
              </TabsList>

              <TabsContent value="overview" className="space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Performance Overview</CardTitle>
                    <CardDescription>
                      Your automation and portfolio performance
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium">Task Success Rate</span>
                          <span className="text-sm text-muted-foreground">94%</span>
                        </div>
                        <Progress value={94} className="h-2" />
                      </div>
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium">Browser Automation</span>
                          <span className="text-sm text-muted-foreground">87%</span>
                        </div>
                        <Progress value={87} className="h-2" />
                      </div>
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium">Web3 Transactions</span>
                          <span className="text-sm text-muted-foreground">98%</span>
                        </div>
                        <Progress value={98} className="h-2" />
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Quick Actions</CardTitle>
                    <CardDescription>
                      Common tasks and shortcuts
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-2 gap-4">
                      <Button className="h-20 flex-col gap-2">
                        <MessageSquare className="h-6 w-6" />
                        <span>New AI Chat</span>
                      </Button>
                      <Button variant="outline" className="h-20 flex-col gap-2">
                        <Globe className="h-6 w-6" />
                        <span>Browser Session</span>
                      </Button>
                      <Button variant="outline" className="h-20 flex-col gap-2">
                        <Wallet className="h-6 w-6" />
                        <span>Connect Wallet</span>
                      </Button>
                      <Button variant="outline" className="h-20 flex-col gap-2">
                        <BarChart3 className="h-6 w-6" />
                        <span>View Analytics</span>
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="ai" className="space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>AI Agent Status</CardTitle>
                    <CardDescription>
                      Current AI tasks and capabilities
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <Bot className="h-8 w-8 text-primary" />
                          <div>
                            <p className="font-medium">GPT-4 Turbo</p>
                            <p className="text-sm text-muted-foreground">Primary AI Model</p>
                          </div>
                        </div>
                        <Badge variant="outline">Active</Badge>
                      </div>
                      <div className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <Eye className="h-8 w-8 text-secondary" />
                          <div>
                            <p className="font-medium">Vision Analysis</p>
                            <p className="text-sm text-muted-foreground">Screenshot understanding</p>
                          </div>
                        </div>
                        <Badge variant="secondary">Available</Badge>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="browser" className="space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Browser Sessions</CardTitle>
                    <CardDescription>
                      Active and recent browser automation sessions
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {[1, 2, 3].map((session) => (
                        <div key={session} className="flex items-center justify-between p-4 border rounded-lg">
                          <div className="flex items-center gap-3">
                            <Globe className="h-8 w-8 text-blue-500" />
                            <div>
                              <p className="font-medium">Session {session}</p>
                              <p className="text-sm text-muted-foreground">
                                {session === 1 ? 'google.com' : session === 2 ? 'github.com' : 'example.com'}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            <Badge variant={session === 1 ? "default" : "secondary"}>
                              {session === 1 ? 'Active' : 'Idle'}
                            </Badge>
                            <Button variant="ghost" size="sm">
                              <Eye className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="web3" className="space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Web3 Portfolio</CardTitle>
                    <CardDescription>
                      Your cryptocurrency and DeFi positions
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-bold">
                            ETH
                          </div>
                          <div>
                            <p className="font-medium">Ethereum</p>
                            <p className="text-sm text-muted-foreground">2.5 ETH</p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatCurrency(6250)}</p>
                          <p className="text-sm text-green-600">+5.2%</p>
                        </div>
                      </div>
                      <div className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center text-white text-sm font-bold">
                            USDC
                          </div>
                          <div>
                            <p className="font-medium">USD Coin</p>
                            <p className="text-sm text-muted-foreground">5,000 USDC</p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatCurrency(5000)}</p>
                          <p className="text-sm text-muted-foreground">0.0%</p>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>

          {/* Right Column - Activity Feed */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Recent Activity</CardTitle>
                <CardDescription>
                  Latest actions and updates
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {recentActivity.map((activity) => (
                    <motion.div
                      key={activity.id}
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      className="flex items-start gap-3 p-3 rounded-lg border"
                    >
                      <div className="flex-shrink-0 mt-1">
                        {getActivityIcon(activity.type)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-sm">{activity.title}</p>
                        <p className="text-xs text-muted-foreground mb-2">
                          {activity.description}
                        </p>
                        <div className="flex items-center gap-2">
                          <div className={`w-2 h-2 rounded-full ${getStatusColor(activity.status)}`} />
                          <span className="text-xs text-muted-foreground">
                            {formatRelativeTime(activity.timestamp)}
                          </span>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>System Status</CardTitle>
                <CardDescription>
                  Service health and performance
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm">AI Agent</span>
                    <Badge variant="outline" className="text-green-600 border-green-600">
                      Operational
                    </Badge>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm">Browser Service</span>
                    <Badge variant="outline" className="text-green-600 border-green-600">
                      Operational
                    </Badge>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm">Web3 Service</span>
                    <Badge variant="outline" className="text-green-600 border-green-600">
                      Operational
                    </Badge>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm">API Gateway</span>
                    <Badge variant="outline" className="text-yellow-600 border-yellow-600">
                      Degraded
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
